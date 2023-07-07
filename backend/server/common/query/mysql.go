package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
)

type MySqlApiClient struct {
	Username     string
	Password     string
	DatabaseName string
	Host         string
}

type mysqlIterator struct {
	queryResult *sql.Rows
	schema      data.Schema
}

func (it *mysqlIterator) Next(_ context.Context) (data.Row, error) {
	if it.queryResult.Next() {
		numColumns := len(it.schema)
		values := make([]any, numColumns)
		valuePtrs := make([]any, numColumns)
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}

		err := it.queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.mysqlIterator.Next)")
		}

		return convertMySqlRow(values, it.schema), nil
	}

	defer it.queryResult.Close()
	err := it.queryResult.Err()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.mysqlIterator.Next) iterating over query results")
	}

	return nil, data.ErrDone
}

// TODO: this must be in order
func (it *mysqlIterator) Schema() data.Schema {
	return it.schema
}

func (mc MySqlApiClient) openConnection(ctx context.Context) (*sql.DB, error) {
	params := url.Values{}
	params.Add("tls", "true")

	// Can't use url.Url because mysql does not accept a scheme
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mc.Username, mc.Password, mc.Host, mc.DatabaseName, params.Encode())

	return sql.Open("mysql", dsn)
}

func (mc MySqlApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT(table_name) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '%s'", namespace)

	queryResult, err := mc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetTables) error running query")
	}

	var tableNames []string
	for _, row := range queryResult.Data {
		tableNames = append(tableNames, row[0].(string))
	}

	return tableNames, nil
}

func (mc MySqlApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	queryString := fmt.Sprintf("SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '%s' AND table_name = '%s'", namespace, tableName)

	queryResult, err := mc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrapf(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetTables) getting schema for %s.%s", namespace, tableName)
	}

	schema := data.Schema{}
	for _, row := range queryResult.Data {
		dataType := getMySqlFieldType(row[1].(string))
		schema = append(schema, data.Field{Name: row[0].(string), Type: dataType})
	}

	return schema, nil
}

func (mc MySqlApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT %s FROM %s.%s LIMIT 100", fieldName, namespace, tableName)

	queryResult, err := mc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetFieldValues) error running query")
	}

	values := []any{}
	for _, row := range queryResult.Data {
		values = append(values, row[0])
	}

	return values, nil
}

func (mc MySqlApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	queryString := "SELECT schema_name FROM INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('_vt', 'sys', 'performance_schema', 'information_schema')"
	queryResult, err := mc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetNamespaces) error running query")
	}

	var namespaces []string
	for _, row := range queryResult.Data {
		namespaces = append(namespaces, row[0].(string))
	}

	return namespaces, nil
}

func (mc MySqlApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.RunQuery) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.RunQuery) running query")
	}
	defer queryResult.Close()

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.RunQuery) getting column types")
	}
	numColumns := len(columns)
	schema := convertMySqlSchema(columns)

	var rows []data.Row
	values := make([]any, numColumns)
	valuePtrs := make([]any, numColumns)
	for queryResult.Next() {
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}
		err := queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.RunQuery) scanning row")
		}

		rows = append(rows, convertMySqlRow(values, schema))
	}

	return &data.QueryResults{
		Schema: schema,
		Data:   rows,
	}, nil
}

func (mc MySqlApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetQueryIterator) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetQueryIterator) running query")
	}

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MySqlApiClient.GetQueryIterator) getting column types")
	}

	return &mysqlIterator{
		queryResult: queryResult,
		schema:      convertMySqlSchema(columns),
	}, nil
}

func getMySqlFieldType(mysqlType string) data.FieldType {
	uppemcased := strings.ToUpper(mysqlType)
	switch uppemcased {
	case "BOOL", "BOOLEAN":
		return data.FieldTypeBoolean
	case "INT", "INT2", "INT4", "INT8", "BIGINT":
		return data.FieldTypeInteger
	case "FLOAT", "FLOAT4", "FLOAT8", "NUMERIC", "DOUBLE":
		return data.FieldTypeNumber
	case "DATE":
		return data.FieldTypeDate
	case "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE":
		return data.FieldTypeDateTimeTz
	case "TIMESTAMP", "TIMESTAMP WITHOUT TIME ZONE":
		return data.FieldTypeDateTimeNtz
	case "TIME", "TIME WITHOUT TIME ZONE":
		return data.FieldTypeTimeNtz
	case "TIMETZ", "TIME WITH TIME ZONE":
		return data.FieldTypeTimeTz
	case "":
		// Objects from MySql will have an empty type
		return data.FieldTypeJson
	default:
		// Everything can always be treated as a string
		return data.FieldTypeString
	}
}

func convertMySqlRow(mysqlRow []any, schema data.Schema) data.Row {
	row := make(data.Row, len(mysqlRow))
	for i, value := range mysqlRow {
		row[i] = convertMySqlValue(value, schema[i].Type)
	}

	return row
}

func convertMySqlValue(mysqlValue any, fieldType data.FieldType) any {
	// Don't try to convert value that is nil
	if mysqlValue == nil {
		return nil
	}

	switch fieldType {
	case data.FieldTypeDateTimeTz:
		return mysqlValue.(time.Time).Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case data.FieldTypeDateTimeNtz:
		return mysqlValue.(time.Time).Format(FABRA_TIMESTAMP_NTZ_FORMAT)
	case data.FieldTypeString:
		// MySQL strings are sometimes returned as uint8 slices
		if v, ok := mysqlValue.([]uint8); ok {
			return string([]byte(v))
		}
		return string([]byte(mysqlValue.(string)))
	case data.FieldTypeJson:
		var strValue string
		if v, ok := mysqlValue.([]uint8); ok {
			strValue = string([]byte(v))
		} else {
			strValue = string([]byte(mysqlValue.(string)))
		}

		// TODO: handle error
		unquoted, _ := strconv.Unquote(strValue)
		jsonValue := map[string]any{}
		json.Unmarshal([]byte(unquoted), &jsonValue)
		return jsonValue
	default:
		return mysqlValue
	}
}

func convertMySqlSchema(columns []*sql.ColumnType) data.Schema {
	schema := data.Schema{}

	for _, column := range columns {
		field := data.Field{
			Name: column.Name(),
			Type: getMySqlFieldType(column.DatabaseTypeName()),
		}

		schema = append(schema, field)
	}

	return schema
}
