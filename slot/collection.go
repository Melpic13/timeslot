package slot

import (
	"sort"
	"time"
)

// SlotCollection is an immutable, sorted collection of non-overlapping slots.
type SlotCollection struct {
	slots    []TimeSlot
	location *time.Location
}

func NewCollection(slots ...TimeSlot) SlotCollection {
	c := SlotCollection{slots: append([]TimeSlot(nil), slots...)}
	if len(c.slots) > 0 {
		c.location = c.slots[0].locationOrUTC()
	}
	return c.Merge()
}

func (c SlotCollection) Add(slots ...TimeSlot) SlotCollection {
	combined := append(c.Slots(), slots...)
	return NewCollection(combined...)
}

func (c SlotCollection) Remove(slots ...TimeSlot) SlotCollection {
	out := c
	for _, r := range slots {
		out = out.Subtract(NewCollection(r))
	}
	return out
}

func (c SlotCollection) Merge() SlotCollection {
	if len(c.slots) == 0 {
		return SlotCollection{location: c.location}
	}
	sorted := Sort(c.slots)
	merged := make([]TimeSlot, 0, len(sorted))
	cur := sorted[0]
	for i := 1; i < len(sorted); i++ {
		next := sorted[i]
		if cur.Overlaps(next) || cur.End.Equal(next.Start) {
			u, _ := cur.Union(next)
			cur = u
			continue
		}
		merged = append(merged, cur)
		cur = next
	}
	merged = append(merged, cur)
	return SlotCollection{slots: merged, location: merged[0].locationOrUTC()}
}

func (c SlotCollection) Subtract(other SlotCollection) SlotCollection {
	remaining := c.Slots()
	for _, cut := range other.slots {
		next := make([]TimeSlot, 0, len(remaining))
		for _, s := range remaining {
			if !s.Overlaps(cut) {
				next = append(next, s)
				continue
			}
			if cut.Start.After(s.Start) {
				next = append(next, TimeSlot{Start: s.Start, End: cut.Start, Location: s.locationOrUTC()})
			}
			if cut.End.Before(s.End) {
				next = append(next, TimeSlot{Start: cut.End, End: s.End, Location: s.locationOrUTC()})
			}
		}
		remaining = next
	}
	return NewCollection(remaining...)
}

func (c SlotCollection) Intersect(other SlotCollection) SlotCollection {
	out := make([]TimeSlot, 0)
	i, j := 0, 0
	a := c.Merge().slots
	b := other.Merge().slots
	for i < len(a) && j < len(b) {
		if inter, ok := a[i].Intersection(b[j]); ok {
			out = append(out, inter)
		}
		if a[i].End.Before(b[j].End) {
			i++
		} else {
			j++
		}
	}
	return NewCollection(out...)
}

func (c SlotCollection) Union(other SlotCollection) SlotCollection {
	return c.Add(other.slots...)
}

func (c SlotCollection) Filter(fn func(TimeSlot) bool) SlotCollection {
	out := make([]TimeSlot, 0, len(c.slots))
	for _, s := range c.slots {
		if fn(s) {
			out = append(out, s)
		}
	}
	return SlotCollection{slots: out, location: c.location}
}

func (c SlotCollection) FindOverlaps(slot TimeSlot) []TimeSlot {
	if len(c.slots) == 0 {
		return nil
	}
	idx := sort.Search(len(c.slots), func(i int) bool {
		return !c.slots[i].End.Before(slot.Start)
	})
	out := make([]TimeSlot, 0)
	for i := idx; i < len(c.slots); i++ {
		if !c.slots[i].Start.Before(slot.End) {
			break
		}
		if c.slots[i].Overlaps(slot) {
			out = append(out, c.slots[i])
		}
	}
	return out
}

func (c SlotCollection) TotalDuration() time.Duration {
	var d time.Duration
	for _, s := range c.slots {
		d += s.Duration()
	}
	return d
}

func (c SlotCollection) Gaps(within TimeSlot) SlotCollection {
	base := NewCollection(within)
	covered := base.Intersect(c)
	return base.Subtract(covered)
}

func (c SlotCollection) Slots() []TimeSlot {
	out := make([]TimeSlot, len(c.slots))
	copy(out, c.slots)
	return out
}

func (c SlotCollection) Len() int {
	return len(c.slots)
}

func (c SlotCollection) IsEmpty() bool {
	return len(c.slots) == 0
}

func (c SlotCollection) First() (TimeSlot, bool) {
	if len(c.slots) == 0 {
		return TimeSlot{}, false
	}
	return c.slots[0], true
}

func (c SlotCollection) Last() (TimeSlot, bool) {
	if len(c.slots) == 0 {
		return TimeSlot{}, false
	}
	return c.slots[len(c.slots)-1], true
}
