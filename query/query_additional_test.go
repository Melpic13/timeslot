package query

import (
	"testing"
	"time"
)

func TestQueryBuilderMethodsAndBuildDefaults(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.UTC
	}
	builder := NewQuery().
		Duration(30*time.Minute).
		From(time.Date(2025, 1, 1, 8, 0, 0, 0, time.UTC)).
		To(time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)).
		OnWeekdays(time.Monday, time.Tuesday).
		OnlyMornings().
		OnlyAfternoons().
		OnlyEvenings().
		NotBefore(NewTimeOfDay(8, 0, 0)).
		NotAfter(NewTimeOfDay(19, 0, 0)).
		WithConstraint(nil).
		WithConstraint(NewWeekdayConstraint(time.Monday)).
		PreferEarlier().
		PreferLater().
		PreferTime(NewTimeOfDay(10, 0, 0)).
		Limit(7).
		InTimezone(loc)

	q := builder.Build()
	if q.Limit != 7 {
		t.Fatalf("limit not set")
	}
	if q.Location == nil {
		t.Fatalf("location should be set")
	}
	if len(q.Constraints) < 1 || len(q.Preferences) < 1 {
		t.Fatalf("expected constraints and preferences")
	}
	if err := builder.Validate(); err != nil {
		t.Fatalf("builder validate failed: %v", err)
	}

	withNext := NewQuery().Duration(time.Hour).InNext(2 * time.Hour).Build()
	if withNext.To.Sub(withNext.From) < time.Hour {
		t.Fatalf("in-next should set a future range")
	}

	defaultBuild := (&QueryBuilder{query: Query{Duration: time.Hour}}).Build()
	if defaultBuild.From.IsZero() || defaultBuild.To.IsZero() {
		t.Fatalf("expected default from/to")
	}

	qb := &QueryBuilder{query: Query{}}
	if qb.locationOrUTC() != time.UTC {
		t.Fatalf("expected UTC fallback")
	}
}

func TestQueryValidationErrors(t *testing.T) {
	badDuration := Query{Duration: 0, From: time.Now(), To: time.Now().Add(time.Hour)}
	if err := badDuration.Validate(); err == nil {
		t.Fatalf("expected bad duration")
	}
	now := time.Now()
	badRange := Query{Duration: time.Hour, From: now, To: now}
	if err := badRange.Validate(); err == nil {
		t.Fatalf("expected bad range")
	}
}
