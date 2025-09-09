package types

import (
	"bytes"
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

func (c CIDR) String() any {
	return (*net.IPNet)(&c).String()
}

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

// CIDR network operations

// Contains checks if the CIDR contains the given IP address
func (c CIDR) Contains(ip net.IP) bool {
	return (*net.IPNet)(&c).Contains(ip)
}

// IsValid checks if the CIDR is valid
func (c CIDR) IsValid() bool {
	n := (*net.IPNet)(&c)
	return n.IP != nil && n.Mask != nil
}

// IsIPv4 checks if the CIDR is IPv4
func (c CIDR) IsIPv4() bool {
	return (*net.IPNet)(&c).IP.To4() != nil
}

// IsIPv6 checks if the CIDR is IPv6
func (c CIDR) IsIPv6() bool {
	return (*net.IPNet)(&c).IP.To4() == nil && (*net.IPNet)(&c).IP.To16() != nil
}

// IsPrivate checks if the CIDR is in a private IP range
func (c CIDR) IsPrivate() bool {
	return (*net.IPNet)(&c).IP.IsPrivate()
}

// IsPublic checks if the CIDR is in a public IP range
func (c CIDR) IsPublic() bool {
	ip := (*net.IPNet)(&c).IP
	return !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsMulticast()
}

// IsLoopback checks if the CIDR is a loopback address
func (c CIDR) IsLoopback() bool {
	return (*net.IPNet)(&c).IP.IsLoopback()
}

// IsMulticast checks if the CIDR is a multicast address
func (c CIDR) IsMulticast() bool {
	return (*net.IPNet)(&c).IP.IsMulticast()
}

// NetworkAddress returns the network address (first IP in the range)
func (c CIDR) NetworkAddress() net.IP {
	return (*net.IPNet)(&c).IP
}

// BroadcastAddress returns the broadcast address (last IP in the range, IPv4 only)
func (c CIDR) BroadcastAddress() net.IP {
	n := (*net.IPNet)(&c)
	if n.IP.To4() == nil {
		return nil // IPv6 doesn't have broadcast
	}

	ip := make(net.IP, len(n.IP))
	copy(ip, n.IP)

	// Set host bits to 1
	for i := 0; i < len(ip); i++ {
		ip[i] |= ^n.Mask[i]
	}

	return ip
}

// Overlaps checks if this CIDR overlaps with another CIDR
func (c CIDR) Overlaps(other CIDR) bool {
	n1 := (*net.IPNet)(&c)
	n2 := (*net.IPNet)(&other)

	return n1.Contains(n2.IP) || n2.Contains(n1.IP)
}

// Clone creates a copy of the CIDR
func (c CIDR) Clone() CIDR {
	n := (*net.IPNet)(&c)
	cloned := &net.IPNet{
		IP:   make(net.IP, len(n.IP)),
		Mask: make(net.IPMask, len(n.Mask)),
	}
	copy(cloned.IP, n.IP)
	copy(cloned.Mask, n.Mask)
	return CIDR(*cloned)
}

// Equals checks if two CIDRs are equal
func (c CIDR) Equals(other CIDR) bool {
	n1 := (*net.IPNet)(&c)
	n2 := (*net.IPNet)(&other)
	return n1.IP.Equal(n2.IP) && bytes.Equal(n1.Mask, n2.Mask)
}
