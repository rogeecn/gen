package types

import (
    "context"
    "database/sql/driver"
    "errors"
    "fmt"
    "strconv"
    "strings"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/schema"
)

type Point struct{ X, Y float64 }

func (Point) GormDataType() string                          { return "point" }
func (Point) GormDBDataType(*gorm.DB, *schema.Field) string { return "POINT" }
func (p *Point) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported point scan type %T", value)
	}
	// format: (x,y)
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return errors.New("invalid point")
	}
	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return err
	}
	*p = Point{X: x, Y: y}
	return nil
}
func (p Point) Value() (driver.Value, error) { return fmt.Sprintf("(%g,%g)", p.X, p.Y), nil }

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    v, _ := p.Value(); return gorm.Expr("?", v)
}

// Constructors
func NewPoint(x, y float64) Point { return Point{X: x, Y: y} }

// Edit helpers
func (p *Point) Set(x, y float64)         { p.X, p.Y = x, y }
func (p *Point) Translate(dx, dy float64) { p.X += dx; p.Y += dy }
