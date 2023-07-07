package models

type EndCustomerApiKey struct {
	OrganizationID int64
	EndCustomerID  string
	EncryptedKey   string

	BaseModel
}
