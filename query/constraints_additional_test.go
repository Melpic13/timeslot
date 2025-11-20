package query

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func testSlotAt(hour int) slot.TimeSlot {
	loc := time.UTC
	start := time.Date(2025, 1, 6, hour, 0, 0, 0, loc)
	return slot.TimeSlot{Start: start, End: start.Add(time.Hour), Location: loc}
}

func TestQueryTimeOfDayAndConstraints(t *testing.T) {
	if err := NewTimeOfDay(24, 0, 0).Validate(); err == nil {
		t.Fatalf("expected invalid hour")
	}
	if err := NewTimeOfDay(10, 60, 0).Validate(); err == nil {
		t.Fatalf("expected invalid minute")
	}
	if err := NewTimeOfDay(10, 0, 60).Validate(); err == nil {
		t.Fatalf("expected invalid second")
	}

	day := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	if got := NewTimeOfDay(9, 30, 0).toTime(day, nil); got.Hour() != 9 || got.Minute() != 30 {
		t.Fatalf("toTime mismatch")
	}

	wc := NewWeekdayConstraint(time.Monday)
	if !wc.IsSatisfied(testSlotAt(9)) {
		t.Fatalf("weekday should match")
	}
	if wc.String() == "" {
		t.Fatalf("weekday string empty")
	}

	startTOD := NewTimeOfDay(9, 0, 0)
	endTOD := NewTimeOfDay(17, 0, 0)
	tod := TimeOfDayConstraint{Start: &startTOD, End: &endTOD}
	if !tod.IsSatisfied(testSlotAt(9)) {
		t.Fatalf("time-of-day should satisfy lower bound")
	}
	if tod.IsSatisfied(testSlotAt(18)) {
		t.Fatalf("time-of-day should fail upper bound")
	}
	if tod.String() == "" {
		t.Fatalf("time-of-day string empty")
	}

	nb := NotBeforeConstraint{Time: NewTimeOfDay(10, 0, 0)}
	if nb.IsSatisfied(testSlotAt(9)) {
		t.Fatalf("not-before should fail")
	}
	if nb.String() == "" {
		t.Fatalf("not-before string empty")
	}

	na := NotAfterConstraint{Time: NewTimeOfDay(10, 0, 0)}
	if na.IsSatisfied(testSlotAt(11)) {
		t.Fatalf("not-after should fail")
	}
	if na.String() == "" {
		t.Fatalf("not-after string empty")
	}

	booking := testSlotAt(12)
	gap := MinGapConstraint{Gap: 30 * time.Minute, Bookings: []slot.TimeSlot{booking}}
	if gap.IsSatisfied(slot.TimeSlot{Start: booking.Start.Add(15 * time.Minute), End: booking.End.Add(15 * time.Minute), Location: time.UTC}) {
		t.Fatalf("min gap should fail overlap")
	}
	if !gap.IsSatisfied(testSlotAt(14)) {
		t.Fatalf("min gap should pass")
	}
	if gap.String() == "" {
		t.Fatalf("min gap string empty")
	}
}

func TestQueryPreferences(t *testing.T) {
	early := testSlotAt(9)
	late := testSlotAt(11)
	pe := preferEarlier{}
	if pe.Score(early) <= pe.Score(late) {
		t.Fatalf("preferEarlier should score earlier higher")
	}
	if pe.String() == "" {
		t.Fatalf("preferEarlier string")
	}

	pl := preferLater{}
	if pl.Score(late) <= pl.Score(early) {
		t.Fatalf("preferLater should score later higher")
	}
	if pl.String() == "" {
		t.Fatalf("preferLater string")
	}

	pt := preferTime{Target: NewTimeOfDay(10, 0, 0)}
	if pt.Score(testSlotAt(10)) <= pt.Score(testSlotAt(12)) {
		t.Fatalf("preferTime should reward closeness")
	}
	if pt.String() == "" {
		t.Fatalf("preferTime string")
	}
}
