package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Int4Range represents a PostgreSQL int4range field
type Int4Range Field

func NewInt4Range(table, column string, opts ...Option) Int4Range {
	return Int4Range{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Int4Range) Eq(v types.Int4Range) Expr {
	return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}}
}

func (f Int4Range) Neq(v types.Int4Range) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

func (f Int4Range) Overlaps(v types.Int4Range) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests range contains ( @> ) another range
func (f Int4Range) Contains(v types.Int4Range) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests range is contained by ( <@ ) another range
func (f Int4Range) ContainedBy(v types.Int4Range) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictLeft tests range strictly left of ( << ) another range
func (f Int4Range) StrictLeft(v types.Int4Range) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictRight tests range strictly right of ( >> ) another range
func (f Int4Range) StrictRight(v types.Int4Range) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Adjacent tests ranges are adjacent ( -|- )
func (f Int4Range) Adjacent(v types.Int4Range) Expr {
	return expr{e: clause.Expr{SQL: "? -|- ?", Vars: []interface{}{f.RawExpr(), v}}}
}
