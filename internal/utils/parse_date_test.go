package utils

import (
	"testing"
	"time"
)

// mustYM parses a year-month string (YYYY-MM or MM-YYYY) and fails the test on error.
func mustYM(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := ParseYearMonth(s)
	if err != nil {
		t.Fatalf("parse %s: %v", s, err)
	}
	return v
}

func TestParseYearMonth_ValidFormats(t *testing.T) {
	got1 := mustYM(t, "2025-07")
	if got1.Year() != 2025 || got1.Month() != time.July || got1.Day() != 1 {
		t.Fatalf("want 2025-07-01, got %v", got1)
	}
	// Alternative format MM-YYYY must give the same result
	got2 := mustYM(t, "07-2025")
	if !got1.Equal(got2) {
		t.Fatalf("formats must match: %v vs %v", got1, got2)
	}
}

func TestParseYearMonth_InvalidFormats(t *testing.T) {
	// Wrong delimiter
	if _, err := ParseYearMonth("2025/07"); err == nil {
		t.Fatal("want error for bad delimiter")
	}
	// Too short / wrong length
	if _, err := ParseYearMonth("2025-7"); err == nil {
		t.Fatal("want error for wrong length")
	}
	// Month out of range (upper)
	if _, err := ParseYearMonth("13-2025"); err == nil {
		t.Fatal("want error for month=13")
	}
	// Month out of range (lower)
	if _, err := ParseYearMonth("00-2025"); err == nil {
		t.Fatal("want error for month=00")
	}
	if _, err := ParseYearMonth("2025-00"); err == nil {
		t.Fatal("want error for month=00 (YYYY-MM)")
	}
}

func TestYMString(t *testing.T) {
	d := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC)
	if s := ymString(d); s != "2025-07" {
		t.Fatalf("ymString: got %q want %q", s, "2025-07")
	}
}

func TestMonthsOverlapInclusive_Basics(t *testing.T) {
	// Subscription: Jul..Sep 2025 (3 months)
	aStart := mustYM(t, "2025-07")
	aEnd := mustYM(t, "2025-09")

	// Period fully covers subscription
	pStart := mustYM(t, "2025-07")
	pEnd := mustYM(t, "2025-09")
	if n := MonthsOverlapInclusive(aStart, &aEnd, pStart, &pEnd); n != 3 {
		t.Fatalf("full overlap: got %d want 3", n)
	}

	// Period inside subscription → 1 month (Aug only)
	p2Start := mustYM(t, "2025-08")
	p2End := mustYM(t, "2025-08")
	if n := MonthsOverlapInclusive(aStart, &aEnd, p2Start, &p2End); n != 1 {
		t.Fatalf("single-month overlap: got %d want 1", n)
	}

	p3Start := mustYM(t, "2025-10")
	p3End := mustYM(t, "2025-12")
	if n := MonthsOverlapInclusive(aStart, &aEnd, p3Start, &p3End); n != 0 {
		t.Fatalf("no overlap: got %d want 0", n)
	}
}

func TestMonthsOverlapInclusive_OpenEnded(t *testing.T) {
	// Open-ended subscription starting Jul 2025
	aStart := mustYM(t, "2025-07")
	// aEnd = nil (open-ended)

	// Period Jul..Aug 2025 → 2 months
	pStart := mustYM(t, "2025-07")
	pEnd := mustYM(t, "2025-08")
	if n := MonthsOverlapInclusive(aStart, nil, pStart, &pEnd); n != 2 {
		t.Fatalf("open-ended overlap: got %d want 2", n)
	}
}

func TestMonthsOverlapInclusive_BoundariesInclusive(t *testing.T) {
	// Subscription: exactly one month (Jul 2025)
	aStart := mustYM(t, "2025-07")
	aEnd := mustYM(t, "2025-07")

	// Period also Jul 2025 → inclusive boundaries = 1 month
	pStart := mustYM(t, "2025-07")
	pEnd := mustYM(t, "2025-07")
	if n := MonthsOverlapInclusive(aStart, &aEnd, pStart, &pEnd); n != 1 {
		t.Fatalf("inclusive boundary should count as 1, got %d", n)
	}
}
