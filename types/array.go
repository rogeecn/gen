package types

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Array is a generic PostgreSQL array field type.
// Supported element types include: string, int, int32, int64, uint, uint32, uint64, float32, float64, bool.
// For other types, StringArray/IntArray/FloatArray aliases are recommended.
type Array[T any] []T

func NewArray[T any](s []T) Array[T] { return Array[T](s) }

// Set replaces the slice with given value
func (a *Array[T]) Set(s []T) { *a = Array[T](s) }

// Append appends elements to the slice
func (a *Array[T]) Append(elems ...T) { *a = append(*a, elems...) }

// Value implements driver.Valuer using PostgreSQL array literal syntax.
func (a Array[T]) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	var b strings.Builder
	b.WriteByte('{')
	for i := range a {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(arrayElemToLiteral(a[i]))
	}
	b.WriteByte('}')
	return b.String(), nil
}

// Scan implements sql.Scanner from PostgreSQL array literal.
func (a *Array[T]) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	s, ok := toString(value)
	if !ok {
		return driver.ErrBadConn
	}
	tokens := parsePgArrayElements(s)
	out := make([]T, 0, len(tokens))
	for _, t := range tokens {
		var v T
		if err := parseInto(&v, t); err != nil {
			// skip invalid element
			continue
		}
		out = append(out, v)
	}
	*a = out
	return nil
}

// GORM data type mapping
func (Array[T]) GormDataType() string { return elementDBType[T]() }

func (Array[T]) GormDBDataType(*gorm.DB, *schema.Field) string {
	return strings.ToUpper(elementDBType[T]())
}

func (a Array[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := a.Value()
	return gorm.Expr("?", v)
}

// Helpers

func arrayElemToLiteral[T any](v T) string {
	switch x := any(v).(type) {
	case string:
		// quote and escape
		var b strings.Builder
		b.WriteByte('"')
		for _, r := range x {
			switch r {
			case '\\', '"':
				b.WriteByte('\\')
				b.WriteRune(r)
			default:
				b.WriteRune(r)
			}
		}
		b.WriteByte('"')
		return b.String()
	case bool:
		if x {
			return "t"
		}
		return "f"
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	default:
		// fallback to quoted string
		s := strings.ReplaceAll(strings.ReplaceAll(toStringOrSprint(x), "\\", "\\\\"), "\"", "\\\"")
		return "\"" + s + "\""
	}
}

func toStringOrSprint(v any) string { return fmt.Sprint(v) }

func parseInto[T any](dst *T, token string) error {
	switch any(*dst).(type) {
	case string:
		*dst = any(token).(T)
		return nil
	case bool:
		switch strings.ToLower(token) {
		case "t", "true":
			*dst = any(true).(T)
		case "f", "false":
			*dst = any(false).(T)
		default:
			return strconv.ErrSyntax
		}
		return nil
	case int:
		n, err := strconv.Atoi(token)
		if err != nil {
			return err
		}
		*dst = any(n).(T)
		return nil
	case int32:
		n, err := strconv.ParseInt(token, 10, 32)
		if err != nil {
			return err
		}
		*dst = any(int32(n)).(T)
		return nil
	case int64:
		n, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return err
		}
		*dst = any(n).(T)
		return nil
	case uint:
		n, err := strconv.ParseUint(token, 10, 0)
		if err != nil {
			return err
		}
		*dst = any(uint(n)).(T)
		return nil
	case uint32:
		n, err := strconv.ParseUint(token, 10, 32)
		if err != nil {
			return err
		}
		*dst = any(uint32(n)).(T)
		return nil
	case uint64:
		n, err := strconv.ParseUint(token, 10, 64)
		if err != nil {
			return err
		}
		*dst = any(n).(T)
		return nil
	case float32:
		n, err := strconv.ParseFloat(token, 32)
		if err != nil {
			return err
		}
		*dst = any(float32(n)).(T)
		return nil
	case float64:
		n, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return err
		}
		*dst = any(n).(T)
		return nil
	default:
		return strconv.ErrSyntax
	}
}

