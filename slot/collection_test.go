package slot

import (
	"testing"
	"time"
)

func TestCollectionMerge(t *testing.T) {
	loc := time.UTC
	a := TimeSlot{Start: time.Date(2024, 1, 1, 9, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 10, 0, 0, 0, loc), Location: loc}
	b := TimeSlot{Start: time.Date(2024, 1, 1, 10, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 11, 0, 0, 0, loc), Location: loc}
	c := NewCollection(a, b)
	if c.Len() != 1 {
		t.Fatalf("expected merged slot, got %d", c.Len())
	}
}

func TestCollectionSubtract(t *testing.T) {
	loc := time.UTC
	base := NewCollection(TimeSlot{Start: time.Date(2024, 1, 1, 9, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 12, 0, 0, 0, loc), Location: loc})
	cut := NewCollection(TimeSlot{Start: time.Date(2024, 1, 1, 10, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 11, 0, 0, 0, loc), Location: loc})
	result := base.Subtract(cut)
	if result.Len() != 2 {
		t.Fatalf("expected 2 slots, got %d", result.Len())
	}
}
