package models

type Source struct {
	OrganizationID int64  `json:"organization_id"`
	DisplayName    string `json:"display_name"`
	EndCustomerID  string `json:"end_customer_id"`
	ConnectionID   int64  `json:"connection_id"`

	BaseModel
}

type SourceConnection struct {
	ID             int64
	OrganizationID int64
	EndCustomerID  string
	DisplayName    string
	ConnectionID   int64
	ConnectionType ConnectionType
}
