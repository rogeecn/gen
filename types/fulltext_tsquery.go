package types

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TSQuery string

func (TSQuery) GormDataType() string                          { return "tsquery" }
func (TSQuery) GormDBDataType(*gorm.DB, *schema.Field) string { return "TSQUERY" }
func (t *TSQuery) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported tsquery scan type %T", value)
	}
	*t = TSQuery(s)
	return nil
}
func (t TSQuery) Value() (driver.Value, error) { return string(t), nil }

// Constructors
func NewTSQuery(s string) TSQuery { return TSQuery(s) }

// Edit helpers
func (q *TSQuery) Set(s string) { *q = TSQuery(s) }
