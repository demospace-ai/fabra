package sources

import (
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"

	"gorm.io/gorm"
)

func CreateSource(
	db *gorm.DB,
	organizationID int64,
	displayName string,
	endCustomerID string,
	connectionID int64,
) (*models.Source, error) {

	source := models.Source{
		OrganizationID: organizationID,
		DisplayName:    displayName,
		EndCustomerID:  endCustomerID,
		ConnectionID:   connectionID,
	}

	result := db.Create(&source)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sources.CreateSource)")
	}

	return &source, nil
}

// TODO: test that connection credentials are not exposed
func LoadSourceByID(db *gorm.DB, organizationID int64, endCustomerID string, sourceID int64) (*models.Source, error) {
	var source models.Source
	result := db.Table("sources").
		Select("sources.*").
		Where("sources.id = ?", sourceID).
		Where("sources.organization_id = ?", organizationID).
		Where("sources.end_customer_id = ?", endCustomerID).
		Where("sources.deactivated_at IS NULL").
		Take(&source)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sources.LoadSourceByID)")
	}

	return &source, nil
}

func LoadSourcesByIDs(db *gorm.DB, organizationID int64, sourceIDs []int64) ([]models.SourceConnection, error) {
	var sources []models.SourceConnection
	result := db.Table("sources").
		Select("sources.*, connections.connection_type").
		Joins("JOIN connections ON sources.connection_id = connections.id").
		Where("sources.id IN ?", sourceIDs).
		Where("sources.organization_id = ?", organizationID).
		Where("sources.deactivated_at IS NULL").
		Order("sources.created_at ASC").
		Where("connections.deactivated_at IS NULL").
		Find(&sources)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sources.LoadSourcesByIDs)")
	}

	return sources, nil
}

func LoadSourcesByIDsForCustomer(db *gorm.DB, organizationID int64, endCustomerID string, sourceIDs []int64) ([]models.SourceConnection, error) {
	var sources []models.SourceConnection
	result := db.Table("sources").
		Select("sources.*, connections.connection_type").
		Joins("JOIN connections ON sources.connection_id = connections.id").
		Where("sources.id IN ?", sourceIDs).
		Where("sources.organization_id = ?", organizationID).
		Where("sources.end_customer_id = ?", endCustomerID).
		Where("sources.deactivated_at IS NULL").
		Order("sources.created_at ASC").
		Where("connections.deactivated_at IS NULL").
		Find(&sources)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sources.LoadSourcesByIDsForCustomer)")
	}

	return sources, nil
}

func LoadAllSources(
	db *gorm.DB,
	organizationID int64,
	endCustomerID string,
) ([]models.SourceConnection, error) {
	var sources []models.SourceConnection
	result := db.Table("sources").
		Select("sources.*, connections.connection_type").
		Joins("JOIN connections ON sources.connection_id = connections.id").
		Where("sources.organization_id = ?", organizationID).
		Where("sources.end_customer_id = ?", endCustomerID).
		Where("sources.deactivated_at IS NULL").
		Order("sources.created_at ASC").
		Where("connections.deactivated_at IS NULL").
		Find(&sources)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(sources.LoadAllSources)")
	}

	return sources, nil
}
