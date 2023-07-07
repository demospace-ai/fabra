package query

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/bigquery"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
)

const FABRA_TIMESTAMP_TZ_FORMAT = "2006-01-02 15:04:05.000-07:00"
const FABRA_TIMESTAMP_NTZ_FORMAT = "2006-01-02 15:04:05.000"

type StagingOptions struct {
	Bucket string
	Object string
}

type LoadOptions struct {
	GcsReference   string
	BigQuerySchema bigquery.Schema
	WriteMode      bigquery.TableWriteDisposition
}

type QueryService interface {
	GetNamespaces(ctx context.Context, connection *models.Connection) ([]string, error)
	GetTables(ctx context.Context, connection *models.Connection, namespace string) ([]string, error)
	GetSchema(ctx context.Context, connection *models.Connection, namespace string, tableName string) ([]data.Field, error)
	GetFieldValues(ctx context.Context, connection *models.Connection, namespace string, tableName string, fieldName string) ([]any, error)
	RunQuery(ctx context.Context, connection *models.Connection, queryString string) (*data.QueryResults, error)
	GetQueryIterator(ctx context.Context, connection *models.Connection, queryString string) (data.RowIterator, error)
	GetClient(ctx context.Context, connection *models.Connection) (ConnectorClient, error)
	GetWarehouseClient(ctx context.Context, connection *models.Connection) (WarehouseClient, error)
}

type QueryServiceImpl struct {
	cryptoService crypto.CryptoService
}

func NewQueryService(cryptoService crypto.CryptoService) QueryService {
	return QueryServiceImpl{
		cryptoService: cryptoService,
	}
}

type ConnectorClient interface {
	GetTables(ctx context.Context, namespace string) ([]string, error)
	GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error)
	GetNamespaces(ctx context.Context) ([]string, error)
	GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error)
	RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error)
	GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error)
}

type WarehouseClient interface {
	ConnectorClient
	StageData(ctx context.Context, csvData string, stagingOptions StagingOptions) error
	LoadFromStaging(ctx context.Context, namespace string, tableName string, loadOptions LoadOptions) error
	CleanUpStagingData(ctx context.Context, stagingOptions StagingOptions) error
}

type DatabaseClient interface {
	ConnectorClient
	LoadData(ctx context.Context, namespace string, tableName string, rows []data.Row) error
}

func (qs QueryServiceImpl) GetClient(ctx context.Context, connection *models.Connection) (ConnectorClient, error) {
	switch connection.ConnectionType {
	case models.ConnectionTypeBigQuery:
		bigQueryCredentialsString, err := qs.cryptoService.DecryptConnectionCredentials(connection.Credentials.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting BigQuery credentials")
		}

		var bigQueryCredentials models.BigQueryCredentials
		err = json.Unmarshal([]byte(*bigQueryCredentialsString), &bigQueryCredentials)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) unmarshalling BigQuery credentials")
		}

		if !connection.Location.Valid {
			return nil, errors.NewCustomerVisibleError("BigQuery connection must have location defined")
		}

		return BigQueryApiClient{
			ProjectID:   &bigQueryCredentials.ProjectID,
			Credentials: bigQueryCredentialsString,
			Location:    &connection.Location.String,
		}, nil
	case models.ConnectionTypeDynamoDb:
		password, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting DynamoDB credentials")
		}

		return DynamoDbClient{
			KeyID:     connection.Username.String,
			AccessKey: *password,
			Location:  connection.Location.String,
		}, nil

	case models.ConnectionTypeSnowflake:
		snowflakePassword, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting Snowflake password")
		}

		// TODO: validate all connection params
		return SnowflakeApiClient{
			Username:      connection.Username.String,
			Password:      *snowflakePassword,
			WarehouseName: connection.WarehouseName.String,
			DatabaseName:  connection.DatabaseName.String,
			Role:          connection.Role.String,
			Host:          connection.Host.String,
		}, nil
	case models.ConnectionTypeRedshift:
		redshiftPassword, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting Redshift password")
		}

		// TODO: validate all connection params
		return RedshiftApiClient{
			Username:     connection.Username.String,
			Password:     *redshiftPassword,
			DatabaseName: connection.DatabaseName.String,
			Host:         connection.Host.String,
		}, nil
	case models.ConnectionTypeSynapse:
		synapsePassword, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting Synapse password")
		}

		// TODO: validate all connection params
		return SynapseApiClient{
			Username:     connection.Username.String,
			Password:     *synapsePassword,
			DatabaseName: connection.DatabaseName.String,
			Host:         connection.Host.String,
		}, nil
	case models.ConnectionTypeMongoDb:
		mongodbPassword, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting MongoDB password")
		}

		// TODO: validate all connection params
		return MongoDbApiClient{
			Username:          connection.Username.String,
			Password:          *mongodbPassword,
			Host:              connection.Host.String,
			ConnectionOptions: connection.ConnectionOptions.String,
		}, nil
	case models.ConnectionTypePostgres:
		postgresPassword, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting Postgres password")
		}

		// TODO: validate all connection params
		return PostgresApiClient{
			Username:     connection.Username.String,
			Password:     *postgresPassword,
			DatabaseName: connection.DatabaseName.String,
			Host:         connection.Host.String,
		}, nil
	case models.ConnectionTypeMySQL:
		mysqlPassword, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) decrypting MySQL password")
		}

		// TODO: validate all connection params
		return MySqlApiClient{
			Username:     connection.Username.String,
			Password:     *mysqlPassword,
			DatabaseName: connection.DatabaseName.String,
			Host:         connection.Host.String,
		}, nil
	default:
		return nil, errors.Newf("(query.QueryServiceImpl.GetClient) unrecognized warehouse type %v", connection.ConnectionType)
	}
}

