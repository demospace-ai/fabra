package models

import "go.fabra.io/server/common/database"

type ConnectionType string

const (
	ConnectionTypeSynapse   ConnectionType = "synapse"
	ConnectionTypeBigQuery  ConnectionType = "bigquery"
	ConnectionTypeMongoDb   ConnectionType = "mongodb"
	ConnectionTypeRedshift  ConnectionType = "redshift"
	ConnectionTypeSnowflake ConnectionType = "snowflake"
	ConnectionTypePostgres  ConnectionType = "postgres"
	ConnectionTypeMySQL     ConnectionType = "mysql"
	ConnectionTypeWebhook   ConnectionType = "webhook"
	ConnectionTypeDynamoDb  ConnectionType = "dynamodb"
)

type BigQueryCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

type Connection struct {
	OrganizationID    int64
	ConnectionType    ConnectionType      `json:"connection_type"`
	Credentials       database.NullString `json:"-"`
	Username          database.NullString `json:"-"`
	Password          database.NullString `json:"-"`
	Location          database.NullString
	WarehouseName     database.NullString
	DatabaseName      database.NullString
	Role              database.NullString
	Host              database.NullString
	Port              database.NullString
	ConnectionOptions database.NullString

	BaseModel
}
