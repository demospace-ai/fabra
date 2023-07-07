package query

import (
	"context"
	"fmt"
	"strings"

	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbApiClient struct {
	Username          string
	Password          string
	Host              string
	ConnectionOptions string
}

type MongoQuery struct {
	Database   string               `json:"database"`
	Collection string               `json:"collection"`
	Filter     bson.D               `json:"filter"`
	Options    *options.FindOptions `json:"options"`
}

type mongoDbIterator struct {
	schema data.Schema
	cursor *mongo.Cursor
	client *mongo.Client
}

func (it *mongoDbIterator) Next(ctx context.Context) (data.Row, error) {
	if it.cursor.Next(ctx) {
		var row bson.D
		err := it.cursor.Decode(&row)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.mongoDbIterator.Next) decoding row")
		}

		return convertMongoDbRow(row, it.schema), nil
	}

	defer it.cursor.Close(ctx)
	defer it.client.Disconnect(ctx)
	err := it.cursor.Err()
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.mongoDbIterator.Next) cursor error")
	}

	return nil, data.ErrDone
}

func (it *mongoDbIterator) Schema() data.Schema {
	return it.schema
}

func (mc MongoDbApiClient) openConnection(ctx context.Context) (*mongo.Client, error) {
	connectionString := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?%s",
		mc.Username,
		mc.Password,
		mc.Host,
		mc.ConnectionOptions,
	)
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetServerAPIOptions(serverAPIOptions)
	return mongo.Connect(ctx, clientOptions)
}

func (mc MongoDbApiClient) GetTables(ctx context.Context, namespace string) ([]string, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetTables) opening connection")
	}

	defer client.Disconnect(ctx)

	db := client.Database(namespace)
	tables, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetTables) getting tables")
	}

	return tables, nil
}

func (mc MongoDbApiClient) GetSchema(ctx context.Context, namespace string, tableName string) (data.Schema, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetSchema) opening connection")
	}

	defer client.Disconnect(ctx)

	db := client.Database(namespace)
	collection := db.Collection(tableName)
	fields, err := getFields(collection)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetSchema) getting fields")
	}

	fieldTypes, err := getFieldTypes(collection, fields)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetSchema) getting field types")
	}

	return convertMongoDbSchema(fieldTypes), nil
}

func (mc MongoDbApiClient) GetFieldValues(ctx context.Context, namespace string, tableName string, fieldName string) ([]any, error) {
	// TODO
	return nil, nil
}

func (mc MongoDbApiClient) GetNamespaces(ctx context.Context) ([]string, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetNamespaces) opening connection")
	}

	defer client.Disconnect(ctx)

	databaseNames, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetNamespaces) listing database names")
	}

	return databaseNames, nil
}

func (mc MongoDbApiClient) RunQuery(ctx context.Context, queryString string, args ...any) (*data.QueryResults, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.RunQuery) opening connection")
	}
	defer client.Disconnect(ctx)

	var mongoQuery MongoQuery
	err = bson.Unmarshal([]byte(queryString), &mongoQuery)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.RunQuery) unmarshalling query")
	}

	schemaC := make(chan data.Schema)
	errC := make(chan error)
	go func() {
		schema, err := mc.GetSchema(ctx, mongoQuery.Database, mongoQuery.Collection)
		schemaC <- schema
		errC <- err

		close(schemaC)
		close(errC)
	}()

	db := client.Database(mongoQuery.Database)
	collection := db.Collection(mongoQuery.Collection)
	cursor, err := collection.Find(
		ctx,
		mongoQuery.Filter,
		mongoQuery.Options,
	)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.RunQuery) running query")
	}
	defer cursor.Close(ctx)

	var rows bson.A
	err = cursor.All(ctx, &rows)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.RunQuery) getting rows")
	}

	schema := <-schemaC
	err = <-errC
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.RunQuery) getting schema")
	}

	return &data.QueryResults{
		Schema: schema,
		Data:   convertMongoDbRows(rows, schema),
	}, nil
}

func (mc MongoDbApiClient) GetQueryIterator(ctx context.Context, queryString string) (data.RowIterator, error) {
	client, err := mc.openConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetQueryIterator) opening connection")
	}

	var mongoQuery MongoQuery
	err = bson.Unmarshal([]byte(queryString), &mongoQuery)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetQueryIterator) unmarshalling query")
	}

	schemaC := make(chan data.Schema)
	errC := make(chan error)
	go func() {
		schema, err := mc.GetSchema(ctx, mongoQuery.Database, mongoQuery.Collection)
		schemaC <- schema
		errC <- err
		close(schemaC)
		close(errC)
	}()

	db := client.Database(mongoQuery.Database)
	collection := db.Collection(mongoQuery.Collection)
	cursor, err := collection.Find(
		ctx,
		mongoQuery.Filter,
		mongoQuery.Options,
	)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetQueryIterator) running query")
	}

	schema := <-schemaC
	err = <-errC
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.GetQueryIterator) getting schema")
	}

	return &mongoDbIterator{
		schema: schema,
		cursor: cursor,
		client: client,
	}, nil
}

func getFields(collection *mongo.Collection) ([]string, error) {
	ctx := context.TODO()
	cursor, err := collection.Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{{Key: "$limit", Value: 10000}},
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "data", Value: bson.D{
					{Key: "$objectToArray", Value: "$$ROOT"},
				}},
			}}},
			bson.D{{Key: "$unwind", Value: "$data"}},
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "fields", Value: bson.D{
					{Key: "$addToSet", Value: "$data.k"},
				}},
			}}},
		},
	)
	if err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.getFields) running query")
	}

	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.getFields) decoding results")
	}

	var fields []string
	for _, field := range results[0]["fields"].(bson.A) {
		fields = append(fields, field.(string))
	}

	return fields, nil
}

