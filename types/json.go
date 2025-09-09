package types

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// JSON defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
type JSON json.RawMessage

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("null")
		return nil
	}
	var bytes []byte
	if s, ok := value.(fmt.Stringer); ok {
		bytes = []byte(s.String())
	} else {
		switch v := value.(type) {
		case []byte:
			if len(v) > 0 {
				bytes = make([]byte, len(v))
				copy(bytes, v)
			}
		case string:
			bytes = []byte(v)
		default:
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}
	}

	result := json.RawMessage(bytes)
	*j = JSON(result)
	return nil
}

// MarshalJSON to output non base64 encoded []byte
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (j *JSON) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*j = JSON(result)
	return err
}

func (j JSON) String() string {
	return string(j)
}

// GormDataType gorm common data type
func (JSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// Postgres only
	return "JSONB"
}

func (js JSON) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(js) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := js.MarshalJSON()
	return gorm.Expr("?", string(data))
}

// Convenience methods for JSON operations

// IsValid checks if the JSON is valid
func (j JSON) IsValid() bool {
	return json.Valid([]byte(j))
}

// IsNull checks if the JSON represents null
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// IsEmpty checks if the JSON is empty or null
func (j JSON) IsEmpty() bool {
	return len(j) == 0 || string(j) == "null" || string(j) == "{}" || string(j) == "[]"
}

// Pretty returns a formatted JSON string with indentation
func (j JSON) Pretty() string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(j), "", "  "); err != nil {
		return string(j) // fallback to original if formatting fails
	}
	return buf.String()
}

// Compact removes unnecessary whitespace from JSON
func (j JSON) Compact() JSON {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(j)); err != nil {
		return j // fallback to original if compacting fails
	}
	return JSON(buf.Bytes())
}

// Size returns the byte length of the JSON
func (j JSON) Size() int {
	return len(j)
}

// Clone creates a copy of the JSON
func (j JSON) Clone() JSON {
	if len(j) == 0 {
		return nil
	}
	clone := make([]byte, len(j))
	copy(clone, j)
	return JSON(clone)
}

// Equals compares two JSON values for equality
func (j JSON) Equals(other JSON) bool {
	return bytes.Equal([]byte(j), []byte(other))
}

// GetKeys returns the keys of a JSON object (returns empty slice for non-objects)
func (j JSON) GetKeys() []string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(j), &obj); err != nil {
		return []string{}
	}

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	return keys
}

// HasKey checks if the JSON object has the specified key
func (j JSON) HasKey(key string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(j), &obj); err != nil {
		return false
	}
	_, exists := obj[key]
	return exists
}

// GetValue extracts a value from JSON object by key
func (j JSON) GetValue(key string) (interface{}, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(j), &obj); err != nil {
		return nil, err
	}
	return obj[key], nil
}

// SetValue sets a value in the JSON object (only works for JSON objects)
func (j *JSON) SetValue(key string, value interface{}) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(*j), &obj); err != nil {
		// If it's not an object, create a new one
		obj = make(map[string]interface{})
	}

	obj[key] = value
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	*j = JSON(data)
	return nil
}

// DeleteKey removes a key from the JSON object
func (j *JSON) DeleteKey(key string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(*j), &obj); err != nil {
		return err
	}

	delete(obj, key)
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	*j = JSON(data)
	return nil
}

// Merge merges another JSON object into this one (only works for JSON objects)
func (j *JSON) Merge(other JSON) error {
	var thisObj, otherObj map[string]interface{}

	if err := json.Unmarshal([]byte(*j), &thisObj); err != nil {
		thisObj = make(map[string]interface{})
	}

	if err := json.Unmarshal([]byte(other), &otherObj); err != nil {
		return err
	}

	for key, value := range otherObj {
		thisObj[key] = value
	}

	data, err := json.Marshal(thisObj)
	if err != nil {
		return err
	}

	*j = JSON(data)
	return nil
}
