package api

import (
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/repositories/sessions"
)

func (s ApiService) Logout(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if !auth.IsAuthenticated {
		return nil
	}

	return sessions.Clear(s.db, auth.Session)
}
