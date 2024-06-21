package provider

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/query"
)

func TestProviderFindSlots(t *testing.T) {
	ws := availability.NewWeeklySchedule(time.UTC).SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(11, 0, 0)})
	p := NewProvider("p1", WithWeeklySchedule(ws))
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	q := query.NewQuery().Duration(time.Hour).Between(from, to).Build()
	slots, err := p.FindSlots(q)
	if err != nil {
		t.Fatalf("find slots: %v", err)
	}
	if len(slots) != 2 {
		t.Fatalf("expected 2, got %d", len(slots))
	}
}
