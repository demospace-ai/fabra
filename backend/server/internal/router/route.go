package router

import (
	"net/http"

	"go.fabra.io/server/common/auth"
)

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
	PATCH
)

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	default:
		return ""
	}
}

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request) error
type AuthenticatedHandlerFunc func(auth.Authentication, http.ResponseWriter, *http.Request) error
type AuthenticatedRoute struct {
	Name        string
	Method      Method
	Pattern     string
	HandlerFunc AuthenticatedHandlerFunc
}

type UnauthenticatedRoute struct {
	Name        string
	Method      Method
	Pattern     string
	HandlerFunc ErrorHandlerFunc
}

type LinkAuthenticatedRoute struct {
	Name        string
	Method      Method
	Pattern     string
	HandlerFunc AuthenticatedHandlerFunc
}
