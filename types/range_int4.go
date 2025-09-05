package types

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Int4Range Range[int32]

func (Int4Range) GormDBDataType(*gorm.DB, *schema.Field) string { return "INT4RANGE" }

// Constructors
func NewInt4Range(lower, upper int32, lowerInclusive, upperInclusive bool) Int4Range {
	return Int4Range(Range[int32]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive})
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
