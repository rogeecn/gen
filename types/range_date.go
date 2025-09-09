package types

import (
    "context"
    "database/sql/driver"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/schema"
)

type DateRange Range[time.Time]

func (DateRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "DATERANGE" }

func (r *DateRange) Scan(value interface{}) error { return (*Range[time.Time])(r).Scan(value) }
func (r DateRange) Value() (driver.Value, error) { return (Range[time.Time])(r).Value() }
func (r DateRange) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    return (Range[time.Time])(r).GormValue(ctx, db)
}

// Constructors
func NewDateRange(lower, upper time.Time, lowerInclusive, upperInclusive bool) DateRange {
	return DateRange(
		Range[time.Time]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive},
	)
}

// Edit wrappers
func (r *DateRange) SetBounds(lower, upper time.Time) {
	rr := (*Range[time.Time])(r)
	rr.SetBounds(lower, upper)
}

func (r *DateRange) SetInclusivity(lowerInclusive, upperInclusive bool) {
	rr := (*Range[time.Time])(r)
	rr.SetInclusivity(lowerInclusive, upperInclusive)
}
func (r *DateRange) SetEmpty(empty bool) { rr := (*Range[time.Time])(r); rr.SetEmpty(empty) }
