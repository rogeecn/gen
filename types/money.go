package types

import (
    "context"
    "database/sql/driver"
    "fmt"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "gorm.io/gorm/schema"
)

// Money stored as string to preserve formatting/precision.
type Money string

func (Money) GormDataType() string                          { return "money" }
func (Money) GormDBDataType(*gorm.DB, *schema.Field) string { return "MONEY" }
func (m *Money) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported money scan type %T", value)
	}
	*m = Money(s)
	return nil
}
func (m Money) Value() (driver.Value, error) { return string(m), nil }

// GormValuer
func (m Money) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    v, _ := m.Value()
    return gorm.Expr("?", v)
}

// Constructors
func NewMoney(s string) Money { return Money(s) }

// Edit helpers
func (m *Money) Set(s string) { *m = Money(s) }
