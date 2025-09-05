package types

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TSVector string

func (TSVector) GormDataType() string                          { return "tsvector" }
func (TSVector) GormDBDataType(*gorm.DB, *schema.Field) string { return "TSVECTOR" }
func (t *TSVector) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported tsvector scan type %T", value)
	}
	*t = TSVector(s)
	return nil
}
func (t TSVector) Value() (driver.Value, error) { return string(t), nil }

// Constructors
func NewTSVector(s string) TSVector { return TSVector(s) }

// Edit helpers
func (t *TSVector) Set(s string)       { *t = TSVector(s) }
func (t *TSVector) AppendRaw(s string) { *t = TSVector(string(*t) + " " + s) }
