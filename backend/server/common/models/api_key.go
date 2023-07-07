package models

type ApiKey struct {
	OrganizationID int64
	EncryptedKey   string
	HashedKey      string

	BaseModel
}
