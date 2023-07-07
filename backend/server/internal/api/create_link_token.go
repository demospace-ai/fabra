package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/link_tokens"
	"go.fabra.io/server/common/repositories/webhooks"
	"gorm.io/gorm"
)

type CreateLinkTokenRequest struct {
	EndCustomerID  string             `json:"end_customer_id" validate:"required"`
	DestinationIDs []int64            `json:"destination_ids"`
	WebhookData    *input.WebhookData `json:"webhook_data,omitempty"`
}

type CreateLinkTokenResponse struct {
	LinkToken string `json:"link_token"`
}

func (s ApiService) CreateLinkToken(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("cannot request users without organization"), "(api.CreateLinkToken)")
	}

	decoder := json.NewDecoder(r.Body)
	var createLinkTokenRequest CreateLinkTokenRequest
	err := decoder.Decode(&createLinkTokenRequest)
	if err != nil {
		return errors.Wrap(err, "(api.CreateLinkToken)")
	}

	signedToken, err := link_tokens.CreateLinkToken(link_tokens.TokenInfo{
		OrganizationID: auth.Organization.ID,
		EndCustomerID:  createLinkTokenRequest.EndCustomerID,
		DestinationIDs: createLinkTokenRequest.DestinationIDs,
	})
	if err != nil {
		return errors.Wrap(err, "(api.CreateLinkToken)")
	}

	if createLinkTokenRequest.WebhookData != nil {
		if createLinkTokenRequest.WebhookData.EndCustomerApiKey != nil {
			encryptedEndCustomerApiKey, err := s.cryptoService.EncryptEndCustomerApiKey(*createLinkTokenRequest.WebhookData.EndCustomerApiKey)
			if err != nil {
				return errors.Wrap(err, "(api.CreateLinkToken)")
			}

			// this operation always replaces the existing api key
			err = s.db.Transaction(func(tx *gorm.DB) error {
				err = webhooks.DeactivateExistingEndCustomerApiKey(tx, auth.Organization.ID, createLinkTokenRequest.EndCustomerID)
				if err != nil {
					return errors.Wrap(err, "(api.CreateLinkToken)")
				}

				err = webhooks.CreateEndCustomerApiKey(tx, auth.Organization.ID, createLinkTokenRequest.EndCustomerID, *encryptedEndCustomerApiKey)
				if err != nil {
					return errors.Wrap(err, "(api.CreateLinkToken)")
				}

				return nil
			})

			if err != nil {
				return errors.Wrap(err, "(api.CreateLinkToken)")
			}
		}
	}

	return json.NewEncoder(w).Encode(CreateLinkTokenResponse{
		LinkToken: *signedToken,
	})
}
