package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// NumRange represents a PostgreSQL numrange field
type NumRange Field

func NewNumRange(table, column string, opts ...Option) NumRange {
	return NumRange{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f NumRange) Eq(v types.NumRange) Expr { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f NumRange) Neq(v types.NumRange) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

func (f NumRange) Overlaps(v types.NumRange) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests range contains ( @> ) another range
func (f NumRange) Contains(v types.NumRange) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests range is contained by ( <@ ) another range
func (f NumRange) ContainedBy(v types.NumRange) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictLeft tests range strictly left of ( << ) another range
func (f NumRange) StrictLeft(v types.NumRange) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictRight tests range strictly right of ( >> ) another range
func (f NumRange) StrictRight(v types.NumRange) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Adjacent tests ranges are adjacent ( -|- )
func (f NumRange) Adjacent(v types.NumRange) Expr {
	return expr{e: clause.Expr{SQL: "? -|- ?", Vars: []interface{}{f.RawExpr(), v}}}
}
