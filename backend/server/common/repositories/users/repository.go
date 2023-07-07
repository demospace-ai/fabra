package users

import (
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/events"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/oauth"
	"go.fabra.io/server/common/repositories/external_profiles"
	"gorm.io/gorm"
)

// Maximum of 62^8 guarantees number will be at most 8 digits in base
const MAX_RANDOM = 218340105584896

func LoadByExternalID(db *gorm.DB, externalID string) (*models.User, error) {
	var user models.User
	result := db.Table("users").
		Joins("JOIN external_profiles ON external_profiles.user_id = users.id").
		Where("external_profiles.external_id = ?", externalID).
		Where("users.deactivated_at IS NULL").
		Take(&user)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(users.LoadByExternalID)")
	}

	return &user, nil
}

func LoadByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	result := db.Table("users").
		Joins("JOIN emails ON emails.user_id = users.id").
		Where("emails.email = ?", email).
		Where("users.deactivated_at IS NULL").
		Take(&user)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(users.LoadByEmail)")
	}

	return &user, nil
}

func LoadUserByID(db *gorm.DB, userID int64) (*models.User, error) {
	var user models.User
	result := db.Table("users").
		Select("users.*").
		Where("users.id = ?", userID).
		Where("users.deactivated_at IS NULL").
		Take(&user)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(users.LoadUserByID)")
	}

	return &user, nil
}

func create(db *gorm.DB, name string, email string) (*models.User, error) {
	user := models.User{
		Name:  name,
		Email: email,
	}

	result := db.Create(&user)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(users.create)")
	}

	return &user, nil
}

func CreateUserForExternalInfo(db *gorm.DB, externalUserInfo *oauth.ExternalUserInfo) (*models.User, error) {
	user, err := create(db, externalUserInfo.Name, externalUserInfo.Email)
	if err != nil {
		return nil, errors.Wrap(err, "(users.CreateUserForExternalInfo)")
	}

	_, err = external_profiles.Create(db, externalUserInfo.ExternalID, externalUserInfo.OauthProvider, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "(users.CreateUserForExternalInfo)")
	}

	events.TrackSignup(user.ID, user.Name, user.Email)

	return user, nil
}

func SetOrganization(db *gorm.DB, user *models.User, organizationID int64) (*models.User, error) {
	result := db.Model(user).Update("organization_id", organizationID)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(users.SetOrganization)")
	}

	return user, nil
}

func GetOrCreateForExternalInfo(db *gorm.DB, externalUserInfo *oauth.ExternalUserInfo) (*models.User, error) {
	existingUser, err := LoadByExternalID(db, externalUserInfo.ExternalID)
	if err != nil && !errors.IsRecordNotFound(err) {
		return nil, errors.Wrap(err, "(users.GetOrCreateForExternalInfo)")
	} else if err == nil {
		return existingUser, nil
	}

	user, err := CreateUserForExternalInfo(db, externalUserInfo)
	if err != nil {
		return nil, errors.Wrap(err, "(users.GetOrCreateForExternalInfo)")
	}

	return user, nil
}

func LoadAllByOrganizationID(db *gorm.DB, organizationID int64) ([]models.User, error) {
	var users []models.User
	result := db.Table("users").
		Where("users.organization_id = ?", organizationID).
		Where("users.deactivated_at IS NULL").
		Find(&users)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(users.LoadAllByOrganizationID)")
	}

	return users, nil

}
