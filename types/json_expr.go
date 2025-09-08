package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JSONQueryExpression json query expression, implements clause.Expression interface to use as querier
type JSONQueryExpression struct {
	column      string
	keys        []string
	hasKeys     bool
	equals      bool
	likes       bool
	equalsValue interface{}
	extract     bool
	path        string
}

// JSONQuery query column as json
func JSONQuery(column string) *JSONQueryExpression {
	return &JSONQueryExpression{column: column}
}

// Extract extract json with path
func (jsonQuery *JSONQueryExpression) Extract(path string) *JSONQueryExpression {
	jsonQuery.extract = true
	jsonQuery.path = path
	return jsonQuery
}

// HasKey returns clause.Expression
func (jsonQuery *JSONQueryExpression) HasKey(keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.hasKeys = true
	return jsonQuery
}

// Keys returns clause.Expression
func (jsonQuery *JSONQueryExpression) Equals(value interface{}, keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.equals = true
	jsonQuery.equalsValue = value
	return jsonQuery
}

// Likes return clause.Expression
func (jsonQuery *JSONQueryExpression) Likes(value interface{}, keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.likes = true
	jsonQuery.equalsValue = value
	return jsonQuery
}

// Build implements clause.Expression
func (jsonQuery *JSONQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch {
		case jsonQuery.extract:
			builder.WriteString(fmt.Sprintf("json_extract_path_text(%v::json,", stmt.Quote(jsonQuery.column)))
			stmt.AddVar(builder, jsonQuery.path)
			builder.WriteByte(')')
		case jsonQuery.hasKeys:
			if len(jsonQuery.keys) > 0 {
				stmt.WriteQuoted(jsonQuery.column)
				stmt.WriteString("::jsonb")
				for _, key := range jsonQuery.keys[0 : len(jsonQuery.keys)-1] {
					stmt.WriteString(" -> ")
					stmt.AddVar(builder, key)
				}
				stmt.WriteString(" ? ")
				stmt.AddVar(builder, jsonQuery.keys[len(jsonQuery.keys)-1])
			}
		case jsonQuery.equals:
			if len(jsonQuery.keys) > 0 {
				builder.WriteString(fmt.Sprintf("json_extract_path_text(%v::json,", stmt.Quote(jsonQuery.column)))
				for idx, key := range jsonQuery.keys {
					if idx > 0 {
						builder.WriteByte(',')
					}
					stmt.AddVar(builder, key)
				}
				builder.WriteString(") = ")
				if _, ok := jsonQuery.equalsValue.(string); ok {
					stmt.AddVar(builder, jsonQuery.equalsValue)
				} else {
					stmt.AddVar(builder, fmt.Sprint(jsonQuery.equalsValue))
				}
			}
		case jsonQuery.likes:
			if len(jsonQuery.keys) > 0 {
				builder.WriteString(fmt.Sprintf("json_extract_path_text(%v::json,", stmt.Quote(jsonQuery.column)))
				for idx, key := range jsonQuery.keys {
					if idx > 0 {
						builder.WriteByte(',')
					}
					stmt.AddVar(builder, key)
				}
				builder.WriteString(") LIKE ")
				if _, ok := jsonQuery.equalsValue.(string); ok {
					stmt.AddVar(builder, jsonQuery.equalsValue)
				} else {
					stmt.AddVar(builder, fmt.Sprint(jsonQuery.equalsValue))
				}
			}
		}
	}
}

// JSONOverlapsExpression JSON_OVERLAPS expression, implements clause.Expression interface to use as querier
type JSONOverlapsExpression struct {
	column clause.Expression
	val    string
}

// JSONOverlaps query column as json
func JSONOverlaps(column clause.Expression, value string) *JSONOverlapsExpression {
	return &JSONOverlapsExpression{
		column: column,
		val:    value,
	}
}

// Build implements clause.Expression
// only mysql support JSON_OVERLAPS
func (json *JSONOverlapsExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		builder.WriteString("(")
		json.column.Build(builder)
		builder.WriteString("::jsonb && ")
		builder.AddVar(stmt, json.val)
		builder.WriteString("::jsonb)")
	}
}

type columnExpression string

func Column(col string) columnExpression {
	return columnExpression(col)
}

func (col columnExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		builder.WriteString(stmt.Quote(string(col)))
	}
}

const prefix = "$."

