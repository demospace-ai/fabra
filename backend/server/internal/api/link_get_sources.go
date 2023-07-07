package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/sources"
	"go.fabra.io/server/common/views"
)

type GetSourcesResponse struct {
	Sources []views.Source `json:"sources"`
}

func (s ApiService) LinkGetSources(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.LinkGetSources)")
	}

	if auth.LinkToken == nil {
		return errors.Wrap(errors.NewBadRequest("must send link token"), "(api.LinkGetSources)")
	}

	// TODO: write test to make sure only authorized users can use the data connection
	// Needed to ensure end customer ID encoded by the link token owns the source/connection
	sources, err := sources.LoadAllSources(s.db, auth.Organization.ID, auth.LinkToken.EndCustomerID)
	if err != nil {
		return errors.Wrap(err, "(api.LinkGetSources)")
	}

	return json.NewEncoder(w).Encode(GetSourcesResponse{
		Sources: views.ConvertSourceConnections(sources),
	})
}
