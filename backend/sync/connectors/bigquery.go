package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/views"
)

type BigQueryImpl struct {
	client query.WarehouseClient
}

func NewBigQueryConnector(client query.WarehouseClient) Connector {
	return BigQueryImpl{
		client: client,
	}
}

func (bq BigQueryImpl) Read(
	ctx context.Context,
	sourceConnection views.FullConnection,
	sync views.Sync,
	fieldMappings []views.FieldMapping,
	rowsC chan<- []data.Row,
	readOutputC chan<- ReadOutput,
	errC chan<- error,
) {
	readQuery := bq.getReadQuery(sourceConnection, sync, fieldMappings)
	iterator, err := bq.client.GetQueryIterator(ctx, readQuery)
	if err != nil {
		errC <- errors.Wrap(err, "(connectors.BigQueryImpl.Read) getting iterator")
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
				errC <- errors.Wrap(err, "(connectors.BigQueryImpl.Read) iterating data")
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

	newCursorPosition := bq.getNewCursorPosition(lastRow, iterator.Schema(), sync)
	readOutputC <- ReadOutput{
		CursorPosition: newCursorPosition,
	}

	close(rowsC)
	close(errC)
}

func (bq BigQueryImpl) getReadQuery(sourceConnection views.FullConnection, sync views.Sync, fieldMappings []views.FieldMapping) string {
	var queryString string
	if sync.CustomJoin != nil {
		queryString = *sync.CustomJoin
	} else {
		selectString := bq.getSelectString(fieldMappings)
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

func (bq BigQueryImpl) getSelectString(fieldMappings []views.FieldMapping) string {
	fields := []string{}
	for _, fieldMapping := range fieldMappings {
		fields = append(fields, fieldMapping.SourceFieldName)
	}

	return strings.Join(fields, ",")
}

func (bq BigQueryImpl) getNewCursorPosition(lastRow data.Row, schema data.Schema, sync views.Sync) *string {
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
	case data.FieldTypeInteger:
		newCursorPos = fmt.Sprintf("%v", lastRow[cursorFieldPos])
	default:
		newCursorPos = fmt.Sprintf("'%v'", lastRow[cursorFieldPos])
	}

	return &newCursorPos
}

func (bq BigQueryImpl) Write(
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
	// always clean up the data in the storage bucket
	objectPrefix := uuid.New().String()
	wildcardObject := fmt.Sprintf("%s-*", objectPrefix)
	gcsReference := fmt.Sprintf("gs://%s/%s", destinationOptions.StagingBucket, wildcardObject)

	batchNum := 0
	rowsWritten := 0
	for {
		rows, more := <-rowsC
		if !more {
			break
		}

		rowsWritten += len(rows)
		objectName := fmt.Sprintf("%s-%d", objectPrefix, batchNum)
		err := bq.stageBatch(ctx, rows, fieldMappings, object, sync, destinationOptions, bq.client, objectName)
		if err != nil {
			errC <- errors.Wrap(err, "(connectors.BigQueryImpl.Write) staging batch")
			return
		}

		// use a separate context for cleanup so it won't get cancelled
		defer bq.client.CleanUpStagingData(context.Background(), query.StagingOptions{Bucket: destinationOptions.StagingBucket, Object: objectName})

		batchNum++
	}

	if rowsWritten > 0 {
		writeMode := bq.toBigQueryWriteMode(sync.SyncMode)
		csvSchema := bq.createCsvSchema(*object.EndCustomerIDField, object.ObjectFields)
		err := bq.client.LoadFromStaging(ctx, *object.Namespace, *object.TableName, query.LoadOptions{
			GcsReference:   gcsReference,
			BigQuerySchema: csvSchema,
			WriteMode:      writeMode,
		})
		if err != nil {
			errC <- errors.Wrap(err, "(connectors.BigQueryImpl.Write) loading data from staging")
			return
		}
	}

	writeOutputC <- WriteOutput{
		RowsWritten: rowsWritten,
	}

	close(errC)
}

func (bq BigQueryImpl) stageBatch(
	ctx context.Context,
	rows []data.Row,
	fieldMappings []views.FieldMapping,
	object views.Object, sync views.Sync,
	destinationOptions DestinationOptions,
	destClient query.WarehouseClient,
	objectName string,
) error {
	// count the fields since there may be multiple mappings for a single JSON object in the destination
	// also track where each field should go in the output row based on the order of object fields.
	// use the count as the index since we want to skip omitted fields in the output
	numFields := 0
	objectFieldsIdToIndex := make(map[int64]int)
	for _, objectField := range object.ObjectFields {
		if !objectField.Omit {
			objectFieldsIdToIndex[objectField.ID] = numFields
			numFields++
		}
	}

	// extra field for end customer ID
	numFields++

	// allocate the arrays and reuse them to save memory
	rowStrings := make([]string, len(rows))
	rowTokens := make([]string, numFields)
	rowTokens[numFields-1] = sync.EndCustomerID // end customer ID will be the same for every row

	// write to temporary table in destination
	for i, row := range rows {
		indexToJsonValueMap := make(map[int]map[string]any)
		for j, value := range row {
			fieldMapping := fieldMappings[j]
			sourceType := fieldMapping.SourceFieldType
			destFieldIdx := objectFieldsIdToIndex[fieldMapping.DestinationFieldId]

			// just collect the raw values into a map
			if fieldMapping.IsJsonField {
				existing, ok := indexToJsonValueMap[destFieldIdx]
				if !ok {
					existing = make(map[string]any)
					indexToJsonValueMap[destFieldIdx] = existing
				}

				existing[fieldMapping.SourceFieldName] = value
			} else {
				if value == nil {
					// empty string for null values will be interpreted as null when loading from csv
					rowTokens[j] = ""
				} else {
					switch sourceType {
					case data.FieldTypeJson:
						jsonStr, err := getBigQueryJsonString(value)
						if err != nil {
							return errors.Wrap(err, "(connectors.BigQueryImpl.stageBatch)")
						}
						rowTokens[destFieldIdx] = jsonStr
					case data.FieldTypeString:
						// escape the string so commas don't break the CSV schema
						rowTokens[destFieldIdx] = fmt.Sprintf("\"%v\"", value)
					default:
						rowTokens[destFieldIdx] = fmt.Sprintf("%v", value)
					}
				}
			}
		}

		// insert the many-to-one mappings into the row tokens slice
		for key, value := range indexToJsonValueMap {
			jsonStr, err := getBigQueryJsonString(value)
			if err != nil {
				return errors.Wrap(err, "(connectors.BigQueryImpl.stageBatch)")
			}

			rowTokens[key] = jsonStr
		}

		rowString := strings.Join(rowTokens, ",")
		rowStrings[i] = rowString
	}

	stagingOptions := query.StagingOptions{Bucket: destinationOptions.StagingBucket, Object: objectName}
	err := destClient.StageData(ctx, strings.Join(rowStrings, "\n"), stagingOptions)
	if err != nil {
		return err
	}

	return nil
}

// JSON-like values need to be escaped according to BigQuery expectations. Even if the destination
// type is not JSON, it is necessary to escape to avoid issues
// https://cloud.google.com/bigquery/docs/reference/standard-sql/json-data#load_from_csv_files
func getBigQueryJsonString(value any) (string, error) {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "(connectors.BigQueryImpl.getBigQueryJsonString)")
	}
	return fmt.Sprintf("\"%s\"", strings.ReplaceAll(string(jsonStr), "\"", "\"\"")), nil
}

