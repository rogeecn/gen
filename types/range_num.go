package types

import (
	"context"
	"database/sql/driver"
	"math/big"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type NumRange Range[*big.Rat]

func (NumRange) GormDBDataType(*gorm.DB, *schema.Field) string { return "NUMRANGE" }

func (r *NumRange) Scan(value interface{}) error { return (*Range[*big.Rat])(r).Scan(value) }
func (r NumRange) Value() (driver.Value, error)  { return (Range[*big.Rat])(r).Value() }
func (r NumRange) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return (Range[*big.Rat])(r).GormValue(ctx, db)
}

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
