package field

import (
	"gorm.io/gorm/clause"
)

// Array represents a PostgreSQL array column
// It provides common array operators for Postgres: contains, contained-by, overlaps.
type Array Field

func NewArray(table, column string, opts ...Option) Array {
	return Array{expr: expr{col: toColumn(table, column, opts...)}}
}

// Eq compares entire array equality. Caller provides a slice or driver.Valuer like pq.Array.
func (f Array) Eq(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? = ?", Vars: []interface{}{f.RawExpr(), v}}}
}

func (f Array) Neq(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? <> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains uses @> operator (left contains right)
func (f Array) Contains(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy uses <@ operator (left is contained by right)
func (f Array) ContainedBy(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Overlaps uses && operator (arrays have elements in common)
func (f Array) Overlaps(v interface{}) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}
