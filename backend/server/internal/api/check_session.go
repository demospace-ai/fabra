package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/intercom"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/organizations"
	"go.fabra.io/server/common/views"
)

type CheckSessionResponse struct {
	User                   views.User            `json:"user"`
	Organization           *models.Organization  `json:"organization"`
	SuggestedOrganizations []models.Organization `json:"suggested_organizations"`
}

func (s ApiService) CheckSession(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	var suggestedOrganizations []models.Organization
	var err error
	if !auth.User.OrganizationID.Valid {
		var userEmailDomain = strings.Split(auth.User.Email, "@")[1]

		// Don't suggest organizations via domain for common domain names
		if _, unauthorized := UNAUTHORIZED_DOMAINS[userEmailDomain]; !unauthorized {

			// TODO: hack to associate watershed.com emails with watershedclimate.com (since that's how it was created)
			if userEmailDomain == "watershed.com" {
				userEmailDomain = "watershedclimate.com"
			}

			suggestedOrganizations, err = organizations.LoadOrganizationsByEmailDomain(s.db, userEmailDomain)
			if err != nil {
				return errors.Wrap(err, "(api.CheckSession)")
			}
		}
	}

	intercomHash, err := intercom.GenerateIntercomHash(*auth.User)
	if err != nil {
		return errors.Wrap(err, "(api.CheckSession)")
	}

	return json.NewEncoder(w).Encode(CheckSessionResponse{
		User:                   views.ConvertUser(*auth.User, *intercomHash),
		Organization:           auth.Organization,
		SuggestedOrganizations: suggestedOrganizations,
	})
}
