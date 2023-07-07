package query

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
)

type SynapseApiClient struct {
	Username     string
	Password     string
	DatabaseName string
	Host         string
}

type synapseIterator struct {
	queryResult *sql.Rows
	schema      data.Schema
}

func (it *synapseIterator) Next(_ context.Context) (data.Row, error) {
	if it.queryResult.Next() {
		numColumns := len(it.schema)
		values := make([]any, numColumns)
		valuePtrs := make([]any, numColumns)
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}

		err := it.queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.synapseIterator.Next)")
		}
		return convertSynapseRow(values, it.schema), nil
	}

	defer it.queryResult.Close()
	err := it.queryResult.Err()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.synapseIterator.Next) iterating over query results")
	}

	return nil, data.ErrDone
}

// TODO: this must be in order
func (it *synapseIterator) Schema() data.Schema {
	return it.schema
}

func (sc SynapseApiClient) openConnection(ctx context.Context) (*sql.DB, error) {
	params := url.Values{}
	params.Add("database", sc.DatabaseName)
	params.Add("sslmode", "encrypt")
	params.Add("TrustServerCertificate", "true")
	dsn := url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(sc.Username, sc.Password),
		Host:     sc.Host,
		RawQuery: params.Encode(),
	}

	return sql.Open("sqlserver", dsn.String())
}

func (sc SynapseApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	queryString := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_SCHEMA = '%s'", namespace)

	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetTables)")
	}

	var tableNames []string
	for _, row := range queryResult.Data {
		tableNames = append(tableNames, row[0].(string))
	}

	return tableNames, nil
}

func (sc SynapseApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	queryString := fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND table_name = '%s'", namespace, tableName)

	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrapf(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetSchema) getting schema for %s.%s", namespace, tableName)
	}

	schema := data.Schema{}
	for _, row := range queryResult.Data {
		dataType := getSynapseFieldType(row[1].(string))
		schema = append(schema, data.Field{Name: row[0].(string), Type: dataType})
	}

	return schema, nil
}

func (sc SynapseApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT %s FROM %s.%s LIMIT 100", fieldName, namespace, tableName)

	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetFieldValues)")
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

func (sc SynapseApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	queryString := "SELECT TABLE_SCHEMA FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'"
	queryResult, err := sc.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetNamespaces)")
	}

	var namespaces []string
	for _, row := range queryResult.Data {
		if row[0] == nil {
			continue
		}

		namespaces = append(namespaces, row[0].(string))
	}

	return namespaces, nil
}

func (sc SynapseApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := sc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.RunQuery) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.RunQuery) running query")
	}
	defer queryResult.Close()

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.RunQuery) getting column types")
	}
	numColumns := len(columns)
	schema := convertSynapseSchema(columns)

	var rows []data.Row
	values := make([]any, numColumns)
	valuePtrs := make([]any, numColumns)
	for queryResult.Next() {
		for i := 0; i < numColumns; i++ {
			valuePtrs[i] = &values[i]
		}
		err := queryResult.Scan(valuePtrs...)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.RunQuery) scanning row")
		}

		rows = append(rows, convertSynapseRow(values, schema))
	}

	return &data.QueryResults{
		Schema: schema,
		Data:   rows,
	}, nil
}

func (sc SynapseApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := sc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetQueryIterator) opening connection")
	}
	defer client.Close()

	queryResult, err := client.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetQueryIterator) running query")
	}

	columns, err := queryResult.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.SynapseApiClient.GetQueryIterator) getting column types")
	}

	return &synapseIterator{
		queryResult: queryResult,
		schema:      convertSynapseSchema(columns),
	}, nil
}

func convertSynapseRow(synapseRow []any, schema data.Schema) data.Row {
	row := make(data.Row, len(synapseRow))
	for i, value := range synapseRow {
		row[i] = convertSynapseValue(value, schema[i].Type)
	}

	return row
}

func convertSynapseValue(synapseValue any, fieldType data.FieldType) any {
	// Don't try to convert value that is nil
	if synapseValue == nil {
		return nil
	}

	switch fieldType {
	case data.FieldTypeDateTimeTz:
		return synapseValue.(time.Time).Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case data.FieldTypeDateTimeNtz:
		return synapseValue.(time.Time).Format(FABRA_TIMESTAMP_NTZ_FORMAT)
	case data.FieldTypeString:
		return synapseValue.(string)
	default:
		return synapseValue
	}
}

func getSynapseFieldType(synapseType string) data.FieldType {
	uppercased := strings.ToUpper(synapseType)
	switch uppercased {
	case "BIT":
		return data.FieldTypeBoolean
	case "INT", "BIGINT", "SMALLINT", "TINYINT":
		return data.FieldTypeInteger
	case "REAL", "DECIMAL", "NUMERIC", "FLOAT":
		return data.FieldTypeNumber
	case "DATE":
		return data.FieldTypeDate
	case "TIME":
		return data.FieldTypeTimeNtz
	case "DATETIME", "DATETIME2", "SMALLDATETIME":
		return data.FieldTypeDateTimeNtz
	case "DATETIMEOFFSET":
		return data.FieldTypeDateTimeTz
	default:
		// Everything can always be treated as a string
		return data.FieldTypeString
	}
}

func convertSynapseSchema(columns []*sql.ColumnType) data.Schema {
	schema := data.Schema{}

	for _, column := range columns {
		field := data.Field{
			Name: column.Name(),
			Type: getSynapseFieldType(column.DatabaseTypeName()),
		}

		schema = append(schema, field)
	}

	return schema
}
