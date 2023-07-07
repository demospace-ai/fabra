package temporal

import (
	"context"
	"time"

	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/views"
	"go.fabra.io/sync/connectors"
	"go.temporal.io/sdk/activity"
)

const FABRA_STAGING_BUCKET = "fabra-staging"

type ReplicateInput = SyncConfig

type ReplicateOutput struct {
	RowsWritten    int
	CursorPosition *string
}

type FormatToken struct {
	Format string
	Index  int
}

func (a *Activities) Replicate(ctx context.Context, input ReplicateInput) (*ReplicateOutput, error) {
	cryptoService := crypto.NewCryptoService()
	queryService := query.NewQueryService(cryptoService)

	rowsC := make(chan []data.Row)
	readOutputC := make(chan connectors.ReadOutput)
	writeOutputC := make(chan connectors.WriteOutput)
	readErrC := make(chan error)
	writeErrC := make(chan error)
	doneC := make(chan bool)

	sourceConnector, err := getSourceConnector(ctx, input.SourceConnection, queryService)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.Replicate) getSourceConnector")
	}

	destConnector, err := getDestinationConnector(ctx, input.DestinationConnection, queryService, cryptoService, input.EncryptedEndCustomerApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "(temporal.Replicate) getDestinationConnector")
	}

	go safeCall(func() {
		sourceConnector.Read(ctx, input.SourceConnection, input.Sync, input.FieldMappings, rowsC, readOutputC, readErrC)
	}, readErrC)

	go safeCall(func() {
		destConnector.Write(ctx, input.DestinationConnection, input.DestinationOptions, input.Object, input.Sync, input.FieldMappings, rowsC, writeOutputC, writeErrC)
	}, writeErrC)

	go heartbeat(ctx, doneC) // TODO: heartbeat from the write/read methods to ensure the worker is making progress

	var readOutput connectors.ReadOutput
	var writeOutput connectors.WriteOutput
	var readDone, writeDone bool
	for {
		if readDone && writeDone {
			break
		}

		// wait for both error channels in any order, immediately exiting if an error is returned
		select {
		case err = <-readErrC:
			if err != nil {
				return nil, errors.Wrap(err, "(temporal.Replicate) readErrC")
			}
		case err = <-writeErrC:
			if err != nil {
				return nil, errors.Wrap(err, "(temporal.Replicate) writeErrC")
			}
		case readOutput = <-readOutputC:
			readDone = true
		case writeOutput = <-writeOutputC:
			writeDone = true
		}
	}

	// signal the heartbeat worker that the replication is finished
	doneC <- true

	return &ReplicateOutput{
		RowsWritten:    writeOutput.RowsWritten,
		CursorPosition: readOutput.CursorPosition,
	}, nil
}

func getSourceConnector(ctx context.Context, connection views.FullConnection, queryService query.QueryService) (connectors.Connector, error) {
	connectionModel := views.ConvertConnectionView(connection)
	switch connection.ConnectionType {
	case models.ConnectionTypeBigQuery:
		warehouseClient, err := queryService.GetWarehouseClient(ctx, connectionModel)
		if err != nil {
			return nil, errors.Wrap(err, "(temporal.getSourceConnector)")
		}
		return connectors.NewBigQueryConnector(warehouseClient), nil
	case models.ConnectionTypeSnowflake:
		return connectors.NewSnowflakeConnector(queryService), nil
	case models.ConnectionTypeRedshift:
		return connectors.NewRedshiftConnector(queryService), nil
	case models.ConnectionTypeSynapse:
		return connectors.NewSynapseConnector(queryService), nil
	case models.ConnectionTypeMongoDb:
		return connectors.NewMongoDbConnector(queryService), nil
	case models.ConnectionTypePostgres:
		return connectors.NewPostgresConnector(queryService), nil
	case models.ConnectionTypeMySQL:
		return connectors.NewMySqlConnector(queryService), nil
	default:
		return nil, errors.Newf("(temporal.getSourceConnector) source not implemented for %s", connection.ConnectionType)
	}
}

func getDestinationConnector(ctx context.Context, connection views.FullConnection, queryService query.QueryService, cryptoService crypto.CryptoService, encryptedEndCustomerApiKey *string) (connectors.Connector, error) {
	connectionModel := views.ConvertConnectionView(connection)
	switch connection.ConnectionType {
	case models.ConnectionTypeBigQuery:
		warehouseClient, err := queryService.GetWarehouseClient(ctx, connectionModel)
		if err != nil {
			return nil, errors.Wrap(err, "(temporal.getDestinationConnector)")
		}
		return connectors.NewBigQueryConnector(warehouseClient), nil
	case models.ConnectionTypeWebhook:
		// TODO: does end customer api key belong here?
		return connectors.NewWebhookConnector(queryService, cryptoService, encryptedEndCustomerApiKey), nil
	default:
		return nil, errors.Newf("(temporal.getDestinationConnector) destination not implemented for %s", connection.ConnectionType)
	}
}

// Any new goroutines can crash the whole worker unless we recover from panics
func safeCall(fn func(), errC chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			errC <- errors.Newf("panic: %v", r)
		}
	}()

	fn()
}

func heartbeat(ctx context.Context, doneC <-chan bool) {
	timeChan := time.NewTicker(time.Minute).C
	for {
		activity.RecordHeartbeat(ctx)
		select {
		case <-doneC:
			return
		case <-timeChan:
			continue
		}
	}
}
