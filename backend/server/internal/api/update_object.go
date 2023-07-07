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

type UpdateObjectResponse struct {
	Object views.Object `json:"object"`
}

func (s ApiService) UpdateObject(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
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
	var updateObjectRequest input.PartialUpdateObjectInput

	if err := decoder.Decode(&updateObjectRequest); err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(updateObjectRequest)
	if err != nil {
		return err
	}

	object, err := objects.PartialUpdateObject(
		s.db,
		auth.Organization.ID,
		objectID,
		updateObjectRequest,
	)
	if err != nil {
		return err
	}

	objectFields, err := objects.LoadObjectFieldsByID(s.db, object.ID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(UpdateObjectResponse{
		views.ConvertObject(object, objectFields),
	})
}
