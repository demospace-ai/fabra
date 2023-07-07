package models

import (
	"time"
)

type Verification struct {
	Code       string
	UserID     int64
	Expiration time.Time
	Used       bool

	BaseModel
}
