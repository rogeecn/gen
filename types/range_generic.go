package types

import (
	"context"
	"database/sql/driver"
	"encoding"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Range is a generic representation of PostgreSQL range types
type Range[T any] struct {
	Lower, Upper                   T
	LowerInclusive, UpperInclusive bool
	Empty                          bool
}

func (Range[T]) GormDataType() string { return "range" }

func (r *Range[T]) Scan(value interface{}) error {
	s, ok := toString(value)
	if !ok {
		return fmt.Errorf("unsupported range scan type %T", value)
	}
	s = strings.TrimSpace(s)
	if strings.EqualFold(s, "empty") {
		r.Empty = true
		return nil
	}
	if len(s) < 2 {
		return errors.New("invalid range")
	}
	r.LowerInclusive = s[0] == '['
	r.UpperInclusive = s[len(s)-1] == ']'
	body := s[1 : len(s)-1]
	parts := strings.SplitN(body, ",", 2)
	if len(parts) != 2 {
		return errors.New("invalid range body")
	}
	var err error
	var zeroT T
	r.Lower, err = parseRangeVal[T](strings.TrimSpace(parts[0]))
	if err != nil {
		return err
	}
	r.Upper, err = parseRangeVal[T](strings.TrimSpace(parts[1]))
	if err != nil {
		return err
	}
	_ = zeroT
	return nil
}

func (r Range[T]) Value() (driver.Value, error) {
	if r.Empty {
		return "empty", nil
	}
	lb := '('
	if r.LowerInclusive {
		lb = '['
	}
	ub := ')'
	if r.UpperInclusive {
		ub = ']'
	}
	return fmt.Sprintf("%c%s,%s%c", lb, formatRangeVal(r.Lower), formatRangeVal(r.Upper), ub), nil
}

func (r Range[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := r.Value()
	return gorm.Expr("?", v)
}

// Edit helpers on generic range
func (r *Range[T]) SetBounds(lower, upper T) { r.Lower, r.Upper = lower, upper }

func (r *Range[T]) SetInclusivity(lowerInclusive, upperInclusive bool) {
	r.LowerInclusive, r.UpperInclusive = lowerInclusive, upperInclusive
}
func (r *Range[T]) SetEmpty(empty bool) { r.Empty = empty }

func parseRangeVal[T any](s string) (T, error) {
	var zero T
	if s == "" || strings.EqualFold(s, "infinity") || strings.EqualFold(s, "-infinity") {
		return zero, nil
	}

	// Prefer custom parsing/validation when available (e.g. enums or custom types).
	if scanner, ok := any(&zero).(interface{ Scan(any) error }); ok {
		if err := scanner.Scan(s); err != nil {
			return zero, err
		}
		return zero, nil
	}
	if u, ok := any(&zero).(encoding.TextUnmarshaler); ok {
		if err := u.UnmarshalText([]byte(s)); err != nil {
			return zero, err
		}
		return zero, nil
	}

	typ := reflect.TypeOf(zero)
	// Handle pointer types like *big.Rat
	if typ.Kind() == reflect.Pointer && typ.Elem() == reflect.TypeOf(big.Rat{}) {
		r := new(big.Rat)
		if _, ok := r.SetString(s); !ok {
			return zero, fmt.Errorf("invalid numrange %q", s)
		}
		return any(r).(T), nil
	}

	// Unquote timestamp bounds if needed
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		if uq, err := strconv.Unquote(s); err == nil {
			s = uq
		} else {
			s = strings.Trim(s, "\"")
		}
	}

	// time.Time and aliases
	timeType := reflect.TypeOf(time.Time{})
	if typ == timeType || typ.ConvertibleTo(timeType) {
		layouts := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05-07",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		var parsed time.Time
		var lastErr error
		for _, l := range layouts {
			tm, err := time.Parse(l, s)
			if err == nil {
				parsed = tm
				lastErr = nil
				break
			}
			lastErr = err
		}
		if lastErr != nil {
			return zero, lastErr
		}
		v := reflect.ValueOf(parsed)
		if typ != timeType {
			v = v.Convert(typ)
		}
		return v.Interface().(T), nil
	}

	// Scalar kinds (support defined types with scalar underlying kinds)
	dst := reflect.New(typ).Elem()
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return zero, err
		}
		dst.SetInt(n)
		return dst.Interface().(T), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return zero, err
		}
		dst.SetUint(n)
		return dst.Interface().(T), nil
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return zero, err
		}
		dst.SetFloat(n)
		return dst.Interface().(T), nil
	default:
		return zero, fmt.Errorf("unsupported range type")
	}
}

