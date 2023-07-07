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

	_ "github.com/lib/pq"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
)

type RedshiftApiClient struct {
	Username     string
	Password     string
	DatabaseName string
	Host         string
}

type redshiftIterator struct {
	queryResult *sql.Rows
	schema      data.Schema
}

func (it *redshiftIterator) Next(_ context.Context) (data.Row, error) {
	if it.queryResult.Next() {
		numColumns := len(it.schema)
		values := make([]any, numColumns)
		valuePtrs := make([]any, numColumns)
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}

		err := it.queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.redshiftIterator.Next)")
		}

		return convertRedshiftRow(values, it.schema), nil
	}

	defer it.queryResult.Close()
	err := it.queryResult.Err()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.redshiftIterator.Next) iterating over query results")
	}

	return nil, data.ErrDone
}

// TODO: this must be in order
func (it *redshiftIterator) Schema() data.Schema {
	return it.schema
}

func (rc RedshiftApiClient) openConnection(ctx context.Context) (*sql.DB, error) {
	params := url.Values{}
	params.Add("sslmode", "require")
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(rc.Username, rc.Password),
		Host:     rc.Host,
		Path:     rc.DatabaseName,
		RawQuery: params.Encode(),
	}

	return sql.Open("postgres", dsn.String())
}

func (rc RedshiftApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT(tablename) FROM pg_table_def WHERE schemaname = '%s'", namespace)

	queryResult, err := rc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetTables) running query")
	}

	var tableNames []string
	for _, row := range queryResult.Data {
		tableNames = append(tableNames, row[0].(string))
	}

	return tableNames, nil
}

func (rc RedshiftApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	queryString := fmt.Sprintf("SELECT pg_table_def.column, pg_table_def.type FROM pg_table_def WHERE schemaname = '%s' AND tablename = '%s'", namespace, tableName)

	queryResult, err := rc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrapf(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetSchema) getting schema for %s.%s", namespace, tableName)
	}

	schema := data.Schema{}
	for _, row := range queryResult.Data {
		dataType := getRedshiftFieldType(row[1].(string))
		schema = append(schema, data.Field{Name: row[0].(string), Type: dataType})
	}

	return schema, nil
}

func (rc RedshiftApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT %s FROM %s.%s LIMIT 100", fieldName, namespace, tableName)

	queryResult, err := rc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetFieldValues) running query")
	}

	values := []any{}
	for _, row := range queryResult.Data {
		values = append(values, row[0])
	}

	return values, nil
}

func (rc RedshiftApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	queryString := "SELECT nspname FROM pg_namespace WHERE nspname NOT IN ('pg_toast', 'pg_internal', 'catalog_history', 'pg_automv', 'pg_temp_1', 'pg_catalog', 'information_schema')"
	queryResult, err := rc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetNamespaces) running query")
	}

	var namespaces []string
	for _, row := range queryResult.Data {
		namespaces = append(namespaces, row[0].(string))
	}

	return namespaces, nil
}

func (rc RedshiftApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := rc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.RunQuery) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.RunQuery) running query")
	}
	defer queryResult.Close()

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.RunQuery) getting column types")
	}
	numColumns := len(columns)
	schema := convertRedshiftSchema(columns)

	var rows []data.Row
	values := make([]any, numColumns)
	valuePtrs := make([]any, numColumns)
	for queryResult.Next() {
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}
		err := queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.RunQuery) scanning row")
		}

		rows = append(rows, convertRedshiftRow(values, schema))
	}

	return &data.QueryResults{
		Schema: schema,
		Data:   rows,
	}, nil
}

func (rc RedshiftApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := rc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetQueryIterator) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetQueryIterator) running query")
	}

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.RedshiftApiClient.GetQueryIterator) getting column types")
	}

	return &redshiftIterator{
		queryResult: queryResult,
		schema:      convertRedshiftSchema(columns),
	}, nil
}

func getRedshiftFieldType(redshiftType string) data.FieldType {
	uppercased := strings.ToUpper(redshiftType)
	switch uppercased {
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
		// Objects from Redshift will have an empty type
		return data.FieldTypeJson
	default:
		// Everything can always be treated as a string
		return data.FieldTypeString
	}
}

func convertRedshiftRow(redshiftRow []any, schema data.Schema) data.Row {
	row := make(data.Row, len(redshiftRow))
	for i, value := range redshiftRow {
		row[i] = convertRedshiftValue(value, schema[i].Type)
	}

	return row
}

func convertRedshiftValue(redshiftValue any, fieldType data.FieldType) any {
	// Don't try to convert value that is nil
	if redshiftValue == nil {
		return nil
	}

	switch fieldType {
	case data.FieldTypeDateTimeTz:
		return redshiftValue.(time.Time).Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case data.FieldTypeDateTimeNtz:
		return redshiftValue.(time.Time).Format(FABRA_TIMESTAMP_NTZ_FORMAT)
	case data.FieldTypeString:
		// Redshift strings are sometimes returned as uint8 slices
		if v, ok := redshiftValue.([]uint8); ok {
			return string([]byte(v))
		}
		return string([]byte(redshiftValue.(string)))
	case data.FieldTypeJson:
		var strValue string
		if v, ok := redshiftValue.([]uint8); ok {
			strValue = string([]byte(v))
		} else {
			strValue = string([]byte(redshiftValue.(string)))
		}

		// TODO: handle error
		unquoted, _ := strconv.Unquote(strValue)
		jsonValue := map[string]any{}
		json.Unmarshal([]byte(unquoted), &jsonValue)
		return jsonValue
	default:
		return redshiftValue
	}
}

func convertRedshiftSchema(columns []*sql.ColumnType) data.Schema {
	schema := data.Schema{}

	for _, column := range columns {
		field := data.Field{
			Name: column.Name(),
			Type: getRedshiftFieldType(column.DatabaseTypeName()),
		}

		schema = append(schema, field)
	}

	return schema
}
