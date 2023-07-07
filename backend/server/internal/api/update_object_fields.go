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

type UpdateObjectFieldsRequest = struct {
	ObjectFields []input.PartialUpdateObjectField `json:"object_fields" validate:"required"`
}

type UpdateObjectFieldsResponse struct {
	ObjectFields []views.ObjectField `json:"object_fields"`
	Failures     []int64             `json:"failures"`
}

func (s ApiService) UpdateObjectFields(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.NewBadRequest("must setup organization first")
	}

	vars := mux.Vars(r)
	strObjectID, ok := vars["objectID"]
	if !ok {
		return errors.Newf("missing object ID request URL: %s", r.URL.RequestURI())
	}

	objectID, err := strconv.ParseInt(strObjectID, 10, 64)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)
	var requestBody UpdateObjectFieldsRequest
	if err := decoder.Decode(&requestBody); err != nil {
		return err
	}

	for _, requestItem := range requestBody.ObjectFields {
		validate := validator.New()
		err := validate.Struct(requestItem)
		if err != nil {
			return err
		}
	}

	objectFieldViews := []views.ObjectField{}
	failures := []int64{}
	for _, objectField := range requestBody.ObjectFields {
		updatedObjectField, err := objects.PartialUpdateObjectField(
			s.db,
			auth.Organization.ID,
			objectID,
			objectField,
		)
		if err == nil {
			updated := views.ConvertObjectField(updatedObjectField)
			objectFieldViews = append(objectFieldViews, updated)
		} else {
			failures = append(failures, objectField.ID)
		}
	}

	return json.NewEncoder(w).Encode(UpdateObjectFieldsResponse{
		ObjectFields: objectFieldViews,
		Failures:     failures,
	})
}
