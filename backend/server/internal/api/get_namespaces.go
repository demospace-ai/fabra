package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/connections"
)

type GetNamespacesResponse struct {
	Namespaces []string `json:"namespaces"`
}

func (s ApiService) GetNamespaces(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.GetNamespaces)")
	}

	strConnectionID := r.URL.Query().Get("connectionID")
	if len(strConnectionID) == 0 {
		return errors.Newf("(api.GetNamespaces) missing connection ID from GetNamespaces request URL: %s", r.URL.RequestURI())
	}

	connectionID, err := strconv.ParseInt(strConnectionID, 10, 64)
	if err != nil {
		return errors.Wrap(err, "(api.GetNamespaces)")
	}

	// TODO: write test to make sure only authorized users can use the data connection
	connection, err := connections.LoadConnectionByID(s.db, auth.Organization.ID, connectionID)
	if err != nil {
		return errors.Wrap(err, "(api.GetNamespaces)")
	}

	namespaces, err := s.queryService.GetNamespaces(ctx, connection)
	if err != nil {
		return errors.Wrap(err, "(api.GetNamespaces)")
	}

	return json.NewEncoder(w).Encode(GetNamespacesResponse{
		Namespaces: namespaces,
	})
}
