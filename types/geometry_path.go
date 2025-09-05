package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Path struct {
	Closed bool
	Points []Point
}

func (Path) GormDataType() string                          { return "path" }
func (Path) GormDBDataType(*gorm.DB, *schema.Field) string { return "PATH" }
func (p *Path) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported path scan type %T", value)
	}
	s = strings.TrimSpace(s)
	closed := strings.HasPrefix(s, "(")
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, ")")
	s = strings.TrimSuffix(s, "]")
	pts := strings.Split(s, "),(")
	res := make([]Point, 0, len(pts))
	for _, ps := range pts {
		var pt Point
		if err := pt.Scan(strings.TrimSuffix(strings.TrimPrefix(ps, "("), ")")); err != nil {
			return err
		}
		res = append(res, pt)
	}
	*p = Path{Closed: closed, Points: res}
	return nil
}

func (p Path) Value() (driver.Value, error) {
	var buf bytes.Buffer
	if p.Closed {
		buf.WriteByte('(')
	} else {
		buf.WriteByte('[')
	}
	for i, pt := range p.Points {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("(%g,%g)", pt.X, pt.Y))
	}
	if p.Closed {
		buf.WriteByte(')')
	} else {
		buf.WriteByte(']')
	}
	return buf.String(), nil
}

// Constructors
func NewPath(points []Point, closed bool) Path { return Path{Closed: closed, Points: points} }

// Edit helpers
func (p *Path) Set(points []Point, closed bool) { p.Points, p.Closed = points, closed }
func (p *Path) Append(points ...Point)          { p.Points = append(p.Points, points...) }
func (p *Path) Close()                          { p.Closed = true }
func (p *Path) Open()                           { p.Closed = false }
