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

// Date calculation methods

// IsValid checks if the date is valid
func (date Date) IsValid() bool {
	return !time.Time(date).IsZero()
}

// IsZero checks if the date is zero
func (date Date) IsZero() bool {
	return time.Time(date).IsZero()
}

// AddDays adds the specified number of days
func (date Date) AddDays(days int) Date {
	return Date(time.Time(date).AddDate(0, 0, days))
}

// AddMonths adds the specified number of months
func (date Date) AddMonths(months int) Date {
	return Date(time.Time(date).AddDate(0, months, 0))
}

// AddYears adds the specified number of years
func (date Date) AddYears(years int) Date {
	return Date(time.Time(date).AddDate(years, 0, 0))
}

// DayOfWeek returns the day of the week
func (date Date) DayOfWeek() time.Weekday {
	return time.Time(date).Weekday()
}

// IsWeekend checks if the date falls on a weekend
func (date Date) IsWeekend() bool {
	weekday := date.DayOfWeek()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsBusinessDay checks if the date falls on a business day (Monday-Friday)
func (date Date) IsBusinessDay() bool {
	return !date.IsWeekend()
}

// Quarter returns the quarter of the year (1-4)
func (date Date) Quarter() int {
	month := time.Time(date).Month()
	return int((month-1)/3) + 1
}

// DayOfYear returns the day of the year (1-366)
func (date Date) DayOfYear() int {
	return time.Time(date).YearDay()
}

// WeekOfYear returns the ISO week number and year
func (date Date) WeekOfYear() (year, week int) {
	return time.Time(date).ISOWeek()
}

// StartOfWeek returns the date of the start of the week (Monday)
func (date Date) StartOfWeek() Date {
	t := time.Time(date)
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	return Date(t.AddDate(0, 0, -(weekday - 1)))
}

// EndOfWeek returns the date of the end of the week (Sunday)
func (date Date) EndOfWeek() Date {
	return date.StartOfWeek().AddDays(6)
}

// StartOfMonth returns the first day of the month
func (date Date) StartOfMonth() Date {
	t := time.Time(date)
	return Date(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()))
}

// EndOfMonth returns the last day of the month
func (date Date) EndOfMonth() Date {
	t := time.Time(date)
	return Date(time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1))
}

// StartOfYear returns the first day of the year
func (date Date) StartOfYear() Date {
	t := time.Time(date)
	return Date(time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()))
}

// EndOfYear returns the last day of the year
func (date Date) EndOfYear() Date {
	t := time.Time(date)
	return Date(time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location()))
}

// DaysUntil returns the number of days until the specified date
func (date Date) DaysUntil(other Date) int {
	d1 := time.Time(date)
	d2 := time.Time(other)
	return int(d2.Sub(d1).Hours() / 24)
}

// DaysSince returns the number of days since the specified date
func (date Date) DaysSince(other Date) int {
	d1 := time.Time(date)
	d2 := time.Time(other)
	return int(d1.Sub(d2).Hours() / 24)
}

// Before reports whether the date is before other
func (date Date) Before(other Date) bool {
	return time.Time(date).Before(time.Time(other))
}

// After reports whether the date is after other
func (date Date) After(other Date) bool {
	return time.Time(date).After(time.Time(other))
}

// Equal reports whether the date is equal to other
func (date Date) Equal(other Date) bool {
	return time.Time(date).Equal(time.Time(other))
}

// Clone creates a copy of the date
func (date Date) Clone() Date {
	return Date(time.Time(date))
}
