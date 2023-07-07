package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/snowflakedb/gosnowflake"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
)

const SNOWFLAKE_TZ_FORMAT = "2006-01-02T15:04:05.000-07:00"

type SnowflakeApiClient struct {
	Username      string
	Password      string
	WarehouseName string
	DatabaseName  string
	Role          string
	Host          string
}

type snowflakeSchema struct {
	Type string `json:"type"`
}

type snowflakeIterator struct {
	queryResult *sql.Rows
	schema      data.Schema
}

func (it *snowflakeIterator) Next(_ context.Context) (data.Row, error) {
	if it.queryResult.Next() {
		numColumns := len(it.schema)
		values := make([]any, numColumns)
		valuePtrs := make([]any, numColumns)
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}

		err := it.queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.snowflakeIterator.Next)")
		}

		return convertSnowflakeRow(values, it.schema), nil
	}

	defer it.queryResult.Close()
	err := it.queryResult.Err()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.snowflakeIterator.Next) iterating over query results")
	}

	return nil, data.ErrDone
}

// TODO: this must be in order
func (it *snowflakeIterator) Schema() data.Schema {
	return it.schema
}

func (sc SnowflakeApiClient) openConnection(ctx context.Context) (*sql.DB, error) {
	account := strings.Split(sc.Host, ".")[0] // TODO: remove the https/http
	config := gosnowflake.Config{
		Account:   account,
		User:      sc.Username,
		Password:  sc.Password,
		Warehouse: sc.WarehouseName,
		Database:  sc.DatabaseName,
		Role:      sc.Role,
		Host:      sc.Host,
		Port:      443,
	}

	dsn, err := gosnowflake.DSN(&config)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.openConnection)")
	}

	return sql.Open("snowflake", dsn)
}

func (sc SnowflakeApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	client, err := sc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetTables) opening connection")
	}

	defer client.Close()

	queryString := fmt.Sprintf("SHOW TERSE TABLES IN %s", namespace)
	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetTables) running query")
	}
	defer queryResult.Close()

	columns, err := queryResult.Columns()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetTables) getting columns")
	}
	numColumns := len(columns)

	// just scan into a string list, everything can be a string
	var tableNames []string
	values := make([]any, numColumns)
	valuePtrs := make([]any, numColumns)
	for queryResult.Next() {
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}
		err := queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetTables) scanning row")
		}

		tableNames = append(tableNames, values[1].(string))
	}

	return tableNames, nil
}

func (sc SnowflakeApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	queryString := fmt.Sprintf("SHOW COLUMNS IN %s.%s", namespace, tableName)

	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrapf(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetSchema) getting schema for %s.%s", namespace, tableName)
	}

	schema := data.Schema{}
	for _, row := range queryResult.Data {
		if row[0] == nil {
			continue
		}

		var snowflakeSchema snowflakeSchema
		err := json.Unmarshal([]byte(row[3].(string)), &snowflakeSchema)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetSchema) unmarshalling columns")
		}

		dataType := getSnowflakeFieldType(snowflakeSchema.Type)
		schema = append(schema, data.Field{Name: row[2].(string), Type: dataType})
	}

	return schema, nil
}

func (sc SnowflakeApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT %s FROM %s.%s LIMIT 100", fieldName, namespace, tableName)

	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetFieldValues)")
	}

	values := []any{}
	for _, row := range queryResult.Data {
		if row[0] == nil {
			continue
		}

		values = append(values, row[0])
	}

	return values, nil
}

func (sc SnowflakeApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	queryString := "SHOW TERSE SCHEMAS"
	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetNamespaces)")
	}

	var namespaces []string
	for _, row := range queryResult.Data {
		if row[0] == nil {
			continue
		}

		namespaces = append(namespaces, row[1].(string))
	}

	return namespaces, nil
}

func (sc SnowflakeApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := sc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.RunQuery) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.RunQuery) running query")
	}
	defer queryResult.Close()

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.RunQuery) getting columns")
	}
	numColumns := len(columns)
	schema := convertSnowflakeSchema(columns)

	var rows []data.Row
	values := make([]any, numColumns)
	valuePtrs := make([]any, numColumns)
	for queryResult.Next() {
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}
		err := queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.RunQuery) scanning row")
		}

		rows = append(rows, convertSnowflakeRow(values, schema))
	}

	return &data.QueryResults{
		Schema: schema,
		Data:   rows,
	}, nil
}

func (sc SnowflakeApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := sc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetQueryIterator) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetQueryIterator) running query")
	}

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SnowflakeApiClient.GetQueryIterator) getting columns")
	}

	return &snowflakeIterator{
		queryResult: queryResult,
		schema:      convertSnowflakeSchema(columns),
	}, nil
}

func convertSnowflakeRow(snowflakeRow []any, schema data.Schema) data.Row {
	row := make(data.Row, len(snowflakeRow))
	for i, value := range snowflakeRow {
		row[i] = convertSnowflakeValue(value, schema[i].Type)
	}

	return row
}

func convertSnowflakeValue(snowflakeValue any, fieldType data.FieldType) any {
	if snowflakeValue == nil {
		return nil
	}

	// TODO: convert remaining types to the expected Fabra Golang types
	switch fieldType {
	case data.FieldTypeJson:
		jsonValue := map[string]any{}
		json.Unmarshal([]byte(snowflakeValue.(string)), &jsonValue)
		return jsonValue
	case data.FieldTypeDateTimeTz:
		return snowflakeValue.(time.Time).Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case data.FieldTypeDateTimeNtz:
		return snowflakeValue.(time.Time).Format(FABRA_TIMESTAMP_NTZ_FORMAT)
	default:
		return snowflakeValue
	}
}

func getSnowflakeFieldType(snowflakeType string) data.FieldType {
	uppercased := strings.ToUpper(snowflakeType)
	switch uppercased {
	case "BIT", "BOOLEAN":
		return data.FieldTypeBoolean
	case "INTEGER", "BIGINT", "SMALLINT", "TINYINT":
		return data.FieldTypeInteger
	case "REAL", "DOUBLE", "DECIMAL", "NUMERIC", "FLOAT", "FIXED":
		return data.FieldTypeNumber
	case "TIMESTAMP_TZ":
		return data.FieldTypeDateTimeTz
	case "DATETIME", "TIMESTAMP", "TIMESTAMP_NTZ":
		return data.FieldTypeDateTimeNtz
	case "VARIANT":
		return data.FieldTypeJson
	default:
		// Everything can always be treated as a string
		return data.FieldTypeString
	}
}

func convertSnowflakeSchema(columns []*sql.ColumnType) data.Schema {
	schema := data.Schema{}

	for _, column := range columns {
		field := data.Field{
			Name: column.Name(),
			Type: getSnowflakeFieldType(column.DatabaseTypeName()),
		}

		schema = append(schema, field)
	}

	return schema
}
