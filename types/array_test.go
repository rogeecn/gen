package types

import (
	"fmt"
	"testing"
)

type testRole string

const (
	testRoleUser       testRole = "user"
	testRoleSuperAdmin testRole = "super_admin"
)

func parseTestRole(s string) (testRole, error) {
	switch s {
	case string(testRoleUser):
		return testRoleUser, nil
	case string(testRoleSuperAdmin):
		return testRoleSuperAdmin, nil
	default:
		return "", fmt.Errorf("invalid role: %q", s)
	}
}

// Scan lets parseInto use the Scanner fast-path (mirrors go-enum generated enums).
func (r *testRole) Scan(value any) error {
	switch v := value.(type) {
	case string:
		out, err := parseTestRole(v)
		if err != nil {
			return err
		}
		*r = out
		return nil
	default:
		return fmt.Errorf("invalid type: %T", value)
	}
}

type stringAlias string

func TestArrayScan_EnumScanner(t *testing.T) {
	var a Array[testRole]
	if err := a.Scan(`{user,super_admin}`); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(a) != 2 || a[0] != testRoleUser || a[1] != testRoleSuperAdmin {
		t.Fatalf("unexpected values: %#v", []testRole(a))
	}
}

func TestArrayScan_StringAlias(t *testing.T) {
	var a Array[stringAlias]
	if err := a.Scan(`{"a","b"}`); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(a) != 2 || a[0] != "a" || a[1] != "b" {
		t.Fatalf("unexpected values: %#v", []stringAlias(a))
	}
}