func (qs QueryServiceImpl) GetWarehouseClient(ctx context.Context, connection *models.Connection) (WarehouseClient, error) {
	switch connection.ConnectionType {
	case models.ConnectionTypeBigQuery:
		bigQueryCredentialsString, err := qs.cryptoService.DecryptConnectionCredentials(connection.Credentials.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetWarehouseClient) decrypting BigQuery credentials")
		}

		var bigQueryCredentials models.BigQueryCredentials
		err = json.Unmarshal([]byte(*bigQueryCredentialsString), &bigQueryCredentials)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetClient) unmarshalling BigQuery credentials")
		}

		if !connection.Location.Valid {
			return nil, errors.NewCustomerVisibleError("BigQuery connection must have location defined")
		}

		return BigQueryApiClient{
			ProjectID:   &bigQueryCredentials.ProjectID,
			Credentials: bigQueryCredentialsString,
			Location:    &connection.Location.String,
		}, nil
	default:
		return nil, errors.Newf("(query.QueryServiceImpl.GetWarehouseClient) unrecognized warehouse type %v", connection.ConnectionType)
	}
}

func (qs QueryServiceImpl) GetDatabaseClient(ctx context.Context, connection *models.Connection) (DatabaseClient, error) {
	switch connection.ConnectionType {
	case models.ConnectionTypeDynamoDb:
		dynamoDbAccessKey, err := qs.cryptoService.DecryptConnectionCredentials(connection.Password.String)
		if err != nil {
			return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetDatabaseClient) decrypting DynamoDB credentials")
		}

		if !connection.Location.Valid {
			return nil, errors.NewCustomerVisibleError("DynamoDB connection must have location defined")
		}

		return DynamoDbClient{
			KeyID:     connection.Username.String,
			AccessKey: *dynamoDbAccessKey,
			Location:  connection.Location.String,
		}, nil
	default:
		return nil, errors.Newf("(query.QueryServiceImpl.GetDatabaseClient) unrecognized database type %v", connection.ConnectionType)
	}
}

func (qs QueryServiceImpl) RunQuery(ctx context.Context, connection *models.Connection, queryString string) (*data.QueryResults, error) {
	client, err := qs.GetClient(ctx, connection)
	if err != nil {
		return nil, errors.Wrap(err, "(query.QueryServiceImpl.RunQuery)")
	}

	return client.RunQuery(ctx, queryString)
}

func (qs QueryServiceImpl) GetQueryIterator(ctx context.Context, connection *models.Connection, queryString string) (data.RowIterator, error) {
	client, err := qs.GetClient(ctx, connection)
	if err != nil {
		return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetQueryIterator)")
	}

	return client.GetQueryIterator(ctx, queryString)
}

func (qs QueryServiceImpl) GetNamespaces(ctx context.Context, connection *models.Connection) ([]string, error) {
	client, err := qs.GetClient(ctx, connection)
	if err != nil {
		return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetNamespaces)")
	}

	return client.GetNamespaces(ctx)
}

func (qs QueryServiceImpl) GetTables(ctx context.Context, connection *models.Connection, namespace string) ([]string, error) {
	client, err := qs.GetClient(ctx, connection)
	if err != nil {
		return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetTables)")
	}

	return client.GetTables(ctx, namespace)
}

func (qs QueryServiceImpl) GetSchema(ctx context.Context, connection *models.Connection, namespace string, tableName string) ([]data.Field, error) {
	client, err := qs.GetClient(ctx, connection)
	if err != nil {
		return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetSchema)")
	}

	return client.GetSchema(ctx, namespace, tableName)
}

func (qs QueryServiceImpl) GetFieldValues(ctx context.Context, connection *models.Connection, namespace string, tableName string, fieldName string) ([]any, error) {
	client, err := qs.GetClient(ctx, connection)
	if err != nil {
		return nil, errors.Wrap(err, "(query.QueryServiceImpl.GetFieldValues)")
	}

	return client.GetFieldValues(ctx, namespace, tableName, fieldName)
}
