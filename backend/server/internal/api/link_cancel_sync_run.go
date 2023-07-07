package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/sync_runs"
	"go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/sync/temporal"
)

func (s ApiService) LinkCancelSyncRun(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.LinkCancelSync)")
	}

	if auth.LinkToken == nil {
		return errors.Wrap(errors.NewBadRequest("must send link token"), "(api.LinkCancelSync)")
	}

	vars := mux.Vars(r)
	strSyncId, ok := vars["syncID"]
	if !ok {
		return errors.Newf("(api.LinkCancelSync) missing sync ID from CancelSync request URL: %s", r.URL.RequestURI())
	}

	syncId, err := strconv.ParseInt(strSyncId, 10, 64)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCancelSync)")
	}

	// check the sync belongs to the right organization and customer
	sync, err := syncs.LoadSyncByIDAndCustomer(s.db, auth.Organization.ID, auth.LinkToken.EndCustomerID, syncId)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCancelSync) loading sync")
	}

	syncRun, err := sync_runs.LoadActiveRunBySyncID(s.db, sync.ID)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCancelSync) loading sync run")
	}

	c, err := temporal.CreateClient(CLIENT_PEM_KEY, CLIENT_KEY_KEY)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCancelSync) creating client")
	}
	defer c.Close()

	ctx := context.TODO()
	err = c.CancelWorkflow(
		ctx,
		syncRun.WorkflowID,
		"", // Empty RunID will result in the currently running workflow to be cancelled
	)
	if err != nil {
		return errors.Wrap(err, "(api.LinkCancelSync) cancelling workflow")
	}

	return nil
}
