package availability

import (
	"testing"
	"time"
)

func TestWeeklySetGetValidateAndMerge(t *testing.T) {
	ws := NewWeeklySchedule(nil)
	allDays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
	for _, d := range allDays {
		ws = ws.SetDay(d, TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(10, 0, 0)})
		if len(ws.GetDay(d)) != 1 {
			t.Fatalf("expected day configured")
		}
	}
	if got := ws.GetDay(time.Weekday(99)); got != nil {
		t.Fatalf("invalid weekday should return nil")
	}
	if err := ws.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	other := NewWeeklySchedule(time.UTC).SetDay(time.Monday,
		TimeRange{Start: NewTimeOfDay(10, 0, 0), End: NewTimeOfDay(11, 0, 0)},
	)
	merged := ws.MergeWith(other)
	if len(merged.GetDay(time.Monday)) == 0 {
		t.Fatalf("expected merged monday ranges")
	}
}

func TestWeeklyTimeOfDayAndRanges(t *testing.T) {
	if err := (TimeOfDay{Hour: -1}).Validate(); err == nil {
		t.Fatalf("expected invalid hour")
	}
	if err := (TimeOfDay{Hour: 1, Minute: 60}).Validate(); err == nil {
		t.Fatalf("expected invalid minute")
	}
	if err := (TimeOfDay{Hour: 1, Minute: 1, Second: 60}).Validate(); err == nil {
		t.Fatalf("expected invalid second")
	}

	if !(TimeOfDay{Hour: 2}).after(TimeOfDay{Hour: 1}) {
		t.Fatalf("after hour failed")
	}
	if !(TimeOfDay{Hour: 1, Minute: 2}).after(TimeOfDay{Hour: 1, Minute: 1}) {
		t.Fatalf("after minute failed")
	}
	if !(TimeOfDay{Hour: 1, Minute: 1, Second: 2}).after(TimeOfDay{Hour: 1, Minute: 1, Second: 1}) {
		t.Fatalf("after second failed")
	}

	if err := (TimeRange{Start: NewTimeOfDay(10, 0, 0), End: NewTimeOfDay(9, 0, 0)}).Validate(); err == nil {
		t.Fatalf("expected invalid time range")
	}

	day := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	converted := NewTimeOfDay(9, 30, 0).ToTime(day, nil)
	if converted.Hour() != 9 || converted.Minute() != 30 {
		t.Fatalf("to time failed: %v", converted)
	}
}

func TestWeeklyGenerateNextAvailableAndNormalize(t *testing.T) {
	ws := NewWeeklySchedule(time.UTC).SetDay(time.Monday,
		TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(10, 0, 0)},
		TimeRange{Start: NewTimeOfDay(9, 30, 0), End: NewTimeOfDay(11, 0, 0)},
		TimeRange{Start: NewTimeOfDay(12, 0, 0), End: NewTimeOfDay(13, 0, 0)},
	)

	from := time.Date(2025, 1, 6, 8, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 6, 14, 0, 0, 0, time.UTC)
	slots := ws.GenerateSlots(from, to)
	if slots.Len() != 2 {
		t.Fatalf("expected merged ranges to become two slots, got %d", slots.Len())
	}
	if got := ws.GenerateSlots(to, from); got.Len() != 0 {
		t.Fatalf("expected empty when to <= from")
	}

	if next, ok := ws.NextAvailable(time.Date(2025, 1, 6, 8, 30, 0, 0, time.UTC)); !ok || next.Hour() != 9 {
		t.Fatalf("unexpected next available: %v ok=%v", next, ok)
	}
	if next, ok := ws.NextAvailable(time.Date(2025, 1, 6, 9, 30, 0, 0, time.UTC)); !ok || next.Hour() != 9 {
		t.Fatalf("expected in-range return of current time")
	}

	none := NewWeeklySchedule(time.UTC)
	if _, ok := none.NextAvailable(time.Date(2025, 1, 6, 8, 30, 0, 0, time.UTC)); ok {
		t.Fatalf("expected no availability")
	}

	if !ws.IsAvailable(time.Date(2025, 1, 6, 9, 15, 0, 0, time.UTC)) {
		t.Fatalf("expected available during range")
	}
	if ws.IsAvailable(time.Date(2025, 1, 6, 11, 30, 0, 0, time.UTC)) {
		t.Fatalf("expected unavailable outside range")
	}

	if normalizeRanges([]TimeRange{{Start: NewTimeOfDay(10, 0, 0), End: NewTimeOfDay(9, 0, 0)}}) != nil {
		t.Fatalf("invalid ranges should be dropped")
	}
}
