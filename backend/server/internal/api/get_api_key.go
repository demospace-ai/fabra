package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/api_keys"
)

func (s ApiService) GetApiKey(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("cannot request users without organization"), "(api.GetApiKey)")
	}

	apiKey, err := s.getOrCreateApiKey(auth.Organization.ID)
	if err != nil {
		return errors.Wrap(err, "(api.GetApiKey)")
	}

	_, err = fmt.Fprintf(w, *apiKey)
	if err != nil {
		return errors.Wrap(err, "(api.GetApiKey)")
	}

	return nil
}

func (s ApiService) getOrCreateApiKey(organizationID int64) (*string, error) {
	apiKey, err := api_keys.LoadApiKeyForOrganization(s.db, organizationID)
	if err != nil {
		// no api key found, so just generate one now
		if errors.IsRecordNotFound(err) {
			rawApiKey := generateKey()
			encryptedApiKey, err := s.cryptoService.EncryptApiKey(rawApiKey)
			if err != nil {
				return nil, errors.Wrap(err, "(api.GetOrCreateApiKey)")
			}

			_, err = api_keys.CreateApiKey(s.db, organizationID, *encryptedApiKey, crypto.HashString(rawApiKey))
			if err != nil {
				return nil, errors.Wrap(err, "(api.GetOrCreateApiKey)")
			}

			return &rawApiKey, nil
		} else {
			return nil, errors.Wrap(err, "(api.GetOrCreateApiKey)")
		}
	}

	return s.cryptoService.DecryptApiKey(apiKey.EncryptedKey)
}

func generateKey() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(randomBytes)
}
