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

// Money operations

// IsValid checks if the money value is a valid money format
func (m Money) IsValid() bool {
	s := string(m)
	if s == "" {
		return false
	}

	// Remove currency symbol and whitespace for validation
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "$") {
		s = strings.TrimPrefix(s, "$")
	}
	if strings.HasPrefix(s, "£") {
		s = strings.TrimPrefix(s, "£")
	}
	if strings.HasPrefix(s, "€") {
		s = strings.TrimPrefix(s, "€")
	}

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "") // Remove thousands separators

	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// ToFloat64 converts the money value to float64
func (m Money) ToFloat64() (float64, error) {
	s := string(m)
	if s == "" {
		return 0, fmt.Errorf("empty money value")
	}

	// Remove currency symbols and formatting
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "$") {
		s = strings.TrimPrefix(s, "$")
	}
	if strings.HasPrefix(s, "£") {
		s = strings.TrimPrefix(s, "£")
	}
	if strings.HasPrefix(s, "€") {
		s = strings.TrimPrefix(s, "€")
	}

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "") // Remove thousands separators

	return strconv.ParseFloat(s, 64)
}

// IsZero checks if the money value is zero
func (m Money) IsZero() bool {
	val, err := m.ToFloat64()
	if err != nil {
		return false
	}
	return val == 0.0
}

// IsPositive checks if the money value is positive
func (m Money) IsPositive() bool {
	val, err := m.ToFloat64()
	if err != nil {
		return false
	}
	return val > 0.0
}

// IsNegative checks if the money value is negative
func (m Money) IsNegative() bool {
	val, err := m.ToFloat64()
	if err != nil {
		return false
	}
	return val < 0.0
}

// Abs returns the absolute value of the money
func (m Money) Abs() Money {
	val, err := m.ToFloat64()
	if err != nil {
		return m
	}

	if val < 0 {
		// Extract the format and make it positive
		s := string(m)
		if strings.Contains(s, "-") {
			s = strings.ReplaceAll(s, "-", "")
		}
		return Money(s)
	}

	return m
}

// Compare compares two money values (-1 if less, 0 if equal, 1 if greater)
func (m Money) Compare(other Money) (int, error) {
	val1, err := m.ToFloat64()
	if err != nil {
		return 0, err
	}

	val2, err := other.ToFloat64()
	if err != nil {
		return 0, err
	}

	if val1 < val2 {
		return -1, nil
	} else if val1 > val2 {
		return 1, nil
	}
	return 0, nil
}

// Add adds another money value to this one
func (m Money) Add(other Money) (Money, error) {
	val1, err := m.ToFloat64()
	if err != nil {
		return "", err
	}

	val2, err := other.ToFloat64()
	if err != nil {
		return "", err
	}

	result := val1 + val2
	return Money(fmt.Sprintf("%.2f", result)), nil
}

// Subtract subtracts another money value from this one
func (m Money) Subtract(other Money) (Money, error) {
	val1, err := m.ToFloat64()
	if err != nil {
		return "", err
	}

	val2, err := other.ToFloat64()
	if err != nil {
		return "", err
	}

	result := val1 - val2
	return Money(fmt.Sprintf("%.2f", result)), nil
}

// Multiply multiplies the money value by a factor
func (m Money) Multiply(factor float64) (Money, error) {
	val, err := m.ToFloat64()
	if err != nil {
		return "", err
	}

	result := val * factor
	return Money(fmt.Sprintf("%.2f", result)), nil
}

// Clone creates a copy of the money value
func (m Money) Clone() Money {
	return Money(string(m))
}

// Equals checks if two money values are equal
func (m Money) Equals(other Money) bool {
	cmp, err := m.Compare(other)
	if err != nil {
		return false
	}
	return cmp == 0
}
