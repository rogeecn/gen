package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Point represents a PostgreSQL point field
type Point Field

func NewPoint(table, column string, opts ...Option) Point {
	return Point{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Point) Eq(v types.Point) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Point) Neq(v types.Point) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// DistanceTo returns the distance between two points using the <-> operator.
func (f Point) DistanceTo(p types.Point) Float64 {
	return Float64{expr{e: clause.Expr{SQL: "? <-> ?", Vars: []interface{}{f.RawExpr(), p}}}}
}

// WithinCircle tests if point is within a circle: point <@ circle
func (f Point) WithinCircle(c types.Circle) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), c}}}
}

// WithinBox tests if point is within a box: point <@ box
func (f Point) WithinBox(b types.Box) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), b}}}
}

// WithinPolygon tests if point is within a polygon: point <@ polygon
func (f Point) WithinPolygon(poly types.Polygon) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), poly}}}
}
