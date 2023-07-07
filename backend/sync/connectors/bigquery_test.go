package connectors_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/golang/mock/gomock"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/input"
	mock_query "go.fabra.io/server/common/mocks"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/test"
	"go.fabra.io/server/common/views"
	"go.fabra.io/sync/connectors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BigQueryConnector", func() {
	var (
		sourceConnection      views.FullConnection
		destinationConnection views.FullConnection
		sync                  views.Sync
		fieldMappings         []views.FieldMapping
		object                views.Object
	)

	BeforeEach(func() {
		org := test.CreateOrganization(db)
		endCustomerID := "abc123"
		source, sourceConn := test.CreateSource(db, org.ID, endCustomerID)
		sourceConnection = views.ConvertFullConnection(sourceConn)
		destination, destConn := test.CreateDestination(db, org.ID)
		destinationConnection = views.ConvertFullConnection(destConn)

		objectModel := test.CreateObject(db, org.ID, destination.ID, models.SyncModeFullAppend)
		objectFields := test.CreateObjectFields(db, objectModel.ID, []input.ObjectField{
			{Name: "string", Type: data.FieldTypeString},
			{Name: "integer", Type: data.FieldTypeInteger},
			{Name: "boolean", Type: data.FieldTypeBoolean},
			{Name: "datetime_tz", Type: data.FieldTypeDateTimeTz},
			{Name: "datetime_ntz", Type: data.FieldTypeDateTimeNtz},
			{Name: "json", Type: data.FieldTypeJson, Optional: true},
		})
		object = views.ConvertObject(objectModel, objectFields)
		sync = views.ConvertSync(test.CreateSync(db, org.ID, endCustomerID, source.ID, objectModel.ID, models.SyncModeFullAppend))
		fieldMappings = views.ConvertFieldMappings(test.CreateFieldMappings(db, sync.ID, []input.FieldMapping{
			{SourceFieldName: "source_string", SourceFieldType: data.FieldTypeString, DestinationFieldId: objectFields[0].ID},
			{SourceFieldName: "source_integer", SourceFieldType: data.FieldTypeInteger, DestinationFieldId: objectFields[1].ID},
			{SourceFieldName: "source_boolean", SourceFieldType: data.FieldTypeBoolean, DestinationFieldId: objectFields[2].ID},
			{SourceFieldName: "source_datetime_tz", SourceFieldType: data.FieldTypeDateTimeTz, DestinationFieldId: objectFields[3].ID},
			{SourceFieldName: "source_datetime_ntz", SourceFieldType: data.FieldTypeDateTimeNtz, DestinationFieldId: objectFields[4].ID},
			{SourceFieldName: "source_json", SourceFieldType: data.FieldTypeJson, DestinationFieldId: objectFields[5].ID},
		}), objectFields)
	})

	Describe("Read", func() {

		It("reads in batches", func() {
			ctrl := gomock.NewController(GinkgoT())
			client := mock_query.NewMockWarehouseClient(ctrl)
			defer ctrl.Finish()

			rows := make([]data.Row, 2_000_000)
			for i := 0; i < 2_000_000; i++ {
				rows[i] = data.Row{"string", 1, false, "2006-01-02 15:04:05.000-07:00", "2006-01-02 15:04:05.000", map[string]int{"hello": 123}}
			}
			iterator := test.NewMockIterator(
				rows,
				data.Schema{
					{Name: "source_string", Type: data.FieldTypeString},
					{Name: "source_integer", Type: data.FieldTypeInteger},
					{Name: "source_boolean", Type: data.FieldTypeBoolean},
					{Name: "source_datetime_tz", Type: data.FieldTypeDateTimeTz},
					{Name: "source_datetime_ntz", Type: data.FieldTypeDateTimeNtz},
					{Name: "source_json", Type: data.FieldTypeJson},
				},
			)
			client.EXPECT().GetQueryIterator(
				gomock.Any(),
				"SELECT source_string,source_integer,source_boolean,source_datetime_tz,source_datetime_ntz,source_json FROM namespace.table;",
			).Return(iterator, nil)

			connector := connectors.NewBigQueryConnector(client)
			rowsC := make(chan []data.Row)
			readOutputC := make(chan connectors.ReadOutput)
			errC := make(chan error)

			go func() {
				defer GinkgoRecover()
				defer func() { close(readOutputC) }() // close the output channel so the test completes in case of an error
				connector.Read(context.TODO(), sourceConnection, sync, fieldMappings, rowsC, readOutputC, errC)
			}()
			readOutput, resultRows, numBatches, err := waitForRead(rowsC, readOutputC, errC)

			Expect(err).To(BeNil())
			Expect(readOutput.CursorPosition).To(BeNil())
			Expect(resultRows).To(Equal(rows))
			Expect(numBatches).To(Equal(2))
		})

		It("queries for first cursor correctly", func() {
			ctrl := gomock.NewController(GinkgoT())
			client := mock_query.NewMockWarehouseClient(ctrl)
			defer ctrl.Finish()

			sync.SyncMode = models.SyncModeIncrementalAppend
			cursorField := "source_datetime_tz"
			sync.SourceCursorField = &cursorField

			rows := make([]data.Row, 10)
			for i := 0; i < 10; i++ {
				rows[i] = data.Row{"string", 1, false, "2006-01-02 15:04:05.000-07:00", "2006-01-02 15:04:05.000", map[string]int{"hello": 123}}
			}
			iterator := test.NewMockIterator(
				rows,
				data.Schema{
					{Name: "source_string", Type: data.FieldTypeString},
					{Name: "source_integer", Type: data.FieldTypeInteger},
					{Name: "source_boolean", Type: data.FieldTypeBoolean},
					{Name: "source_datetime_tz", Type: data.FieldTypeDateTimeTz},
					{Name: "source_datetime_ntz", Type: data.FieldTypeDateTimeNtz},
					{Name: "source_json", Type: data.FieldTypeJson},
				},
			)
			client.EXPECT().GetQueryIterator(
				gomock.Any(),
				"SELECT source_string,source_integer,source_boolean,source_datetime_tz,source_datetime_ntz,source_json FROM namespace.table ORDER BY source_datetime_tz ASC;",
			).Return(iterator, nil)

			connector := connectors.NewBigQueryConnector(client)
			rowsC := make(chan []data.Row)
			readOutputC := make(chan connectors.ReadOutput)
			errC := make(chan error)

			go func() {
				defer GinkgoRecover()
				defer func() { close(readOutputC) }() // close the output channel so the test completes in case of an error
				connector.Read(context.TODO(), sourceConnection, sync, fieldMappings, rowsC, readOutputC, errC)
			}()
			readOutput, resultRows, numBatches, err := waitForRead(rowsC, readOutputC, errC)

			Expect(err).To(BeNil())
			Expect(*readOutput.CursorPosition).To(Equal("'2006-01-02 15:04:05.000-07:00'"))
			Expect(resultRows).To(Equal(rows))
			Expect(numBatches).To(Equal(1))
		})

		It("queries for successive cursors correctly", func() {
			ctrl := gomock.NewController(GinkgoT())
			client := mock_query.NewMockWarehouseClient(ctrl)
			defer ctrl.Finish()

			sync.SyncMode = models.SyncModeIncrementalAppend
			cursorPosition := "'2007-01-02 15:04:05.000-07:00'"
			cursorField := "source_datetime_tz"
			sync.CursorPosition = &cursorPosition
			sync.SourceCursorField = &cursorField

			rows := []data.Row{
				{"string", 1, false, "2008-01-02 15:04:05.000-07:00", "2006-01-02 15:04:05.000", map[string]int{"hello": 123}},
			}

			iterator := test.NewMockIterator(
				rows,
				data.Schema{
					{Name: "source_string", Type: data.FieldTypeString},
					{Name: "source_integer", Type: data.FieldTypeInteger},
					{Name: "source_boolean", Type: data.FieldTypeBoolean},
					{Name: "source_datetime_tz", Type: data.FieldTypeDateTimeTz},
					{Name: "source_datetime_ntz", Type: data.FieldTypeDateTimeNtz},
					{Name: "source_json", Type: data.FieldTypeJson},
				},
			)
			client.EXPECT().GetQueryIterator(
				gomock.Any(),
				"SELECT source_string,source_integer,source_boolean,source_datetime_tz,source_datetime_ntz,source_json FROM namespace.table WHERE source_datetime_tz > '2007-01-02 15:04:05.000-07:00' ORDER BY source_datetime_tz ASC;",
			).Return(iterator, nil)

			connector := connectors.NewBigQueryConnector(client)
			rowsC := make(chan []data.Row)
			readOutputC := make(chan connectors.ReadOutput)
			errC := make(chan error)

			go func() {
				defer GinkgoRecover()
				defer func() { close(readOutputC) }() // close the output channel so the test completes in case of an error
				connector.Read(context.TODO(), sourceConnection, sync, fieldMappings, rowsC, readOutputC, errC)
			}()
			readOutput, resultRows, numBatches, err := waitForRead(rowsC, readOutputC, errC)

			Expect(err).To(BeNil())
			Expect(*readOutput.CursorPosition).To(Equal("'2008-01-02 15:04:05.000-07:00'"))
			Expect(len(resultRows)).To(Equal(1))
			Expect(resultRows).To(Equal(rows))
			Expect(numBatches).To(Equal(1))
		})

		It("does not wrap integer cursor with ''", func() {
			ctrl := gomock.NewController(GinkgoT())
			client := mock_query.NewMockWarehouseClient(ctrl)
			defer ctrl.Finish()

			sync.SyncMode = models.SyncModeIncrementalAppend
			cursorPosition := "1"
			cursorField := "source_integer"
			sync.CursorPosition = &cursorPosition
			sync.SourceCursorField = &cursorField

			rows := []data.Row{
				{"string", 2, false, "2006-01-02 15:04:05.000-07:00", "2006-01-02 15:04:05.000", map[string]int{"hello": 123}},
			}

			iterator := test.NewMockIterator(
				rows,
				data.Schema{
					{Name: "source_string", Type: data.FieldTypeString},
					{Name: "source_integer", Type: data.FieldTypeInteger},
					{Name: "source_boolean", Type: data.FieldTypeBoolean},
					{Name: "source_datetime_tz", Type: data.FieldTypeDateTimeTz},
					{Name: "source_datetime_ntz", Type: data.FieldTypeDateTimeNtz},
					{Name: "source_json", Type: data.FieldTypeJson},
				},
			)
			client.EXPECT().GetQueryIterator(
				gomock.Any(),
				"SELECT source_string,source_integer,source_boolean,source_datetime_tz,source_datetime_ntz,source_json FROM namespace.table WHERE source_integer > 1 ORDER BY source_integer ASC;",
			).Return(iterator, nil)

			connector := connectors.NewBigQueryConnector(client)
			rowsC := make(chan []data.Row)
			readOutputC := make(chan connectors.ReadOutput)
			errC := make(chan error)

			go func() {
				defer GinkgoRecover()
				defer func() { close(readOutputC) }() // close the output channel so the test completes in case of an error
				connector.Read(context.TODO(), sourceConnection, sync, fieldMappings, rowsC, readOutputC, errC)
			}()
			readOutput, resultRows, numBatches, err := waitForRead(rowsC, readOutputC, errC)

			Expect(err).To(BeNil())
			Expect(*readOutput.CursorPosition).To(Equal("2"))
			Expect(len(resultRows)).To(Equal(1))
			Expect(resultRows).To(Equal(rows))
			Expect(numBatches).To(Equal(1))
		})
	})

	Describe("Write", func() {
		/*
			TODO

			Write
			1. Unused object fields interspersed
			2. Out of order object fields
			3. Various types mapped correctly to the bigquery expected format
			4. Mapping multiple JSON fields to a single object

		*/

		It("writes correctly", func() {
			ctrl := gomock.NewController(GinkgoT())
			client := mock_query.NewMockWarehouseClient(ctrl)
			defer ctrl.Finish()

			rows := make([]data.Row, 10)
			for i := 0; i < 10; i++ {
				rows[i] = data.Row{"string", 2, false, "2006-01-02 15:04:05.000-07:00", "2006-01-02 15:04:05.000", map[string]int{"hello": 123}}
			}

			csvRows := make([]string, 10)
			for i := 0; i < 10; i++ {
				// Strings should be quoted, while integers, numbers, and datetimes should not. JSON should be double quoted
				csvRows[i] = "\"string\",2,false,2006-01-02 15:04:05.000-07:00,2006-01-02 15:04:05.000,\"{\"\"hello\"\":123}\",abc123"
			}
			csvData := strings.Join(csvRows, "\n")

			client.EXPECT().StageData(
				gomock.Any(),
				csvData,
				MockStagingOptions{Bucket: "staging"},
			).Return(nil)

			client.EXPECT().LoadFromStaging(
				gomock.Any(),
				"namespace",
				"table",
				MockLoadOptions{
					"staging",
					bigquery.Schema{
						{Name: "string", Type: bigquery.StringFieldType, Required: true},
						{Name: "integer", Type: bigquery.IntegerFieldType, Required: true},
						{Name: "boolean", Type: bigquery.BooleanFieldType, Required: true},
						{Name: "datetime_tz", Type: bigquery.TimestampFieldType, Required: true},
						{Name: "datetime_ntz", Type: bigquery.DateTimeFieldType, Required: true},
						{Name: "json", Type: bigquery.JSONFieldType, Required: false},
						{Name: "end_customer_id", Type: bigquery.StringFieldType, Required: true},
					},
					bigquery.WriteAppend,
				},
			).Return(nil)

			client.EXPECT().CleanUpStagingData(
				gomock.Any(),
				MockStagingOptions{Bucket: "staging"},
			).Return(nil)

			connector := connectors.NewBigQueryConnector(client)
			rowsC := make(chan []data.Row)
			writeOutputC := make(chan connectors.WriteOutput)
			errC := make(chan error)

			go func() {
				defer GinkgoRecover()
				defer func() { close(writeOutputC) }() // close the output channel so the test completes in case of an error
				connector.Write(context.TODO(), destinationConnection, connectors.DestinationOptions{StagingBucket: "staging"}, object, sync, fieldMappings, rowsC, writeOutputC, errC)
			}()

			rowsC <- rows
			close(rowsC)

			writeOutput, err := waitForWrite(writeOutputC, errC)

			Expect(err).To(BeNil())
			Expect(writeOutput.RowsWritten).To(Equal(10))
		})
	})
})

