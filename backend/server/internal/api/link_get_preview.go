package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/repositories/connections"
	"go.fabra.io/server/common/repositories/sources"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LinkGetPreviewRequest struct {
	SourceID  int64  `json:"source_id"`
	Namespace string `json:"namespace"`
	TableName string `json:"table_name"`
}

func (s ApiService) LinkGetPreview(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.LinkGetPreview)")
	}

	if auth.LinkToken == nil {
		return errors.Wrap(errors.NewBadRequest("must send link token"), "(api.LinkGetPreview)")
	}

	decoder := json.NewDecoder(r.Body)
	var getPreviewRequest LinkGetPreviewRequest
	err := decoder.Decode(&getPreviewRequest)
	if err != nil {
		return errors.Wrap(err, "(api.LinkGetPreview)")
	}

	// Needed to ensure end customer ID encoded by the link token owns the source/connection
	source, err := sources.LoadSourceByID(s.db, auth.Organization.ID, auth.LinkToken.EndCustomerID, getPreviewRequest.SourceID)
	if err != nil {
		return errors.Wrap(err, "(api.LinkGetPreview)")
	}

	connection, err := connections.LoadConnectionByID(s.db, auth.Organization.ID, source.ConnectionID)
	if err != nil {
		return errors.Wrap(err, "(api.LinkGetPreview)")
	}

	query, err := getPreviewQuery(connection.ConnectionType, getPreviewRequest.Namespace, getPreviewRequest.TableName)
	if err != nil {
		return errors.Wrap(err, "(api.LinkGetPreview)")
	}

	queryResults, err := s.queryService.RunQuery(context.TODO(), connection, *query)
	if err != nil {
		return errors.Wrap(err, "(api.LinkGetPreview)")
	}

	return json.NewEncoder(w).Encode(queryResults)
}

func getPreviewQuery(connectionType models.ConnectionType, namespace string, tableName string) (*string, error) {
	switch connectionType {
	case models.ConnectionTypeMongoDb:
		mongoQuery := query.MongoQuery{
			Database:   namespace,
			Collection: tableName,
			Filter:     bson.D{},
			Options:    options.Find().SetLimit(100),
		}

		queryString := query.CreateMongoQueryString(mongoQuery)
		return &queryString, nil
	case models.ConnectionTypeBigQuery:
		fallthrough
	case models.ConnectionTypeRedshift:
		fallthrough
	case models.ConnectionTypePostgres:
		fallthrough
	case models.ConnectionTypeMySQL:
		fallthrough
	case models.ConnectionTypeSnowflake:
		queryStr := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 100;", namespace, tableName)
		return &queryStr, nil
	case models.ConnectionTypeSynapse:
		queryStr := fmt.Sprintf("SELECT TOP(100) * FROM %s.%s;", namespace, tableName)
		return &queryStr, nil
	default:
		return nil, errors.Wrap(errors.Newf("unexpected connection type: %s", connectionType), "(api.getPreviewQuery)")
	}
}
