package organizations

import (
	"time"

	"go.fabra.io/server/common/database"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"

	"gorm.io/gorm"
)

func Create(db *gorm.DB, organizationName string, emailDomain string) (*models.Organization, error) {
	// 30-day free trial
	freeTrialEnd := time.Now().Add(time.Hour * 24 * 30)

	organization := models.Organization{
		Name:         organizationName,
		EmailDomain:  emailDomain,
		FreeTrialEnd: database.NewNullTime(freeTrialEnd),
	}

	result := db.Create(&organization)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(organizations.Create)")
	}

	return &organization, nil
}

func LoadOrganizationByID(db *gorm.DB, organizationID int64) (*models.Organization, error) {
	var organization models.Organization
	result := db.Table("organizations").
		Select("organizations.*").
		Where("organizations.id = ?", organizationID).
		Where("organizations.deactivated_at IS NULL").
		Take(&organization)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(organizations.LoadOrganizationByID)")
	}

	return &organization, nil
}

func LoadOrganizationsByEmailDomain(db *gorm.DB, emailDomain string) ([]models.Organization, error) {
	var organizations []models.Organization
	result := db.Table("organizations").
		Where("organizations.email_domain = ?", emailDomain).
		Where("organizations.deactivated_at IS NULL").
		Find(&organizations)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(organizations.LoadOrganizationsByEmailDomain)")
	}

	return organizations, nil
}

func LoadOrganizationByApiKey(db *gorm.DB, hashedKey string) (*models.Organization, error) {
	var organization models.Organization
	result := db.Table("organizations").
		Select("organizations.*").
		Joins("JOIN api_keys ON api_keys.organization_id = organizations.id").
		Where("api_keys.hashed_key = ?", hashedKey).
		Where("organizations.deactivated_at IS NULL").
		Where("api_keys.deactivated_at IS NULL").
		Take(&organization)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(organizations.LoadOrganizationByApiKey)")
	}

	return &organization, nil
}
