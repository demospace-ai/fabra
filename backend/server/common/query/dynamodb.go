package query

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
)

type DynamoDbClient struct {
	KeyID     string
	AccessKey string
	Location  string
}

// TODO
// type dynamodbIterator struct {
// 	queryResult sql.Rows
// 	schema      data.Schema
// }

// func (it *dynamodbIterator) Next(_ context.Context) (data.Row, error) {
// 	// TODO
// 	return nil, errors.New("not implemented")
// }

// // TODO: this must be in order
// func (it *dynamodbIterator) Schema() data.Schema {
// 	return it.schema
// }

func (dc DynamoDbClient) openConnection(ctx context.Context) (*dynamodb.Client, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(dc.KeyID, dc.AccessKey, "")), // empty token is ok
		config.WithRegion(dc.Location),
	)

	if err != nil {
		return nil, errors.Wrap(err, "(query.DynamoDbClient.openConnection) loading AWS config")
	}

	return dynamodb.NewFromConfig(awsConfig), nil
}

func (dc DynamoDbClient) GetTables(ctx context.Context, _namespace string) ([]string, error) {
	client, err := dc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.DynamoDbClient.GetTables) opening connection")
	}

	tables, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.DynamoDbClient.GetTables) list tables")
	}

	return tables.TableNames, nil
}

func (dc DynamoDbClient) LoadData(ctx context.Context, namespace string, tableName string, rows []data.Row) error {
	// TODO
	return errors.New("not implemented")
}

func (dc DynamoDbClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	// TODO
	return nil, errors.New("not implemented")
}

func (dc DynamoDbClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	// TODO
	return nil, errors.New("not implemented")
}

func (dc DynamoDbClient) GetNamespaces(ctx context.Context) ([]string, error) {
	return nil, errors.New("DynamoDB does not support namespaces")
}

func (dc DynamoDbClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	// TODO
	return nil, errors.New("not implemented")
}

func (dc DynamoDbClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	// TODO
	return nil, errors.New("not implemented")
}
