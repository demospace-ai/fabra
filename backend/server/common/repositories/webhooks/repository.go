package webhooks

import (
	"time"

	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"gorm.io/gorm"
)

func CreateEndCustomerApiKey(db *gorm.DB, organizationID int64, endCustomerID string, encryptedKey string) error {
	endCustomerApiKey := models.EndCustomerApiKey{
		OrganizationID: organizationID,
		EndCustomerID:  endCustomerID,
		EncryptedKey:   encryptedKey,
	}

	result := db.Create(&endCustomerApiKey)
	if result.Error != nil {
		return errors.Wrap(result.Error, "(api_keys.CreateEndCustomerApiKey)")
	}

	return nil
}

func LoadEndCustomerApiKey(db *gorm.DB, organizationID int64, endCustomerID string) (*string, error) {
	var endCustomerApiKey models.EndCustomerApiKey
	result := db.Table("end_customer_api_keys").
		Select("end_customer_api_keys.*").
		Where("end_customer_api_keys.organization_id = ?", organizationID).
		Where("end_customer_api_keys.end_customer_id = ?", endCustomerID).
		Where("end_customer_api_keys.deactivated_at IS NULL").
		Take(&endCustomerApiKey)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(api_keys.LoadEndCustomerApiKey)")
	}

	return &endCustomerApiKey.EncryptedKey, nil
}

func DeactivateExistingEndCustomerApiKey(db *gorm.DB, organizationID int64, endCustomerID string) error {
	result := db.Table("end_customer_api_keys").
		Where("end_customer_api_keys.organization_id = ?", organizationID).
		Where("end_customer_api_keys.end_customer_id = ?", endCustomerID).
		Update("deactivated_at", time.Now())

	if result.Error != nil {
		return errors.Wrap(result.Error, "(api_keys.DeactivateExistingEndCustomerApiKey)")
	}

	return nil
}
