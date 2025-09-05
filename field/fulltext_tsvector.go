package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// TSVector represents a PostgreSQL tsvector field
type TSVector Field

func NewTSVector(table, column string, opts ...Option) TSVector {
	return TSVector{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f TSVector) Eq(v types.TSVector) Expr { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f TSVector) Neq(v types.TSVector) Expr {
	return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}}
}

// Matches to_tsquery
func (f TSVector) Matches(q types.TSQuery) Expr {
	return expr{e: clause.Expr{SQL: "? @@ ?", Vars: []interface{}{f.RawExpr(), q}}}
}
