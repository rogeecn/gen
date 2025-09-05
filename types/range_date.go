package types

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DateRange Range[time.Time]

func (DateRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "DATERANGE" }

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
