package types

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// MACAddr represents PostgreSQL MACADDR type
type MACAddr net.HardwareAddr

func (m MACAddr) String() any {
	return net.HardwareAddr(m).String()
}

func (MACAddr) GormDataType() string                          { return "macaddr" }
func (MACAddr) GormDBDataType(*gorm.DB, *schema.Field) string { return "MACADDR" }
func (m *MACAddr) Scan(value interface{}) error {
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("unsupported macaddr scan type %T", value)
	}
	hw, err := net.ParseMAC(s)
	if err != nil {
		return err
	}
	*m = MACAddr(hw)
	return nil
}
func (m MACAddr) Value() (driver.Value, error) { return net.HardwareAddr(m).String(), nil }

func (m MACAddr) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := m.Value()
	return gorm.Expr("?", v)
}

// Constructors
func NewMACAddr(s string) (MACAddr, error) {
	hw, err := net.ParseMAC(s)
	if err != nil {
		return MACAddr(nil), err
	}
	return MACAddr(hw), nil
}

func MustMACAddr(s string) MACAddr {
	v, err := NewMACAddr(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Edit applies a mutator to the underlying MAC address bytes
func (m *MACAddr) Edit(mutator func(b []byte) []byte) {
	if mutator == nil {
		return
	}
	*m = MACAddr(mutator([]byte(net.HardwareAddr(*m))))
}
