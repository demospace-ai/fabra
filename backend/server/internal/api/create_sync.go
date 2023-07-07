package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/objects"
	"go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/server/common/timeutils"
	"go.fabra.io/server/common/views"
	"go.fabra.io/sync/temporal"
	"go.temporal.io/sdk/client"

	"github.com/go-playground/validator/v10"
)

const CLIENT_PEM_KEY = "projects/932264813910/secrets/temporal-client-pem/versions/latest"
const CLIENT_KEY_KEY = "projects/932264813910/secrets/temporal-client-key/versions/latest"

type CreateSyncRequest struct {
	DisplayName       string                 `json:"display_name"`
	EndCustomerID     *string                `json:"end_customer_id,omitempty"`
	SourceID          int64                  `json:"source_id"`
	ObjectID          int64                  `json:"object_id"`
	Namespace         *string                `json:"namespace,omitempty"`
	TableName         *string                `json:"table_name,omitempty"`
	CustomJoin        *string                `json:"custom_join,omitempty"`
	SourceCursorField *string                `json:"source_cursor_field,omitempty"`
	SourcePrimaryKey  *string                `json:"source_primary_key,omitempty"`
	SyncMode          *models.SyncMode       `json:"sync_mode,omitempty"`
	Recurring         *bool                  `json:"recurring,omitempty"`
	Frequency         *int64                 `json:"frequency,omitempty"`
	FrequencyUnits    *models.FrequencyUnits `json:"frequency_units,omitempty"`
	FieldMappings     []input.FieldMapping   `json:"field_mappings"`
}

type CreateSyncResponse struct {
	Sync          views.Sync           `json:"sync"`
	FieldMappings []views.FieldMapping `json:"field_mappings"`
}

func (s ApiService) CreateSync(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {

	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.CreateSync)")
	}

	decoder := json.NewDecoder(r.Body)
	var createSyncRequest CreateSyncRequest
	err := decoder.Decode(&createSyncRequest)
	if err != nil {
		return errors.Wrap(err, "(api.CreateSync)")
	}

	// TODO: validate connection parameters
	validate := validator.New()
	err = validate.Struct(createSyncRequest)
	if err != nil {
		return errors.Wrap(err, "(api.CreateSync)")
	}

	sync, fieldMappings, err := s.createSync(auth, createSyncRequest, *createSyncRequest.EndCustomerID)
	if err != nil {
		return errors.Wrap(err, "(api.CreateSync)")
	}

	return json.NewEncoder(w).Encode(CreateSyncResponse{
		Sync:          *sync,
		FieldMappings: fieldMappings,
	})
}

func (s ApiService) createSync(auth auth.Authentication, createSyncRequest CreateSyncRequest, endCustomerID string) (*views.Sync, []views.FieldMapping, error) {
	if (createSyncRequest.TableName == nil || createSyncRequest.Namespace == nil) && createSyncRequest.CustomJoin == nil {
		return nil, nil, errors.Wrap(errors.NewBadRequest("must have table_name and namespace or custom_join"), "(api.createSync)")
	}

	// this also serves to check that this organization owns the object
	object, err := objects.LoadObjectByID(s.db, auth.Organization.ID, createSyncRequest.ObjectID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSync)")
	}

	objectFields, err := objects.LoadObjectFieldsByID(s.db, createSyncRequest.ObjectID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSync)")
	}

	// default values for the sync come from the object
	sourceCursorField := getSourceCursorField(object, objectFields, createSyncRequest.FieldMappings)
	sourcePrimaryKey := getSourcePrimaryKey(object, objectFields, createSyncRequest.FieldMappings)
	syncMode := object.SyncMode
	recurring := object.Recurring
	frequency := object.Frequency
	frequencyUnits := object.FrequencyUnits

	// TODO: validate that the organization allows customizing sync settings
	if true {
		if createSyncRequest.SourceCursorField != nil {
			sourceCursorField = createSyncRequest.SourceCursorField
		}
		if createSyncRequest.SourcePrimaryKey != nil {
			sourcePrimaryKey = createSyncRequest.SourcePrimaryKey
		}
		if createSyncRequest.SyncMode != nil {
			syncMode = *createSyncRequest.SyncMode
		}
		if createSyncRequest.Recurring != nil {
			recurring = *createSyncRequest.Recurring
		}
		if createSyncRequest.Frequency != nil {
			frequency = createSyncRequest.Frequency
		}
		if createSyncRequest.FrequencyUnits != nil {
			frequencyUnits = createSyncRequest.FrequencyUnits
		}
	}

	if recurring {
		if *frequency <= 0 || (*frequency < 30 && *frequencyUnits == models.FrequencyUnitsMinutes) {
			return nil, nil, errors.NewBadRequest("Frequency must be greater than 30 minutes")
		}
	}

	err = validateFieldsMapped(objectFields, createSyncRequest.FieldMappings)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSync)")
	}

	// TODO: create via schedule in Temporal once GA
	// TODO: create field mappings in DB using transaction
	sync, err := syncs.CreateSync(
		s.db,
		auth.Organization.ID,
		createSyncRequest.DisplayName,
		endCustomerID,
		createSyncRequest.SourceID,
		createSyncRequest.ObjectID,
		createSyncRequest.Namespace,
		createSyncRequest.TableName,
		createSyncRequest.CustomJoin,
		sourceCursorField,
		sourcePrimaryKey,
		syncMode,
		recurring,
		frequency,
		frequencyUnits,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSync)")
	}

	fieldMappings, err := syncs.CreateFieldMappings(
		s.db, auth.Organization.ID, sync.ID, createSyncRequest.FieldMappings,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSync)")
	}

	err = createTemporalWorkflow(sync.ID, auth.Organization.ID, recurring, sync.WorkflowID, frequency, frequencyUnits)
	if err != nil {
		return nil, nil, errors.Wrap(err, "(api.createSync)")
	}

	syncView := views.ConvertSync(sync)
	return &syncView, views.ConvertFieldMappings(fieldMappings, objectFields), nil
}

