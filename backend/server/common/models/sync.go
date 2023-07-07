package models

import (
	"go.fabra.io/server/common/database"
)

type SyncMode string

const (
	SyncModeFullOverwrite     SyncMode = "full_overwrite"
	SyncModeFullAppend        SyncMode = "full_append" // for testing only: do not expose to customers in UI
	SyncModeIncrementalAppend SyncMode = "incremental_append"
	SyncModeIncrementalUpdate SyncMode = "incremental_update"
)

func (sm SyncMode) UsesCursor() bool {
	return sm == SyncModeIncrementalAppend || sm == SyncModeIncrementalUpdate
}

type FrequencyUnits string

const (
	FrequencyUnitsMinutes FrequencyUnits = "minutes"
	FrequencyUnitsHours   FrequencyUnits = "hours"
	FrequencyUnitsDays    FrequencyUnits = "days"
	FrequencyUnitsWeeks   FrequencyUnits = "weeks"
)

type SyncStatus string

const (
	SyncStatusActive SyncStatus = "active"
	SyncStatusPaused SyncStatus = "paused"
)

type Sync struct {
	OrganizationID int64
	DisplayName    string              `json:"display_name"`
	Status         SyncStatus          `json:"status"`
	WorkflowID     string              `json:"workflow_id"`
	EndCustomerID  string              `json:"end_customer_id"`
	SourceID       int64               `json:"source_id"`
	ObjectID       int64               `json:"object_id"`
	Namespace      database.NullString `json:"namespace"`
	TableName      database.NullString `json:"table_name"`
	CustomJoin     database.NullString `json:"custom_join"`

	// These values are used to override the object settings, but default to the same values
	SyncMode          SyncMode            `json:"sync_mode"`
	Recurring         bool                `json:"recurring"`
	Frequency         *int64              `json:"frequency,omitempty"`
	FrequencyUnits    *FrequencyUnits     `json:"frequency_units,omitempty"`
	SourceCursorField database.NullString `json:"source_cursor_field,omitempty"`
	SourcePrimaryKey  database.NullString `json:"source_primary_key,omitempty"`
	CursorPosition    database.NullString `json:"cursor_position"` // current value of the cursor to determine where to start a sync from

	BaseModel
}
