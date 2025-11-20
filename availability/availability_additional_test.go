package availability

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func setupAvailabilityAllDayMonday() Availability {
	a := New(time.UTC)
	a.Weekly = a.Weekly.SetDay(time.Monday, TimeRange{Start: NewTimeOfDay(0, 0, 0), End: NewTimeOfDay(23, 59, 0)})
	return a
}

func TestAvailabilityAddersAndBookingChecks(t *testing.T) {
	base := setupAvailabilityAllDayMonday()
	day := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	base = base.AddBlockedDates(day)
	base = base.AddBlockedRange(day.Add(2*time.Hour), day.Add(3*time.Hour))
	base = base.AddAvailableOverride(day.AddDate(0, 0, 1), TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(10, 0, 0)})

	booking := slot.TimeSlot{Start: day.Add(4 * time.Hour), End: day.Add(5 * time.Hour), Location: time.UTC}
	base = base.AddBooking(booking)
	if !base.IsBooked(day.Add(4*time.Hour + 30*time.Minute)) {
		t.Fatalf("expected booked instant")
	}
	base = base.RemoveBooking(booking)
	if base.IsBooked(day.Add(4*time.Hour + 30*time.Minute)) {
		t.Fatalf("expected booking removed")
	}

	if base.locationOrUTC() == nil {
		t.Fatalf("location must not be nil")
	}
}

func TestAvailabilityGetSlotsOverridesAndValidate(t *testing.T) {
	base := setupAvailabilityAllDayMonday()
	monday := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	base = base.AddBlockedRange(monday.Add(8*time.Hour), monday.Add(10*time.Hour))
	base = base.AddAvailableOverride(monday, TimeRange{Start: NewTimeOfDay(11, 0, 0), End: NewTimeOfDay(12, 0, 0)})

	from := monday
	to := monday.Add(24 * time.Hour)
	slots := base.GetSlots(from, to)
	if slots.Len() == 0 {
		t.Fatalf("expected slots")
	}
	if !base.IsAvailable(monday.Add(11*time.Hour + 30*time.Minute)) {
		t.Fatalf("expected available in override")
	}
	if base.IsAvailable(monday.Add(8*time.Hour + 30*time.Minute)) {
		t.Fatalf("expected unavailable in blocked range")
	}
	if got := base.GetSlots(to, from); got.Len() != 0 {
		t.Fatalf("expected empty for invalid bounds")
	}

	invalid := base
	invalid.Bookings = slot.NewCollection(slot.TimeSlot{Start: monday.Add(2 * time.Hour), End: monday.Add(time.Hour), Location: time.UTC})
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected invalid booking validation error")
	}

	badWeekly := base
	badWeekly.Weekly.Monday = []TimeRange{{Start: NewTimeOfDay(10, 0, 0), End: NewTimeOfDay(9, 0, 0)}}
	if err := badWeekly.Validate(); err == nil {
		t.Fatalf("expected weekly validation error")
	}
}
