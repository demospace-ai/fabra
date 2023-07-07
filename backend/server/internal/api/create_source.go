package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/connections"
	"go.fabra.io/server/common/repositories/sources"
	"go.fabra.io/server/common/views"

	"github.com/go-playground/validator/v10"
)

type CreateSourceRequest struct {
	DisplayName     string                 `json:"display_name" validate:"required"`
	ConnectionType  models.ConnectionType  `json:"connection_type"`
	BigQueryConfig  *input.BigQueryConfig  `json:"bigquery_config,omitempty"`
	SnowflakeConfig *input.SnowflakeConfig `json:"snowflake_config,omitempty"`
	RedshiftConfig  *input.RedshiftConfig  `json:"redshift_config,omitempty"`
	MongoDbConfig   *input.MongoDbConfig   `json:"mongodb_config,omitempty"`
	SynapseConfig   *input.SynapseConfig   `json:"synapse_config,omitempty"`
	PostgresConfig  *input.PostgresConfig  `json:"postgres_config,omitempty"`
	MySqlConfig     *input.MySqlConfig     `json:"mysql_config,omitempty"`
	EndCustomerID   *string                `json:"end_customer_id,omitempty"`
}

type CreateSourceResponse struct {
	Source views.Source `json:"source"`
}

func (s ApiService) CreateSource(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.CreateSource)")
	}

	decoder := json.NewDecoder(r.Body)
	var createSourceRequest CreateSourceRequest
	err := decoder.Decode(&createSourceRequest)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.CreateSource)")
	}

	// TODO: validate connection parameters
	validate := validator.New()
	err = validate.Struct(createSourceRequest)
	if err != nil {
		return errors.Wrap(errors.WrapCustomerVisibleError(err), "(api.CreateSource)")
	}

	if createSourceRequest.EndCustomerID == nil {
		return errors.Wrap(errors.NewBadRequest("must provide end customer ID"), "(api.CreateSource)")
	}

	source, connection, err := s.createSource(auth, createSourceRequest, *createSourceRequest.EndCustomerID)
	if err != nil {
		return errors.Wrap(err, "(api.CreateSource)")
	}

	return json.NewEncoder(w).Encode(CreateSourceResponse{
		views.ConvertSource(*source, *connection),
	})
}

func (s ApiService) createSource(auth auth.Authentication, createSourceRequest CreateSourceRequest, endCustomerID string) (*models.Source, *models.Connection, error) {
	// TODO: Create connection + source in a transaction
	var connection *models.Connection
	var encryptedCredentials *string
	var err error
	switch createSourceRequest.ConnectionType {
	case models.ConnectionTypeBigQuery:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.BigQueryConfig.Credentials)
		if err != nil {
			return nil, nil, err
		}
		connection, err = connections.CreateBigQueryConnection(
			s.db, auth.Organization.ID, *encryptedCredentials, createSourceRequest.BigQueryConfig.Location,
		)
	case models.ConnectionTypeSnowflake:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.SnowflakeConfig.Password)
		if err != nil {
			return nil, nil, errors.Wrap(err, "(api.createSource)")
		}
		connection, err = connections.CreateSnowflakeConnection(
			s.db, auth.Organization.ID, *createSourceRequest.SnowflakeConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeRedshift:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.RedshiftConfig.Password)
		if err != nil {
			return nil, nil, errors.Wrap(err, "(api.createSource)")
		}
		connection, err = connections.CreateRedshiftConnection(
			s.db, auth.Organization.ID, *createSourceRequest.RedshiftConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeMongoDb:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.MongoDbConfig.Password)
		if err != nil {
			return nil, nil, errors.Wrap(err, "(api.createSource)")
		}
		connection, err = connections.CreateMongoDbConnection(
			s.db, auth.Organization.ID, *createSourceRequest.MongoDbConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeSynapse:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.SynapseConfig.Password)
		if err != nil {
			return nil, nil, errors.Wrap(err, "(api.createSource)")
		}
		connection, err = connections.CreateSynapseConnection(
			s.db, auth.Organization.ID, *createSourceRequest.SynapseConfig, *encryptedCredentials,
		)
	case models.ConnectionTypePostgres:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.PostgresConfig.Password)
		if err != nil {
			return nil, nil, errors.Wrap(err, "(api.createSource)")
		}
		connection, err = connections.CreatePostgresConnection(
			s.db, auth.Organization.ID, *createSourceRequest.PostgresConfig, *encryptedCredentials,
		)
	case models.ConnectionTypeMySQL:
		encryptedCredentials, err = s.cryptoService.EncryptConnectionCredentials(createSourceRequest.MySqlConfig.Password)
		if err != nil {
			return nil, nil, errors.Wrap(err, "(api.createSource)")
		}
		connection, err = connections.CreateMySqlConnection(
			s.db, auth.Organization.ID, *createSourceRequest.MySqlConfig, *encryptedCredentials,
		)
	default:
		return nil, nil, errors.Wrap(errors.Newf("unsupported connection type: %s", createSourceRequest.ConnectionType), "(api.createSource)")
	}

	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSource)")
	}

	source, err := sources.CreateSource(
		s.db,
		auth.Organization.ID,
		createSourceRequest.DisplayName,
		endCustomerID,
		connection.ID,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSource)")
	}

	return source, connection, nil
}
