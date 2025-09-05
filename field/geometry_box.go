package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Box represents a PostgreSQL box field
type Box Field

func NewBox(table, column string, opts ...Option) Box {
	return Box{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Box) Eq(v types.Box) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Box) Neq(v types.Box) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// Overlaps tests if two boxes overlap: &&
func (f Box) Overlaps(v types.Box) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests if box contains another box: @>
func (f Box) Contains(v types.Box) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests if box is contained by another box: <@
func (f Box) ContainedBy(v types.Box) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsPoint tests if box contains a point: box @> point
func (f Box) ContainsPoint(p types.Point) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), p}}}
}
