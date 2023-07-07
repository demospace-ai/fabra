package models

import (
	"time"
)

type Session struct {
	Token      string
	UserID     int64
	Expiration time.Time

	BaseModel
}
