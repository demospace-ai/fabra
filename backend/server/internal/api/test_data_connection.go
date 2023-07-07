package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cloud.google.com/go/bigquery"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/snowflakedb/gosnowflake"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type TestDataConnectionRequest struct {
	DisplayName     string                 `json:"display_name"`
	ConnectionType  models.ConnectionType  `json:"connection_type"`
	BigQueryConfig  *input.BigQueryConfig  `json:"bigquery_config,omitempty"`
	SnowflakeConfig *input.SnowflakeConfig `json:"snowflake_config,omitempty"`
	RedshiftConfig  *input.RedshiftConfig  `json:"redshift_config,omitempty"`
	SynapseConfig   *input.SynapseConfig   `json:"synapse_config,omitempty"`
	PostgresConfig  *input.PostgresConfig  `json:"postgres_config,omitempty"`
	MySqlConfig     *input.MySqlConfig     `json:"mysql_config,omitempty"`
	MongoDbConfig   *input.MongoDbConfig   `json:"mongodb_config,omitempty"`
	WebhookConfig   *input.WebhookConfig   `json:"webhook_config,omitempty"`
	DynamoDbConfig  *input.DynamoDbConfig  `json:"dynamodb_config,omitempty"`
}

func (s ApiService) TestDataConnection(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "TestDataConnection")
	}

	decoder := json.NewDecoder(r.Body)
	var testDataConnectionRequest TestDataConnectionRequest
	err := decoder.Decode(&testDataConnectionRequest)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.TestDataConnection)")
	}

	err = validateTestDataConnectionRequest(testDataConnectionRequest)
	if err != nil {
		return errors.Wrap(err, "(api.TestDataConnection)")
	}

	switch testDataConnectionRequest.ConnectionType {
	case models.ConnectionTypeBigQuery:
		err = testBigQueryConnection(*testDataConnectionRequest.BigQueryConfig)
	case models.ConnectionTypeSnowflake:
		err = testSnowflakeConnection(*testDataConnectionRequest.SnowflakeConfig)
	case models.ConnectionTypeMongoDb:
		err = testMongoDbConnection(*testDataConnectionRequest.MongoDbConfig)
	case models.ConnectionTypeRedshift:
		err = testRedshiftConnection(*testDataConnectionRequest.RedshiftConfig)
	case models.ConnectionTypeSynapse:
		err = testSynapseConnection(*testDataConnectionRequest.SynapseConfig)
	case models.ConnectionTypePostgres:
		err = testPostgresConnection(*testDataConnectionRequest.PostgresConfig)
	case models.ConnectionTypeMySQL:
		err = testMySqlConnection(*testDataConnectionRequest.MySqlConfig)
	case models.ConnectionTypeWebhook:
		err = testWebhookConnection(*testDataConnectionRequest.WebhookConfig)
	case models.ConnectionTypeDynamoDb:
		err = testDynamoDbConnection(*testDataConnectionRequest.DynamoDbConfig)
	default:
		err = errors.NewBadRequest(fmt.Sprintf("unknown connection type: %s", testDataConnectionRequest.ConnectionType))
	}

	if err != nil {
		return errors.Wrap(err, "(api.TestDataConnection)")
	}

	return nil
}

func testBigQueryConnection(bigqueryConfig input.BigQueryConfig) error {
	var bigQueryCredentials models.BigQueryCredentials
	err := json.Unmarshal([]byte(bigqueryConfig.Credentials), &bigQueryCredentials)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testBigQueryConnection)")
	}

	credentialOption := option.WithCredentialsJSON([]byte(bigqueryConfig.Credentials))

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, bigQueryCredentials.ProjectID, credentialOption)
	if err != nil {
		return errors.Wrap(err, "(api.testBigQueryConnection)")
	}
	defer client.Close()

	it := client.Datasets(ctx)
	_, err = it.Next()

	if err != nil && err != iterator.Done {
		return errors.Wrap(err, "(api.testBigQueryConnection)")
	}

	return nil
}

