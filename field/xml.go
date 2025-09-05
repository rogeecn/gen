package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// XML represents a PostgreSQL XML field
type XML Field

func NewXML(table, column string, opts ...Option) XML {
	return XML{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f XML) Eq(v types.XML) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f XML) Neq(v types.XML) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// In set membership
func (f XML) In(v ...types.XML) Expr {
	vals := make([]interface{}, len(v))
	for i := range v {
		vals[i] = v[i]
	}
	return expr{e: clause.IN{Column: f.RawExpr(), Values: vals}}
}
func (f XML) NotIn(v ...types.XML) Expr { return expr{e: clause.Not(f.In(v...).expression())} }

// Like textual LIKE on XML string representation
func (f XML) Like(pattern string) Expr {
	return expr{e: clause.Like{Column: f.RawExpr(), Value: pattern}}
}
func (f XML) NotLike(pattern string) Expr { return expr{e: clause.Not(f.Like(pattern).expression())} }

// Regexp regex match on XML textual form
func (f XML) Regexp(pattern string) Expr {
	return expr{e: clause.Expr{SQL: "? REGEXP ?", Vars: []interface{}{f.RawExpr(), pattern}}}
}

func (f XML) NotRegexp(pattern string) Expr {
	return expr{e: clause.Not(f.Regexp(pattern).expression())}
}