func elementDBType[T any]() string {
	switch any(*new(T)).(type) {
	case string:
		return "text[]"
	case bool:
		return "boolean[]"
	case int, int32:
		return "integer[]"
	case int64:
		return "bigint[]"
	case uint, uint32:
		return "integer[]"
	case uint64:
		return "bigint[]"
	case float32:
		return "real[]"
	case float64:
		return "double precision[]"
	default:
		return "text[]"
	}
}

// parsePgArrayElements parses a PostgreSQL array text into unquoted elements.
// It handles simple quoted/unquoted tokens and backslash escapes. Nested arrays are not supported.
func parsePgArrayElements(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" || s == "{}" {
		return nil
	}
	if s[0] == '{' && s[len(s)-1] == '}' {
		s = s[1 : len(s)-1]
	}
	res := make([]string, 0)
	var cur strings.Builder
	inQuotes := false
	esc := false
	flush := func() {
		res = append(res, cur.String())
		cur.Reset()
	}
	for _, ch := range s {
		if inQuotes {
			if esc {
				cur.WriteRune(ch)
				esc = false
				continue
			}
			switch ch {
			case '\\':
				esc = true
			case '"':
				inQuotes = false
			default:
				cur.WriteRune(ch)
			}
			continue
		}
		switch ch {
		case '"':
			inQuotes = true
		case ',':
			flush()
		default:
			if ch != ' ' && ch != '\n' && ch != '\t' { // trim whitespace outside quotes
				cur.WriteRune(ch)
			}
		}
	}
	flush()
	// Convert unquoted NULL to empty string for lack of null-element representation
	for i := range res {
		if res[i] == "NULL" {
			res[i] = ""
		}
	}
	return res
}

// Collection operations

// Contains checks if the array contains the given element
func (a Array[T]) Contains(elem T) bool {
	for _, v := range a {
		if any(v) == any(elem) {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of elem in the array, or -1 if not found
func (a Array[T]) IndexOf(elem T) int {
	for i, v := range a {
		if any(v) == any(elem) {
			return i
		}
	}
	return -1
}

// Remove removes all occurrences of elem from the array
func (a Array[T]) Remove(elem T) Array[T] {
	result := make([]T, 0, len(a))
	for _, v := range a {
		if any(v) != any(elem) {
			result = append(result, v)
		}
	}
	return Array[T](result)
}

// RemoveAt removes the element at the given index
func (a Array[T]) RemoveAt(index int) Array[T] {
	if index < 0 || index >= len(a) {
		return a.Clone()
	}
	result := make([]T, 0, len(a)-1)
	result = append(result, a[:index]...)
	result = append(result, a[index+1:]...)
	return Array[T](result)
}

// Unique returns a new array with duplicate elements removed
func (a Array[T]) Unique() Array[T] {
	seen := make(map[any]bool)
	result := make([]T, 0, len(a))
	for _, v := range a {
		key := any(v)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	return Array[T](result)
}

// Filter returns a new array containing only elements that satisfy the predicate
func (a Array[T]) Filter(predicate func(T) bool) Array[T] {
	result := make([]T, 0)
	for _, v := range a {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return Array[T](result)
}

// Clear removes all elements from the array
func (a *Array[T]) Clear() {
	*a = Array[T](nil)
}

// IsEmpty returns true if the array has no elements
func (a Array[T]) IsEmpty() bool {
	return len(a) == 0
}

// Len returns the length of the array
func (a Array[T]) Len() int {
	return len(a)
}

// Clone creates a deep copy of the array
func (a Array[T]) Clone() Array[T] {
	if a == nil {
		return nil
	}
	result := make([]T, len(a))
	copy(result, a)
	return Array[T](result)
}

// Equals compares two arrays for equality
func (a Array[T]) Equals(other Array[T]) bool {
	if len(a) != len(other) {
		return false
	}
	for i, v := range a {
		if any(v) != any(other[i]) {
			return false
		}
	}
	return true
}

// Reverse returns a new array with elements in reverse order
func (a Array[T]) Reverse() Array[T] {
	result := make([]T, len(a))
	for i, v := range a {
		result[len(a)-1-i] = v
	}
	return Array[T](result)
}
