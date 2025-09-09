package types

import (
    "context"
    "database/sql/driver"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/schema"
)

type TsRange Range[time.Time]

func (TsRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "TSRANGE" }

func (r *TsRange) Scan(value interface{}) error { return (*Range[time.Time])(r).Scan(value) }
func (r TsRange) Value() (driver.Value, error) { return (Range[time.Time])(r).Value() }
func (r TsRange) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    return (Range[time.Time])(r).GormValue(ctx, db)
}

// Constructors
func NewTsRange(lower, upper time.Time, lowerInclusive, upperInclusive bool) TsRange {
	return TsRange(Range[time.Time]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive})
}

// Edit wrappers
func (r *TsRange) SetBounds(lower, upper time.Time) {
	rr := (*Range[time.Time])(r)
	rr.SetBounds(lower, upper)
}
func (r *TsRange) SetInclusivity(lowerInclusive, upperInclusive bool) {
	rr := (*Range[time.Time])(r)
	rr.SetInclusivity(lowerInclusive, upperInclusive)
}
func (r *TsRange) SetEmpty(empty bool) { rr := (*Range[time.Time])(r); rr.SetEmpty(empty) }
