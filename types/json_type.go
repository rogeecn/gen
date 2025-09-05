package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// JSONType give a generic data type for json encoded data.
type JSONType[T any] struct {
	data T
}

func NewJSONType[T any](data T) JSONType[T] {
	return JSONType[T]{
		data: data,
	}
}

// Data return data with generic Type T
func (j JSONType[T]) Data() T {
	return j.data
}

// Set replaces the underlying data to given value
func (j *JSONType[T]) Set(data T) { j.data = data }

// Edit applies a mutator to the underlying data in-place
func (j *JSONType[T]) Edit(mutator func(*T)) {
	if mutator != nil {
		mutator(&j.data)
	}
}

// Value return json value, implement driver.Valuer interface
func (j JSONType[T]) Value() (driver.Value, error) {
	return json.Marshal(j.data)
}

// Scan scan value into JSONType[T], implements sql.Scanner interface
func (j *JSONType[T]) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &j.data)
}

// MarshalJSON to output non base64 encoded []byte
func (j JSONType[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.data)
}

// UnmarshalJSON to deserialize []byte
func (j *JSONType[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.data)
}

// GormDataType gorm common data type
func (JSONType[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSONType[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// Postgres only
	return "JSONB"
}

func (js JSONType[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := js.MarshalJSON()
	return gorm.Expr("?", string(data))
}

// JSONSlice give a generic data type for json encoded slice data.
type JSONSlice[T any] []T

func NewJSONSlice[T any](s []T) JSONSlice[T] {
	return JSONSlice[T](s)
}

// Set replaces the slice with given value
func (j *JSONSlice[T]) Set(s []T) { *j = JSONSlice[T](s) }

// Append appends elements to the slice
func (j *JSONSlice[T]) Append(elems ...T) { *j = append(*j, elems...) }

// Edit applies a mutator to the underlying slice in-place
func (j *JSONSlice[T]) Edit(mutator func(*[]T)) {
	if mutator != nil {
		s := ([]T)(*j)
		mutator(&s)
		*j = JSONSlice[T](s)
	}
}

// Value return json value, implement driver.Valuer interface
func (j JSONSlice[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan scan value into JSONType[T], implements sql.Scanner interface
func (j *JSONSlice[T]) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &j)
}

// GormDataType gorm common data type
func (JSONSlice[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSONSlice[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// Postgres only
	return "JSONB"
}

func (j JSONSlice[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := json.Marshal(j)
	return gorm.Expr("?", string(data))
}
