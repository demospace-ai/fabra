package api

import (
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/internal/router"

	"gorm.io/gorm"
)

type ApiService struct {
	db            *gorm.DB
	authService   auth.AuthService
	cryptoService crypto.CryptoService
	queryService  query.QueryService
}

func NewApiService(db *gorm.DB, authService auth.AuthService, cryptoService crypto.CryptoService, queryService query.QueryService) ApiService {
	return ApiService{
		db:            db,
		authService:   authService,
		cryptoService: cryptoService,
		queryService:  queryService,
	}
}

func (s ApiService) AuthenticatedRoutes() []router.AuthenticatedRoute {
	return []router.AuthenticatedRoute{
		{
			Name:        "Check session",
			Method:      router.GET,
			Pattern:     "/check_session",
			HandlerFunc: s.CheckSession,
		},
		{
			Name:        "Logout",
			Method:      router.DELETE,
			Pattern:     "/logout",
			HandlerFunc: s.Logout,
		},
		{
			Name:        "Get API key",
			Method:      router.GET,
			Pattern:     "/api_key",
			HandlerFunc: s.GetApiKey,
		},
		{
			Name:        "Get all destinations",
			Method:      router.GET,
			Pattern:     "/destinations",
			HandlerFunc: s.GetDestinations,
		},
		{
			Name:        "Get destination",
			Method:      router.GET,
			Pattern:     "/destination/{destinationID}",
			HandlerFunc: s.GetDestination,
		},
		{
			Name:        "Get namespaces",
			Method:      router.GET,
			Pattern:     "/connection/namespaces",
			HandlerFunc: s.GetNamespaces,
		},
		{
			Name:        "Get tables for a given namespace",
			Method:      router.GET,
			Pattern:     "/connection/tables",
			HandlerFunc: s.GetTables,
		},
		{
			Name:        "Get schema for a given table",
			Method:      router.GET,
			Pattern:     "/connection/schema",
			HandlerFunc: s.GetSchema,
		},
		{
			Name:        "Get all syncs",
			Method:      router.GET,
			Pattern:     "/syncs",
			HandlerFunc: s.GetSyncs,
		},
		{
			Name:        "Get all syncs for a customer",
			Method:      router.GET,
			Pattern:     "/customer/{endCustomerId}/syncs",
			HandlerFunc: s.GetSyncsForCustomer,
		},
		{
			Name:        "Query object record for customer",
			Method:      router.POST,
			Pattern:     "/customer/{endCustomerId}/object/{objectId}/record",
			HandlerFunc: s.QueryObjectRecord,
		},
		{
			Name:        "Get all users",
			Method:      router.GET,
			Pattern:     "/users",
			HandlerFunc: s.GetAllUsers,
		},
		{
			Name:        "Create destination for sync",
			Method:      router.POST,
			Pattern:     "/destination",
			HandlerFunc: s.CreateDestination,
		},
		{
			Name:        "Create source for sync",
			Method:      router.POST,
			Pattern:     "/source",
			HandlerFunc: s.CreateSource,
		},
		{
			Name:        "Create object for sync",
			Method:      router.POST,
			Pattern:     "/object",
			HandlerFunc: s.CreateObject,
		},
		{
			Name:        "Update an object for sync",
			Method:      router.PATCH,
			Pattern:     "/object/{objectID}",
			HandlerFunc: s.UpdateObject,
		},
		{
			Name:        "Update object fields for sync",
			Method:      router.PATCH,
			Pattern:     "/object/{objectID}/object_fields",
			HandlerFunc: s.UpdateObjectFields,
		},
		{
			Name:        "Create sync",
			Method:      router.POST,
			Pattern:     "/sync",
			HandlerFunc: s.CreateSync,
		},
		{
			Name:        "Delete sync",
			Method:      router.DELETE,
			Pattern:     "/sync/{syncID}",
			HandlerFunc: s.DeleteSync,
		},
		{
			Name:        "Update sync",
			Method:      router.PATCH,
			Pattern:     "/sync/{syncID}",
			HandlerFunc: s.UpdateSync,
		},
		{
			Name:        "Run sync",
			Method:      router.POST,
			Pattern:     "/sync/{syncID}/run",
			HandlerFunc: s.RunSync,
		},
		{
			Name:        "Cancel sync run",
			Method:      router.DELETE,
			Pattern:     "/sync/{syncID}/run",
			HandlerFunc: s.CancelSyncRun,
		},
		{
			Name:        "Get sync",
			Method:      router.GET,
			Pattern:     "/sync/{syncID}",
			HandlerFunc: s.GetSync,
		},
		{
			Name:        "Create link token",
			Method:      router.POST,
			Pattern:     "/link_token",
			HandlerFunc: s.CreateLinkToken,
		},
		{
			Name:        "Get values for a specified field",
			Method:      router.GET,
			Pattern:     "/connection/field_values",
			HandlerFunc: s.GetFieldValues,
		},
		{
			Name:        "Set organization for user",
			Method:      router.POST,
			Pattern:     "/organization",
			HandlerFunc: s.SetOrganization,
		},
	}
}

