package types

import (
	"context"
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type TstzRange Range[time.Time]

func (TstzRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "TSTZRANGE" }

func (r *TstzRange) Scan(value interface{}) error { return (*Range[time.Time])(r).Scan(value) }
func (r TstzRange) Value() (driver.Value, error)  { return (Range[time.Time])(r).Value() }
func (r TstzRange) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return (Range[time.Time])(r).GormValue(ctx, db)
}

// Constructors
func NewTstzRange(lower, upper time.Time, lowerInclusive, upperInclusive bool) TstzRange {
	return TstzRange(Range[time.Time]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive})
}

// Edit wrappers
func (r *TstzRange) SetBounds(lower, upper time.Time) {
	rr := (*Range[time.Time])(r)
	rr.SetBounds(lower, upper)
}

func (r *TstzRange) SetInclusivity(lowerInclusive, upperInclusive bool) {
	rr := (*Range[time.Time])(r)
	rr.SetInclusivity(lowerInclusive, upperInclusive)
}
func (r *TstzRange) SetEmpty(empty bool) { rr := (*Range[time.Time])(r); rr.SetEmpty(empty) }
