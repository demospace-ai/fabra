package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/connections"
)

type GetTablesResponse struct {
	Tables []string `json:"tables"`
}

func (s ApiService) GetTables(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.GetTables)")
	}

	strConnectionID := r.URL.Query().Get("connectionID")
	if len(strConnectionID) == 0 {
		return errors.Newf("(api.GetTables) missing connection ID from GetDatasets request URL: %s", r.URL.RequestURI())
	}

	connectionID, err := strconv.ParseInt(strConnectionID, 10, 64)
	if err != nil {
		return errors.Wrap(err, "(api.GetTables)")
	}

	// TODO: write test to make sure only authorized users can use the data connection
	connection, err := connections.LoadConnectionByID(s.db, auth.Organization.ID, connectionID)
	if err != nil {
		return errors.Wrap(err, "(api.GetTables)")
	}

	namespace := r.URL.Query().Get("namespace")

	if connection.ConnectionType != models.ConnectionTypeDynamoDb && len(namespace) == 0 {
		return errors.Newf("(api.GetTables) missing namespace from GetTables request URL: %s", r.URL.RequestURI())
	}

	tables, err := s.queryService.GetTables(ctx, connection, namespace)
	if err != nil {
		return errors.Wrap(err, "(api.GetTables)")
	}

	return json.NewEncoder(w).Encode(GetTablesResponse{
		Tables: tables,
	})
}
