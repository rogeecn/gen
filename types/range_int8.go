package types

import (
	"context"
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Int8Range Range[int64]

func (Int8Range) GormDBDataType(*gorm.DB, *schema.Field) string { return "INT8RANGE" }

func (r *Int8Range) Scan(value interface{}) error { return (*Range[int64])(r).Scan(value) }
func (r Int8Range) Value() (driver.Value, error)  { return (Range[int64])(r).Value() }
func (r Int8Range) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return (Range[int64])(r).GormValue(ctx, db)
}

// Constructors
func NewInt8Range(lower, upper int64, lowerInclusive, upperInclusive bool) Int8Range {
	return Int8Range(Range[int64]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive})
}

// Edit wrappers
func (r *Int8Range) SetBounds(lower, upper int64) {
	rr := (*Range[int64])(r)
	rr.SetBounds(lower, upper)
}

func (r *Int8Range) SetInclusivity(lowerInclusive, upperInclusive bool) {
	rr := (*Range[int64])(r)
	rr.SetInclusivity(lowerInclusive, upperInclusive)
}
func (r *Int8Range) SetEmpty(empty bool) { rr := (*Range[int64])(r); rr.SetEmpty(empty) }
