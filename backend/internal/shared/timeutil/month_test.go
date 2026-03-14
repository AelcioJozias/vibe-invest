package timeutil

import (
	"testing"
	"time"
)

func TestParseReferenceMonth(t *testing.T) {
	month, err := ParseReferenceMonth("2026-03")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	if !month.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, month)
	}
}

func TestParseReferenceMonthInvalid(t *testing.T) {
	if _, err := ParseReferenceMonth("03-2026"); err == nil {
		t.Fatal("expected error for invalid format")
	}
}