func createTemporalWorkflow(syncID int64, organizationID int64, recurring bool, workflowID string, frequency *int64, frequencyUnits *models.FrequencyUnits) error {
	c, err := temporal.CreateClient(CLIENT_PEM_KEY, CLIENT_KEY_KEY)
	if err != nil {
		return errors.Wrap(err, "(api.createTemporalWorkflow)")
	}
	defer c.Close()
	ctx := context.TODO()
	scheduleClient := c.ScheduleClient()

	scheduleOptions := client.ScheduleOptions{
		ID:                 workflowID,
		TriggerImmediately: recurring,
		Action: &client.ScheduleWorkflowAction{
			TaskQueue: temporal.SyncTaskQueue,
			Workflow:  temporal.SyncWorkflow,
			Args: []interface{}{temporal.SyncInput{
				SyncID: syncID, OrganizationID: organizationID,
			}},
		},
	}

	if recurring {
		schedule, err := createSchedule(*frequency, *frequencyUnits)
		if err != nil {
			return errors.Wrap(err, "(api.createTemporalWorkflow)")
		}

		scheduleOptions.Spec = client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{
					Every: schedule,
				},
			},
		}
	}

	_, err = scheduleClient.Create(ctx, scheduleOptions)
	if err != nil {
		return errors.Wrap(err, "(api.createTemporalWorkflow)")
	}

	return nil
}

func createSchedule(frequency int64, frequencyUnits models.FrequencyUnits) (time.Duration, error) {
	frequencyDuration := time.Duration(frequency)
	switch frequencyUnits {
	case models.FrequencyUnitsMinutes:
		return frequencyDuration * time.Minute, nil
	case models.FrequencyUnitsHours:
		return frequencyDuration * time.Hour, nil
	case models.FrequencyUnitsDays:
		return frequencyDuration * timeutils.DAY, nil
	case models.FrequencyUnitsWeeks:
		return frequencyDuration * timeutils.WEEK, nil
	default:
		// TODO: this should not happen
		return timeutils.WEEK, errors.Newf("(api.createSchedule) unexpected frequency unit: %s", string(frequencyUnits))
	}
}

func getSourcePrimaryKey(object *models.Object, objectFields []models.ObjectField, fieldMappings []input.FieldMapping) *string {
	if object.PrimaryKey.Valid {
		var destinationPrimaryKey models.ObjectField
		for _, field := range objectFields {
			if field.Name == object.PrimaryKey.String {
				destinationPrimaryKey = field
			}
		}

		for _, fieldMapping := range fieldMappings {
			if fieldMapping.DestinationFieldId == destinationPrimaryKey.ID {
				return &fieldMapping.SourceFieldName
			}
		}
	}

	return nil
}

func getSourceCursorField(object *models.Object, objectFields []models.ObjectField, fieldMappings []input.FieldMapping) *string {
	if object.CursorField.Valid {
		var destinationCursorField models.ObjectField
		for _, field := range objectFields {
			if field.Name == object.CursorField.String {
				destinationCursorField = field
			}
		}

		for _, fieldMapping := range fieldMappings {
			if fieldMapping.DestinationFieldId == destinationCursorField.ID {
				return &fieldMapping.SourceFieldName
			}
		}
	}

	return nil
}

func validateFieldsMapped(objectFields []models.ObjectField, fieldMappings []input.FieldMapping) error {
	mappedObjectFieldIDs := make(map[int64]bool)
	for _, fieldMapping := range fieldMappings {
		mappedObjectFieldIDs[fieldMapping.DestinationFieldId] = true
	}

	for _, objectField := range objectFields {
		if !objectField.Optional && !objectField.Omit && !mappedObjectFieldIDs[objectField.ID] {
			return errors.NewBadRequest(fmt.Sprintf("object field %s is not mapped", objectField.Name))
		}
	}

	return nil
}
