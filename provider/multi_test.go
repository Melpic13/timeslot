package provider

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/query"
)

func TestFindAnyAndFindCommon(t *testing.T) {
	ws := availability.NewWeeklySchedule(time.UTC).SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(12, 0, 0)})
	p1 := NewProvider("p1", WithWeeklySchedule(ws))
	p2 := NewProvider("p2", WithWeeklySchedule(ws))
	from := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	q := query.NewQuery().Duration(time.Hour).Between(from, from.Add(24*time.Hour)).Build()
	any := FindAny([]*Provider{p1, p2}, q)
	if len(any) != 2 {
		t.Fatalf("expected results for all providers")
	}
	common := FindCommon([]*Provider{p1, p2}, q)
	if len(common) == 0 {
		t.Fatalf("expected common slots")
	}

	if out := FindCommon(nil, q); out != nil {
		t.Fatalf("nil providers should return nil")
	}

	bad := query.Query{}
	if out := FindCommon([]*Provider{p1}, bad); out != nil {
		t.Fatalf("invalid query should return nil")
	}
}
