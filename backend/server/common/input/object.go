package input

import (
	"encoding/json"

	"go.fabra.io/server/common/models"
)

type PartialUpdateObjectInput struct {
	// Don't include destination, sync mode, cursor field, primary key, or end customer ID
	// field since we don't have a way to safely update these
	DisplayName    *string                `json:"display_name,omitempty"`
	Recurring      *bool                  `json:"recurring,omitempty"`
	Frequency      *int64                 `json:"frequency,omitempty"`
	FrequencyUnits *models.FrequencyUnits `json:"frequency_units,omitempty"`
}

type PartialUpdateObjectField struct {
	// Don't include name or type since we don't have a way to safely update these
	ID             int64           `json:"id" validate:"required"`
	DisplayNameRaw json.RawMessage `json:"display_name"`
	DescriptionRaw json.RawMessage `json:"description"`
}
