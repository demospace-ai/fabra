package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/connections"
	"go.fabra.io/server/common/repositories/destinations"
	"go.fabra.io/server/common/views"
)

type CreateDestinationRequest struct {
	DisplayName     string                 `json:"display_name"`
	ConnectionType  models.ConnectionType  `json:"connection_type"`
	StagingBucket   *string                `json:"staging_bucket"`
	BigQueryConfig  *input.BigQueryConfig  `json:"bigquery_config,omitempty"`
	SnowflakeConfig *input.SnowflakeConfig `json:"snowflake_config,omitempty"`
	RedshiftConfig  *input.RedshiftConfig  `json:"redshift_config,omitempty"`
	MongoDbConfig   *input.MongoDbConfig   `json:"mongodb_config,omitempty"`
	WebhookConfig   *input.WebhookConfig   `json:"webhook_config,omitempty"`
	DynamoDbConfig  *input.DynamoDbConfig  `json:"dynamodb_config,omitempty"`
}

type CreateDestinationResponse struct {
	Destination views.Destination `json:"destination"`
}

func (s ApiService) CreateDestination(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.CreateDestination)")
	}

	decoder := json.NewDecoder(r.Body)
	var createDestinationRequest CreateDestinationRequest
	err := decoder.Decode(&createDestinationRequest)
	if err != nil {
		return errors.Wrap(err, "(api.CreateDestination)")
	}

	err = validateCreateDestinationRequest(createDestinationRequest)
	if err != nil {
		return errors.Wrap(err, "(api.CreateDestination)")
	}

	// TODO: Create connection + destination in a transaction
	var connection *models.Connection
	var webhookSigningKey string
	switch createDestinationRequest.ConnectionType {
	case models.ConnectionTypeBigQuery:
		encryptedCredentials, encryptionErr := s.cryptoService.EncryptConnectionCredentials(createDestinationRequest.BigQueryConfig.Credentials)
		if encryptionErr != nil {
			return errors.Wrap(encryptionErr, "(api.CreateDestination)")
		}
		connection, err = connections.CreateBigQueryConnection(
			s.db, auth.Organization.ID, *encryptedCredentials, createDestinationRequest.BigQueryConfig.Location,
		)
	case models.ConnectionTypeSnowflake:
		encryptedCredentials, encryptionErr := s.cryptoService.EncryptConnectionCredentials(createDestinationRequest.SnowflakeConfig.Password)
		if encryptionErr != nil {
			return errors.Wrap(encryptionErr, "(api.CreateDestination)")
		}
		connection, err = connections.CreateSnowflakeConnection(
			s.db, auth.Organization.ID, *createDestinationRequest.SnowflakeConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeRedshift:
		encryptedCredentials, encryptionErr := s.cryptoService.EncryptConnectionCredentials(createDestinationRequest.RedshiftConfig.Password)
		if encryptionErr != nil {
			return errors.Wrap(encryptionErr, "(api.CreateDestination)")
		}
		connection, err = connections.CreateRedshiftConnection(
			s.db, auth.Organization.ID, *createDestinationRequest.RedshiftConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeMongoDb:
		encryptedCredentials, encryptionErr := s.cryptoService.EncryptConnectionCredentials(createDestinationRequest.MongoDbConfig.Password)
		if encryptionErr != nil {
			return errors.Wrap(encryptionErr, "(api.CreateDestination)")
		}
		connection, err = connections.CreateMongoDbConnection(
			s.db, auth.Organization.ID, *createDestinationRequest.MongoDbConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeDynamoDb:
		encryptedCredentials, encryptionErr := s.cryptoService.EncryptConnectionCredentials(createDestinationRequest.DynamoDbConfig.SecretKey)
		if encryptionErr != nil {
			return errors.Wrap(encryptionErr, "(api.CreateDestination)")
		}
		connection, err = connections.CreateDynamoDbConnection(
			s.db, auth.Organization.ID, createDestinationRequest.DynamoDbConfig.AccessKey, *encryptedCredentials,
			createDestinationRequest.DynamoDbConfig.Region,
		)
	case models.ConnectionTypeWebhook:
		webhookSigningKey = crypto.GenerateSigningKey()
		encryptedSigningKey, encryptionErr := s.cryptoService.EncryptWebhookSigningKey(webhookSigningKey)
		if encryptionErr != nil {
			return errors.Wrap(encryptionErr, "(api.CreateDestination)")
		}
		connection, err = connections.CreateWebhookConnection(
			s.db, auth.Organization.ID, *createDestinationRequest.WebhookConfig, *encryptedSigningKey,
		)
	default:
		return errors.Newf("(api.CreateDestination) unsupported connection type: %s", createDestinationRequest.ConnectionType)
	}

	if err != nil {
		return errors.Wrap(err, "(api.CreateDestination)")
	}

	destination, err := destinations.CreateDestination(
		s.db,
		auth.Organization.ID,
		createDestinationRequest.DisplayName,
		connection.ID,
		createDestinationRequest.StagingBucket,
	)
	if err != nil {
		return errors.Wrap(err, "(api.CreateDestination)")
	}

	if err != nil {
		return errors.Wrap(err, "(api.CreateDestination)")
	}

	var destinationView views.Destination
	if connection.ConnectionType == models.ConnectionTypeWebhook {
		destinationView = views.ConvertWebhook(*destination, *connection, &webhookSigningKey)
	} else {
		destinationView = views.ConvertDestination(*destination, *connection)
	}

	return json.NewEncoder(w).Encode(CreateDestinationResponse{
		destinationView,
	})
}

func validateCreateDestinationRequest(request CreateDestinationRequest) error {
	switch request.ConnectionType {
	case models.ConnectionTypeBigQuery:
		return validateCreateBigQueryDestination(request)
	case models.ConnectionTypeSnowflake:
		return validateCreateSnowflakeDestination(request)
	case models.ConnectionTypeRedshift:
		return validateCreateRedshiftDestination(request)
	case models.ConnectionTypeMongoDb:
		return validateCreateMongoDbDestination(request)
	case models.ConnectionTypeWebhook:
		return validateCreateWebhookDestination(request)
	case models.ConnectionTypeDynamoDb:
		return validateCreateDynamoDbDestination(request)
	default:
		return errors.Wrap(errors.NewBadRequestf("unknown connection type: %s", request.ConnectionType), "(api.validateCreateDestinationRequest)")
	}
}

func validateCreateBigQueryDestination(request CreateDestinationRequest) error {
	if request.BigQueryConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing BigQuery configuration"), "(api.validateCreateBigQueryDestination)")
	}

	var bigQueryCredentials models.BigQueryCredentials
	err := json.Unmarshal([]byte(request.BigQueryConfig.Credentials), &bigQueryCredentials)
	if err != nil {
		return errors.Wrap(err, "(api.validateCreateBigQueryDestination)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateCreateSnowflakeDestination(request CreateDestinationRequest) error {
	if request.SnowflakeConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Snowflake configuration"), "(api.validateCreateSnowflakeDestination)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateCreateRedshiftDestination(request CreateDestinationRequest) error {
	if request.RedshiftConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Redshift configuration"), "(api.validateCreateRedshiftDestination)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateCreateMongoDbDestination(request CreateDestinationRequest) error {
	if request.MongoDbConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing MongoDB configuration"), "(api.validateCreateMongoDbDestination)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateCreateWebhookDestination(request CreateDestinationRequest) error {
	if request.WebhookConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing Webhook configuration"), "(api.validateCreateWebhookDestination)")
	}

	if request.WebhookConfig.URL == "" {
		return errors.Wrap(errors.NewBadRequest("missing Webhook URL"), "(api.validateCreateWebhookDestination)")
	}

	if !strings.Contains(request.WebhookConfig.URL, "https") {
		return errors.Wrap(errors.NewBadRequest("Webhook must be HTTPS"), "(api.validateCreateWebhookDestination)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}

func validateCreateDynamoDbDestination(request CreateDestinationRequest) error {
	if request.DynamoDbConfig == nil {
		return errors.Wrap(errors.NewBadRequest("missing DynamoDB configuration"), "(api.validateCreateDynamoDbDestination)")
	}

	// TODO: validate the fields all exist in the credentials object

	return nil
}
