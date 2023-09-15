package timezone

import (
	"testing"
	"time"
)

func TestIsDSTTransitionDay(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone data not available")
	}
	day := time.Date(2024, 3, 10, 12, 0, 0, 0, loc)
	if !IsDSTTransitionDay(day) {
		t.Fatalf("expected DST transition day")
	}
}
