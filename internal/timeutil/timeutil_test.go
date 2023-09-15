package timeutil

import (
	"testing"
	"time"
)

func TestStartOfDay(t *testing.T) {
	loc := time.UTC
	in := time.Date(2024, 1, 2, 15, 4, 5, 0, loc)
	got := StartOfDay(in, loc)
	want := time.Date(2024, 1, 2, 0, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSameDay(t *testing.T) {
	a := time.Date(2024, 1, 2, 1, 0, 0, 0, time.UTC)
	b := time.Date(2024, 1, 2, 23, 59, 0, 0, time.UTC)
	if !SameDay(a, b, time.UTC) {
		t.Fatalf("expected same day")
	}
}