func jsonQueryJoin(keys []string) string {
	if len(keys) == 1 {
		return prefix + keys[0]
	}

	n := len(prefix)
	n += len(keys) - 1
	for i := 0; i < len(keys); i++ {
		n += len(keys[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(prefix)
	b.WriteString(keys[0])
	for _, key := range keys[1:] {
		b.WriteString(".")
		b.WriteString(key)
	}
	return b.String()
}

// JSONSetExpression json set expression, implements clause.Expression interface to use as updater
type JSONSetExpression struct {
	column     string
	path2value map[string]interface{}
	mutex      sync.RWMutex
}

// JSONSet update fields of json column
func JSONSet(column string) *JSONSetExpression {
	return &JSONSetExpression{column: column, path2value: make(map[string]interface{})}
}

// Set returns a JSONSetExpression with the specified path and value.
//
// Example JSON:
//
//	{
//		"age": 20,
//		"name": "json-1",
//		"orgs": {"orga": "orgv"},
//		"tags": ["tag1", "tag2"]
//	}
//
// Usage:
//
// For PostgreSQL, use curly braces with comma-separated path: "{age}", "{name}", "{orgs,orga}", "{tags,0}", "{tags,1}".
//
//	DB.UpdateColumn("attr", JSONSet("attr").Set("{orgs,orga}", "bar"))
func (jsonSet *JSONSetExpression) Set(path string, value interface{}) *JSONSetExpression {
	jsonSet.mutex.Lock()
	jsonSet.path2value[path] = value
	jsonSet.mutex.Unlock()
	return jsonSet
}

// Build implements clause.Expression
// support mysql, sqlite and postgres
func (jsonSet *JSONSetExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		var expr clause.Expression = columnExpression(jsonSet.column)
		for path, value := range jsonSet.path2value {
			if _, ok = value.(clause.Expression); ok {
				expr = gorm.Expr("JSONB_SET(?,?,?)", expr, path, value)
				continue
			} else {
				b, _ := json.Marshal(value)
				expr = gorm.Expr("JSONB_SET(?,?,?)", expr, path, string(b))
			}
		}
		stmt.AddVar(builder, expr)
	}
}

func JSONArrayQuery(column string) *JSONArrayExpression {
	return &JSONArrayExpression{
		column: column,
	}
}

type JSONArrayExpression struct {
	contains    bool
	in          bool
	overlap     bool
	containsAll bool
	containedBy bool
	column      string
	keys        []string
	equalsValue interface{}
}

// Contains checks if column[keys] contains the value given. The keys parameter is only supported for MySQL and SQLite.
func (json *JSONArrayExpression) Contains(value interface{}, keys ...string) *JSONArrayExpression {
	json.contains = true
	json.equalsValue = value
	json.keys = keys
	return json
}

// Overlap checks if array column overlaps with the given array (&&)
func (json *JSONArrayExpression) Overlap(value interface{}) *JSONArrayExpression {
	json.overlap = true
	json.equalsValue = value
	return json
}

// ContainsAll checks if array column contains all elements of the given array (@>)
func (json *JSONArrayExpression) ContainsAll(value interface{}) *JSONArrayExpression {
	json.containsAll = true
	json.equalsValue = value
	return json
}

// ContainedBy checks if array column is contained by the given array (<@)
func (json *JSONArrayExpression) ContainedBy(value interface{}) *JSONArrayExpression {
	json.containedBy = true
	json.equalsValue = value
	return json
}

// In checks if columns[keys] is in the array value given. This method is only supported for MySQL and SQLite.
func (json *JSONArrayExpression) In(value interface{}, keys ...string) *JSONArrayExpression {
	json.in = true
	json.keys = keys
	json.equalsValue = value
	return json
}

// Build implements clause.Expression
func (json *JSONArrayExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch {
		case json.contains:
			builder.WriteString(stmt.Quote(json.column))
			builder.WriteString(" ? ")
			builder.AddVar(stmt, json.equalsValue)
		case json.in:
			builder.WriteString(stmt.Quote(json.column))
			builder.WriteString(" IN (")
			switch v := json.equalsValue.(type) {
			case []interface{}:
				for i, val := range v {
					if i > 0 {
						builder.WriteString(", ")
					}
					builder.AddVar(stmt, val)
				}
			case []string:
				for i, val := range v {
					if i > 0 {
						builder.WriteString(", ")
					}
					builder.AddVar(stmt, val)
				}
			default:
				builder.AddVar(stmt, json.equalsValue)
			}
			builder.WriteString(")")
		case json.overlap:
			builder.WriteString(stmt.Quote(json.column))
			builder.WriteString(" && ")
			builder.AddVar(stmt, json.equalsValue)
		case json.containsAll:
			builder.WriteString(stmt.Quote(json.column))
			builder.WriteString(" @> ")
			builder.AddVar(stmt, json.equalsValue)
		case json.containedBy:
			builder.WriteString(stmt.Quote(json.column))
			builder.WriteString(" <@ ")
			builder.AddVar(stmt, json.equalsValue)
		}
	}
}
