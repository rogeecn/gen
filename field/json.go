package field

import (
	"strings"

	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// JSON represents a JSON/JSONB column
type JSON Field

func NewJSON(table, column string, opts ...Option) JSON {
	return JSON{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f JSON) Eq(v types.JSON) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f JSON) Neq(v types.JSON) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// ======================== key-path compare & like ========================
// KeyEq compares the JSON value at dot-separated path equals the given value.
// Path example: "a.b.c". Value type is dynamic (int, bool, string, ...).
func (f JSON) KeyEq(dotKey string, v interface{}) Expr  { return f.keyCmp("=", dotKey, v) }
func (f JSON) KeyNeq(dotKey string, v interface{}) Expr { return f.keyCmp("<>", dotKey, v) }
func (f JSON) KeyGt(dotKey string, v interface{}) Expr  { return f.keyCmp(">", dotKey, v) }
func (f JSON) KeyGte(dotKey string, v interface{}) Expr { return f.keyCmp(">=", dotKey, v) }
func (f JSON) KeyLt(dotKey string, v interface{}) Expr  { return f.keyCmp("<", dotKey, v) }
func (f JSON) KeyLte(dotKey string, v interface{}) Expr { return f.keyCmp("<=", dotKey, v) }

// KeyLike performs a LIKE on the text of the JSON value at the given path.
// Path example: "a.b.c". Pattern should include wildcards (e.g., %foo%).
func (f JSON) KeyLike(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? LIKE ?", Vars: []interface{}{textExpr, pattern}}}
}

// KeyILike performs a case-insensitive LIKE (ILIKE) on the text value.
func (f JSON) KeyILike(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? ILIKE ?", Vars: []interface{}{textExpr, pattern}}}
}

// KeyRegexp performs a regex match using PostgreSQL '~'.
func (f JSON) KeyRegexp(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? ~ ?", Vars: []interface{}{textExpr, pattern}}}
}

// KeyIRegexp performs a case-insensitive regex match using PostgreSQL '~*'.
func (f JSON) KeyIRegexp(dotKey, pattern string) Expr {
	textExpr := f.extractTextByDot(dotKey)
	return expr{e: clause.Expr{SQL: "? ~* ?", Vars: []interface{}{textExpr, pattern}}}
}

// keyCmp builds a comparison for the JSON value at path with dynamic typing.
// - bool  -> cast extracted text to boolean
// - number-> cast extracted text to numeric
// - other -> compare as text
func (f JSON) keyCmp(op, dotKey string, v interface{}) Expr {
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

// extractTextByDot builds json_extract_path_text(?::json, keys...)
// with keys parsed from a dot-separated path like "a.b.c".
func (f JSON) extractTextByDot(dotKey string) clause.Expression {
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
