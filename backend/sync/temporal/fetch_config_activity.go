package temporal

import (
	"context"

	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/connections"
	"go.fabra.io/server/common/repositories/destinations"
	"go.fabra.io/server/common/repositories/objects"
	"go.fabra.io/server/common/repositories/sources"
	"go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/server/common/repositories/webhooks"
	"go.fabra.io/server/common/views"
	"go.fabra.io/sync/connectors"
)

type FetchConfigInput struct {
	OrganizationID int64
	SyncID         int64
}

type SyncConfig struct {
	Sync                       views.Sync
	SourceConnection           views.FullConnection
	DestinationConnection      views.FullConnection
	DestinationOptions         connectors.DestinationOptions
	Object                     views.Object
	ObjectFields               []views.ObjectField
	FieldMappings              []views.FieldMapping
	EncryptedEndCustomerApiKey *string
}

func (a *Activities) FetchConfig(ctx context.Context, input FetchConfigInput) (*SyncConfig, error) {
	sync, err := syncs.LoadSyncByID(a.Db, input.OrganizationID, input.SyncID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load sync")
	}

	source, err := sources.LoadSourceByID(a.Db, input.OrganizationID, sync.EndCustomerID, sync.SourceID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load source")
	}

	sourceConnection, err := connections.LoadConnectionByID(a.Db, input.OrganizationID, source.ConnectionID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load source connection")
	}

	object, err := objects.LoadObjectByID(a.Db, input.OrganizationID, sync.ObjectID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load object")
	}

	destination, err := destinations.LoadDestinationByID(a.Db, input.OrganizationID, object.DestinationID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load destination")
	}

	destinationConnection, err := connections.LoadConnectionByID(a.Db, input.OrganizationID, destination.ConnectionID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load destination connection")
	}

	fieldMappings, err := syncs.LoadFieldMappingsForSync(a.Db, input.SyncID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load field mappings")
	}

	objectFields, err := objects.LoadObjectFieldsByID(a.Db, object.ID)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load object fields")
	}

	encryptedEndCustomerApiKey, err := webhooks.LoadEndCustomerApiKey(a.Db, input.OrganizationID, sync.EndCustomerID)
	// This might be missing, but that's ok-- it isn't required
	if err != nil && !errors.IsRecordNotFound(err) {
		return nil, errors.Wrap(err, "(temporal.FetchConfig) failed to load end customer api key")
	}

	// TODO: encrypt this value before returning it, even though the credentials are already encrypted
	syncConfig := SyncConfig{
		Sync:                       views.ConvertSync(sync),
		SourceConnection:           views.ConvertFullConnection(sourceConnection),
		DestinationConnection:      views.ConvertFullConnection(destinationConnection),
		DestinationOptions:         connectors.DestinationOptions{},
		Object:                     views.ConvertObject(object, objectFields),
		FieldMappings:              views.ConvertFieldMappings(fieldMappings, objectFields),
		EncryptedEndCustomerApiKey: encryptedEndCustomerApiKey,
	}

	if destination.StagingBucket.Valid {
		syncConfig.DestinationOptions.StagingBucket = destination.StagingBucket.String
	}

	return &syncConfig, nil
}
