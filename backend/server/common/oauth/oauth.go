package oauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/github"
	"go.fabra.io/server/common/application"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/secret"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	githuboauth "golang.org/x/oauth2/github"
	googleoauth "golang.org/x/oauth2/google"
)

const GITHUB_PRODUCTION_CLIENT_ID = "7eff3cfd664e1e01e19b"
const GITHUB_PRODUCTION_SECRET_KEY = "projects/932264813910/secrets/github-prod-client-secret/versions/latest"
const GITHUB_DEVELOPMENT_CLIENT_ID = "f84f670b7af18144af4a"
const GITHUB_DEVELOPMENT_SECRET_KEY = "projects/86315250181/secrets/github-dev-client-secret/versions/latest"

const GOOGLE_PRODUCTION_CLIENT_ID = "932264813910-egpk1omo3v2cedd89k8go851uko6djpa.apps.googleusercontent.com"
const GOOGLE_PRODUCTION_SECRET_KEY = "projects/932264813910/secrets/google-prod-client-secret/versions/latest"
const GOOGLE_DEVELOPMENT_CLIENT_ID = "86315250181-v19knnmf486fb5nebm2b47hu454abvet.apps.googleusercontent.com"
const GOOGLE_DEVELOPMENT_SECRET_KEY = "projects/86315250181/secrets/google-dev-client-secret/versions/latest"

type StateClaims struct {
	Provider OauthProvider `json:"provider"`
	jwt.RegisteredClaims
}

type OauthProvider string

const (
	OauthProviderGoogle  OauthProvider = "google"
	OauthProviderGithub  OauthProvider = "github"
	OauthProviderUnknown OauthProvider = "unknown"
)

type ExternalUserInfo struct {
	ExternalID    string
	OauthProvider OauthProvider
	Email         string
	Name          string
}

func FetchGithubInfo(code string) (*ExternalUserInfo, error) {
	secretKey := getGithubSecretKey()
	githubClientSecret, err := secret.FetchSecret(context.TODO(), secretKey)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGithubInfo) fetching secret")
	}

	oauthConf := &oauth2.Config{
		ClientID:     getGithubClientID(),
		ClientSecret: *githubClientSecret,
		Scopes:       []string{"user:email", "read:user"},
		RedirectURL:  getOauthRedirectUrl(),
		Endpoint:     githuboauth.Endpoint,
	}

	token, err := oauthConf.Exchange(context.TODO(), code)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGithubInfo) exchanging code for token")
	}

	oauthClient := oauthConf.Client(context.TODO(), token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(context.TODO(), "") // empty string will get info for authenticated user
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGithubInfo) extracting user")
	}

	emails, _, err := client.Users.ListEmails(context.TODO(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGithubInfo) extracting email")
	}

	var primaryEmail string
	for _, email := range emails {
		if *email.Primary {
			primaryEmail = email.GetEmail()
		}
	}

	if primaryEmail == "" {
		return nil, errors.New("(oauth.FetchGithubInfo) could not find primary email")
	}

	return &ExternalUserInfo{
		ExternalID:    fmt.Sprintf("%d", user.GetID()),
		OauthProvider: OauthProviderGithub,
		Email:         primaryEmail,
		Name:          user.GetName(),
	}, nil
}

func FetchGoogleInfo(code string) (*ExternalUserInfo, error) {
	secretKey := getGoogleSecretKey()
	googleClientSecret, err := secret.FetchSecret(context.TODO(), secretKey)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGoogleInfo) fetching secret")
	}

	clientId := getGoogleClientID()
	oauthConf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: *googleClientSecret,
		Scopes:       []string{"email", "profile", "openid"},
		RedirectURL:  getOauthRedirectUrl(),
		Endpoint:     googleoauth.Endpoint,
	}

	oauth2Token, err := oauthConf.Exchange(context.TODO(), code)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGoogleInfo) exchanging code for token")
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.Newf("(oauth.FetchGoogleInfo) no id_token included in token exchange response: %+v", oauth2Token)
	}

	validator, err := idtoken.NewValidator(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGoogleInfo) creating validator")
	}

	payload, err := validator.Validate(context.TODO(), rawIDToken, clientId)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.FetchGoogleInfo) validating token")
	}

	return &ExternalUserInfo{
		ExternalID:    payload.Subject,
		OauthProvider: OauthProviderGoogle,
		Email:         payload.Claims["email"].(string),
		Name:          payload.Claims["name"].(string),
	}, nil
}

func GetOauthRedirect(strProvider string) (*string, error) {
	provider := getOAuthProvider(strProvider)

	var oauthConf *oauth2.Config
	switch provider {
	case OauthProviderGoogle:
		oauthConf = &oauth2.Config{
			ClientID:    getGoogleClientID(),
			Scopes:      []string{"email", "profile", "openid"},
			RedirectURL: getOauthRedirectUrl(),
			Endpoint:    googleoauth.Endpoint,
		}
	case OauthProviderGithub:
		oauthConf = &oauth2.Config{
			ClientID:    getGithubClientID(),
			Scopes:      []string{"user:email", "read:user"},
			RedirectURL: getOauthRedirectUrl(),
			Endpoint:    githuboauth.Endpoint,
		}
	default:
		return nil, errors.Newf("unsupported login method: %s", strProvider)
	}

	token := jwt.NewWithClaims(crypto.SigningMethodKMSHS256, StateClaims{
		provider,
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})

	signedString, err := token.SignedString(nil)
	if err != nil {
		return nil, errors.Wrap(err, "(oauth.GetOauthRedirect) signing token")
	}

	url := oauthConf.AuthCodeURL(signedString, oauth2.AccessTypeOnline, oauth2.ApprovalForce)

	return &url, nil
}

func ValidateState(state string) (*OauthProvider, error) {
	token, err := jwt.ParseWithClaims(state, &StateClaims{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil // no key needs to be fetched— we just call the GCP KMS endpoint
	})

	if err != nil {
		return nil, errors.Wrap(err, "(oauth.ValidateState) parsing token")
	}

	if !token.Valid {
		return nil, errors.Newf("token invalid: %v", token.Raw)
	}

	claims, ok := token.Claims.(*StateClaims)
	if !ok {
		return nil, errors.Newf("token invalid: %v", token.Raw)
	}

	return &claims.Provider, nil
}

func getGithubSecretKey() string {
	if application.IsProd() {
		return GITHUB_PRODUCTION_SECRET_KEY
	} else {
		return GITHUB_DEVELOPMENT_SECRET_KEY
	}
}

func getGithubClientID() string {
	if application.IsProd() {
		return GITHUB_PRODUCTION_CLIENT_ID
	} else {
		return GITHUB_DEVELOPMENT_CLIENT_ID
	}
}

func getGoogleSecretKey() string {
	if application.IsProd() {
		return GOOGLE_PRODUCTION_SECRET_KEY
	} else {
		return GOOGLE_DEVELOPMENT_SECRET_KEY
	}
}

func getGoogleClientID() string {
	if application.IsProd() {
		return GOOGLE_PRODUCTION_CLIENT_ID
	} else {
		return GOOGLE_DEVELOPMENT_CLIENT_ID
	}
}

func getOauthRedirectUrl() string {
	if application.IsProd() {
		return "https://api.fabra.io/oauth_login"
	} else {
		return "http://localhost:8080/oauth_login"
	}
}

func getOAuthProvider(strProvider string) OauthProvider {
	switch strings.ToLower(strProvider) {
	case "google":
		return OauthProviderGoogle
	case "github":
		return OauthProviderGithub
	default:
		return OauthProviderUnknown
	}
}
