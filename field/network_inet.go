package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Inet represents a PostgreSQL INET field
type Inet Field

func NewInet(table, column string, opts ...Option) Inet {
	return Inet{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Inet) Eq(v types.Inet) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Inet) Neq(v types.Inet) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }
func (f Inet) In(v ...types.Inet) Expr {
	vals := make([]interface{}, len(v))
	for i := range v {
		vals[i] = v[i]
	}
	return expr{e: clause.IN{Column: f.RawExpr(), Values: vals}}
}
func (f Inet) NotIn(v ...types.Inet) Expr { return expr{e: clause.Not(f.In(v...).expression())} }

// Contains tests network contains (>>) another network
func (f Inet) Contains(v types.Inet) Expr {
	return expr{e: clause.Expr{SQL: "? >> ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainsEq tests network contains or equals (>>=)
func (f Inet) ContainsEq(v types.Inet) Expr {
	return expr{e: clause.Expr{SQL: "? >>= ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedBy tests network is contained by (<<)
func (f Inet) ContainedBy(v types.Inet) Expr {
	return expr{e: clause.Expr{SQL: "? << ?", Vars: []interface{}{f.RawExpr(), v}}}
}

// ContainedByEq tests network is contained by or equals (<<=)
func (f Inet) ContainedByEq(v types.Inet) Expr {
	return expr{e: clause.Expr{SQL: "? <<= ?", Vars: []interface{}{f.RawExpr(), v}}}
}
