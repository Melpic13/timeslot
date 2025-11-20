package timeslot

import (
	"testing"
	"time"
)

func TestRootConstructors(t *testing.T) {
	start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	if _, err := NewSlot(start, end); err != nil {
		t.Fatalf("new slot failed: %v", err)
	}
	if NewCollection().Len() != 0 {
		t.Fatalf("expected empty collection")
	}
	if NewWeeklySchedule(nil).Location == nil {
		t.Fatalf("weekly schedule should set location")
	}
	if NewAvailability(nil).Location == nil {
		t.Fatalf("availability should set location")
	}
	if NewProvider("id") == nil {
		t.Fatalf("new provider should not be nil")
	}
	if NewQuery() == nil {
		t.Fatalf("new query builder should not be nil")
	}
}
