package models

import "go.fabra.io/server/common/database"

type Destination struct {
	OrganizationID int64               `json:"organization_id"`
	DisplayName    string              `json:"display_name"`
	ConnectionID   int64               `json:"connection_id"`
	StagingBucket  database.NullString `json:"staging_bucket"`

	BaseModel
}

type DestinationConnection struct {
	ID             int64
	OrganizationID int64
	DisplayName    string
	ConnectionID   int64
	StagingBucket  database.NullString
	ConnectionType ConnectionType
}
