package types

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Date time.Time

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

func (date Date) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := date.Value()
	return gorm.Expr("?", v)
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

// String
func (date Date) String() string {
	return time.Time(date).Format("2006-01-02")
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date Date) MarshalJSON() ([]byte, error) {
	return time.Time(date).MarshalJSON()
}

func (date *Date) UnmarshalJSON(b []byte) error {
	return (*time.Time)(date).UnmarshalJSON(b)
}
