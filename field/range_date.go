package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// DateRange represents a PostgreSQL daterange field
type DateRange Field

func NewDateRange(table, column string, opts ...Option) DateRange {
	return DateRange{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f DateRange) Eq(v types.DateRange) Expr {
	return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}}
}

func (f DateRange) Neq(v types.DateRange) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

func (f DateRange) Overlaps(v types.DateRange) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests range contains ( @> ) another range
func (f DateRange) Contains(v types.DateRange) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests range is contained by ( <@ ) another range
func (f DateRange) ContainedBy(v types.DateRange) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictLeft tests range strictly left of ( << ) another range
func (f DateRange) StrictLeft(v types.DateRange) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictRight tests range strictly right of ( >> ) another range
func (f DateRange) StrictRight(v types.DateRange) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Adjacent tests ranges are adjacent ( -|- )
func (f DateRange) Adjacent(v types.DateRange) Expr {
	return expr{e: clause.Expr{SQL: "? -|- ?", Vars: []interface{}{f.RawExpr(), v}}}
}
