package availability

import (
	"testing"
	"time"
)

func TestWeeklyGenerateSlots(t *testing.T) {
	ws := NewWeeklySchedule(time.UTC).SetDay(time.Monday, TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(10, 0, 0)})
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	to := from.Add(24 * time.Hour)
	slots := ws.GenerateSlots(from, to)
	if slots.Len() != 1 {
		t.Fatalf("got %d slots", slots.Len())
	}
}

func TestWeeklyIsAvailable(t *testing.T) {
	ws := NewWeeklySchedule(time.UTC).SetDay(time.Tuesday, TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(17, 0, 0)})
	tm := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
	if !ws.IsAvailable(tm) {
		t.Fatalf("expected available")
	}
}
