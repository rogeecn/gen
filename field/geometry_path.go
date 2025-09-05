package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Path represents a PostgreSQL path field
type Path Field

func NewPath(table, column string, opts ...Option) Path {
	return Path{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Path) Eq(v types.Path) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Path) Neq(v types.Path) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// Overlaps tests if paths overlap (bounding boxes intersect): &&
func (f Path) Overlaps(v types.Path) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}
