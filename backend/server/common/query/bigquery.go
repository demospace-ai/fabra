package query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/storage"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type BigQueryApiClient struct {
	ProjectID   *string
	Credentials *string
	Location    *string
}

type bigQueryIterator struct {
	iterator *bigquery.RowIterator
	schema   data.Schema
}

// use a pointer receiver because schema field is updated
func (it *bigQueryIterator) Next(_ context.Context) (data.Row, error) {
	var row []bigquery.Value
	err := it.iterator.Next(&row)
	if err == iterator.Done {
		return nil, data.ErrDone
	}

	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.bigQueryIterator.Next)")
	}

	return convertBigQueryRow(row, it.iterator.Schema), nil
}

func (it *bigQueryIterator) Schema() data.Schema {
	if it.schema == nil {
		// TODO: this must be in order
		it.schema = convertBigQuerySchema(it.iterator.Schema)
	}

	return it.schema
}

func (ac BigQueryApiClient) openConnection(ctx context.Context) (*bigquery.Client, error) {
	if ac.ProjectID == nil {
		return nil, errors.Newf("missing project ID")
	}

	var credentialOption option.ClientOption
	if ac.Credentials != nil {
		credentialOption = option.WithCredentialsJSON([]byte(*ac.Credentials))
	}

	return bigquery.NewClient(ctx, *ac.ProjectID, credentialOption)
}

func (ac BigQueryApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	client, err := ac.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetTables)")
	}

	defer client.Close()

	ts := client.Dataset(namespace).Tables(ctx)
	var results []string
	for {
		table, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetTables) getting tables for namespace %s", namespace)
		}

		results = append(results, table.TableID)
	}

	return results, nil
}

func (ac BigQueryApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	queryString := fmt.Sprintf("SELECT * FROM %s.INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s'", namespace, tableName)

	queryResults, err := ac.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetSchema)")
	}

	schema := data.Schema{}
	for _, row := range queryResults.Data {
		if row[0] == nil {
			continue
		}

		schema = append(schema, data.Field{Name: row[3].(string), Type: getBigQueryFieldType(row[6].(string))})
	}

	return schema, nil
}

func (ac BigQueryApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	queryString := fmt.Sprintf("SELECT DISTINCT %s FROM %s.%s LIMIT 100", fieldName, namespace, tableName)

	queryResults, err := ac.RunQuery(ctx, queryString)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetFieldValues)")
	}

	values := []any{}
	for _, row := range queryResults.Data {
		if row[0] == nil {
			continue
		}

		values = append(values, row[0])
	}

	return values, nil
}

func (ac BigQueryApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	client, err := ac.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetNamespaces)")
	}

	defer client.Close()

	ts := client.Datasets(ctx)
	var results []string
	for {
		dataset, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetNamespaces) iterating datasets")
		}

		results = append(results, dataset.DatasetID)
	}

	return results, nil
}

func (ac BigQueryApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := ac.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.RunQuery) opening connection")
	}
	defer client.Close()

	q := client.Query(queryString)
	for arg := range args {
		q.Parameters = append(q.Parameters, bigquery.QueryParameter{Value: arg})
	}

	// Location must match that of the dataset(s) referenced in the query.
	q.Location = *ac.Location

	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.RunQuery) running query")
	}

	// If an error happens here it isn't actually a failure, the query was just wrong. Send the details back.
	// TODO: make special error type for this
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.RunQuery) waiting for query to complete")
	}
	if err := status.Err(); err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.RunQuery) query failed")
	}

	it, err := job.Read(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.RunQuery) reading query results")
	}

	var results []data.Row
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.RunQuery) iterating query results")
		}
		results = append(results, convertBigQueryRow(row, it.Schema))
	}

	return &data.QueryResults{
		Schema: convertBigQuerySchema(it.Schema),
		Data:   results,
	}, nil
}

func (ac BigQueryApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := ac.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetQueryIterator) opening connection")
	}
	defer client.Close()

	q := client.Query(queryString)

	// Location must match that of the dataset(s) referenced in the query.
	q.Location = *ac.Location

	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetQueryIterator) running query")
	}

	// Both of these are not actually a failure, the query was just wrong. Send the details back to them.
	_, err = job.Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetQueryIterator) waiting for query to complete")
	}

	it, err := job.Read(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.GetQueryIterator) reading query results")
	}

	return &bigQueryIterator{
		iterator: it,
	}, nil
}

