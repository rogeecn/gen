package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Circle represents a PostgreSQL circle field
type Circle Field

func NewCircle(table, column string, opts ...Option) Circle {
	return Circle{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Circle) Eq(v types.Circle) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Circle) Neq(v types.Circle) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// Overlaps tests if two circles overlap: &&
func (f Circle) Overlaps(v types.Circle) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsCircle tests if circle contains another circle: @>
func (f Circle) ContainsCircle(v types.Circle) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedByCircle tests if circle is contained by another circle: <@
func (f Circle) ContainedByCircle(v types.Circle) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsPoint tests if circle contains a point: circle @> point
func (f Circle) ContainsPoint(p types.Point) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), p}}}
}
