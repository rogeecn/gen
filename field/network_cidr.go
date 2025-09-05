package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// CIDR represents a PostgreSQL CIDR field
type CIDR Field

func NewCIDR(table, column string, opts ...Option) CIDR {
	return CIDR{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f CIDR) Eq(v types.CIDR) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f CIDR) Neq(v types.CIDR) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }
func (f CIDR) In(v ...types.CIDR) Expr {
	vals := make([]interface{}, len(v))
	for i := range v {
		vals[i] = v[i]
	}
	return expr{e: clause.IN{Column: f.RawExpr(), Values: vals}}
}
func (f CIDR) NotIn(v ...types.CIDR) Expr { return expr{e: clause.Not(f.In(v...).expression())} }

// Contains tests network contains (>>) another network
func (f CIDR) Contains(v types.CIDR) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsEq tests network contains or equals (>>=)
func (f CIDR) ContainsEq(v types.CIDR) Expr {
	return expr{e: clause.Expr{SQL: "? >>= ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests network is contained by (<<)
func (f CIDR) ContainedBy(v types.CIDR) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedByEq tests network is contained by or equals (<<=)
func (f CIDR) ContainedByEq(v types.CIDR) Expr {
	return expr{e: clause.Expr{SQL: "? <<= ?", Vars: []interface{}{f.RawExpr(), v}}}
}
