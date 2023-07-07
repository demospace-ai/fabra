package temporal

import (
	"context"

	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/sync_runs"
)

type UpdateType string

const (
	UpdateTypeCreate   UpdateType = "create"
	UpdateTypeComplete UpdateType = "complete"
)

type RecordStatusInput struct {
	OrganizationID int64
	SyncID         int64
	SyncRun        models.SyncRun
	WorkflowID     string
	UpdateType     UpdateType
	NewStatus      models.SyncRunStatus
	RowsWritten    int
	Error          *string
}

func (a *Activities) RecordStatus(ctx context.Context, input RecordStatusInput) (*models.SyncRun, error) {
	switch input.UpdateType {
	case UpdateTypeCreate:
		// This is a no-op if the sync run already exists
		return sync_runs.CreateOrStartSyncRun(a.Db, input.OrganizationID, input.SyncID, input.WorkflowID)
	case UpdateTypeComplete:
		return sync_runs.UpdateSyncRun(a.Db, &input.SyncRun, input.NewStatus, input.Error, &input.RowsWritten)
	default:
		return nil, errors.Newf("unexpected update type: %s", input.UpdateType)
	}
}
