package types

import (
	"testing"
	"time"
)

type myInt64 int64
type myTime time.Time

func TestRangeScan_Int64Alias(t *testing.T) {
	var r Range[myInt64]
	if err := r.Scan(`[1,2)`); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if r.Empty || r.Lower != 1 || r.Upper != 2 || !r.LowerInclusive || r.UpperInclusive {
		t.Fatalf("unexpected range: %#v", r)
	}
}

func TestRangeScan_TimeAlias(t *testing.T) {
	var r Range[myTime]
	if err := r.Scan(`["2025-01-01T00:00:00Z","2025-01-02T00:00:00Z")`); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	lower := time.Time(r.Lower)
	upper := time.Time(r.Upper)
	if lower.IsZero() || upper.IsZero() || !lower.Before(upper) {
		t.Fatalf("unexpected bounds: lower=%v upper=%v", lower, upper)
	}
}

