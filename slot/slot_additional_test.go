package slot

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeSlot(startHour, endHour int, loc *time.Location) TimeSlot {
	return TimeSlot{
		Start:    time.Date(2025, 1, 6, startHour, 0, 0, 0, loc),
		End:      time.Date(2025, 1, 6, endHour, 0, 0, 0, loc),
		Location: loc,
	}
}

func TestSlotCoreMethods(t *testing.T) {
	loc := time.UTC
	s := makeSlot(9, 11, loc)
	if got := s.Duration(); got != 2*time.Hour {
		t.Fatalf("duration got %v", got)
	}
	if !s.Contains(s.Start) {
		t.Fatalf("expected start boundary included")
	}
	if s.Contains(s.End) {
		t.Fatalf("expected end boundary excluded")
	}

	other := makeSlot(10, 12, loc)
	inter, ok := s.Intersection(other)
	if !ok {
		t.Fatalf("expected intersection")
	}
	if inter.Duration() != time.Hour {
		t.Fatalf("expected 1 hour intersection")
	}

	shifted := s.Shift(30 * time.Minute)
	if !shifted.Start.Equal(s.Start.Add(30 * time.Minute)) {
		t.Fatalf("shift failed")
	}

	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone DB unavailable")
	}
	inNY := s.InTimezone(ny)
	if inNY.Location.String() != ny.String() {
		t.Fatalf("timezone conversion failed")
	}
	if s.String() == "" || !strings.Contains(s.String(), "T") {
		t.Fatalf("string format unexpected: %s", s.String())
	}
}

func TestSlotNewValidateAndUnion(t *testing.T) {
	loc := time.UTC
	if _, err := New(time.Time{}, time.Now()); err == nil {
		t.Fatalf("expected error for zero start")
	}
	_, err := New(
		time.Date(2025, 1, 1, 12, 0, 0, 0, loc),
		time.Date(2025, 1, 1, 11, 0, 0, 0, loc),
	)
	if err == nil {
		t.Fatalf("expected invalid range")
	}

	a := makeSlot(9, 10, loc)
	b := makeSlot(10, 11, loc)
	u, err := a.Union(b)
	if err != nil {
		t.Fatalf("expected union for adjacent slots: %v", err)
	}
	if u.Duration() != 2*time.Hour {
		t.Fatalf("union duration wrong: %v", u.Duration())
	}

	c := makeSlot(12, 13, loc)
	if _, ok := a.Intersection(c); ok {
		t.Fatalf("expected no intersection for disjoint slots")
	}
	if _, err := a.Union(c); err == nil {
		t.Fatalf("expected disjoint union error")
	}
}

func TestSplitEqualAndHelpers(t *testing.T) {
	loc := time.UTC
	s := makeSlot(9, 11, loc)
	s.Metadata = map[string]any{"k": "v"}
	parts := s.Split(30 * time.Minute)
	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d", len(parts))
	}
	if parts[0].Metadata["k"] != "v" {
		t.Fatalf("expected metadata clone")
	}
	parts[0].Metadata["k"] = "changed"
	if s.Metadata["k"] != "v" {
		t.Fatalf("metadata should not mutate original")
	}

	if got := s.Split(0); got != nil {
		t.Fatalf("expected nil split for zero duration")
	}
	if got := (TimeSlot{}).Split(time.Minute); got != nil {
		t.Fatalf("expected nil split for zero slot")
	}

	s2 := makeSlot(9, 11, loc)
	s2.Metadata = map[string]any{"k": "v"}
	if !s.Equal(s2) {
		t.Fatalf("expected equal")
	}
	s2.Metadata["k"] = "x"
	if s.Equal(s2) {
		t.Fatalf("expected metadata mismatch")
	}
	if mapsEqual(map[string]any{"a": 1}, map[string]any{"a": 2}) {
		t.Fatalf("expected value mismatch")
	}
	if mapsEqual(map[string]any{"a": 1}, map[string]any{"a": 1, "b": 2}) {
		t.Fatalf("expected len mismatch")
	}
	if !mapsEqual(map[string]any{"a": 1}, map[string]any{"a": 1}) {
		t.Fatalf("expected equal map")
	}
}

func TestSlotJSONPathsAndSort(t *testing.T) {
	raw := []byte(`{"start":"2025-01-01T09:00:00Z","end":"2025-01-01T10:00:00Z","location":"UTC","metadata":{"a":"b"}}`)
	var s TimeSlot
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if s.locationOrUTC().String() != "UTC" {
		t.Fatalf("unexpected location")
	}

	badLoc := []byte(`{"start":"2025-01-01T09:00:00Z","end":"2025-01-01T10:00:00Z","location":"Invalid/Zone"}`)
	if err := json.Unmarshal(badLoc, &s); err == nil {
		t.Fatalf("expected invalid location error")
	}

	badRange := []byte(`{"start":"2025-01-01T10:00:00Z","end":"2025-01-01T09:00:00Z","location":"UTC"}`)
	if err := json.Unmarshal(badRange, &s); err == nil {
		t.Fatalf("expected invalid range error")
	}

	b, err := json.Marshal(makeSlot(9, 10, time.UTC))
	if err != nil || !strings.Contains(string(b), "location") {
		t.Fatalf("marshal failed: %v %s", err, string(b))
	}

	sorted := Sort([]TimeSlot{makeSlot(11, 12, time.UTC), makeSlot(9, 10, time.UTC), makeSlot(9, 9+1, time.UTC)})
	if !sorted[0].Start.Equal(makeSlot(9, 10, time.UTC).Start) {
		t.Fatalf("sort failed")
	}

	if minTime(time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC), time.Date(2025, 1, 1, 2, 0, 0, 0, time.UTC)).Hour() != 1 {
		t.Fatalf("minTime failed")
	}
	if maxTime(time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC), time.Date(2025, 1, 1, 2, 0, 0, 0, time.UTC)).Hour() != 2 {
		t.Fatalf("maxTime failed")
	}

	z := TimeSlot{}
	if z.locationOrUTC() == nil {
		t.Fatalf("locationOrUTC should never return nil")
	}
	if cloneMetadata(nil) != nil {
		t.Fatalf("nil metadata clone should be nil")
	}
}
