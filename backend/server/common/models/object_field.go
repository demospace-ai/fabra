package models

import (
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/database"
)

type ObjectField struct {
	ObjectID    int64               `json:"object_id"`
	Name        string              `json:"name"`
	Type        data.FieldType      `json:"type"`
	Omit        bool                `json:"omit"`
	Optional    bool                `json:"optional"`
	DisplayName database.NullString `json:"display_name"`
	Description database.NullString `json:"description"`

	BaseModel
}
