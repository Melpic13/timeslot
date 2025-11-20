package conflict

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

func makeProviderForConflict() *provider.Provider {
	ws := availability.NewWeeklySchedule(time.UTC)
	for _, d := range []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday} {
		ws = ws.SetDay(d, availability.TimeRange{Start: availability.NewTimeOfDay(0, 0, 0), End: availability.NewTimeOfDay(23, 59, 0)})
	}
	return provider.NewProvider("conflict", provider.WithWeeklySchedule(ws), provider.WithBuffer(30*time.Minute))
}

func TestConflictBufferHelpers(t *testing.T) {
	s := slot.TimeSlot{Start: time.Date(2025, 1, 6, 10, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 6, 11, 0, 0, 0, time.UTC), Location: time.UTC}
	buffered := ApplyBuffer(s, 15*time.Minute, 10*time.Minute)
	if !buffered.Start.Equal(s.Start.Add(-15*time.Minute)) || !buffered.End.Equal(s.End.Add(10*time.Minute)) {
		t.Fatalf("apply buffer failed")
	}
	unbuffered := RemoveBuffer(buffered, 15*time.Minute, 10*time.Minute)
	if !unbuffered.Start.Equal(s.Start) || !unbuffered.End.Equal(s.End) {
		t.Fatalf("remove buffer failed")
	}
}

func TestConflictDetectorPaths(t *testing.T) {
	p := makeProviderForConflict()
	d := NewDetector().AddProvider(p)

	if got := d.Check(slot.TimeSlot{Start: time.Now().UTC(), End: time.Now().UTC().Add(time.Hour), Location: time.UTC}, nil); len(got) != 0 {
		t.Fatalf("nil provider should yield no conflicts")
	}

	start := time.Now().UTC().Add(4 * time.Hour).Truncate(time.Minute)
	candidate := slot.TimeSlot{Start: start, End: start.Add(time.Hour), Location: time.UTC}
	booked, err := p.Book(candidate)
	if err != nil {
		t.Fatalf("setup booking failed: %v", err)
	}

	d.providers = []*provider.Provider{booked}
	d.options.AllowDoubleBooking = false
	d.options.IncludeBuffers = true
	conflicts := d.Check(candidate, booked)
	if len(conflicts) == 0 {
		t.Fatalf("expected conflicts")
	}

	all := d.CheckAll(candidate)
	if len(all) == 0 {
		t.Fatalf("expected check all conflicts")
	}

	q := query.NewQuery().Duration(30*time.Minute).Between(start, start.Add(3*time.Hour)).Build()
	_ = d.FindConflictFree(q)
	availableByProvider := d.FindAvailableSlots(q)
	if len(availableByProvider) == 0 {
		t.Fatalf("expected available slots map")
	}
}

func TestConflictResolutionHelpers(t *testing.T) {
	s := slot.TimeSlot{Start: time.Date(2025, 1, 6, 10, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 6, 11, 0, 0, 0, time.UTC), Location: time.UTC}
	c := Conflict{Slot: s}
	if _, err := resolveSlot(c, ResolutionStrategy(999)); err == nil {
		t.Fatalf("expected unknown strategy error")
	}
	if opts := defaultResolutionOptions(s); len(opts) != 3 {
		t.Fatalf("expected default options")
	}
	w := slotWindowAround(s, time.Hour)
	if !w.Start.Equal(s.Start.Add(-time.Hour)) || !w.End.Equal(s.End.Add(time.Hour)) {
		t.Fatalf("slot window mismatch")
	}
}