func (ac BigQueryApiClient) StageData(ctx context.Context, csvData string, stagingOptions StagingOptions) error {
	var credentialOption option.ClientOption
	if ac.Credentials != nil {
		credentialOption = option.WithCredentialsJSON([]byte(*ac.Credentials))
	}

	gcsClient, err := storage.NewClient(ctx, credentialOption)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.StageData) creating client")
	}
	defer gcsClient.Close()

	w := gcsClient.Bucket(stagingOptions.Bucket).Object(stagingOptions.Object).Retryer(
		storage.WithPolicy(storage.RetryAlways),
	).NewWriter(ctx)
	if _, err := fmt.Fprint(w, csvData); err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.StageData) writing data")
	}

	if err := w.Close(); err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.StageData) closing writer")
	}

	return nil
}

func (ac BigQueryApiClient) LoadFromStaging(ctx context.Context, namespace string, tableName string, loadOptions LoadOptions) error {
	client, err := ac.openConnection(ctx)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.LoadFromStaging) opening connection")
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference(loadOptions.GcsReference)
	gcsRef.SourceFormat = bigquery.CSV
	gcsRef.Schema = loadOptions.BigQuerySchema

	loader := client.Dataset(namespace).Table(tableName).LoaderFrom(gcsRef)
	loader.WriteDisposition = loadOptions.WriteMode

	job, err := loader.Run(ctx)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.LoadFromStaging) running job")
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.LoadFromStaging) waiting for job")
	}

	if status.Err() != nil {
		errors.Wrap(errors.WrapCustomerVisibleError(status.Err()), "(query.BigQueryApiClient.LoadFromStaging) status error")
	}

	return nil
}

func (ac BigQueryApiClient) CleanUpStagingData(ctx context.Context, stagingOptions StagingOptions) error {
	var credentialOption option.ClientOption
	if ac.Credentials != nil {
		credentialOption = option.WithCredentialsJSON([]byte(*ac.Credentials))
	}

	gcsClient, err := storage.NewClient(ctx, credentialOption)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.CleanUpStagingData)")
	}
	defer gcsClient.Close()

	object := gcsClient.Bucket(stagingOptions.Bucket).Object(stagingOptions.Object).Retryer(
		storage.WithPolicy(storage.RetryAlways),
	)

	err = object.Delete(ctx)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.BigQueryApiClient.CleanUpStagingData) deleting object")
	}

	return nil
}

func convertBigQueryRow(bigQueryRow []bigquery.Value, schema bigquery.Schema) data.Row {
	row := make(data.Row, len(bigQueryRow))
	for i, value := range bigQueryRow {
		row[i] = convertBigQueryValue(value, schema[i].Type)
	}

	return row
}

func convertBigQueryValue(bigqueryValue any, fieldType bigquery.FieldType) any {
	if bigqueryValue == nil {
		return nil
	}

	switch fieldType {
	case bigquery.TimestampFieldType:
		return bigqueryValue.(time.Time).Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case bigquery.DateFieldType:
		return bigqueryValue.(civil.Date).String()
	case bigquery.TimeFieldType:
		return bigqueryValue.(civil.Time).String()
	case bigquery.DateTimeFieldType:
		return bigqueryValue.(civil.DateTime).String()
	case bigquery.JSONFieldType:
		jsonValue := map[string]any{}
		json.Unmarshal([]byte(bigqueryValue.(string)), &jsonValue)
		return jsonValue
	default:
		return bigqueryValue
	}
}

func getBigQueryFieldType(bigQueryType string) data.FieldType {
	uppercased := strings.ToUpper(bigQueryType)
	switch uppercased {
	case "INTEGER", "INT64":
		return data.FieldTypeInteger
	case "FLOAT", "NUMERIC", "BIGNUMERIC":
		return data.FieldTypeNumber
	case "BOOLEAN":
		return data.FieldTypeBoolean
	case "TIMESTAMP":
		// BigQuery timestamps are actually datetimes with optional timezone information
		return data.FieldTypeDateTimeTz
	case "JSON":
		return data.FieldTypeJson
	case "DATE":
		return data.FieldTypeDate
	case "TIME":
		// BigQuery times do not contain timezone information
		return data.FieldTypeTimeNtz
	case "DATETIME":
		// BigQuery datetimes do not contain timezone information
		return data.FieldTypeDateTimeNtz
	default:
		return data.FieldTypeString
	}
}

func convertBigQuerySchema(bigQuerySchema bigquery.Schema) data.Schema {
	schema := data.Schema{}

	for _, bigQuerySchemaField := range bigQuerySchema {
		field := data.Field{
			Name: bigQuerySchemaField.Name,
			Type: getBigQueryFieldType(string(bigQuerySchemaField.Type)),
		}

		schema = append(schema, field)
	}

	return schema
}
