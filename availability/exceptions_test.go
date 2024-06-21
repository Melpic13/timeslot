package availability

import (
	"testing"
	"time"
)

func TestExceptionSetBlockedDates(t *testing.T) {
	day := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	e := ExceptionSet{}.AddBlockedDates(day)
	if !e.IsBlocked(day) {
		t.Fatalf("expected blocked")
	}
}
