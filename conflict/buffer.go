package conflict

import (
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func ApplyBuffer(s slot.TimeSlot, before, after time.Duration) slot.TimeSlot {
	return slot.TimeSlot{Start: s.Start.Add(-before), End: s.End.Add(after), Location: s.Location, Metadata: s.Metadata}
}

func RemoveBuffer(s slot.TimeSlot, before, after time.Duration) slot.TimeSlot {
	return slot.TimeSlot{Start: s.Start.Add(before), End: s.End.Add(-after), Location: s.Location, Metadata: s.Metadata}
}
