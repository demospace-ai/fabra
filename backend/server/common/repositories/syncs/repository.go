package syncs

import (
	"time"

	"github.com/google/uuid"
	"go.fabra.io/server/common/database"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/models"

	"gorm.io/gorm"
)

func CreateSync(
	db *gorm.DB,
	organizationID int64,
	displayName string,
	endCustomerID string,
	sourceID int64,
	objectID int64,
	namespace *string,
	tableName *string,
	customJoin *string,
	sourceCursorField *string,
	sourcePrimaryKey *string,
	syncMode models.SyncMode,
	recurring bool,
	frequency *int64,
	frequencyUnits *models.FrequencyUnits,
) (*models.Sync, error) {

	sync := models.Sync{
		OrganizationID: organizationID,
		DisplayName:    displayName,
		WorkflowID:     uuid.NewString(),
		EndCustomerID:  endCustomerID,
		SourceID:       sourceID,
		ObjectID:       objectID,
		SyncMode:       syncMode,
		Frequency:      frequency,
		FrequencyUnits: frequencyUnits,
		Status:         models.SyncStatusActive,
	}

	if tableName != nil && namespace != nil {
		sync.Namespace = database.NewNullString(*namespace)
		sync.TableName = database.NewNullString(*tableName)
	}

	if customJoin != nil {
		sync.CustomJoin = database.NewNullString(*customJoin)
	}

	if sourceCursorField != nil {
		sync.SourceCursorField = database.NewNullString(*sourceCursorField)
	}

	if sourcePrimaryKey != nil {
		sync.SourcePrimaryKey = database.NewNullString(*sourcePrimaryKey)
	}

	result := db.Create(&sync)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.CreateSync)")
	}

	return &sync, nil
}

func CreateFieldMappings(
	db *gorm.DB,
	organizationID int64,
	syncID int64,
	fieldMappings []input.FieldMapping,
) ([]models.FieldMapping, error) {
	// TODO: validate that the mapped object fields belong to the right object
	var createdFieldMappings []models.FieldMapping
	for _, fieldMapping := range fieldMappings {
		fieldMappingModel := models.FieldMapping{
			SyncID:             syncID,
			SourceFieldName:    fieldMapping.SourceFieldName,
			SourceFieldType:    fieldMapping.SourceFieldType,
			DestinationFieldId: fieldMapping.DestinationFieldId,
			IsJsonField:        fieldMapping.IsJsonField,
		}

		result := db.Create(&fieldMappingModel)
		if result.Error != nil {
			return nil, errors.Wrap(result.Error, "(syncs.CreateFieldMappings)")
		}
		createdFieldMappings = append(createdFieldMappings, fieldMappingModel)
	}

	return createdFieldMappings, nil
}

func LoadSyncByID(db *gorm.DB, organizationID int64, syncID int64) (*models.Sync, error) {
	var sync models.Sync
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.id = ?", syncID).
		Where("syncs.organization_id = ?", organizationID).
		Where("syncs.deactivated_at IS NULL").
		Take(&sync)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.LoadSyncByID)")
	}

	return &sync, nil
}

func DeactivateSyncByID(db *gorm.DB, syncID int64) error {
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.id = ?", syncID).
		Update("deactivated_at", time.Now())

	if result.Error != nil {
		return errors.Wrap(result.Error, "(syncs.DeactivateSyncByID)")
	}

	return nil
}

func UpdateSyncStatusByID(db *gorm.DB, syncID int64, status models.SyncStatus) error {
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.id = ?", syncID).
		Update("status", status)

	if result.Error != nil {
		return errors.Wrap(result.Error, "(syncs.UpdateSyncStatusByID)")
	}

	return nil
}

func LoadSyncByIDAndCustomer(db *gorm.DB, organizationID int64, endCustomerID string, syncID int64) (*models.Sync, error) {
	var sync models.Sync
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.id = ?", syncID).
		Where("syncs.organization_id = ?", organizationID).
		Where("syncs.end_customer_id = ?", endCustomerID).
		Where("syncs.deactivated_at IS NULL").
		Take(&sync)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.LoadSyncByIDAndCustomer)")
	}

	return &sync, nil
}

func LoadAllSyncs(
	db *gorm.DB,
	organizationID int64,
) ([]models.Sync, error) {
	var sync []models.Sync
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.organization_id = ?", organizationID).
		Where("syncs.deactivated_at IS NULL").
		Order("syncs.created_at DESC").
		Find(&sync)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.LoadAllSyncs)")
	}

	return sync, nil
}

func LoadAllSyncsForCustomer(
	db *gorm.DB,
	organizationID int64,
	endCustomerID string,
) ([]models.Sync, error) {
	var syncs []models.Sync
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.organization_id = ?", organizationID).
		Where("syncs.end_customer_id = ?", endCustomerID).
		Where("syncs.deactivated_at IS NULL").
		Order("syncs.created_at DESC").
		Find(&syncs)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.LoadAllSyncsForCustomer)")
	}

	return syncs, nil
}

func LoadSyncsForCustomerAndObject(
	db *gorm.DB,
	organizationID int64,
	endCustomerID string,
	objectID int64,
) ([]models.Sync, error) {
	var syncs []models.Sync
	result := db.Table("syncs").
		Select("syncs.*").
		Where("syncs.organization_id = ?", organizationID).
		Where("syncs.end_customer_id = ?", endCustomerID).
		Where("syncs.object_id = ?", objectID).
		Where("syncs.deactivated_at IS NULL").
		Order("syncs.created_at DESC").
		Find(&syncs)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.LoadSyncsForCustomerAndObject)")
	}

	return syncs, nil
}

func LoadFieldMappingsForSync(
	db *gorm.DB,
	syncID int64,
) ([]models.FieldMapping, error) {
	// TODO: validate that the mapped object fields belong to the right object
	var fieldMappings []models.FieldMapping
	result := db.Table("field_mappings").
		Select("field_mappings.*").
		Where("field_mappings.sync_id = ?", syncID).
		Where("field_mappings.deactivated_at IS NULL").
		Find(&fieldMappings)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.LoadFieldMappingsForSync)")
	}

	return fieldMappings, nil
}

func UpdateCursor(
	db *gorm.DB,
	sync *models.Sync,
	cursorPosition string,
) (*models.Sync, error) {
	updates := models.Sync{
		CursorPosition: database.NewNullString(cursorPosition),
	}

	result := db.Model(sync).Updates(updates)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "(syncs.UpdateCursor)")
	}

	return sync, nil
}
