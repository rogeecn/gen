package field

import (
	"strings"

	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// JSONB represents a PostgreSQL JSONB column
type JSONB Field

func NewJSONB(table, column string, opts ...Option) JSONB {
	return JSONB{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f JSONB) Eq(v types.JSON) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f JSONB) Neq(v types.JSON) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// Contains uses @> (left contains right)
func (f JSONB) Contains(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy uses <@ (left is contained by right)
func (f JSONB) ContainedBy(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// HasKey tests if jsonb has the given key or array contains the element
func (f JSONB) HasKey(key string) Expr {
	return expr{e: clause.Expr{SQL: "? ? ?", Vars: []interface{}{f.RawExpr(), key}}}
}

// HasAnyKeys tests if any keys exist: column ?| ARRAY[keys]
func (f JSONB) HasAnyKeys(keys ...string) Expr {
	if len(keys) == 0 {
		return expr{e: clause.Expr{SQL: "1=0"}}
	}
	ph := "?" + strings.Repeat(",?", len(keys)-1)
	sql := "? ?| ARRAY[" + ph + "]"
	vars := make([]interface{}, 0, 1+len(keys))
	vars = append(vars, f.RawExpr())
	for i := range keys {
		vars = append(vars, keys[i])
	}
	return expr{e: clause.Expr{SQL: sql, Vars: vars}}
}

// HasAllKeys tests if all keys exist: column ?& ARRAY[keys]
func (f JSONB) HasAllKeys(keys ...string) Expr {
	if len(keys) == 0 {
		return expr{e: clause.Expr{SQL: "1=1"}}
	}
	ph := "?" + strings.Repeat(",?", len(keys)-1)
	sql := "? ?& ARRAY[" + ph + "]"
	vars := make([]interface{}, 0, 1+len(keys))
	vars = append(vars, f.RawExpr())
	for i := range keys {
		vars = append(vars, keys[i])
	}
	return expr{e: clause.Expr{SQL: sql, Vars: vars}}
}

// ExtractText builds json_extract_path_text(column::json, ...keys)
func (f JSONB) ExtractText(keys ...string) String {
	if len(keys) == 0 {
		return String{expr{e: clause.Expr{SQL: "?::text", Vars: []interface{}{f.RawExpr()}}}}
	}
	ph := "?" + strings.Repeat(",?", len(keys)-1)
	sql := "json_extract_path_text(?::json," + ph + ")"
	vars := make([]interface{}, 0, 1+len(keys))
	vars = append(vars, f.RawExpr())
	for i := range keys {
		vars = append(vars, keys[i])
	}
	return String{expr{e: clause.Expr{SQL: sql, Vars: vars}}}
}

// ======================== key-path compare & like ========================
// KeyEq compares the JSON value at dot-separated path equals the given value.
// Path example: "a.b.c". Value type is dynamic (int, bool, string, ...).
func (f JSONB) KeyEq(dotKey string, v interface{}) Expr  { return f.keyCmp("=", dotKey, v) }
func (f JSONB) KeyNeq(dotKey string, v interface{}) Expr { return f.keyCmp("<>", dotKey, v) }
func (f JSONB) KeyGt(dotKey string, v interface{}) Expr  { return f.keyCmp(">", dotKey, v) }
func (f JSONB) KeyGte(dotKey string, v interface{}) Expr { return f.keyCmp(">=", dotKey, v) }
func (f JSONB) KeyLt(dotKey string, v interface{}) Expr  { return f.keyCmp("<", dotKey, v) }
func (f JSONB) KeyLte(dotKey string, v interface{}) Expr { return f.keyCmp("<=", dotKey, v) }

// KeyLike performs a LIKE on the text of the JSON value at the given path.
// Path example: "a.b.c". Pattern should include wildcards (e.g., %foo%).
func (f JSONB) KeyLike(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? LIKE ?", Vars: []interface{}{textExpr, pattern}}}
}

// KeyILike performs a case-insensitive LIKE (ILIKE) on the text value.
func (f JSONB) KeyILike(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? ILIKE ?", Vars: []interface{}{textExpr, pattern}}}
}

// KeyRegexp performs a regex match using PostgreSQL '~'.
func (f JSONB) KeyRegexp(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? ~ ?", Vars: []interface{}{textExpr, pattern}}}
}

// KeyIRegexp performs a case-insensitive regex match using PostgreSQL '~*'.
func (f JSONB) KeyIRegexp(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? ~* ?", Vars: []interface{}{textExpr, pattern}}}
}

// keyCmp builds a comparison for the JSON value at path with dynamic typing.
// - bool  -> cast extracted text to boolean
// - number-> cast extracted text to numeric
// - other -> compare as text
func (f JSONB) keyCmp(op, dotKey string, v interface{}) Expr {
	textExpr := f.extractTextByDot(dotKey)
	switch v.(type) {
	case bool:
		return expr{e: clause.Expr{SQL: "CAST(? AS boolean) " + op + " ?", Vars: []interface{}{textExpr, v}}}
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return expr{e: clause.Expr{SQL: "CAST(? AS numeric) " + op + " ?", Vars: []interface{}{textExpr, v}}}
	default:
		// default to text compare
		return expr{e: clause.Expr{SQL: "? " + op + " ?", Vars: []interface{}{textExpr, v}}}
	}
}

// extractTextByDot builds json_extract_path_text(column::json, keys...)
// with keys parsed from a dot-separated path like "a.b.c".
func (f JSONB) extractTextByDot(dotKey string) clause.Expression {
	if dotKey == "" {
		return clause.Expr{SQL: "?::text", Vars: []interface{}{f.RawExpr()}}
	}
	parts := make([]string, 0, 4)
	for _, k := range strings.Split(dotKey, ".") {
		if k == "" {
			continue
		}
		parts = append(parts, k)
	}
	if len(parts) == 0 {
		return clause.Expr{SQL: "?::text", Vars: []interface{}{f.RawExpr()}}
	}
	ph := "?" + strings.Repeat(",?", len(parts)-1)
	sql := "json_extract_path_text(?::json," + ph + ")"
	vars := make([]interface{}, 0, 1+len(parts))
	vars = append(vars, f.RawExpr())
	for i := range parts {
		vars = append(vars, parts[i])
	}
	return clause.Expr{SQL: sql, Vars: vars}
}
