package types

import (
    "context"
    "database/sql/driver"
    "errors"
    "fmt"
    "math/big"
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
    v, _ := r.Value(); return gorm.Expr("?", v)
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
	switch any(zero).(type) {
	case int32:
		v, err := strconv.ParseInt(s, 10, 32)
		return any(int32(v)).(T), err
	case int64:
		v, err := strconv.ParseInt(s, 10, 64)
		return any(int64(v)).(T), err
	case *big.Rat:
		r := new(big.Rat)
		if _, ok := r.SetString(s); !ok {
			return zero, fmt.Errorf("invalid numrange %q", s)
		}
		return any(r).(T), nil
	case time.Time:
		layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05-07", "2006-01-02"}
		var t time.Time
		var err error
		for _, l := range layouts {
			t, err = time.Parse(l, s)
			if err == nil {
				return any(t).(T), nil
			}
		}
		return zero, err
	default:
		return zero, fmt.Errorf("unsupported range type")
	}
}

func formatRangeVal[T any](v T) string {
	switch x := any(v).(type) {
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
    case *big.Rat:
        if x == nil {
            return ""
        }
        return x.RatString()
	case time.Time:
		return x.Format(time.RFC3339Nano)
	default:
		return ""
	}
}
