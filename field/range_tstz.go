package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// TstzRange represents a PostgreSQL tstzrange field
type TstzRange Field

func NewTstzRange(table, column string, opts ...Option) TstzRange {
	return TstzRange{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f TstzRange) Eq(v types.TstzRange) Expr {
	return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}}
}

func (f TstzRange) Neq(v types.TstzRange) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

func (f TstzRange) Overlaps(v types.TstzRange) Expr {
	return expr{e: clause.Expr{SQL: "? && ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Contains tests range contains ( @> ) another range
func (f TstzRange) Contains(v types.TstzRange) Expr {
	return expr{e: clause.Expr{SQL: "? @> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests range is contained by ( <@ ) another range
func (f TstzRange) ContainedBy(v types.TstzRange) Expr {
	return expr{e: clause.Expr{SQL: "? <@ ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictLeft tests range strictly left of ( << ) another range
func (f TstzRange) StrictLeft(v types.TstzRange) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// StrictRight tests range strictly right of ( >> ) another range
func (f TstzRange) StrictRight(v types.TstzRange) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// Adjacent tests ranges are adjacent ( -|- )
func (f TstzRange) Adjacent(v types.TstzRange) Expr {
	return expr{e: clause.Expr{SQL: "? -|- ?", Vars: []interface{}{f.RawExpr(), v}}}
}
