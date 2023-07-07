package sync_runs

import (
	"time"

	"go.fabra.io/server/common/database"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"

	"gorm.io/gorm"
)

func createSyncRun(
	db *gorm.DB,
	organizationID int64,
	syncID int64,
	workflowID string,
) (*models.SyncRun, error) {
	newSyncRun := models.SyncRun{
		OrganizationID: organizationID,
		SyncID:         syncID,
		Status:         models.SyncRunStatusRunning,
		StartedAt:      time.Now(),
		WorkflowID:     workflowID,
	}

	result := db.Create(&newSyncRun)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sync_runs.CreateSyncRun)")
	}

	return &newSyncRun, nil
}

// This function exists so that RecordStatusActivity can be idempotent
func CreateOrStartSyncRun(
	db *gorm.DB,
	organizationID int64,
	syncID int64,
	workflowID string,
) (*models.SyncRun, error) {
	syncRun, err := LoadActiveByWorkflowID(db, workflowID)
	if err != nil && !errors.IsRecordNotFound(err) {
		return nil, errors.Wrap(err, "CreateOrStartSyncRun")
	} else if err == nil {
		return UpdateSyncRun(db, syncRun, models.SyncRunStatusRunning, nil, nil)
	} else {
		// Didn't find an active sync run, so create a new one
		return createSyncRun(db, organizationID, syncID, workflowID)
	}
}

func UpdateSyncRun(db *gorm.DB, syncRun *models.SyncRun, newStatus models.SyncRunStatus, syncError *string, rowsWritten *int) (*models.SyncRun, error) {
	updates := models.SyncRun{
		CompletedAt: time.Now(),
		Status:      newStatus,
	}

	if rowsWritten != nil {
		updates.RowsWritten = *rowsWritten
	}

	if syncError != nil {
		updates.Error = database.NewNullString(*syncError)
	}

	result := db.Model(syncRun).Updates(updates)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sync_runs.CompleteSyncRun)")
	}

	return syncRun, nil
}

// Temporal guarantees that only one active workflow execution will have a given workflow ID, so we check the status
// to be doubly sure we have the right sync run even though the workflow IDs should always be unique
func LoadActiveByWorkflowID(db *gorm.DB, workflowID string) (*models.SyncRun, error) {
	var syncRun models.SyncRun
	var sync models.Sync
	result := db.Table("sync_runs").
		Select("sync_runs.*").
		Where("sync_runs.workflow_id = ?", workflowID).
		// We filter for active sync runs to be sure we have the right one
		Where("sync_runs.status = ?", string(models.SyncRunStatusRunning)).
		Where("sync_runs.deactivated_at IS NULL").
		Take(&sync)
	if result.Error != nil {
		return nil, result.Error
	}

	return &syncRun, nil
}

func LoadActiveRunBySyncID(db *gorm.DB, syncID int64) (*models.SyncRun, error) {
	var syncRun models.SyncRun
	var sync models.Sync
	result := db.Table("sync_runs").
		Select("sync_runs.*").
		Where("sync_runs.sync_id = ?", syncID).
		// We filter for active sync runs to be sure we have the right one
		Where("sync_runs.status = ?", string(models.SyncRunStatusRunning)).
		Where("sync_runs.deactivated_at IS NULL").
		Take(&sync)
	if result.Error != nil {
		return nil, result.Error
	}

	return &syncRun, nil
}

func LoadAllRunsForSync(db *gorm.DB, organizationID int64, syncID int64) ([]models.SyncRun, error) {
	var syncRuns []models.SyncRun
	result := db.Table("sync_runs").
		Select("sync_runs.*").
		Where("sync_runs.organization_id = ?", organizationID).
		Where("sync_runs.sync_id = ?", syncID).
		Where("sync_runs.deactivated_at IS NULL").
		Order("sync_runs.created_at DESC").
		Find(&syncRuns)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sync_runs.LoadAllRunsForSync)")
	}

	return syncRuns, nil
}
