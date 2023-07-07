package temporal

import (
	"time"

	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const SyncTaskQueue = "SYNC_TASK_QUEUE"

type SyncInput struct {
	OrganizationID int64
	SyncID         int64
}

var FETCH_OPTIONS = workflow.ActivityOptions{
	StartToCloseTimeout: time.Minute * 2,
	RetryPolicy: &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	},
}

var RECORD_OPTIONS = workflow.ActivityOptions{
	StartToCloseTimeout: time.Minute * 3,
	RetryPolicy: &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	},
}

var REPLICATE_OPTIONS = workflow.ActivityOptions{
	StartToCloseTimeout: time.Hour * 24,
	HeartbeatTimeout:    time.Minute * 5,
	RetryPolicy: &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        time.Minute,
		MaximumAttempts:        3,
		NonRetryableErrorTypes: []string{"CustomerVisibleError"},
	},
}

var CURSOR_OPTIONS = workflow.ActivityOptions{
	StartToCloseTimeout: time.Minute * 3,
	RetryPolicy: &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	},
}

func SyncWorkflow(ctx workflow.Context, input SyncInput) error {
	var a *Activities // Temporal handles calling the registered activity object

	recordCtx := workflow.WithActivityOptions(ctx, RECORD_OPTIONS)
	fetchCtx := workflow.WithActivityOptions(ctx, FETCH_OPTIONS)
	replicateCtx := workflow.WithActivityOptions(ctx, REPLICATE_OPTIONS)
	cursorCtx := workflow.WithActivityOptions(ctx, CURSOR_OPTIONS)

	var syncRun models.SyncRun
	err := workflow.ExecuteActivity(recordCtx, a.RecordStatus, RecordStatusInput{
		OrganizationID: input.OrganizationID,
		SyncID:         input.SyncID,
		WorkflowID:     workflow.GetInfo(ctx).WorkflowExecution.ID,
		UpdateType:     UpdateTypeCreate,
	}).Get(recordCtx, &syncRun)
	if err != nil {
		return errors.Wrap(err, "(workflow.RecordStatus)")
	}

	var syncConfig SyncConfig
	fetchInput := FetchConfigInput(input)
	err = workflow.ExecuteActivity(fetchCtx, a.FetchConfig, fetchInput).Get(fetchCtx, &syncConfig)
	if err != nil {
		// Ignore the error returned here. It is logged by Temporal as the activity task
		// failing, and the reason for the workflow failing is the original error
		recordFailure(recordCtx, err, syncRun)
		return errors.Wrap(err, "(workflow.FetchConfig)")
	}

	var replicateOutput ReplicateOutput
	replicateInput := ReplicateInput(syncConfig)
	err = workflow.ExecuteActivity(replicateCtx, a.Replicate, replicateInput).Get(replicateCtx, &replicateOutput)
	if err != nil {
		// Ignore the error returned here. It is logged by Temporal as the activity task
		// failing, and the reason for the workflow failing is the original error
		recordFailure(recordCtx, err, syncRun)
		return errors.Wrap(err, "(workflow.Replicate)")
	}

	if syncConfig.Sync.SyncMode.UsesCursor() && replicateOutput.CursorPosition != nil {
		cursorInput := UpdateCursorInput{
			Sync:           syncConfig.Sync,
			CursorPosition: *replicateOutput.CursorPosition,
		}
		err = workflow.ExecuteActivity(cursorCtx, a.UpdateCursor, cursorInput).Get(cursorCtx, nil)
		if err != nil {
			// Ignore the error returned here. It is logged by Temporal as the activity task
			// failing, and the reason for the workflow failing is the original error
			recordFailure(recordCtx, err, syncRun)
			return errors.Wrap(err, "(workflow.UpdateCursor)")
		}
	}

	return recordSuccess(recordCtx, syncRun, replicateOutput.RowsWritten)
}

func recordFailure(ctx workflow.Context, err error, syncRun models.SyncRun) error {
	var applicationErr *temporal.ApplicationError
	var errString string
	if errors.As(err, &applicationErr) && applicationErr.Type() == "CustomerVisibleError" {
		// Interceptor will update the error message to only include the CustomerVisisbleError message
		errString = applicationErr.Message()
	} else {
		errString = "unexpected error"
	}

	var a *Activities // Temporal handles calling the registered activity object
	return workflow.ExecuteActivity(ctx, a.RecordStatus, RecordStatusInput{
		UpdateType: UpdateTypeComplete,
		SyncRun:    syncRun,
		NewStatus:  models.SyncRunStatusFailed,
		Error:      &errString,
	}).Get(ctx, nil)
}

func recordSuccess(ctx workflow.Context, syncRun models.SyncRun, rowsWritten int) error {
	var a *Activities // Temporal handles calling the registered activity object
	return workflow.ExecuteActivity(ctx, a.RecordStatus, RecordStatusInput{
		UpdateType:  UpdateTypeComplete,
		SyncRun:     syncRun,
		NewStatus:   models.SyncRunStatusCompleted,
		RowsWritten: rowsWritten,
	}).Get(ctx, nil)
}
