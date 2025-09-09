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

// CIDR represents PostgreSQL CIDR type
type CIDR net.IPNet

func (CIDR) GormDataType() string                          { return "cidr" }
func (CIDR) GormDBDataType(*gorm.DB, *schema.Field) string { return "CIDR" }
func (c *CIDR) Scan(value interface{}) error {
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("unsupported cidr scan type %T", value)
	}
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}
	ipnet.IP = ip
	*c = CIDR(*ipnet)
	return nil
}
func (c CIDR) Value() (driver.Value, error) { return (*net.IPNet)(&c).String(), nil }

func (c CIDR) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
    v, _ := c.Value()
    return gorm.Expr("?", v)
}

// Constructors
func NewCIDR(s string) (CIDR, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return CIDR{}, err
	}
	ipnet.IP = ip
	return CIDR(*ipnet), nil
}

func MustCIDR(s string) CIDR {
	v, err := NewCIDR(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Edit applies a mutator to the underlying IPNet
func (c *CIDR) Edit(mutator func(n *net.IPNet)) {
	if mutator == nil {
		return
	}
	n := (*net.IPNet)(c)
	mutator(n)
	*c = CIDR(*n)
}