func waitForRead(
	rowsC <-chan []data.Row,
	readOutputC <-chan connectors.ReadOutput,
	errC <-chan error,
) (*connectors.ReadOutput, []data.Row, int, error) {
	var readOutput connectors.ReadOutput
	var rows []data.Row
	var readDone bool
	numBatches := 0
	for {
		if readDone {
			break
		}

		select {
		case err := <-errC:
			if err != nil {
				return nil, nil, numBatches, err
			}
		case rowBatch := <-rowsC:
			rows = append(rows, rowBatch...)
			numBatches++
		case readOutput = <-readOutputC:
			readDone = true
		}
	}

	return &readOutput, rows, numBatches, nil
}

func waitForWrite(
	writeOutputC <-chan connectors.WriteOutput,
	errC <-chan error,
) (*connectors.WriteOutput, error) {
	var writeOutput connectors.WriteOutput
	for {
		select {
		case err := <-errC:
			if err != nil {
				return nil, errors.Wrap(err, "bigquery test waitForWrite")
			}
		case writeOutput = <-writeOutputC:
			return &writeOutput, nil
		}
	}

}

type MockStagingOptions struct {
	Bucket string
}

func (so MockStagingOptions) Matches(x interface{}) bool {
	actual, ok := x.(query.StagingOptions)
	if !ok {
		return false
	}

	return so.Bucket == actual.Bucket
}