func getFieldTypes(collection *mongo.Collection, fields []string) (map[string]string, error) {
	ctx := context.TODO()
	fieldTypes := make(map[string]string)
	for _, field := range fields {
		cursor, err := collection.Aggregate(
			ctx,
			mongo.Pipeline{
				bson.D{{Key: "$limit", Value: 10000}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "fieldType", Value: bson.D{
						{Key: "$type", Value: "$" + field},
					}},
				}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "fieldType", Value: "$fieldType"},
					}},
					{Key: "count", Value: bson.D{
						{Key: "$sum", Value: "1"},
					}},
				}}},
			},
		)
		if err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.getFields) running query")
		}

		defer cursor.Close(ctx)

		var typesForField []string
		for cursor.Next(ctx) {
			var result bson.M
			err = cursor.Decode(&result)
			if err != nil {
				return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.getFields) decoding results")
			}

			// Even if most of the fields are missing/null, we use the most common non-null type if one exists
			fieldType := result["_id"].(bson.M)["fieldType"]
			if fieldType == "missing" || fieldType == "null" {
				continue
			}

			typesForField = append(typesForField, fieldType.(string))
		}

		// If there are truly no types, then the field type is null
		if len(typesForField) == 0 {
			typesForField = append(typesForField, "null")
		}

		if err = cursor.Err(); err != nil {
			return nil, errors.Wrap(errors.WrapCustomerVisibleError(err), "(query.MongoDbApiClient.getFields) iterating results")
		}

		fieldTypes[field] = typesForField[0]
	}

	return fieldTypes, nil
}

func convertMongoDbRows(mongoDbRows bson.A, schema data.Schema) []data.Row {
	var rows []data.Row
	for _, mongoDbRow := range mongoDbRows {
		rows = append(rows, convertMongoDbRow(mongoDbRow.(bson.D), schema))
	}

	return rows
}

func convertMongoDbRow(mongoDbRow bson.D, schema data.Schema) data.Row {
	// TODO: convert the values to the expected Fabra Golang types
	valueMap := make(map[string]any)
	for _, keyPair := range mongoDbRow {
		valueMap[keyPair.Key] = keyPair.Value
	}

	// make sure every result is in the same order by looping through schema
	row := make(data.Row, len(schema))
	for i, field := range schema {
		// if the value is missing, the map will return nil because that is the default value for "any"
		value := convertMongoDbValue(valueMap[field.Name], field.Type)
		row[i] = value
	}

	return row
}

func convertMongoDbValue(mongoDbValue any, fieldType data.FieldType) any {
	if mongoDbValue == nil {
		return nil
	}

	switch fieldType {
	case data.FieldTypeDateTimeTz:
		return mongoDbValue.(primitive.DateTime).Time().UTC().Format(FABRA_TIMESTAMP_TZ_FORMAT)
	case data.FieldTypeJson:
		return ToMap(mongoDbValue.(bson.D))
	case data.FieldTypeArray:
		return ToArray(mongoDbValue.(bson.A))
	default:
		return mongoDbValue
	}
}

func convertMongoDbSchema(fieldTypes map[string]string) data.Schema {
	var schema data.Schema
	for fieldName, fieldType := range fieldTypes {
		schema = append(schema, data.Field{
			Name: fieldName,
			Type: getMongoDbFieldType(fieldType),
		})
	}

	return schema
}

func getMongoDbFieldType(mongoDbType string) data.FieldType {
	uppercased := strings.ToUpper(mongoDbType)
	switch uppercased {
	case "INT", "INT32", "LONG":
		return data.FieldTypeInteger
	case "DATE", "DATETIME":
		// MongoDB dates/datetimes are in UTC
		return data.FieldTypeDateTimeTz
	case "TIMESTAMP":
		return data.FieldTypeTimestamp
	case "DECIMAL", "DOUBLE", "FLOAT64":
		return data.FieldTypeNumber
	case "ARRAY":
		return data.FieldTypeArray
	case "OBJECT":
		return data.FieldTypeJson
	case "BOOL":
		return data.FieldTypeBoolean
	default:
		return data.FieldTypeString
	}
}

func ToMap(bsonMap bson.D) map[string]any {
	output := map[string]any{}
	for _, pair := range bsonMap {
		if pair.Value == nil {
			continue
		}

		switch value := pair.Value.(type) {
		case bson.D:
			output[pair.Key] = ToMap(value)
		case bson.A:
			output[pair.Key] = ToArray(value)
		default:
			output[pair.Key] = value
		}
	}
	return output
}

func ToArray(bsonArray bson.A) []any {
	output := []any{}
	for _, rawValue := range bsonArray {
		switch value := rawValue.(type) {
		case bson.D:
			output = append(output, ToMap(value))
		case bson.A:
			output = append(output, ToArray(value))
		default:
			output = append(output, value)
		}
	}

	return output
}

func CreateMongoQueryString(mongoQuery MongoQuery) string {
	mongoQueryBytes, err := bson.Marshal(mongoQuery)
	if err != nil {
		// this should never happen anyway so just panic
		panic(err)
	}

	return string(mongoQueryBytes)
}
