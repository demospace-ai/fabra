package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/sync/temporal"
)

func (s ApiService) DeleteSync(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.DeleteSync)")
	}

	vars := mux.Vars(r)
	strSyncId, ok := vars["syncID"]
	if !ok {
		return errors.Wrap(errors.NewBadRequestf("missing sync ID from DeleteSync request URL: %s", r.URL.RequestURI()), "(api.DeleteSync)")
	}

	syncId, err := strconv.ParseInt(strSyncId, 10, 64)
	if err != nil {
		return errors.Wrap(err, "(api.DeleteSync)")
	}

	// check the sync belongs to the right organization
	sync, err := syncs.LoadSyncByID(s.db, auth.Organization.ID, syncId)
	if err != nil {
		return errors.Wrap(err, "(api.DeleteSync) loading sync")
	}

	err = syncs.DeactivateSyncByID(s.db, sync.ID)
	if err != nil {
		return errors.Wrap(err, "(api.DeleteSync) deactivating sync")
	}

	c, err := temporal.CreateClient(CLIENT_PEM_KEY, CLIENT_KEY_KEY)
	if err != nil {
		return errors.Wrap(err, "(api.DeleteSync) creating temporal client")
	}
	defer c.Close()

	ctx := context.TODO()
	scheduleClient := c.ScheduleClient()
	schedule := scheduleClient.GetHandle(ctx, sync.WorkflowID)
	err = schedule.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "(api.DeleteSync) deleting temporal schedule")
	}

	return nil
}
