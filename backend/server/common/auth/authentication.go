package auth

import (
	"go.fabra.io/server/common/link_tokens"
	"go.fabra.io/server/common/models"
)

type Authentication struct {
	Session         *models.Session
	User            *models.User
	Organization    *models.Organization
	LinkToken       *link_tokens.TokenInfo
	IsAuthenticated bool
}
