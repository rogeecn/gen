package types

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TstzRange Range[time.Time]

func (TstzRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "TSTZRANGE" }

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
