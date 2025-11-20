package provider

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
)

func TestProviderOptions(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.UTC
	}
	day := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	ws := availability.NewWeeklySchedule(time.UTC).SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(17, 0, 0)})
	p := NewProvider("p",
		WithWeeklySchedule(ws),
		WithBufferBefore(5*time.Minute),
		WithBufferAfter(10*time.Minute),
		WithBuffer(15*time.Minute),
		WithMinNotice(time.Hour),
		WithMaxAdvance(30*24*time.Hour),
		WithTimezone(loc),
		WithMetadata("k", "v"),
		WithBlockedDates(day),
	)
	if p.BufferBefore != 15*time.Minute || p.BufferAfter != 15*time.Minute {
		t.Fatalf("with buffer should set both")
	}
	if p.MinNotice != time.Hour || p.MaxAdvance == 0 {
		t.Fatalf("notice/advance options not applied")
	}
	if p.Availability.Location == nil || p.Availability.Weekly.Location == nil {
		t.Fatalf("timezone should be set")
	}
	if p.Metadata["k"] != "v" {
		t.Fatalf("metadata not set")
	}
}
