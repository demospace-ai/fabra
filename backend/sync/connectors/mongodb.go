package connectors

import (
	"context"
	"fmt"
	"time"

	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/views"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbImpl struct {
	queryService query.QueryService
}

func NewMongoDbConnector(queryService query.QueryService) Connector {
	return MongoDbImpl{
		queryService: queryService,
	}
}

func (md MongoDbImpl) Read(
	ctx context.Context,
	sourceConnection views.FullConnection,
	sync views.Sync,
	fieldMappings []views.FieldMapping,
	rowsC chan<- []data.Row,
	readOutputC chan<- ReadOutput,
	errC chan<- error,
) {
	connectionModel := views.ConvertConnectionView(sourceConnection)

	sourceClient, err := md.queryService.GetClient(ctx, connectionModel)
	if err != nil {
		errC <- err
		return
	}

	readQuery, err := md.getReadQuery(connectionModel, sync, fieldMappings)
	if err != nil {
		errC <- err
		return
	}

	queryString := query.CreateMongoQueryString(*readQuery)

	iterator, err := sourceClient.GetQueryIterator(ctx, queryString)
	if err != nil {
		errC <- err
		return
	}

	currentIndex := 0
	var rowBatch []data.Row
	var lastRow data.Row
	schema := iterator.Schema()
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

		reordered := reorderMongoRow(row, schema, fieldMappings)
		rowBatch = append(rowBatch, reordered)
		lastRow = reordered
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

	// pass field mappings not schema since the row will be reordered
	newCursorPosition, err := md.getNewCursorPosition(lastRow, fieldMappings, sync)
	if err != nil {
		errC <- err
		return
	}

	readOutputC <- ReadOutput{
		CursorPosition: newCursorPosition,
	}

	close(rowsC)
	close(errC)
}

// TODO: only read 10,000 rows at once or something
func (md MongoDbImpl) getReadQuery(sourceConnection *models.Connection, sync views.Sync, fieldMappings []views.FieldMapping) (*query.MongoQuery, error) {
	projection := createProjection(fieldMappings)
	mongoQuery := query.MongoQuery{
		Database:   *sync.Namespace,
		Collection: *sync.TableName,
		Filter:     bson.D{},
		Options:    options.Find(),
	}

	mongoQuery.Options.SetProjection(projection)

	if sync.SyncMode.UsesCursor() {
		// order by cursor field to simplify
		mongoQuery.Options.SetSort(bson.D{
			bson.E{
				Key:   *sync.SourceCursorField,
				Value: 1, // 1 indicates ascending order
			},
		})

		if sync.CursorPosition != nil {
			sourceCursorFieldType, err := getSourceCursorFieldType(*sync.SourceCursorField, fieldMappings)
			if err != nil {
				return nil, errors.Wrap(err, "(connectors.MongoDbImpl.getReadQuery) error getting source cursor field type")
			}

			var comparisonValue any
			switch *sourceCursorFieldType {
			case data.FieldTypeDateTimeTz:
				timeCursor, err := time.Parse(query.FABRA_TIMESTAMP_TZ_FORMAT, *sync.CursorPosition)
				if err != nil {
					return nil, errors.Wrap(err, "(connectors.MongoDbImpl.getReadQuery) error parsing cursor position")
				}
				comparisonValue = primitive.NewDateTimeFromTime(timeCursor)
			default:
				comparisonValue = *sync.CursorPosition
			}

			// TODO: allow choosing other operators (rows smaller than current cursor, etc.)
			mongoQuery.Filter = bson.D{
				bson.E{
					Key: *sync.SourceCursorField,
					Value: bson.D{
						bson.E{
							Key:   "$gt",
							Value: comparisonValue,
						},
					},
				},
			}
		}
	}

	return &mongoQuery, nil
}

// Used to ensure every field is in the correct order, and to omit the _id field
func createProjection(fieldMappings []views.FieldMapping) bson.D {
	projection := bson.D{
		bson.E{
			Key:   "_id",
			Value: 0,
		},
	}

	for _, fieldMapping := range fieldMappings {
		projection = append(projection, bson.E{
			Key:   fieldMapping.SourceFieldName,
			Value: 1,
		})
	}

	return projection
}

func (md MongoDbImpl) getNewCursorPosition(lastRow data.Row, fieldMappings []views.FieldMapping, sync views.Sync) (*string, error) {
	if sync.SourceCursorField == nil {
		return nil, nil
	}

	if lastRow == nil {
		return nil, nil
	}

	var cursorFieldPos int
	var cursorFieldType data.FieldType
	for i := range fieldMappings {
		if fieldMappings[i].SourceFieldName == *sync.SourceCursorField {
			cursorFieldPos = i
			cursorFieldType = fieldMappings[i].SourceFieldType
		}
	}

	// TODO: make sure we don't miss any rows
	// we sort rows by cursor field so just take the last row
	var newCursorPos string
	switch cursorFieldType {
	case data.FieldTypeInteger, data.FieldTypeNumber, data.FieldTypeTimestamp, data.FieldTypeDateTimeTz:
		newCursorPos = fmt.Sprintf("%v", lastRow[cursorFieldPos])
	default:
		newCursorPos = fmt.Sprintf("'%v'", lastRow[cursorFieldPos])
	}

	return &newCursorPos, nil
}

func reorderMongoRow(unorderedRow data.Row, schema data.Schema, fieldMappings []views.FieldMapping) data.Row {
	fieldNameToValue := make(map[string]any)
	for i, field := range schema {
		fieldNameToValue[field.Name] = unorderedRow[i]
	}

	orderedRow := make(data.Row, len(fieldMappings))
	for i, fieldMapping := range fieldMappings {
		orderedRow[i] = fieldNameToValue[fieldMapping.SourceFieldName]
	}

	return orderedRow
}

func (md MongoDbImpl) Write(
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
	errC <- errors.New("mongodb destination not implemented")
}
