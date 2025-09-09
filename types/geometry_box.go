package types

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Box struct{ P1, P2 Point }

func (Box) GormDataType() string                          { return "box" }
func (Box) GormDBDataType(*gorm.DB, *schema.Field) string { return "BOX" }
func (b *Box) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported box scan type %T", value)
	}
	// format: (x1,y1),(x2,y2)
	parts := strings.SplitN(strings.TrimSpace(s), "),(", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid box")
	}
	var p1, p2 Point
	if err := p1.Scan(strings.TrimSuffix(strings.TrimPrefix(parts[0], "("), ")")); err != nil {
		return err
	}
	if err := p2.Scan(strings.TrimSuffix(strings.TrimPrefix(parts[1], "("), ")")); err != nil {
		return err
	}
	*b = Box{P1: p1, P2: p2}
	return nil
}

func (b Box) Value() (driver.Value, error) {
	return fmt.Sprintf("(%g,%g),(%g,%g)", b.P1.X, b.P1.Y, b.P2.X, b.P2.Y), nil
}

func (b Box) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := b.Value()
	return gorm.Expr("?", v)
}

// Constructors
func NewBox(p1, p2 Point) Box { return Box{P1: p1, P2: p2} }

// Edit helpers
func (b *Box) Set(p1, p2 Point) { b.P1, b.P2 = p1, p2 }

func (b Box) P1Point() Point { return b.P1 }
func (b Box) P2Point() Point { return b.P2 }
