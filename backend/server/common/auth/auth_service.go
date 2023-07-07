package auth

import (
	"net/http"

	"go.fabra.io/server/common/application"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/link_tokens"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/organizations"
	"go.fabra.io/server/common/repositories/sessions"
	"go.fabra.io/server/common/repositories/users"

	"gorm.io/gorm"
)

const SESSION_COOKIE_NAME = "X-Session-Token"

type AuthService interface {
	GetAuthentication(r *http.Request) (*Authentication, error)
	GetLinkAuthentication(r *http.Request) (*Authentication, error)
}

type AuthServiceImpl struct {
	db            *gorm.DB
	cryptoService crypto.CryptoService
}

func NewAuthService(db *gorm.DB, cryptoService crypto.CryptoService) AuthService {
	return AuthServiceImpl{
		db:            db,
		cryptoService: cryptoService,
	}
}

func addCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   application.IsProd(), // disable secure for local testing
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

func AddSessionCookie(w http.ResponseWriter, token string) {
	addCookie(w, SESSION_COOKIE_NAME, token)
}

func (as AuthServiceImpl) authenticateCookie(r *http.Request) (*Authentication, error) {
	cookie, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		if errors.IsCookieNotFound(err) {
			// this just means that the authentication failed, not unexpected
			return &Authentication{
				IsAuthenticated: false,
			}, nil
		} else {
			return nil, errors.Wrap(err, "(auth.authenticateCookie)")
		}
	}

	session, err := sessions.LoadValidByToken(as.db, cookie.Value)
	if err != nil {
		if errors.IsRecordNotFound(err) {
			// this just means that the authentication failed, not unexpected
			return &Authentication{
				IsAuthenticated: false,
			}, nil
		} else {
			return nil, errors.Wrap(err, "(auth.authenticateCookie)")
		}
	}

	refreshed, err := sessions.Refresh(as.db, session)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.authenticateCookie)")
	}

	user, err := users.LoadUserByID(as.db, session.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.authenticateCookie) Unexpected error fetching user")
	}

	if user.Blocked {
		return nil, errors.Wrap(errors.Forbidden, "(auth.authenticateCookie)")
	}

	// If organization is null, this means the user still needs to set their organization
	var organization *models.Organization
	if user.OrganizationID.Valid {
		organization, err = organizations.LoadOrganizationByID(as.db, user.OrganizationID.Int64)
		if err != nil {
			return nil, errors.Wrap(err, "(auth.authenticateCookie) Unexpected error fetching organization")
		}
	}

	return &Authentication{
		Session:         refreshed,
		User:            user,
		Organization:    organization,
		IsAuthenticated: true,
	}, nil
}

func (as AuthServiceImpl) authApiKey(r *http.Request) (*Authentication, error) {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		return &Authentication{
			IsAuthenticated: false,
		}, nil
	}

	hashedKey := crypto.HashString(apiKey)
	organization, err := organizations.LoadOrganizationByApiKey(as.db, hashedKey)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.authApiKey)")
	}

	return &Authentication{
		Organization:    organization,
		IsAuthenticated: true,
	}, nil
}

func (as AuthServiceImpl) authLinkToken(r *http.Request) (*Authentication, error) {
	linkToken := r.Header.Get("X-LINK-TOKEN")
	if linkToken == "" {
		return &Authentication{
			IsAuthenticated: false,
		}, nil
	}

	tokenInfo, err := link_tokens.ValidateLinkToken(linkToken)
	if err != nil {
		if errors.IsInvalidLinkToken(err) {
			return &Authentication{
				IsAuthenticated: false,
			}, nil
		} else {
			return nil, errors.Wrap(err, "(auth.authLinkToken)")
		}
	}

	organization, err := organizations.LoadOrganizationByID(as.db, tokenInfo.OrganizationID)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.authLinkToken)")
	}

	return &Authentication{
		LinkToken:       tokenInfo,
		Organization:    organization,
		IsAuthenticated: true,
	}, nil
}

func (as AuthServiceImpl) GetAuthentication(r *http.Request) (*Authentication, error) {
	authentication, err := as.authenticateCookie(r)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.GetAuthentication)")
	}

	if authentication.IsAuthenticated {
		return authentication, nil
	}

	// no session found check for API key authentication first
	authentication, err = as.authApiKey(r)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.GetAuthentication)")
	}
	if authentication.IsAuthenticated {
		return authentication, nil
	}

	return &Authentication{
		IsAuthenticated: false,
	}, nil
}

func (as AuthServiceImpl) GetLinkAuthentication(r *http.Request) (*Authentication, error) {
	// try link token first since some methods depends on it
	authentication, err := as.authLinkToken(r)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.GetLinkAuthentication)")
	}
	if authentication.IsAuthenticated {
		return authentication, nil
	}

	// some link authenticated routes should also work with regular authentication
	authentication, err = as.GetAuthentication(r)
	if err != nil {
		return nil, errors.Wrap(err, "(auth.GetLinkAuthentication)")
	}
	if authentication.IsAuthenticated {
		return authentication, nil
	}

	return &Authentication{
		IsAuthenticated: false,
	}, nil
}
