package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// Money represents a PostgreSQL money field
type Money Field

func NewMoney(table, column string, opts ...Option) Money {
	return Money{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f Money) Eq(v types.Money) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f Money) Neq(v types.Money) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }

// Gt greater than
func (f Money) Gt(v types.Money) Expr { return expr{e: clause.Gt{Column: f.RawExpr(), Value: v}} }

// Gte greater or equal to
func (f Money) Gte(v types.Money) Expr { return expr{e: clause.Gte{Column: f.RawExpr(), Value: v}} }

// Lt less than
func (f Money) Lt(v types.Money) Expr { return expr{e: clause.Lt{Column: f.RawExpr(), Value: v}} }

// Lte less or equal to
func (f Money) Lte(v types.Money) Expr { return expr{e: clause.Lte{Column: f.RawExpr(), Value: v}} }

// In set membership
func (f Money) In(values ...types.Money) Expr {
	return expr{e: clause.IN{Column: f.RawExpr(), Values: f.toSlice(values...)}}
}

// NotIn negated set membership
func (f Money) NotIn(values ...types.Money) Expr {
	return expr{e: clause.Not(f.In(values...).expression())}
}

// Between inclusive range
func (f Money) Between(left, right types.Money) Expr {
	return f.between([]interface{}{left, right})
}

// NotBetween negated range
func (f Money) NotBetween(left, right types.Money) Expr {
	return Not(f.Between(left, right))
}

// Like textual LIKE on money literal representation
func (f Money) Like(v types.Money) Expr { return expr{e: clause.Like{Column: f.RawExpr(), Value: v}} }

// NotLike negated LIKE on textual form
func (f Money) NotLike(v types.Money) Expr { return expr{e: clause.Not(f.Like(v).expression())} }

func (f Money) toSlice(values ...types.Money) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
