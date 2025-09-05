package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Polygon represents a PostgreSQL polygon field
type Polygon Field

func NewPolygon(table, column string, opts ...Option) Polygon {
	return Polygon{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Polygon) Eq(v types.Polygon) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Polygon) Neq(v types.Polygon) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// Overlaps tests if polygons overlap: &&
func (f Polygon) Overlaps(v types.Polygon) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsPolygon tests if polygon contains another polygon: @>
func (f Polygon) ContainsPolygon(v types.Polygon) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedByPolygon tests if polygon is contained by another polygon: <@
func (f Polygon) ContainedByPolygon(v types.Polygon) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsPoint tests if polygon contains a point: polygon @> point
func (f Polygon) ContainsPoint(p types.Point) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), p}}}
}
