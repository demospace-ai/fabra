package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
)

func (s ApiService) LinkGetSyncs(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.LinkGetSyncs)")
	}

	if auth.LinkToken == nil {
		return errors.Wrap(errors.NewBadRequest("must send link token"), "(api.LinkGetSyncs)")
	}

	syncs, sources, objects, err := s.getSyncsForCustomer(auth, auth.LinkToken.EndCustomerID)
	if err != nil {
		return errors.Wrap(err, "(api.GetSyncsForCustomer)")
	}

	return json.NewEncoder(w).Encode(GetSyncsResponse{
		Syncs:   syncs,
		Sources: sources,
		Objects: objects,
	})
}
