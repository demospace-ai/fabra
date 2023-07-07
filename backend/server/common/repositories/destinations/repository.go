package destinations

import (
	"go.fabra.io/server/common/database"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"

	"gorm.io/gorm"
)

func CreateDestination(
	db *gorm.DB,
	organizationID int64,
	displayName string,
	connectionID int64,
	stagingBucket *string,
) (*models.Destination, error) {

	destination := models.Destination{
		OrganizationID: organizationID,
		DisplayName:    displayName,
		ConnectionID:   connectionID,
	}

	if stagingBucket != nil {
		destination.StagingBucket = database.NewNullString(*stagingBucket)
	}

	result := db.Create(&destination)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(destinations.CreateDestination)")
	}

	return &destination, nil
}

// TODO: test that connection credentials are not exposed
func LoadDestinationByID(db *gorm.DB, organizationID int64, destinationID int64) (*models.Destination, error) {
	var destination models.Destination
	result := db.Table("destinations").
		Select("destinations.*").
		Where("destinations.id = ?", destinationID).
		Where("destinations.organization_id = ?", organizationID).
		Where("destinations.deactivated_at IS NULL").
		Take(&destination)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(destinations.LoadDestinationByID)")
	}

	return &destination, nil
}

func LoadAllDestinations(
	db *gorm.DB,
	organizationID int64,
) ([]models.DestinationConnection, error) {
	var destinations []models.DestinationConnection
	result := db.Table("destinations").
		Select("destinations.*, connections.connection_type").
		Joins("JOIN connections ON destinations.connection_id = connections.id").
		Where("destinations.organization_id = ?", organizationID).
		Where("destinations.deactivated_at IS NULL").
		Order("destinations.created_at ASC").
		Where("connections.deactivated_at IS NULL").
		Find(&destinations)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(destinations.LoadAllDestinations)")
	}

	return destinations, nil
}
