package models

import "go.fabra.io/server/common/database"

type User struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	ProfilePictureURL string
	OrganizationID    database.NullInt64
	Blocked           bool `json:"-"`

	BaseModel
}
