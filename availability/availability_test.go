package availability

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func TestAvailabilityGetSlotsSubtractBookings(t *testing.T) {
	base := New(time.UTC)
	base.Weekly = base.Weekly.SetDay(time.Monday, TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(12, 0, 0)})
	booking := slot.TimeSlot{Start: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC), Location: time.UTC}
	base = base.AddBooking(booking)
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	free := base.GetSlots(from, to)
	if free.Len() != 2 {
		t.Fatalf("expected 2 free slots, got %d", free.Len())
	}
}

func TestFindAvailableSlots(t *testing.T) {
	base := New(time.UTC)
	base.Weekly = base.Weekly.SetDay(time.Monday, TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(11, 0, 0)})
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	got := base.FindAvailableSlots(time.Hour, from, to)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
