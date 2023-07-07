package models

import "go.fabra.io/server/common/oauth"

type ExternalProfile struct {
	ExternalID    string
	OauthProvider oauth.OauthProvider
	UserID        int64

	BaseModel
}
