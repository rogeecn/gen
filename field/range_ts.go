package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// TsRange represents a PostgreSQL tsrange field
type TsRange Field

func NewTsRange(table, column string, opts ...Option) TsRange {
	return TsRange{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f TsRange) Eq(v types.TsRange) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f TsRange) Neq(v types.TsRange) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }
func (f TsRange) Overlaps(v types.TsRange) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests range contains ( @> ) another range
func (f TsRange) Contains(v types.TsRange) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests range is contained by ( <@ ) another range
func (f TsRange) ContainedBy(v types.TsRange) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictLeft tests range strictly left of ( << ) another range
func (f TsRange) StrictLeft(v types.TsRange) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictRight tests range strictly right of ( >> ) another range
func (f TsRange) StrictRight(v types.TsRange) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Adjacent tests ranges are adjacent ( -|- )
func (f TsRange) Adjacent(v types.TsRange) Expr {
	return expr{e: clause.Expr{SQL: "? -|- ?", Vars: []interface{}{f.RawExpr(), v}}}
}
