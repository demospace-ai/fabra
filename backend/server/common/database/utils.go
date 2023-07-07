package database

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullString struct{ sql.NullString }

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

func NewNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: true}}
}

func NewNullStringFromPtr(s *string) NullString {
	if s == nil {
		return NullString{}
	} else {
		return NullString{sql.NullString{String: *s, Valid: true}}
	}
}

type NullInt64 struct{ sql.NullInt64 }

func (i NullInt64) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int64)
	}
	return []byte(`null`), nil
}

func NewNullInt64(i int64) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: i, Valid: true}}
}

type NullTime struct{ sql.NullTime }

func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time)
	}
	return []byte(`null`), nil
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{sql.NullTime{Time: t, Valid: true}}
}

var EmptyNullInt64 = NullInt64{sql.NullInt64{Valid: false}}

// Assigns the database.NullString based on the key is null, "", or does not exist.
// For example:
// { input: null } sets stringVal to null
// { input: "" } sets stringVal to ""
// { } leaves the stringVal unchanged
func SetNullStringFromRaw(input json.RawMessage, stringVal *NullString) error {
	if len(input) > 0 { // if key exists in JSON input
		if string(input) == "null" { // value is null
			*stringVal = NullString{}
		} else {
			var nativeString string
			err := json.Unmarshal(input, &nativeString)
			if err != nil {
				return err
			}
			*stringVal = NewNullString(nativeString)
		}
	}
	return nil
}
