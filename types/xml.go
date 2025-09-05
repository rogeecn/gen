package types

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type XML string

func (XML) GormDataType() string                          { return "xml" }
func (XML) GormDBDataType(*gorm.DB, *schema.Field) string { return "XML" }
func (x *XML) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported xml scan type %T", value)
	}
	*x = XML(s)
	return nil
}
func (x XML) Value() (driver.Value, error) { return string(x), nil }

// Constructors
func NewXML(s string) XML { return XML(s) }

// Edit helpers
func (x *XML) Set(s string) { *x = XML(s) }
