package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/repositories/objects"
	"go.fabra.io/server/common/views"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CreateObjectFieldsRequest = struct {
	ObjectFields []input.ObjectField `json:"object_fields" validate:"required"`
}

type CreateObjectFieldsResponse struct {
	ObjectFields []views.ObjectField `json:"object_fields"`
	Failures     []input.ObjectField `json:"failures"`
}

func (s ApiService) CreateObjectFields(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.NewBadRequest("must setup organization first")
	}

	vars := mux.Vars(r)
	strObjectId, ok := vars["objectID"]
	if !ok {
		return errors.Newf("missing objectID request URL: %s", r.URL.RequestURI())
	}

	objectId, err := strconv.ParseInt(strObjectId, 10, 64)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)
	var requestBody CreateObjectFieldsRequest
	if err := decoder.Decode(&requestBody); err != nil {
		return err
	}

	validate := validator.New()
	for _, objectField := range requestBody.ObjectFields {
		err = validate.Struct(objectField)
		if err != nil {
			return err
		}
	}

	objectFieldsView := []views.ObjectField{}
	failures := []input.ObjectField{}
	for _, objectField := range requestBody.ObjectFields {
		field, err := objects.CreateObjectField(
			s.db,
			auth.Organization.ID,
			objectId,
			objectField,
		)
		if err == nil {
			objectFieldsView = append(objectFieldsView, views.ConvertObjectField(field))
		} else {
			failures = append(failures, objectField)
		}
	}

	return json.NewEncoder(w).Encode(CreateObjectFieldsResponse{
		ObjectFields: objectFieldsView,
		Failures:     failures,
	})
}
