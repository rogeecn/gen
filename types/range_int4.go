package types

import (
    "context"
    "database/sql"
    "database/sql/driver"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/schema"
)

var (
	_ gorm.Valuer   = (*Int4Range)(nil)
	_ driver.Valuer = (*Int4Range)(nil)
	_ sql.Scanner   = (*Int4Range)(nil)
)

type Int4Range Range[int32]

func (Int4Range) GormDBDataType(*gorm.DB, *schema.Field) string { return "INT4RANGE" }

// Interface implementations via delegation to Range[int32]
func (r *Int4Range) Scan(value interface{}) error { return (*Range[int32])(r).Scan(value) }
func (r Int4Range) Value() (driver.Value, error) { return (Range[int32])(r).Value() }
func (r Int4Range) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    return (Range[int32])(r).GormValue(ctx, db)
}

// Constructors
func NewInt4Range(lower, upper int32, lowerInclusive, upperInclusive bool) Int4Range {
	return Int4Range(
		Range[int32]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive},
	)
}

// Edit wrappers
func (r *Int4Range) SetBounds(lower, upper int32) {
	rr := (*Range[int32])(r)
	rr.SetBounds(lower, upper)
}

func (r *Int4Range) SetInclusivity(lowerInclusive, upperInclusive bool) {
	rr := (*Range[int32])(r)
	rr.SetInclusivity(lowerInclusive, upperInclusive)
}
func (r *Int4Range) SetEmpty(empty bool) { rr := (*Range[int32])(r); rr.SetEmpty(empty) }
