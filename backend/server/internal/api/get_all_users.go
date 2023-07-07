package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	user_repository "go.fabra.io/server/common/repositories/users"
)

type GetAllUsersResponse struct {
	Users []models.User `json:"users"`
}

func (s ApiService) GetAllUsers(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("cannot request users without organization"), "(api.GetAllUsers)")
	}

	users, err := user_repository.LoadAllByOrganizationID(s.db, auth.Organization.ID)
	if err != nil {
		return errors.Wrap(err, "(api.GetAllUsers)")
	}

	return json.NewEncoder(w).Encode(GetAllUsersResponse{
		Users: users,
	})
}
