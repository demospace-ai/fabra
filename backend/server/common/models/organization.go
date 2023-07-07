package models

import "go.fabra.io/server/common/database"

type Organization struct {
	Name         string            `json:"name"`
	EmailDomain  string            `json:"email_domain"`
	FreeTrialEnd database.NullTime `json:"free_trial_end,omitempty"`

	BaseModel
}
