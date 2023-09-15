package slot

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeSlotOverlaps(t *testing.T) {
	loc := time.UTC
	s1 := TimeSlot{Start: time.Date(2024, 1, 1, 9, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 10, 0, 0, 0, loc), Location: loc}
	s2 := TimeSlot{Start: time.Date(2024, 1, 1, 9, 30, 0, 0, loc), End: time.Date(2024, 1, 1, 10, 30, 0, 0, loc), Location: loc}
	if !s1.Overlaps(s2) {
		t.Fatalf("expected overlap")
	}
}

func TestTimeSlotUnionDisjoint(t *testing.T) {
	loc := time.UTC
	a := TimeSlot{Start: time.Date(2024, 1, 1, 9, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 10, 0, 0, 0, loc), Location: loc}
	b := TimeSlot{Start: time.Date(2024, 1, 1, 11, 0, 0, 0, loc), End: time.Date(2024, 1, 1, 12, 0, 0, 0, loc), Location: loc}
	if _, err := a.Union(b); err == nil {
		t.Fatalf("expected error")
	}
}

func TestTimeSlotJSONRoundTrip(t *testing.T) {
	orig := TimeSlot{Start: time.Now().UTC().Truncate(time.Second), End: time.Now().UTC().Add(time.Hour).Truncate(time.Second), Location: time.UTC, Metadata: map[string]any{"kind": "test"}}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got TimeSlot
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !orig.Start.Equal(got.Start) || !orig.End.Equal(got.End) {
		t.Fatalf("round trip mismatch")
	}
}
