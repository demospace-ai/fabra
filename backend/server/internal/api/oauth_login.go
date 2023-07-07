package api

import (
	"log"
	"net/http"
	"strings"

	"go.fabra.io/server/common/application"
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/oauth"
	"go.fabra.io/server/common/repositories/sessions"
	"go.fabra.io/server/common/repositories/users"
)

var UNAUTHORIZED_DOMAINS = map[string]bool{
	"gmail.com":     true,
	"outlook.com":   true,
	"icloud.com":    true,
	"yahoo.com":     true,
	"aol.com":       true,
	"hotmail.com":   true,
	"hey.com":       true,
	"supaglue.com":  true,
	"merge.dev":     true,
	"vessel.land":   true,
	"hightouch.com": true,
	"getcensus.com": true,
	"airbyte.io":    true,
	"airbyte.com":   true,
	"fivetran.com":  true,
}

func (s ApiService) OAuthLogin(w http.ResponseWriter, r *http.Request) error {
	if !r.URL.Query().Has("state") {
		return errors.Newf("(api.OAuthLogin) missing state from OAuth Login request URL: %s", r.URL.RequestURI())
	}

	if !r.URL.Query().Has("code") {
		return errors.Newf("(api.OAuthLogin) missing code from OAuth Login request URL: %s", r.URL.RequestURI())
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	provider, err := oauth.ValidateState(state)
	if err != nil {
		return errors.Wrap(err, "(api.OAuthLogin)")
	}

	var externalUserInfo *oauth.ExternalUserInfo
	switch *provider {
	case oauth.OauthProviderGoogle:
		externalUserInfo, err = oauth.FetchGoogleInfo(code)
	case oauth.OauthProviderGithub:
		externalUserInfo, err = oauth.FetchGithubInfo(code)
	default:
		return errors.Newf("(api.OAuthLogin) unexpected provider %s", *provider)
	}
	if err != nil {
		return errors.Wrap(err, "(api.OAuthLogin)")
	}

	// separately check for existing user to bypass allowlist for domains
	user, err := users.LoadByExternalID(s.db, externalUserInfo.ExternalID)
	if err != nil && !errors.IsRecordNotFound(err) {
		return errors.Wrap(err, "(api.OAuthLogin)")
	}

	// no user exists yet, so if the domain is not allowed then redirect
	if user == nil {
		var userEmailDomain = strings.Split(externalUserInfo.Email, "@")[1]
		// allow unauthorized domains in development
		if _, unauthorized := UNAUTHORIZED_DOMAINS[userEmailDomain]; unauthorized && application.IsProd() {
			log.Printf("Unauthorized login: %v", externalUserInfo)
			http.Redirect(w, r, getUnauthorizedRedirect(), http.StatusFound)
			return nil
		}

		user, err = users.CreateUserForExternalInfo(s.db, externalUserInfo)
		if err != nil {
			return errors.Wrap(err, "(api.OAuthLogin)")
		}
	}

	sessionToken, err := sessions.Create(s.db, user.ID)
	if err != nil {
		return errors.Wrap(err, "(api.OAuthLogin)")
	}

	auth.AddSessionCookie(w, *sessionToken)
	http.Redirect(w, r, getOauthSuccessRedirect(), http.StatusFound)

	return nil
}

func getOauthSuccessRedirect() string {
	if application.IsProd() {
		return "https://app.fabra.io"
	} else {
		return "http://localhost:3000"
	}
}

func getUnauthorizedRedirect() string {
	if application.IsProd() {
		return "https://app.fabra.io/unauthorized"
	} else {
		return "http://localhost:3000/unauthorized"
	}
}