func (s ApiService) UnauthenticatedRoutes() []router.UnauthenticatedRoute {
	return []router.UnauthenticatedRoute{
		{
			Name:        "Hello",
			Method:      router.GET,
			Pattern:     "/hello",
			HandlerFunc: s.Hello,
		},
		{
			Name:        "OAuth Redirect",
			Method:      router.GET,
			Pattern:     "/oauth_redirect",
			HandlerFunc: s.OAuthRedirect,
		},
		{
			Name:        "OAuth Login",
			Method:      router.GET,
			Pattern:     "/oauth_login",
			HandlerFunc: s.OAuthLogin,
		},
	}
}

// TODO: find a better way of doing this
// Only a subset of the APIs should be accessible with Link Token Authentication
func (s ApiService) LinkAuthenticatedRoutes() []router.LinkAuthenticatedRoute {
	return []router.LinkAuthenticatedRoute{
		{
			Name:        "Get all objects",
			Method:      router.GET,
			Pattern:     "/objects",
			HandlerFunc: s.GetObjects,
		},
		{
			Name:        "Get object",
			Method:      router.GET,
			Pattern:     "/object/{objectID}",
			HandlerFunc: s.GetObject,
		},
		{
			Name:        "Get all sources",
			Method:      router.GET,
			Pattern:     "/link/sources",
			HandlerFunc: s.LinkGetSources,
		},
		{
			Name:        "Get source namespaces",
			Method:      router.GET,
			Pattern:     "/link/namespaces",
			HandlerFunc: s.LinkGetNamespaces,
		},
		{
			Name:        "Get tables for a given source and namespace",
			Method:      router.GET,
			Pattern:     "/link/tables",
			HandlerFunc: s.LinkGetTables,
		},
		{
			Name:        "Get schema for a given table",
			Method:      router.GET,
			Pattern:     "/link/schema",
			HandlerFunc: s.LinkGetSchema,
		},
		{
			Name:        "Create source for sync",
			Method:      router.POST,
			Pattern:     "/link/source",
			HandlerFunc: s.LinkCreateSource,
		},
		{
			Name:        "Create sync",
			Method:      router.POST,
			Pattern:     "/link/sync",
			HandlerFunc: s.LinkCreateSync,
		},
		{
			Name:        "Delete sync",
			Method:      router.DELETE,
			Pattern:     "/link/sync/{syncID}",
			HandlerFunc: s.LinkDeleteSync,
		},
		{
			Name:        "Update sync",
			Method:      router.PATCH,
			Pattern:     "/sync/{syncID}",
			HandlerFunc: s.LinkUpdateSync,
		},
		{
			Name:        "Run sync",
			Method:      router.POST,
			Pattern:     "/link/sync/{syncID}/run",
			HandlerFunc: s.LinkRunSync,
		},
		{
			Name:        "Cancel sync run",
			Method:      router.DELETE,
			Pattern:     "/link/sync/{syncID}/run",
			HandlerFunc: s.LinkCancelSyncRun,
		},
		{
			Name:        "Get sync",
			Method:      router.GET,
			Pattern:     "/link/sync/{syncID}",
			HandlerFunc: s.LinkGetSync,
		},
		{
			Name:        "Get all syncs",
			Method:      router.GET,
			Pattern:     "/link/syncs",
			HandlerFunc: s.LinkGetSyncs,
		},
		{
			Name:        "Get preview",
			Method:      router.POST,
			Pattern:     "/link/preview",
			HandlerFunc: s.LinkGetPreview,
		},
		{
			Name:        "Test data connection",
			Method:      router.POST,
			Pattern:     "/connection/test",
			HandlerFunc: s.TestDataConnection,
		},
	}
}
