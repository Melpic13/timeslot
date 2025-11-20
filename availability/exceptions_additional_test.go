package availability

import (
	"testing"
	"time"
)

func TestExceptionsRangeAndOverrides(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	es := ExceptionSet{}.
		AddBlockedRange(start, end).
		AddAvailableOverride(start)

	if !es.IsBlocked(start.Add(time.Hour)) {
		t.Fatalf("expected blocked")
	}
	if !es.HasAvailableOverride(start.Add(time.Hour)) {
		t.Fatalf("expected available override")
	}

	es = es.AddAvailableOverride(start,
		TimeRange{Start: NewTimeOfDay(9, 0, 0), End: NewTimeOfDay(10, 0, 0)},
	)
	ranges, ok := es.ModifiedForDate(start)
	if !ok || len(ranges) != 1 {
		t.Fatalf("expected modified ranges")
	}
	if _, ok := es.ModifiedForDate(start.AddDate(0, 0, 1)); ok {
		t.Fatalf("unexpected modified date")
	}
}

func TestExceptionsValidate(t *testing.T) {
	invalid := ExceptionSet{Blocked: []DateRange{{Start: time.Now(), End: time.Now().Add(-time.Hour)}}}
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected blocked range validation error")
	}

	now := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	invalidAvail := ExceptionSet{Available: []DateRange{{Start: now, End: now.Add(-time.Hour)}}}
	if err := invalidAvail.Validate(); err == nil {
		t.Fatalf("expected available range validation error")
	}

	invalidModified := ExceptionSet{Modified: []DateOverride{{Date: now, Ranges: []TimeRange{{Start: NewTimeOfDay(10, 0, 0), End: NewTimeOfDay(9, 0, 0)}}}}}
	if err := invalidModified.Validate(); err == nil {
		t.Fatalf("expected modified range validation error")
	}
}
