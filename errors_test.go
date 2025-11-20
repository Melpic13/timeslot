package timeslot

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func TestSlotErrorFormattingAndUnwrap(t *testing.T) {
	err := &SlotError{}
	if !strings.Contains(err.Error(), "failed") {
		t.Fatalf("unexpected error format")
	}

	underlying := errors.New("boom")
	e := &SlotError{
		Op: "book",
		Slot: slot.TimeSlot{
			Start:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			End:      time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			Location: time.UTC,
		},
		Err: underlying,
	}
	if !strings.Contains(e.Error(), "book") || !strings.Contains(e.Error(), "boom") {
		t.Fatalf("unexpected wrapped error message: %s", e.Error())
	}
	if !errors.Is(e, underlying) {
		t.Fatalf("expected unwrap")
	}
	if (*SlotError)(nil).Unwrap() != nil {
		t.Fatalf("nil unwrap should be nil")
	}
}
