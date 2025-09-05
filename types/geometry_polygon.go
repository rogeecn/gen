package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Polygon struct{ Points []Point }

func (Polygon) GormDataType() string                          { return "polygon" }
func (Polygon) GormDBDataType(*gorm.DB, *schema.Field) string { return "POLYGON" }
func (p *Polygon) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported polygon scan type %T", value)
	}
	s = strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(s), ")"), "(")
	parts := strings.Split(s, "),(")
	res := make([]Point, 0, len(parts))
	for _, ps := range parts {
		var pt Point
		if err := pt.Scan(strings.TrimSuffix(strings.TrimPrefix(ps, "("), ")")); err != nil {
			return err
		}
		res = append(res, pt)
	}
	*p = Polygon{Points: res}
	return nil
}

func (p Polygon) Value() (driver.Value, error) {
	var buf bytes.Buffer
	buf.WriteByte('(')
	for i, pt := range p.Points {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("(%g,%g)", pt.X, pt.Y))
	}
	buf.WriteByte(')')
	return buf.String(), nil
}

// Constructors
func NewPolygon(points []Point) Polygon { return Polygon{Points: points} }

// Edit helpers
func (p *Polygon) Set(points []Point)     { p.Points = points }
func (p *Polygon) Append(points ...Point) { p.Points = append(p.Points, points...) }
