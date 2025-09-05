package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Int8Range represents a PostgreSQL int8range field
type Int8Range Field

func NewInt8Range(table, column string, opts ...Option) Int8Range {
	return Int8Range{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Int8Range) Eq(v types.Int8Range) Expr {
	return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}}
}

func (f Int8Range) Neq(v types.Int8Range) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

func (f Int8Range) Overlaps(v types.Int8Range) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests range contains ( @> ) another range
func (f Int8Range) Contains(v types.Int8Range) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests range is contained by ( <@ ) another range
func (f Int8Range) ContainedBy(v types.Int8Range) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictLeft tests range strictly left of ( << ) another range
func (f Int8Range) StrictLeft(v types.Int8Range) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictRight tests range strictly right of ( >> ) another range
func (f Int8Range) StrictRight(v types.Int8Range) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Adjacent tests ranges are adjacent ( -|- )
func (f Int8Range) Adjacent(v types.Int8Range) Expr {
	return expr{e: clause.Expr{SQL: "? -|- ?", Vars: []interface{}{f.RawExpr(), v}}}
}
