package types

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Inet represents PostgreSQL INET type
type Inet net.IP

func (i Inet) String() any {
	return net.IP(i).String()
}

func (Inet) GormDataType() string                          { return "inet" }
func (Inet) GormDBDataType(*gorm.DB, *schema.Field) string { return "INET" }
func (i *Inet) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		ip := net.ParseIP(string(v))
		if ip == nil {
			return errors.New("invalid inet")
		}
		*i = Inet(ip)
		return nil
	case string:
		ip := net.ParseIP(v)
		if ip == nil {
			return errors.New("invalid inet")
		}
		*i = Inet(ip)
		return nil
	case net.IP:
		*i = Inet(v)
		return nil
	default:
		return fmt.Errorf("unsupported inet scan type %T", value)
	}
}
func (i Inet) Value() (driver.Value, error) { return net.IP(i).String(), nil }

func (i Inet) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := i.Value()
	return gorm.Expr("?", v)
}

// Constructors
func NewInet(s string) (Inet, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return Inet(nil), fmt.Errorf("invalid inet: %q", s)
	}
	return Inet(ip), nil
}

func MustInet(s string) Inet {
	v, err := NewInet(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Edit applies a mutator to the underlying IP
func (i *Inet) Edit(mutator func(ip net.IP) net.IP) {
	if mutator == nil {
		return
	}
	*i = Inet(mutator(net.IP(*i)))
}

// INET network operations

// IsValid checks if the INET is a valid IP address
func (i Inet) IsValid() bool {
	ip := net.IP(i)
	return len(ip) == net.IPv4len || len(ip) == net.IPv6len
}

// IsIPv4 checks if the INET is IPv4
func (i Inet) IsIPv4() bool {
	return net.IP(i).To4() != nil
}

// IsIPv6 checks if the INET is IPv6
func (i Inet) IsIPv6() bool {
	return net.IP(i).To4() == nil && net.IP(i).To16() != nil
}

// IsLoopback checks if the INET is a loopback address
func (i Inet) IsLoopback() bool {
	return net.IP(i).IsLoopback()
}

// IsPrivate checks if the INET is in a private IP range
func (i Inet) IsPrivate() bool {
	return net.IP(i).IsPrivate()
}

// IsMulticast checks if the INET is a multicast address
func (i Inet) IsMulticast() bool {
	return net.IP(i).IsMulticast()
}

// IsUnspecified checks if the INET is an unspecified address
func (i Inet) IsUnspecified() bool {
	return net.IP(i).IsUnspecified()
}

// IsLinkLocalUnicast checks if the INET is a link-local unicast address
func (i Inet) IsLinkLocalUnicast() bool {
	return net.IP(i).IsLinkLocalUnicast()
}

// IsGlobalUnicast checks if the INET is a global unicast address
func (i Inet) IsGlobalUnicast() bool {
	return net.IP(i).IsGlobalUnicast()
}

// IsPublic checks if the INET is in a public IP range
func (i Inet) IsPublic() bool {
	ip := net.IP(i)
	return !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsMulticast()
}

// Clone creates a copy of the INET
func (i Inet) Clone() Inet {
	clone := make(net.IP, len(net.IP(i)))
	copy(clone, net.IP(i))
	return Inet(clone)
}

// Equals checks if two INETs are equal
func (i Inet) Equals(other Inet) bool {
	return net.IP(i).Equal(net.IP(other))
}
