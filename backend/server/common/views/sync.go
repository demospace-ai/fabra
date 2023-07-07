package views

import (
	"time"

	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/timeutils"
)

const CUSTOMER_VISIBLE_TIME_FORMAT = "01/02/06 at 03:04 PM MST"

type Sync struct {
	ID                int64                  `json:"id"`
	OrganizationID    int64                  `json:"organization_id"`
	Status            models.SyncStatus      `json:"status"`
	EndCustomerID     string                 `json:"end_customer_id"`
	DisplayName       string                 `json:"display_name"`
	SourceID          int64                  `json:"source_id"`
	ObjectID          int64                  `json:"object_id"`
	Namespace         *string                `json:"namespace,omitempty"`
	TableName         *string                `json:"table_name,omitempty"`
	CustomJoin        *string                `json:"custom_join,omitempty"`
	CursorPosition    *string                `json:"cursor_position,omitempty"`
	SourceCursorField *string                `json:"source_cursor_field,omitempty"`
	SourcePrimaryKey  *string                `json:"source_primary_key,omitempty"`
	SyncMode          models.SyncMode        `json:"sync_mode"`
	Recurring         bool                   `json:"recurring"`
	Frequency         *int64                 `json:"frequency,omitempty"`
	FrequencyUnits    *models.FrequencyUnits `json:"frequency_units,omitempty"`
}

type SyncRun struct {
	Status      models.SyncRunStatus `json:"status"`
	StartedAt   string               `json:"started_at"`
	CompletedAt string               `json:"completed_at"`
	Duration    *string              `json:"duration,omitempty"`
	Error       *string              `json:"error,omitempty"`
	RowsWritten int                  `json:"rows_written"`
}

type FieldMapping struct {
	SourceFieldName      string         `json:"source_field_name"`
	SourceFieldType      data.FieldType `json:"source_field_type"`
	DestinationFieldId   int64          `json:"destination_field_id"`
	DestinationFieldName string         `json:"destination_field_name"`
	DestinationFieldType data.FieldType `json:"destination_field_type"`
	IsJsonField          bool           `json:"is_json_field"`
}

func ConvertSync(sync *models.Sync) Sync {
	syncView := Sync{
		ID:             sync.ID,
		OrganizationID: sync.OrganizationID,
		Status:         sync.Status,
		EndCustomerID:  sync.EndCustomerID,
		DisplayName:    sync.DisplayName,
		SourceID:       sync.SourceID,
		ObjectID:       sync.ObjectID,
		SyncMode:       sync.SyncMode,
		Frequency:      sync.Frequency,
	}

	if sync.Namespace.Valid {
		syncView.Namespace = &sync.Namespace.String
	}
	if sync.TableName.Valid {
		syncView.TableName = &sync.TableName.String
	}
	if sync.CustomJoin.Valid {
		syncView.CustomJoin = &sync.CustomJoin.String
	}
	if sync.CursorPosition.Valid {
		syncView.CursorPosition = &sync.CursorPosition.String
	}
	if sync.SourceCursorField.Valid {
		syncView.SourceCursorField = &sync.SourceCursorField.String
	}
	if sync.SourcePrimaryKey.Valid {
		syncView.SourcePrimaryKey = &sync.SourcePrimaryKey.String
	}

	return syncView
}

func ConvertFieldMappings(fieldMappings []models.FieldMapping, objectFields []models.ObjectField) []FieldMapping {
	// Create a map of object fields by id
	objectFieldsById := make(map[int64]models.ObjectField)
	for _, objectField := range objectFields {
		objectFieldsById[objectField.ID] = objectField
	}

	var fieldMappingsView []FieldMapping
	for _, fieldMapping := range fieldMappings {
		destinationField := objectFieldsById[fieldMapping.DestinationFieldId]
		fieldMappingsView = append(fieldMappingsView, FieldMapping{
			SourceFieldName:      fieldMapping.SourceFieldName,
			SourceFieldType:      fieldMapping.SourceFieldType,
			DestinationFieldId:   fieldMapping.DestinationFieldId,
			DestinationFieldName: destinationField.Name,
			DestinationFieldType: destinationField.Type,
			IsJsonField:          fieldMapping.IsJsonField,
		})
	}

	return fieldMappingsView
}

func ConvertSyncRuns(syncRuns []models.SyncRun, timezone *time.Location) ([]SyncRun, error) {
	var syncRunsView []SyncRun
	for _, syncRun := range syncRuns {
		syncRunView := SyncRun{
			Status:      syncRun.Status,
			StartedAt:   syncRun.StartedAt.In(timezone).Format(CUSTOMER_VISIBLE_TIME_FORMAT),
			CompletedAt: syncRun.CompletedAt.In(timezone).Format(CUSTOMER_VISIBLE_TIME_FORMAT),
			RowsWritten: syncRun.RowsWritten,
		}
		if syncRun.Error.Valid {
			syncError := syncRun.Error.String
			syncRunView.Error = &syncError
		}
		if syncRun.Status != models.SyncRunStatusRunning {
			duration, err := timeutils.GetDurationString(syncRun.CompletedAt.Sub(syncRun.StartedAt))
			if err != nil {
				return nil, errors.Wrap(err, "(views.ConvertSyncRuns) getting duration string")
			}
			syncRunView.Duration = duration
		}

		syncRunsView = append(syncRunsView, syncRunView)
	}

	return syncRunsView, nil
}
