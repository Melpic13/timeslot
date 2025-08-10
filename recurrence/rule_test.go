package recurrence

import (
	"testing"
	"time"
)

func TestParseRule(t *testing.T) {
	r, err := ParseRule("FREQ=DAILY;INTERVAL=2;COUNT=3")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	start := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	got := r.Generate(start, 10)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func FuzzParseRule(f *testing.F) {
	f.Add("FREQ=DAILY;INTERVAL=1")
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = ParseRule(input)
	})
}