func testSnowflakeConnection(snowflakeConfig input.SnowflakeConfig) error {
	account := strings.Split(snowflakeConfig.Host, ".")[0] // TODO: remove the https/http
	config := gosnowflake.Config{
		Account:       account,
		User:          snowflakeConfig.Username,
		Password:      snowflakeConfig.Password,
		Warehouse:     snowflakeConfig.WarehouseName,
		Database:      snowflakeConfig.DatabaseName,
		Role:          snowflakeConfig.Role,
		Host:          snowflakeConfig.Host,
		LoginTimeout:  3 * time.Second,
		ClientTimeout: 3 * time.Second,
	}

	dsn, err := gosnowflake.DSN(&config)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSnowflakeConnection)")
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSnowflakeConnection)")
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSnowflakeConnection)")
	}
	defer rows.Close()

	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSnowflakeConnection)")
		}
		if v != 1 {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSnowflakeConnection)")
		}
	}
	if rows.Err() != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSnowflakeConnection)")
	}

	return nil
}

func testRedshiftConnection(redshiftConfig input.RedshiftConfig) error {
	params := url.Values{}
	params.Add("sslmode", "require")
	params.Add("connect_timeout", "5")

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(redshiftConfig.Username, redshiftConfig.Password),
		Host:     redshiftConfig.Endpoint,
		Path:     redshiftConfig.DatabaseName,
		RawQuery: params.Encode(),
	}

	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testRedshiftConnection)")
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testRedshiftConnection)")
	}
	defer rows.Close()

	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testRedshiftConnection)")
		}
		if v != 1 {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testRedshiftConnection)")
		}
	}
	if rows.Err() != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testRedshiftConnection)")
	}

	return nil
}

func testSynapseConnection(synapseConfig input.SynapseConfig) error {
	params := url.Values{}
	params.Add("database", synapseConfig.DatabaseName)
	params.Add("sslmode", "encrypt")
	params.Add("TrustServerCertificate", "true")
	params.Add("dial timeout", "3")

	dsn := url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(synapseConfig.Username, synapseConfig.Password),
		Host:     synapseConfig.Endpoint,
		RawQuery: params.Encode(),
	}

	db, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSynapseConnection)")
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSynapseConnection)")
	}
	defer rows.Close()

	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSynapseConnection)")
		}
		if v != 1 {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSynapseConnection)")
		}
	}
	if rows.Err() != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testSynapseConnection)")
	}

	return nil
}

func testMongoDbConnection(mongodbConfig input.MongoDbConfig) error {
	connectionOptions := ""
	if mongodbConfig.ConnectionOptions != nil {
		connectionOptions = *mongodbConfig.ConnectionOptions
	}

	connectionString := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?%s",
		mongodbConfig.Username,
		mongodbConfig.Password,
		mongodbConfig.Host,
		connectionOptions,
	)
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		SetServerAPIOptions(serverAPIOptions).
		SetConnectTimeout(3 * time.Second).
		ApplyURI(connectionString) // Apply URI last since this contains connection options from the user
	_, err := mongo.Connect(context.TODO(), clientOptions)
	return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testMongoDbConnection)")
}

func testPostgresConnection(postgresConfig input.PostgresConfig) error {
	params := url.Values{}
	params.Add("sslmode", "require")
	params.Add("connect_timeout", "5")

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(postgresConfig.Username, postgresConfig.Password),
		Host:     postgresConfig.Endpoint,
		Path:     postgresConfig.DatabaseName,
		RawQuery: params.Encode(),
	}

	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testPostgresConnection)")
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testPostgresConnection)")
	}
	defer rows.Close()

	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testPostgresConnection)")
		}
		if v != 1 {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testPostgresConnection)")
		}
	}
	if rows.Err() != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testPostgresConnection)")
	}

	return nil
}

func testMySqlConnection(mysqlConfig input.MySqlConfig) error {
	params := url.Values{}
	params.Add("tls", "true")
	params.Add("timeout", "5s")

	// Can't use url.Url because mysql does not accept a scheme
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Endpoint, mysqlConfig.DatabaseName, params.Encode())

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testMySqlConnection)")
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1")
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testMySqlConnection)")
	}
	defer rows.Close()

	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testMySqlConnection)")
		}
		if v != 1 {
			return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testMySqlConnection)")
		}
	}
	if rows.Err() != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testMySqlConnection)")
	}

	return nil
}

