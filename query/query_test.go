package query

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func TestBuildAndValidate(t *testing.T) {
	q := NewQuery().Duration(time.Hour).Between(time.Now(), time.Now().Add(2*time.Hour)).Build()
	if err := q.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOptimizeSlots(t *testing.T) {
	start := time.Now().UTC().Truncate(time.Minute)
	slots := []slot.TimeSlot{
		{Start: start.Add(2 * time.Hour), End: start.Add(3 * time.Hour), Location: time.UTC},
		{Start: start.Add(1 * time.Hour), End: start.Add(2 * time.Hour), Location: time.UTC},
	}
	q := NewQuery().Duration(time.Hour).Between(start, start.Add(24*time.Hour)).PreferEarlier().Build()
	got := OptimizeSlots(slots, q)
	if !got[0].Start.Equal(slots[1].Start) {
		t.Fatalf("expected earlier slot first")
	}
}
