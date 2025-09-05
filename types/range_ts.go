package types

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TsRange Range[time.Time]

func (TsRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "TSRANGE" }

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
