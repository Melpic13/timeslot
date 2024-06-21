package conflict

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

func TestDetectorFindCommon(t *testing.T) {
	ws := availability.NewWeeklySchedule(time.UTC).SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(10, 0, 0)})
	p1 := provider.NewProvider("a", provider.WithWeeklySchedule(ws))
	p2 := provider.NewProvider("b", provider.WithWeeklySchedule(ws))
	d := NewDetector(p1, p2)
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	q := query.NewQuery().Duration(time.Hour).Between(from, to).Build()
	common := d.FindCommonAvailability(q)
	if len(common) == 0 {
		t.Fatalf("expected common availability")
	}
}

func TestResolve(t *testing.T) {
	d := NewDetector()
	c := Conflict{Slot: slot.TimeSlot{Start: time.Now().UTC(), End: time.Now().UTC().Add(time.Hour), Location: time.UTC}}
	resolved, err := d.Resolve(c, StrategyShiftForward)
	if err != nil || resolved == nil {
		t.Fatalf("expected resolved slot")
	}
}
