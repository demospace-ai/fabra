package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/destinations"
	"go.fabra.io/server/common/views"
)

type GetDestinationsResponse struct {
	Destinations []views.Destination `json:"destinations"`
}

func (s ApiService) GetDestinations(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {

	destinations, err := destinations.LoadAllDestinations(s.db, auth.Organization.ID)
	if err != nil {
		return errors.Wrap(err, "(api.GetDestinations)")
	}

	return json.NewEncoder(w).Encode(GetDestinationsResponse{
		views.ConvertDestinationConnections(destinations),
	})
}
