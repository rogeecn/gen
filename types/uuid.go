package types

import (
	"context"
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// This datatype stores the uuid in the database as a string. To store the uuid
// in the database as a binary (byte) array, please refer to types.BinUUID.
type UUID uuid.UUID

// NewUUIDv1 generates a UUID version 1, panics on generation failure.
func NewUUIDv1() UUID {
	return UUID(uuid.Must(uuid.NewUUID()))
}

// NewUUIDv4 generates a UUID version 4, panics on generation failure.
func NewUUIDv4() UUID {
	return UUID(uuid.Must(uuid.NewRandom()))
}

// GormDataType gorm common data type.
func (UUID) GormDataType() string {
	return "string"
}

// GormDBDataType gorm db data type.
func (UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// Postgres only
	return "UUID"
}

// Scan is the scanner function for this datatype.
func (u *UUID) Scan(value interface{}) error {
	var result uuid.UUID
	if err := result.Scan(value); err != nil {
		return err
	}
	*u = UUID(result)
	return nil
}

// Value is the valuer function for this datatype.
func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).Value()
}

func (u UUID) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := u.Value()
	return gorm.Expr("?", v)
}

// String returns the string form of the UUID.
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// Equals returns true if string form of UUID matches other, false otherwise.
func (u UUID) Equals(other UUID) bool {
	return u.String() == other.String()
}

// Length returns the number of characters in string form of UUID.
func (u UUID) Length() int {
	return len(u.String())
}

// IsNil returns true if the UUID is a nil UUID (all zeroes), false otherwise.
func (u UUID) IsNil() bool {
	return uuid.UUID(u) == uuid.Nil
}

// IsEmpty returns true if UUID is nil UUID or of zero length, false otherwise.
func (u UUID) IsEmpty() bool {
	return u.IsNil() || u.Length() == 0
}

// IsNilPtr returns true if caller UUID ptr is nil, false otherwise.
func (u *UUID) IsNilPtr() bool {
	return u == nil
}

// IsEmptyPtr returns true if caller UUID ptr is nil or it's value is empty.
func (u *UUID) IsEmptyPtr() bool {
	return u.IsNilPtr() || u.IsEmpty()
}

// UUID utility methods

// Version returns the UUID version
func (u UUID) Version() int {
	return int(uuid.UUID(u).Version())
}

// Variant returns the UUID variant
func (u UUID) Variant() uuid.Variant {
	return uuid.UUID(u).Variant()
}

// IsV1 checks if the UUID is version 1
func (u UUID) IsV1() bool {
	return u.Version() == 1
}

// IsV4 checks if the UUID is version 4
func (u UUID) IsV4() bool {
	return u.Version() == 4
}

// Compare compares two UUIDs (-1 if less, 0 if equal, 1 if greater)
func (u UUID) Compare(other UUID) int {
	s1, s2 := u.String(), other.String()
	if s1 < s2 {
		return -1
	} else if s1 > s2 {
		return 1
	}
	return 0
}

// Clone creates a copy of the UUID
func (u UUID) Clone() UUID {
	return UUID(uuid.UUID(u))
}

// MarshalBinary implements encoding.BinaryMarshaler
func (u UUID) MarshalBinary() ([]byte, error) {
	return uuid.UUID(u).MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (u *UUID) UnmarshalBinary(data []byte) error {
	var parsed uuid.UUID
	if err := parsed.UnmarshalBinary(data); err != nil {
		return err
	}
	*u = UUID(parsed)
	return nil
}

// IsValid checks if the UUID is valid
func (u UUID) IsValid() bool {
	return uuid.UUID(u) != uuid.Nil && len(u.String()) == 36
}
