package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/slot"
)

func TestExportAvailabilityAndInvalidSlot(t *testing.T) {
	a := availability.New(time.UTC)
	a.Weekly = a.Weekly.SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(10, 0, 0)})
	from := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	b, err := ExportAvailability(a, from, to)
	if err != nil {
		t.Fatalf("export availability failed: %v", err)
	}
	if !strings.Contains(string(b), "X-WR-CALNAME:Availability") {
		t.Fatalf("expected availability calendar name")
	}

	invalid := slot.TimeSlot{Start: from.Add(time.Hour), End: from, Location: time.UTC}
	if _, err := ExportSlots([]slot.TimeSlot{invalid}, "bad"); err == nil {
		t.Fatalf("expected invalid slot error")
	}
}
