package types

import (
	"math/big"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type NumRange Range[*big.Rat]

func (NumRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "NUMRANGE" }

// Constructors
func NewNumRange(lower, upper *big.Rat, lowerInclusive, upperInclusive bool) NumRange {
	return NumRange(Range[*big.Rat]{Lower: lower, Upper: upper, LowerInclusive: lowerInclusive, UpperInclusive: upperInclusive})
}

// Edit wrappers
func (r *NumRange) SetBounds(lower, upper *big.Rat) {
	rr := (*Range[*big.Rat])(r)
	rr.SetBounds(lower, upper)
}
func (r *NumRange) SetInclusivity(lowerInclusive, upperInclusive bool) {
	rr := (*Range[*big.Rat])(r)
	rr.SetInclusivity(lowerInclusive, upperInclusive)
}
func (r *NumRange) SetEmpty(empty bool) { rr := (*Range[*big.Rat])(r); rr.SetEmpty(empty) }