func (bq BigQueryImpl) toBigQueryWriteMode(syncMode models.SyncMode) bigquery.TableWriteDisposition {
	switch syncMode {
	case models.SyncModeFullOverwrite:
		return bigquery.WriteTruncate
	case models.SyncModeFullAppend:
		return bigquery.WriteAppend
	case models.SyncModeIncrementalAppend:
		return bigquery.WriteAppend
	case models.SyncModeIncrementalUpdate:
		// incremental update loads updated/new rows into a temp table, before merging with the destination
		return bigquery.WriteTruncate
	default:
		return bigquery.WriteAppend
	}
}

func (bq BigQueryImpl) createCsvSchema(endCustomerIDColumn string, orderedObjectFields []views.ObjectField) bigquery.Schema {
	var csvSchema bigquery.Schema
	for _, objectField := range orderedObjectFields {
		if !objectField.Omit {
			field := bigquery.FieldSchema{
				Name:     objectField.Name,
				Type:     getBigQueryType(objectField.Type),
				Required: !objectField.Optional,
			}
			csvSchema = append(csvSchema, &field)
		}
	}

	endCustomerIDField := bigquery.FieldSchema{
		Name:     endCustomerIDColumn,
		Type:     bigquery.StringFieldType,
		Required: true,
	}
	csvSchema = append(csvSchema, &endCustomerIDField)

	return csvSchema
}

func getBigQueryType(fieldType data.FieldType) bigquery.FieldType {
	switch fieldType {
	case data.FieldTypeInteger:
		return bigquery.IntegerFieldType
	case data.FieldTypeNumber:
		return bigquery.NumericFieldType
	case data.FieldTypeBoolean:
		return bigquery.BooleanFieldType
	case data.FieldTypeTimestamp, data.FieldTypeDateTimeTz:
		return bigquery.TimestampFieldType
	case data.FieldTypeDateTimeNtz:
		return bigquery.DateTimeFieldType
	case data.FieldTypeJson:
		return bigquery.JSONFieldType
	case data.FieldTypeDate:
		return bigquery.DateFieldType
	case data.FieldTypeTimeTz, data.FieldTypeTimeNtz:
		return bigquery.TimeFieldType
	default:
		return bigquery.StringFieldType
	}
}