func (so MockStagingOptions) String() string {
	return fmt.Sprintf("{%s, ANY}", so.Bucket)
}

type MockLoadOptions struct {
	Bucket    string
	Schema    bigquery.Schema
	WriteMode bigquery.TableWriteDisposition
}

func (so MockLoadOptions) Matches(x interface{}) bool {
	actual, ok := x.(query.LoadOptions)
	if !ok {
		return false
	}

	gcsReference := fmt.Sprintf("gs://%s/", so.Bucket)
	return gcsReference == strings.Join(strings.SplitAfter(actual.GcsReference, "/")[0:3], "") &&
		reflect.DeepEqual(so.Schema, actual.BigQuerySchema) &&
		so.WriteMode == actual.WriteMode
}

func (so MockLoadOptions) Got(got interface{}) string {
	actual, ok := got.(query.LoadOptions)
	if !ok {
		return fmt.Sprintf("unexpected struct: %v", got)
	}

	fieldList := make([]string, len(actual.BigQuerySchema))
	for i, field := range actual.BigQuerySchema {
		value := "nil"
		if field != nil {
			value = fmt.Sprintf("%+v", *field)
		}
		fieldList[i] = value
	}

	outputSchema := fmt.Sprintf("[\n\t%s\n]", strings.Join(fieldList, ",\n\t"))

	return fmt.Sprintf("{\n%s, \n%v, \n%s\n}", strings.Join(strings.SplitAfter(actual.GcsReference, "/")[0:3], ""), outputSchema, actual.WriteMode)
}

func (so MockLoadOptions) String() string {

	fieldList := make([]string, len(so.Schema))
	for i, field := range so.Schema {
		value := "nil"
		if field != nil {
			value = fmt.Sprintf("%+v", *field)
		}
		fieldList[i] = value
	}

	outputSchema := fmt.Sprintf("[\n\t%s\n]", strings.Join(fieldList, ",\n\t"))

	return fmt.Sprintf("{\ngs://%s/, \n%v, \n%s\n}", so.Bucket, outputSchema, so.WriteMode)
}
