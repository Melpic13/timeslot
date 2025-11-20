package validate

import (
	"testing"
	"time"
)

func TestValidateHelpers(t *testing.T) {
	now := time.Now()
	if err := TimeRange(now, now); err == nil {
		t.Fatalf("expected invalid time range")
	}
	if err := TimeRange(now, now.Add(time.Second)); err != nil {
		t.Fatalf("unexpected time range error: %v", err)
	}
	if err := PositiveDuration(0); err == nil {
		t.Fatalf("expected invalid duration")
	}
	if err := PositiveDuration(time.Second); err != nil {
		t.Fatalf("unexpected duration error: %v", err)
	}
	if LocationOrUTC(nil).String() != "UTC" {
		t.Fatalf("expected UTC fallback")
	}
}