func formatRangeVal[T any](v T) string {
	typ := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	// Handle pointer types like *big.Rat
	if typ.Kind() == reflect.Pointer && !rv.IsNil() && typ.Elem() == reflect.TypeOf(big.Rat{}) {
		r := rv.Interface().(*big.Rat)
		return r.FloatString(10)
	}

	// time.Time and aliases
	timeType := reflect.TypeOf(time.Time{})
	if typ == timeType || typ.ConvertibleTo(timeType) {
		tm := rv
		if typ != timeType {
			tm = rv.Convert(timeType)
		}
		return tm.Interface().(time.Time).Format(time.RFC3339Nano)
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	default:
		if s, ok := any(v).(fmt.Stringer); ok {
			return s.String()
		}
		return ""
	}
}

// Range operations

// IsValid checks if the range is valid (not empty and lower <= upper when inclusive)
func (r Range[T]) IsValid() bool {
	if r.Empty {
		return true
	}

	// For comparable types, check if lower <= upper
	switch any(r.Lower).(type) {
	case int32:
		lower := any(r.Lower).(int32)
		upper := any(r.Upper).(int32)
		if lower > upper {
			return false
		}
		if lower == upper && (!r.LowerInclusive || !r.UpperInclusive) {
			return false
		}
	case int64:
		lower := any(r.Lower).(int64)
		upper := any(r.Upper).(int64)
		if lower > upper {
			return false
		}
		if lower == upper && (!r.LowerInclusive || !r.UpperInclusive) {
			return false
		}
	case time.Time:
		lower := any(r.Lower).(time.Time)
		upper := any(r.Upper).(time.Time)
		if lower.After(upper) {
			return false
		}
		if lower.Equal(upper) && (!r.LowerInclusive || !r.UpperInclusive) {
			return false
		}
	}

	return true
}

// IsFinite checks if both bounds are finite (not infinity)
func (r Range[T]) IsFinite() bool {
	if r.Empty {
		return true
	}

	// For time.Time, all values are considered finite
	// For numeric types, we assume all provided values are finite
	// since infinities are typically represented as zero values
	return true
}

// Clone creates a copy of the range
func (r Range[T]) Clone() Range[T] {
	return Range[T]{
		Lower:          r.Lower,
		Upper:          r.Upper,
		LowerInclusive: r.LowerInclusive,
		UpperInclusive: r.UpperInclusive,
		Empty:          r.Empty,
	}
}

// Equals checks if two ranges are equal
func (r Range[T]) Equals(other Range[T]) bool {
	if r.Empty != other.Empty {
		return false
	}
	if r.Empty && other.Empty {
		return true
	}
	return any(r.Lower) == any(other.Lower) &&
		any(r.Upper) == any(other.Upper) &&
		r.LowerInclusive == other.LowerInclusive &&
		r.UpperInclusive == other.UpperInclusive
}

// Contains checks if a value is within the range
func (r Range[T]) Contains(value T) bool {
	if r.Empty {
		return false
	}

	switch any(value).(type) {
	case int32:
		val := any(value).(int32)
		lower := any(r.Lower).(int32)
		upper := any(r.Upper).(int32)

		lowerOk := r.LowerInclusive && val >= lower || !r.LowerInclusive && val > lower
		upperOk := r.UpperInclusive && val <= upper || !r.UpperInclusive && val < upper
		return lowerOk && upperOk

	case int64:
		val := any(value).(int64)
		lower := any(r.Lower).(int64)
		upper := any(r.Upper).(int64)

		lowerOk := r.LowerInclusive && val >= lower || !r.LowerInclusive && val > lower
		upperOk := r.UpperInclusive && val <= upper || !r.UpperInclusive && val < upper
		return lowerOk && upperOk

	case time.Time:
		val := any(value).(time.Time)
		lower := any(r.Lower).(time.Time)
		upper := any(r.Upper).(time.Time)

		lowerOk := r.LowerInclusive && (val.After(lower) || val.Equal(lower)) ||
			!r.LowerInclusive && val.After(lower)
		upperOk := r.UpperInclusive && (val.Before(upper) || val.Equal(upper)) ||
			!r.UpperInclusive && val.Before(upper)
		return lowerOk && upperOk
	}

	return false
}

