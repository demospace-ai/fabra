package views

import "go.fabra.io/server/common/models"

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IntercomHash string `json:"intercom_hash"`
}

func ConvertUser(user models.User, intercomHash string) User {
	return User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IntercomHash: intercomHash,
	}
}
