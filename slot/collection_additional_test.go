package slot

import (
	"testing"
	"time"
)

func TestCollectionOperationsFull(t *testing.T) {
	loc := time.UTC
	a := TimeSlot{Start: time.Date(2025, 1, 6, 9, 0, 0, 0, loc), End: time.Date(2025, 1, 6, 10, 0, 0, 0, loc), Location: loc}
	b := TimeSlot{Start: time.Date(2025, 1, 6, 11, 0, 0, 0, loc), End: time.Date(2025, 1, 6, 12, 0, 0, 0, loc), Location: loc}
	c := TimeSlot{Start: time.Date(2025, 1, 6, 9, 30, 0, 0, loc), End: time.Date(2025, 1, 6, 11, 30, 0, 0, loc), Location: loc}

	col := NewCollection(a).Add(b)
	if col.Len() != 2 {
		t.Fatalf("expected 2 slots")
	}
	unioned := col.Union(NewCollection(c))
	if unioned.Len() != 1 {
		t.Fatalf("expected merged union, got %d", unioned.Len())
	}

	inter := col.Intersect(NewCollection(c))
	if inter.Len() != 2 {
		t.Fatalf("expected 2 intersections, got %d", inter.Len())
	}

	removed := unioned.Remove(TimeSlot{Start: time.Date(2025, 1, 6, 10, 0, 0, 0, loc), End: time.Date(2025, 1, 6, 10, 30, 0, 0, loc), Location: loc})
	if removed.Len() != 2 {
		t.Fatalf("expected split after remove, got %d", removed.Len())
	}

	filtered := unioned.Filter(func(s TimeSlot) bool { return s.Duration() >= 2*time.Hour })
	if filtered.Len() != 1 {
		t.Fatalf("filter failed")
	}

	overlaps := unioned.FindOverlaps(TimeSlot{Start: time.Date(2025, 1, 6, 9, 15, 0, 0, loc), End: time.Date(2025, 1, 6, 9, 45, 0, 0, loc), Location: loc})
	if len(overlaps) != 1 {
		t.Fatalf("find overlaps failed")
	}
	if NewCollection().FindOverlaps(a) != nil {
		t.Fatalf("empty overlaps should be nil")
	}

	if unioned.TotalDuration() != 3*time.Hour {
		t.Fatalf("total duration mismatch")
	}

	gaps := NewCollection(a, b).Gaps(TimeSlot{Start: time.Date(2025, 1, 6, 9, 0, 0, 0, loc), End: time.Date(2025, 1, 6, 12, 0, 0, 0, loc), Location: loc})
	if gaps.Len() != 1 {
		t.Fatalf("expected one gap")
	}

	if NewCollection().IsEmpty() != true {
		t.Fatalf("expected empty")
	}
	if _, ok := NewCollection().First(); ok {
		t.Fatalf("expected no first element")
	}
	if _, ok := NewCollection().Last(); ok {
		t.Fatalf("expected no last element")
	}
	if _, ok := col.First(); !ok {
		t.Fatalf("expected first")
	}
	if _, ok := col.Last(); !ok {
		t.Fatalf("expected last")
	}
}
