package input

import "go.fabra.io/server/common/data"

type FieldMapping struct {
	SourceFieldName    string         `json:"source_field_name,omitempty"`
	SourceFieldType    data.FieldType `json:"source_field_type,omitempty"`
	DestinationFieldId int64          `json:"destination_field_id,omitempty"`
	IsJsonField        bool           `json:"is_json_field,omitempty"`
}
