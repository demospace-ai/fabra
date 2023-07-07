package data

import (
	"context"
	"errors"
)

type Schema []Field

type FieldType string

const (
	FieldTypeString      FieldType = "STRING"
	FieldTypeInteger     FieldType = "INTEGER"
	FieldTypeNumber      FieldType = "NUMBER"
	FieldTypeDateTimeTz  FieldType = "DATETIME_TZ"
	FieldTypeDateTimeNtz FieldType = "DATETIME_NTZ"
	FieldTypeTimestamp   FieldType = "TIMESTAMP"
	FieldTypeTimeTz      FieldType = "TIME_TZ"
	FieldTypeTimeNtz     FieldType = "TIME_NTZ"
	FieldTypeDate        FieldType = "DATE"
	FieldTypeBoolean     FieldType = "BOOLEAN"
	FieldTypeArray       FieldType = "ARRAY"
	FieldTypeJson        FieldType = "JSON"
)

type Field struct {
	Name string    `json:"name"`
	Type FieldType `json:"type"`
}

type Row []any

var ErrDone = errors.New("no more items in fabra iterator")

type RowIterator interface {
	Next(ctx context.Context) (Row, error)
	Schema() Schema
}

type QueryResults struct {
	Data   []Row  `json:"data"`
	Schema Schema `json:"schema"`
}
