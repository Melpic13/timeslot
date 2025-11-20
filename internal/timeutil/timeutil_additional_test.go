package timeutil

import (
	"testing"
	"time"
)

func TestEndOfDayClampAndOverlap(t *testing.T) {
	loc := time.UTC
	tm := time.Date(2025, 1, 1, 10, 0, 0, 0, loc)
	eod := EndOfDay(tm, loc)
	if eod.Hour() != 23 || eod.Minute() != 59 {
		t.Fatalf("unexpected end of day: %v", eod)
	}

	if got := Clamp(tm.Add(-time.Hour), tm, tm.Add(time.Hour)); !got.Equal(tm) {
		t.Fatalf("clamp low failed")
	}
	if got := Clamp(tm.Add(2*time.Hour), tm, tm.Add(time.Hour)); !got.Equal(tm.Add(time.Hour)) {
		t.Fatalf("clamp high failed")
	}
	if got := Clamp(tm, tm.Add(-time.Hour), tm.Add(time.Hour)); !got.Equal(tm) {
		t.Fatalf("clamp mid failed")
	}

	if !Overlap(tm, tm.Add(time.Hour), tm.Add(30*time.Minute), tm.Add(90*time.Minute)) {
		t.Fatalf("expected overlap")
	}
	if Overlap(tm, tm.Add(time.Hour), tm.Add(time.Hour), tm.Add(2*time.Hour)) {
		t.Fatalf("expected no overlap at boundary")
	}
}