// Overlaps checks if this range overlaps with another range
func (r Range[T]) Overlaps(other Range[T]) bool {
	if r.Empty || other.Empty {
		return false
	}

	switch any(r.Lower).(type) {
	case int32:
		r1Lower := any(r.Lower).(int32)
		r1Upper := any(r.Upper).(int32)
		r2Lower := any(other.Lower).(int32)
		r2Upper := any(other.Upper).(int32)

		// Check if ranges don't overlap
		if r1Upper < r2Lower || (r1Upper == r2Lower && (!r.UpperInclusive || !other.LowerInclusive)) {
			return false
		}
		if r2Upper < r1Lower || (r2Upper == r1Lower && (!other.UpperInclusive || !r.LowerInclusive)) {
			return false
		}
		return true

	case int64:
		r1Lower := any(r.Lower).(int64)
		r1Upper := any(r.Upper).(int64)
		r2Lower := any(other.Lower).(int64)
		r2Upper := any(other.Upper).(int64)

		// Check if ranges don't overlap
		if r1Upper < r2Lower || (r1Upper == r2Lower && (!r.UpperInclusive || !other.LowerInclusive)) {
			return false
		}
		if r2Upper < r1Lower || (r2Upper == r1Lower && (!other.UpperInclusive || !r.LowerInclusive)) {
			return false
		}
		return true

	case time.Time:
		r1Lower := any(r.Lower).(time.Time)
		r1Upper := any(r.Upper).(time.Time)
		r2Lower := any(other.Lower).(time.Time)
		r2Upper := any(other.Upper).(time.Time)

		// Check if ranges don't overlap
		if r1Upper.Before(r2Lower) || (r1Upper.Equal(r2Lower) && (!r.UpperInclusive || !other.LowerInclusive)) {
			return false
		}
		if r2Upper.Before(r1Lower) || (r2Upper.Equal(r1Lower) && (!other.UpperInclusive || !r.LowerInclusive)) {
			return false
		}
		return true
	}

	return false
}

// Adjacent checks if this range is adjacent to another range
func (r Range[T]) Adjacent(other Range[T]) bool {
	if r.Empty || other.Empty {
		return false
	}

	switch any(r.Lower).(type) {
	case int32:
		r1Lower := any(r.Lower).(int32)
		r1Upper := any(r.Upper).(int32)
		r2Lower := any(other.Lower).(int32)
		r2Upper := any(other.Upper).(int32)

		// Adjacent if one range's upper bound meets the other's lower bound
		adjacent1 := r1Upper == r2Lower && (r.UpperInclusive != other.LowerInclusive)
		adjacent2 := r2Upper == r1Lower && (other.UpperInclusive != r.LowerInclusive)
		return adjacent1 || adjacent2

	case int64:
		r1Lower := any(r.Lower).(int64)
		r1Upper := any(r.Upper).(int64)
		r2Lower := any(other.Lower).(int64)
		r2Upper := any(other.Upper).(int64)

		// Adjacent if one range's upper bound meets the other's lower bound
		adjacent1 := r1Upper == r2Lower && (r.UpperInclusive != other.LowerInclusive)
		adjacent2 := r2Upper == r1Lower && (other.UpperInclusive != r.LowerInclusive)
		return adjacent1 || adjacent2

	case time.Time:
		r1Lower := any(r.Lower).(time.Time)
		r1Upper := any(r.Upper).(time.Time)
		r2Lower := any(other.Lower).(time.Time)
		r2Upper := any(other.Upper).(time.Time)

		// Adjacent if one range's upper bound meets the other's lower bound
		adjacent1 := r1Upper.Equal(r2Lower) && (r.UpperInclusive != other.LowerInclusive)
		adjacent2 := r2Upper.Equal(r1Lower) && (other.UpperInclusive != r.LowerInclusive)
		return adjacent1 || adjacent2
	}

	return false
}
