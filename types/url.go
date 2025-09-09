package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type URL url.URL

func (u URL) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u URL) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v, _ := u.Value()
	return gorm.Expr("?", v)
}

func (u *URL) Scan(value interface{}) error {
	var us string
	switch v := value.(type) {
	case []byte:
		us = string(v)
	case string:
		us = v
	default:
		return errors.New(fmt.Sprint("Failed to parse URL:", value))
	}
	uu, err := url.Parse(us)
	if err != nil {
		return err
	}
	*u = URL(*uu)
	return nil
}

func (URL) GormDataType() string {
	return "url"
}

func (URL) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "TEXT"
}

func (u *URL) String() string {
	return (*url.URL)(u).String()
}

func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URL) UnmarshalJSON(data []byte) error {
	// ignore null
	if string(data) == "null" {
		return nil
	}
	uu, err := url.Parse(strings.Trim(string(data), `"'`))
	if err != nil {
		return err
	}
	*u = URL(*uu)
	return nil
}

// URL operations

// IsValid checks if the URL is valid
func (u URL) IsValid() bool {
	_, err := url.Parse(u.String())
	return err == nil
}

// Domain returns the domain of the URL
func (u URL) Domain() string {
	return (*url.URL)(&u).Host
}

// IsHTTPS checks if the URL uses HTTPS
func (u URL) IsHTTPS() bool {
	return (*url.URL)(&u).Scheme == "https"
}

// IsHTTP checks if the URL uses HTTP
func (u URL) IsHTTP() bool {
	return (*url.URL)(&u).Scheme == "http"
}

// WithScheme returns a new URL with the specified scheme
func (u URL) WithScheme(scheme string) URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	cloned.Scheme = scheme
	return URL(cloned)
}

// WithHost returns a new URL with the specified host
func (u URL) WithHost(host string) URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	cloned.Host = host
	return URL(cloned)
}

// WithPath returns a new URL with the specified path
func (u URL) WithPath(path string) URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	cloned.Path = path
	return URL(cloned)
}

// AddQuery adds a query parameter to the URL
func (u URL) AddQuery(key, value string) URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	q := cloned.Query()
	q.Add(key, value)
	cloned.RawQuery = q.Encode()
	return URL(cloned)
}

// SetQuery sets a query parameter in the URL (overwrites existing)
func (u URL) SetQuery(key, value string) URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	q := cloned.Query()
	q.Set(key, value)
	cloned.RawQuery = q.Encode()
	return URL(cloned)
}

// RemoveQuery removes a query parameter from the URL
func (u URL) RemoveQuery(key string) URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	q := cloned.Query()
	q.Del(key)
	cloned.RawQuery = q.Encode()
	return URL(cloned)
}

// HasQuery checks if the URL has a specific query parameter
func (u URL) HasQuery(key string) bool {
	return (*url.URL)(&u).Query().Has(key)
}

// GetQuery gets the value of a query parameter
func (u URL) GetQuery(key string) string {
	return (*url.URL)(&u).Query().Get(key)
}

// Clone creates a copy of the URL
func (u URL) Clone() URL {
	newURL := (*url.URL)(&u)
	cloned := *newURL
	return URL(cloned)
}

// Equals checks if two URLs are equal
func (u URL) Equals(other URL) bool {
	return u.String() == other.String()
}

// IsAbsolute checks if the URL is absolute
func (u URL) IsAbsolute() bool {
	return (*url.URL)(&u).IsAbs()
}

// GetPath returns the path component of the URL
func (u URL) GetPath() string {
	return (*url.URL)(&u).Path
}

// Query returns the query component of the URL
func (u URL) Query() url.Values {
	return (*url.URL)(&u).Query()
}

// GetFragment returns the fragment component of the URL
func (u URL) GetFragment() string {
	return (*url.URL)(&u).Fragment
}
