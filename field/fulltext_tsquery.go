package field

import (
	"go.ipao.vip/gen/types"
	"gorm.io/gorm/clause"
)

// TSQuery represents a PostgreSQL tsquery field
type TSQuery Field

func NewTSQuery(table, column string, opts ...Option) TSQuery {
	return TSQuery{expr: expr{col: toColumn(table, column, opts...)}}
}

func (f TSQuery) Eq(v types.TSQuery) Expr  { return expr{e: clause.Eq{Column: f.RawExpr(), Value: v}} }
func (f TSQuery) Neq(v types.TSQuery) Expr { return expr{e: clause.Neq{Column: f.RawExpr(), Value: v}} }
