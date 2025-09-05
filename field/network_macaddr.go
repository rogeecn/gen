package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// MACAddr represents a PostgreSQL MACADDR field
type MACAddr Field

func NewMACAddr(table, column string, opts ...Option) MACAddr {
	return MACAddr{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f MACAddr) Eq(v types.MACAddr) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f MACAddr) Neq(v types.MACAddr) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// In set membership
func (f MACAddr) In(v ...types.MACAddr) Expr {
	vals := make([]interface{}, len(v))
	for i := range v {
		vals[i] = v[i]
	}
	return expr{e: clause.IN{Column: f.RawExpr(), Values: vals}}
}

// NotIn negated set membership
func (f MACAddr) NotIn(v ...types.MACAddr) Expr { return expr{e: clause.Not(f.In(v...).expression())} }
