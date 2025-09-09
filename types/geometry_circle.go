package types

import (
    "context"
    "database/sql/driver"
    "fmt"
    "strconv"
    "strings"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/schema"
)

type Circle struct {
	Center Point
	Radius float64
}

func (Circle) GormDataType() string                          { return "circle" }
func (Circle) GormDBDataType(*gorm.DB, *schema.Field) string { return "CIRCLE" }
func (c *Circle) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported circle scan type %T", value)
	}
	// format: <(x,y),r> or (x,y),r
	s = strings.Trim(s, "<>")
	parts := strings.SplitN(s, "),", 2)
	var pt Point
	if err := pt.Scan(parts[0]); err != nil {
		return err
	}
	r, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return err
	}
	*c = Circle{Center: pt, Radius: r}
	return nil
}
func (c Circle) Value() (driver.Value, error) {
    return fmt.Sprintf("<(%g,%g),%g>", c.Center.X, c.Center.Y, c.Radius), nil
}

func (c Circle) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    v, _ := c.Value(); return gorm.Expr("?", v)
}

// Constructors
func NewCircle(center Point, radius float64) Circle { return Circle{Center: center, Radius: radius} }

// Edit helpers
func (c *Circle) Set(center Point, radius float64) { c.Center, c.Radius = center, radius }
func (c *Circle) SetCenter(center Point)           { c.Center = center }
func (c *Circle) SetRadius(radius float64)         { c.Radius = radius }
