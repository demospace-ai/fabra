package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"

	"github.com/go-playground/validator/v10"
)

func (s ApiService) LinkCreateSync(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.LinkCreateSync)")
	}

	if auth.LinkToken == nil {
		return errors.Wrap(errors.NewBadRequest("must send link token"), "(api.LinkCreateSync)")
	}

	decoder := json.NewDecoder(r.Body)
	var createSyncRequest CreateSyncRequest
	err := decoder.Decode(&createSyncRequest)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCreateSync)")
	}

	// TODO: validate connection parameters
	validate := validator.New()
	err = validate.Struct(createSyncRequest)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCreateSync)")
	}

	// Do NOT use the end customer ID from the request— we must pull it from the link token to ensure
	// the customer is authorized.
	sync, fieldMappings, err := s.createSync(auth, createSyncRequest, auth.LinkToken.EndCustomerID)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCreateSync)")
	}

	return json.NewEncoder(w).Encode(CreateSyncResponse{
		Sync:          *sync,
		FieldMappings: fieldMappings,
	})
}
