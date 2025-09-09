package types

import (
	"context"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// toString converts driver input to string
func toString(v interface{}) (string, bool) {
	switch x := v.(type) {
	case []byte:
		return string(x), true
	case string:
		return x, true
	default:
		return "", false
	}
}

// HexBytes is a helper datatype for bytea hex representations
type HexBytes []byte

func (HexBytes) GormDataType() string                          { return "bytea" }
func (HexBytes) GormDBDataType(*gorm.DB, *schema.Field) string { return "BYTEA" }
func (h *HexBytes) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported bytea scan type %T", value)
	}
	s = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(s)), "\\x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*h = HexBytes(b)
	return nil
}
func (h HexBytes) Value() (driver.Value, error) { return []byte(h), nil }

func (h HexBytes) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := h.Value()
	return gorm.Expr("?", v)
}

// SetHex replaces content from a hex string (with or without leading \\x)
func (h *HexBytes) SetHex(hexStr string) error {
	hexStr = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(hexStr)), "\\x")
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	*h = HexBytes(b)
	return nil
}

// AppendBytes appends raw bytes to the value
func (h *HexBytes) AppendBytes(b []byte) { *h = append(*h, b...) }
