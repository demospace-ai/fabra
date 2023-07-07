package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	// TODO: don't expose common database ID to public
	ID int64 `json:"id"`

	// Don't expose these fields
	CreatedAt     time.Time    `json:"-"`
	UpdatedAt     time.Time    `json:"-"`
	DeactivatedAt sql.NullTime `json:"-"`
}
