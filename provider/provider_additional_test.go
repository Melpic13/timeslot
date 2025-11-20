package provider

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

func providerAllDay() *Provider {
	ws := availability.NewWeeklySchedule(time.UTC)
	for _, d := range []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday} {
		ws = ws.SetDay(d, availability.TimeRange{Start: availability.NewTimeOfDay(0, 0, 0), End: availability.NewTimeOfDay(23, 59, 0)})
	}
	return NewProvider("all-day", WithWeeklySchedule(ws))
}

func TestProviderBehavior(t *testing.T) {
	p := providerAllDay()
	now := time.Now().UTC().Truncate(time.Minute)
	future := slot.TimeSlot{Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Location: time.UTC}

	if !p.IsAvailable(future) {
		t.Fatalf("expected available")
	}

	blockedCopy := p.WithBlockedDates(now)
	if blockedCopy == p {
		t.Fatalf("expected cloned provider")
	}

	booked, err := p.Book(future)
	if err != nil {
		t.Fatalf("expected successful booking: %v", err)
	}
	if len(booked.Availability.Bookings.Slots()) != 1 {
		t.Fatalf("expected booking persisted")
	}
	if p.Availability.Bookings.Len() != 0 {
		t.Fatalf("original provider should stay immutable")
	}

	if p2, err := booked.CancelBooking(future); err != nil || p2.Availability.Bookings.Len() != 0 {
		t.Fatalf("cancel booking failed: %v", err)
	}
	if _, err := booked.CancelBooking(slot.TimeSlot{Start: now.Add(9 * time.Hour), End: now.Add(10 * time.Hour), Location: time.UTC}); err == nil {
		t.Fatalf("expected not found cancel error")
	}

	windowBookings := booked.GetBookings(now, now.Add(24*time.Hour))
	if len(windowBookings) != 1 {
		t.Fatalf("get bookings failed")
	}

	p.BufferBefore = 10 * time.Minute
	p.BufferAfter = 20 * time.Minute
	eff := p.EffectiveAvailability(future)
	if !eff.Start.Equal(future.Start.Add(-10*time.Minute)) || !eff.End.Equal(future.End.Add(20*time.Minute)) {
		t.Fatalf("effective availability failed")
	}

	if !p.passesConstraints(future, nil) {
		t.Fatalf("nil constraints should pass")
	}
	if p.locationOrUTC() == nil {
		t.Fatalf("location fallback should not be nil")
	}
}

func TestProviderBookValidationErrors(t *testing.T) {
	p := providerAllDay()
	now := time.Now().UTC().Truncate(time.Minute)

	if _, err := p.Book(slot.TimeSlot{Start: now.Add(-2 * time.Hour), End: now.Add(-time.Hour), Location: time.UTC}); err == nil {
		t.Fatalf("expected past booking error")
	}

	p.MinNotice = 3 * time.Hour
	if _, err := p.Book(slot.TimeSlot{Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Location: time.UTC}); err == nil {
		t.Fatalf("expected insufficient notice")
	}
	p.MinNotice = 0

	p.MaxAdvance = 24 * time.Hour
	if _, err := p.Book(slot.TimeSlot{Start: now.Add(48 * time.Hour), End: now.Add(49 * time.Hour), Location: time.UTC}); err == nil {
		t.Fatalf("expected too far advance")
	}

	limited := NewProvider("none", WithWeeklySchedule(availability.NewWeeklySchedule(time.UTC)))
	if _, err := limited.Book(slot.TimeSlot{Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Location: time.UTC}); err == nil {
		t.Fatalf("expected unavailable error")
	}
}

func TestProviderFindSlotsLimitAndConstraints(t *testing.T) {
	p := providerAllDay()
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	q := query.NewQuery().Duration(time.Hour).Between(start, start.Add(24*time.Hour)).OnlyMornings().Limit(2).Build()
	slots, err := p.FindSlots(q)
	if err != nil {
		t.Fatalf("find slots failed: %v", err)
	}
	if len(slots) != 2 {
		t.Fatalf("expected limit=2")
	}

	bad := query.Query{}
	if _, err := p.FindSlots(bad); err == nil {
		t.Fatalf("expected invalid query")
	}

	p.BufferBefore = time.Hour
	existing := slot.TimeSlot{Start: start.Add(10 * time.Hour), End: start.Add(11 * time.Hour), Location: time.UTC}
	p.Availability = p.Availability.AddBooking(existing)
	if p.IsAvailable(slot.TimeSlot{Start: start.Add(10*time.Hour + 30*time.Minute), End: start.Add(11*time.Hour + 30*time.Minute), Location: time.UTC}) {
		t.Fatalf("buffer overlap should be unavailable")
	}
}
