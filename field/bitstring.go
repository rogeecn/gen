package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// BitString represents a PostgreSQL BIT/VARBIT field
type BitString Field

func NewBitString(table, column string, opts ...Option) BitString {
	return BitString{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f BitString) Eq(v types.BitString) Expr {
	return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}}
}

func (f BitString) Neq(v types.BitString) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

// In set membership
func (f BitString) In(v ...types.BitString) Expr {
	vals := make([]interface{}, len(v))
	for i := range v {
		vals[i] = v[i]
	}
	return expr{e: clause.IN{Column: f.RawExpr(), Values: vals}}
}

// NotIn negated set membership
func (f BitString) NotIn(v ...types.BitString) Expr {
	return expr{e: clause.Not(f.In(v...).expression())}
}
