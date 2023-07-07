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

type PostgresApiClient struct {
	Username     string
	Password     string
	DatabaseName string
	Host         string
}

type postgresIterator struct {
	queryResult *sql.Rows
	schema      data.Schema
}

func (it *postgresIterator) Next(_ context.Context) (data.Row, error) {
	if it.queryResult.Next() {
		numColumns := len(it.schema)
		values := make([]any, numColumns)
		valuePtrs := make([]any, numColumns)
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}

		err := it.queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.postgresIterator.Next)")
		}

		return convertPostgresRow(values, it.schema), nil
	}

	defer it.queryResult.Close()
	err := it.queryResult.Err()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.postgresIterator.Next) iterating over query results")
	}

	return nil, data.ErrDone
}

// TODO: this must be in order
func (it *postgresIterator) Schema() data.Schema {
	return it.schema
}

func (pc PostgresApiClient) openConnection(ctx context.Context) (*sql.DB, error) {
	params := url.Values{}
	params.Add("sslmode", "require")
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(pc.Username, pc.Password),
		Host:     pc.Host,
		Path:     pc.DatabaseName,
		RawQuery: params.Encode(),
	}

	return sql.Open("postgres", dsn.String())
}

func (pc PostgresApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT(table_name) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '%s'", namespace)

	queryResult, err := pc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetTables) error running query")
	}

	var tableNames []string
	for _, row := range queryResult.Data {
		tableNames = append(tableNames, row[0].(string))
	}

	return tableNames, nil
}

func (pc PostgresApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	queryString := fmt.Sprintf("SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '%s' AND table_name = '%s'", namespace, tableName)

	queryResult, err := pc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrapf(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetTables) getting schema for %s.%s", namespace, tableName)
	}

	schema := data.Schema{}
	for _, row := range queryResult.Data {
		dataType := getPostgresFieldType(row[1].(string))
		schema = append(schema, data.Field{Name: row[0].(string), Type: dataType})
	}

	return schema, nil
}

func (pc PostgresApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT %s FROM %s.%s LIMIT 100", fieldName, namespace, tableName)

	queryResult, err := pc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetFieldValues) error running query")
	}

	values := []any{}
	for _, row := range queryResult.Data {
		values = append(values, row[0])
	}

	return values, nil
}

func (pc PostgresApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	queryString := "SELECT nspname FROM pg_catalog.pg_namespace	WHERE nspname NOT IN ('pg_toast', 'pg_internal', 'catalog_history', 'pg_automv', 'pg_temp_1', 'pg_catalog', 'information_schema')"
	queryResult, err := pc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetNamespaces) error running query")
	}

	var namespaces []string
	for _, row := range queryResult.Data {
		namespaces = append(namespaces, row[0].(string))
	}

	return namespaces, nil
}

func (pc PostgresApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := pc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.RunQuery) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.RunQuery) running query")
	}
	defer queryResult.Close()

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.RunQuery) getting column types")
	}
	numColumns := len(columns)
	schema := convertPostgresSchema(columns)

	var rows []data.Row
	values := make([]any, numColumns)
	valuePtrs := make([]any, numColumns)
	for queryResult.Next() {
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}
		err := queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.RunQuery) scanning row")
		}

		rows = append(rows, convertPostgresRow(values, schema))
	}

	return &data.QueryResults{
		Schema: schema,
		Data:   rows,
	}, nil
}

func (pc PostgresApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := pc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetQueryIterator) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetQueryIterator) running query")
	}

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.PostgresApiClient.GetQueryIterator) getting column types")
	}

	return &postgresIterator{
		queryResult: queryResult,
		schema:      convertPostgresSchema(columns),
	}, nil
}

func getPostgresFieldType(postgresType string) data.FieldType {
	uppepcased := strings.ToUpper(postgresType)
	switch uppepcased {
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
		// Objects from Postgres will have an empty type
		return data.FieldTypeJson
	default:
		// Everything can always be treated as a string
		return data.FieldTypeString
	}
}

func convertPostgresRow(postgresRow []any, schema data.Schema) data.Row {
	row := make(data.Row, len(postgresRow))
	for i, value := range postgresRow {
		row[i] = convertPostgresValue(value, schema[i].Type)
	}

	return row
}

func convertPostgresValue(postgresValue any, fieldType data.FieldType) any {
	// Don't try to convert value that is nil
	if postgresValue == nil {
		return nil
	}

	switch fieldType {
	case data.FieldTypeDateTimeTz:
		return postgresValue.(time.Time).Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case data.FieldTypeDateTimeNtz:
		return postgresValue.(time.Time).Format(FABRA_TIMESTAMP_NTZ_FORMAT)
	case data.FieldTypeString:
		// Postgres strings are sometimes returned as uint8 slices
		if v, ok := postgresValue.([]uint8); ok {
			return string([]byte(v))
		}
		return string([]byte(postgresValue.(string)))
	case data.FieldTypeJson:
		var strValue string
		if v, ok := postgresValue.([]uint8); ok {
			strValue = string([]byte(v))
		} else {
			strValue = string([]byte(postgresValue.(string)))
		}

		// TODO: handle error
		unquoted, _ := strconv.Unquote(strValue)
		jsonValue := map[string]any{}
		json.Unmarshal([]byte(unquoted), &jsonValue)
		return jsonValue
	default:
		return postgresValue
	}
}

func convertPostgresSchema(columns []*sql.ColumnType) data.Schema {
	schema := data.Schema{}

	for _, column := range columns {
		field := data.Field{
			Name: column.Name(),
			Type: getPostgresFieldType(column.DatabaseTypeName()),
		}

		schema = append(schema, field)
	}

	return schema
}
