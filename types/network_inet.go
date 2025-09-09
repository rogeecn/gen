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
