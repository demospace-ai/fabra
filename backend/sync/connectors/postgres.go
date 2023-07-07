package connectors

import (
	"context"
	"fmt"
	"strings"

	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/views"
)

type PostgresImpl struct {
	queryService query.QueryService
}

func NewPostgresConnector(queryService query.QueryService) Connector {
	return PostgresImpl{
		queryService: queryService,
	}
}

func (pg PostgresImpl) Read(
	ctx context.Context,
	sourceConnection views.FullConnection,
	sync views.Sync,
	fieldMappings []views.FieldMapping,
	rowsC chan<- []data.Row,
	readOutputC chan<- ReadOutput,
	errC chan<- error,
) {
	connectionModel := views.ConvertConnectionView(sourceConnection)

	sourceClient, err := pg.queryService.GetClient(ctx, connectionModel)
	if err != nil {
		errC <- err
		return
	}

	readQuery := pg.getReadQuery(connectionModel, sync, fieldMappings)

	iterator, err := sourceClient.GetQueryIterator(ctx, readQuery)
	if err != nil {
		errC <- err
		return
	}

	currentIndex := 0
	var rowBatch []data.Row
	var lastRow data.Row
	for {
		row, err := iterator.Next(ctx)
		if err != nil {
			if err == data.ErrDone {
				break
			} else {
				errC <- err
				return
			}
		}

		rowBatch = append(rowBatch, row)
		lastRow = row
		currentIndex++
		if currentIndex == READ_BATCH_SIZE {
			rowsC <- rowBatch
			currentIndex = 0
			rowBatch = []data.Row{}
		}
	}

	// write any remaining roows
	if currentIndex > 0 {
		rowsC <- rowBatch
	}

	newCursorPosition := pg.getNewCursorPosition(lastRow, iterator.Schema(), sync)
	readOutputC <- ReadOutput{
		CursorPosition: newCursorPosition,
	}

	close(rowsC)
	close(errC)
}

func (pg PostgresImpl) getReadQuery(sourceConnection *models.Connection, sync views.Sync, fieldMappings []views.FieldMapping) string {
	var queryString string
	if sync.CustomJoin != nil {
		queryString = *sync.CustomJoin
	} else {
		selectString := pg.getSelectString(fieldMappings)
		queryString = fmt.Sprintf("SELECT %s FROM %s.%s", selectString, *sync.Namespace, *sync.TableName)
	}

	if sync.SyncMode.UsesCursor() {
		if sync.CursorPosition != nil {
			// TODO: allow choosing other operators (rows smaller than current cursor)
			// order by cursor field to simplify
			return fmt.Sprintf("%s WHERE %s > %s ORDER BY %s ASC;", queryString, *sync.SourceCursorField, *sync.CursorPosition, *sync.SourceCursorField)
		} else {
			return fmt.Sprintf("%s ORDER BY %s ASC;", queryString, *sync.SourceCursorField)
		}
	} else {
		return fmt.Sprintf("%s;", queryString)
	}
}

func (pg PostgresImpl) getSelectString(fieldMappings []views.FieldMapping) string {
	fields := []string{}
	for _, fieldMapping := range fieldMappings {
		fields = append(fields, fieldMapping.SourceFieldName)
	}

	return strings.Join(fields, ",")
}

func (pg PostgresImpl) getNewCursorPosition(lastRow data.Row, schema data.Schema, sync views.Sync) *string {
	if sync.SourceCursorField == nil {
		return nil
	}

	if lastRow == nil {
		return nil
	}

	var cursorFieldPos int
	var cursorFieldType data.FieldType
	for i := range schema {
		if schema[i].Name == *sync.SourceCursorField {
			cursorFieldPos = i
			cursorFieldType = schema[i].Type
		}
	}

	// TODO: make sure we don't miss any rows
	// we sort rows by cursor field so just take the last row
	var newCursorPos string
	switch cursorFieldType {
	case data.FieldTypeInteger, data.FieldTypeNumber:
		newCursorPos = fmt.Sprintf("%v", lastRow[cursorFieldPos])
	default:
		newCursorPos = fmt.Sprintf("'%v'", lastRow[cursorFieldPos])
	}

	return &newCursorPos
}

func (pg PostgresImpl) Write(
	ctx context.Context,
	destinationConnection views.FullConnection,
	destinationOptions DestinationOptions,
	object views.Object,
	sync views.Sync,
	fieldMappings []views.FieldMapping,
	rowsC <-chan []data.Row,
	writeOutputC chan<- WriteOutput,
	errC chan<- error,
) {
	errC <- errors.New("postgres destination not implemented")
}
