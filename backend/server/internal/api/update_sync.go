package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/sync/temporal"
	"go.temporal.io/sdk/client"
)

type UpdateSyncRequest struct {
	Status models.SyncStatus `json:"status"`
}

func (s ApiService) UpdateSync(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.UpdateSync)")
	}

	vars := mux.Vars(r)
	strSyncId, ok := vars["syncID"]
	if !ok {
		return errors.Wrap(errors.NewBadRequestf("missing sync ID from UpdateSync request URL: %s", r.URL.RequestURI()), "(api.UpdateSync)")
	}

	syncId, err := strconv.ParseInt(strSyncId, 10, 64)
	if err != nil {
		return errors.Wrap(err, "(api.UpdateSync)")
	}

	var updateSyncRequest UpdateSyncRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateSyncRequest); err != nil {
		return err
	}

	// check the sync belongs to the right organization
	sync, err := syncs.LoadSyncByID(s.db, auth.Organization.ID, syncId)
	if err != nil {
		return errors.Wrap(err, "(api.UpdateSync) loading sync")
	}

	// TODO: allow updating more than just the status
	err = syncs.UpdateSyncStatusByID(s.db, syncId, updateSyncRequest.Status)
	if err != nil {
		return errors.Wrap(err, "(api.UpdateSync) updating sync status")
	}

	c, err := temporal.CreateClient(CLIENT_PEM_KEY, CLIENT_KEY_KEY)
	if err != nil {
		return errors.Wrap(err, "(api.UpdateSync) creating client")
	}
	defer c.Close()

	ctx := context.TODO()
	scheduleClient := c.ScheduleClient()
	schedule := scheduleClient.GetHandle(ctx, sync.WorkflowID)
	if updateSyncRequest.Status == models.SyncStatusActive {
		err = schedule.Unpause(ctx, client.ScheduleUnpauseOptions{Note: "Updated by customer"})
		if err != nil {
			return errors.Wrap(err, "(api.UpdateSync) resuming temporal schedule")
		}
	} else {
		err = schedule.Pause(ctx, client.SchedulePauseOptions{Note: "Updated by customer"})
		if err != nil {
			return errors.Wrap(err, "(api.UpdateSync) pausing temporal schedule")
		}
	}

	return nil
}
