package timeslot

import (
	"errors"
	"fmt"

	"github.com/Melpic13/timeslot/slot"
)

var (
	ErrInvalidTimeRange   = errors.New("timeslot: end time must be after start time")
	ErrInvalidDuration    = errors.New("timeslot: duration must be positive")
	ErrSlotOverlap        = errors.New("timeslot: slots overlap")
	ErrNoAvailability     = errors.New("timeslot: no availability found")
	ErrConflict           = errors.New("timeslot: booking conflict detected")
	ErrInvalidTimezone    = errors.New("timeslot: invalid timezone")
	ErrPastTime           = errors.New("timeslot: cannot book in the past")
	ErrInsufficientNotice = errors.New("timeslot: insufficient booking notice")
	ErrTooFarAdvance      = errors.New("timeslot: booking too far in advance")
	ErrInvalidRecurrence  = errors.New("timeslot: invalid recurrence rule")
	ErrInvalidICS         = errors.New("timeslot: invalid iCal format")
)

// SlotError wraps a slot-specific failure with operation context.
type SlotError struct {
	Op   string
	Slot slot.TimeSlot
	Err  error
}

func (e *SlotError) Error() string {
	if e == nil {
		return "timeslot: <nil>"
	}
	if e.Err == nil {
		return fmt.Sprintf("timeslot: %s failed for %s", e.Op, e.Slot.String())
	}
	return fmt.Sprintf("timeslot: %s failed for %s: %v", e.Op, e.Slot.String(), e.Err)
}

func (e *SlotError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
