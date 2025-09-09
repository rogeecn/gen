package types

import (
	"context"
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// BitString represents PostgreSQL BIT/VARBIT
type BitString string

func (BitString) GormDataType() string                          { return "bit" }
func (BitString) GormDBDataType(*gorm.DB, *schema.Field) string { return "BIT VARYING" }
func (b *BitString) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*b = BitString(string(v))
		return nil
	case string:
		*b = BitString(v)
		return nil
	default:
		return fmt.Errorf("unsupported bit scan type %T", value)
	}
}
func (b BitString) Value() (driver.Value, error) { return string(b), nil }

func (b BitString) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := b.Value()
	return gorm.Expr("?", v)
}

// Constructor
func NewBitString(s string) BitString { return BitString(s) }

// Edit helpers
func (b *BitString) Set(s string)       { *b = BitString(s) }
func (b *BitString) Append(bits string) { *b = BitString(string(*b) + bits) }
func (b BitString) Len() int            { return len(string(b)) }
