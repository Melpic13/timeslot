package conflict

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

type ResolutionStrategy int

const (
	StrategyShiftForward ResolutionStrategy = iota
	StrategyShiftBackward
	StrategySkip
)

type ResolutionOption struct {
	Strategy    ResolutionStrategy
	Description string
}

func resolveSlot(c Conflict, strategy ResolutionStrategy) (*slot.TimeSlot, error) {
	s := c.Slot
	d := s.Duration()
	switch strategy {
	case StrategyShiftForward:
		out := slot.TimeSlot{Start: s.End, End: s.End.Add(d), Location: s.Location}
		return &out, nil
	case StrategyShiftBackward:
		out := slot.TimeSlot{Start: s.Start.Add(-d), End: s.Start, Location: s.Location}
		return &out, nil
	case StrategySkip:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown resolution strategy")
	}
}

func defaultResolutionOptions(s slot.TimeSlot) []ResolutionOption {
	_ = s
	return []ResolutionOption{
		{Strategy: StrategyShiftForward, Description: "Shift to next slot"},
		{Strategy: StrategyShiftBackward, Description: "Shift to previous slot"},
		{Strategy: StrategySkip, Description: "Skip this slot"},
	}
}

func slotWindowAround(s slot.TimeSlot, d time.Duration) slot.TimeSlot {
	return slot.TimeSlot{Start: s.Start.Add(-d), End: s.End.Add(d), Location: s.Location}
}