func testWebhookConnection(webhookConfig input.WebhookConfig) error {
	_, err := url.ParseRequestURI(webhookConfig.URL)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testWebhookConnection)")
	}
	resp, err := http.Head(webhookConfig.URL)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testWebhookConnection)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.WrapCustomerVisibleError(fmt.Errorf("unexpected status code: %d", resp.StatusCode)), "(api.testWebhookConnection) ")
	}

	return nil
}

func testDynamoDbConnection(dynamoDbConfig input.DynamoDbConfig) error {
	fmt.Print("testDynamoDbConnection", dynamoDbConfig)
	// region := config.WithRegion(dynamoDbConfig.Region)
	creds := credentials.NewStaticCredentialsProvider(
		dynamoDbConfig.AccessKey,
		dynamoDbConfig.SecretKey,
		"",
	)
	credProvider := config.WithCredentialsProvider(creds)
	cfg, err := config.LoadDefaultConfig(context.TODO(), credProvider)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testDynamoDbConnection)")
	}

	svc := dynamodb.NewFromConfig(cfg)
	// Build the request with its input parameters
	_, err = svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})

	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.testDynamoDbConnection)")
	}

	return nil
}

func validateTestDataConnectionRequest(request TestDataConnectionRequest) error {
	switch request.ConnectionType {
	case models.ConnectionTypeBigQuery:
		return validateTestBigQueryConnection(request)
	case models.ConnectionTypeSnowflake:
		return validateTestSnowflakeConnection(request)
	case models.ConnectionTypeRedshift:
		return validateTestRedshiftConnection(request)
	case models.ConnectionTypeMongoDb:
		return validateTestMongoConnection(request)
	case models.ConnectionTypeSynapse:
		return validateTestSynapseConnection(request)
	case models.ConnectionTypePostgres:
		return validateTestPostgresConnection(request)
	case models.ConnectionTypeMySQL:
		return validateTestMySqlConnection(request)
	case models.ConnectionTypeWebhook:
		return validateTestWebhookConnection(request)
	case models.ConnectionTypeDynamoDb:
		return validateTestDynamoDbConnection(request)
	default:
		return errors.Wrap(errors.NewBadRequestf("unknown connection type: %s", request.ConnectionType), "(api.validateTestDataConnectionRequest)")
	}
}

func validateTestBigQueryConnection(request TestDataConnectionRequest) error {
	if request.BigQueryConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing BigQuery configuration"), "(api.validateTestBigQueryConnection)")
	}

	var bigQueryCredentials models.BigQueryCredentials
	err := json.Unmarshal([]byte(request.BigQueryConfig.Credentials), &bigQueryCredentials)
	if err != nil {
		return errors.Wrap(err, "validateTestBigQueryConnection")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestSnowflakeConnection(request TestDataConnectionRequest) error {
	if request.SnowflakeConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Snowflake configuration"), "(api.validateTestSnowflakeConnection)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestRedshiftConnection(request TestDataConnectionRequest) error {
	if request.RedshiftConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Redshift configuration"), "(api.validateTestRedshiftConnection)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestMongoConnection(request TestDataConnectionRequest) error {
	if request.MongoDbConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing MongoDB configuration"), "(api.validateTestMongoConnection)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestSynapseConnection(request TestDataConnectionRequest) error {
	if request.SynapseConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Synapse configuration"), "(api.validateTestSynapseConnection)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestPostgresConnection(request TestDataConnectionRequest) error {
	if request.PostgresConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Postgres configuration"), "(api.validateTestPostgresConnection)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestMySqlConnection(request TestDataConnectionRequest) error {
	if request.MySqlConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing MySQL configuration"), "(api.validateTestMySqlConnection)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateTestWebhookConnection(request TestDataConnectionRequest) error {
	if request.WebhookConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Webhook configuration"), "(api.validateTestWebhookConnection)")
	}

	// TODO: validate the fields all exist

	return nil
}

func validateTestDynamoDbConnection(request TestDataConnectionRequest) error {
	if request.DynamoDbConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing DynamoDB configuration"), "(api.validateTestDynamoDBConnection)")
	}

	// TODO: validate the fields all exist

	return nil
}
