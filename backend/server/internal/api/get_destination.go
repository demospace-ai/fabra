package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/connections"
	"go.fabra.io/server/common/repositories/destinations"
	"go.fabra.io/server/common/views"
)

type GetDestinationResponse struct {
	Destination views.Destination `json:"destination"`
}

func (s ApiService) GetDestination(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	strDestinationId, ok := vars["destinationID"]
	if !ok {
		return errors.Newf("(api.GetDestination) missing destination ID from GetDestination request URL: %s", r.URL.RequestURI())
	}

	destinationId, err := strconv.ParseInt(strDestinationId, 10, 64)
	if err != nil {
		return errors.Wrap(err, "(api.GetDestination)")
	}

	destination, err := destinations.LoadDestinationByID(s.db, auth.Organization.ID, destinationId)
	if err != nil {
		return errors.Wrap(err, "(api.GetDestination)")
	}

	connection, err := connections.LoadConnectionByID(s.db, auth.Organization.ID, destination.ConnectionID)
	if err != nil {
		return errors.Wrap(err, "(api.GetDestination)")
	}

	var destinationView views.Destination
	if connection.ConnectionType == models.ConnectionTypeWebhook {
		webhookSigningKey, err := s.cryptoService.DecryptWebhookSigningKey(connection.Credentials.String)
		if err != nil {
			return errors.Wrap(err, "(api.GetDestination)")
		}

		destinationView = views.ConvertWebhook(*destination, *connection, webhookSigningKey)
	} else {
		destinationView = views.ConvertDestination(*destination, *connection)
	}

	return json.NewEncoder(w).Encode(GetDestinationResponse{
		destinationView,
	})
}
