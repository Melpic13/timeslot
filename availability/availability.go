package availability

import (
	"time"

	"github.com/Melpic13/timeslot/internal/timeutil"
	"github.com/Melpic13/timeslot/slot"
)

// Availability combines weekly schedule with exceptions.
type Availability struct {
	Weekly     WeeklySchedule
	Exceptions ExceptionSet
	Bookings   slot.SlotCollection
	Location   *time.Location
}

func New(loc *time.Location) Availability {
	if loc == nil {
		loc = time.UTC
	}
	return Availability{Weekly: NewWeeklySchedule(loc), Location: loc, Bookings: slot.NewCollection()}
}

func (a Availability) AddBlockedDates(dates ...time.Time) Availability {
	a.Exceptions = a.Exceptions.AddBlockedDates(dates...)
	return a
}

func (a Availability) AddBlockedRange(start, end time.Time) Availability {
	a.Exceptions = a.Exceptions.AddBlockedRange(start, end)
	return a
}

func (a Availability) AddAvailableOverride(date time.Time, ranges ...TimeRange) Availability {
	a.Exceptions = a.Exceptions.AddAvailableOverride(date, ranges...)
	return a
}

func (a Availability) AddBooking(s slot.TimeSlot) Availability {
	a.Bookings = a.Bookings.Add(s)
	return a
}

func (a Availability) RemoveBooking(s slot.TimeSlot) Availability {
	a.Bookings = a.Bookings.Remove(s)
	return a
}

func (a Availability) GetSlots(from, to time.Time) slot.SlotCollection {
	if !to.After(from) {
		return slot.NewCollection()
	}
	loc := a.locationOrUTC()
	from = from.In(loc)
	to = to.In(loc)

	var generated []slot.TimeSlot
	for day := timeutil.StartOfDay(from, loc); day.Before(to) || day.Equal(to); day = day.Add(24 * time.Hour) {
		dayStart := timeutil.StartOfDay(day, loc)
		dayEnd := dayStart.Add(24 * time.Hour)
		if dayEnd.Before(from) || dayStart.After(to) {
			continue
		}

		ranges, modified := a.Exceptions.ModifiedForDate(dayStart)
		if modified {
			for _, r := range ranges {
				start := r.Start.ToTime(dayStart, loc)
				end := r.End.ToTime(dayStart, loc)
				start = timeutil.Clamp(start, from, to)
				end = timeutil.Clamp(end, from, to)
				if end.After(start) {
					generated = append(generated, slot.TimeSlot{Start: start, End: end, Location: loc})
				}
			}
			continue
		}

		weeklyRanges := a.Weekly.GetDay(dayStart.Weekday())
		for _, r := range weeklyRanges {
			start := r.Start.ToTime(dayStart, loc)
			end := r.End.ToTime(dayStart, loc)
			start = timeutil.Clamp(start, from, to)
			end = timeutil.Clamp(end, from, to)
			if end.After(start) {
				generated = append(generated, slot.TimeSlot{Start: start, End: end, Location: loc})
			}
		}
	}

	base := slot.NewCollection(generated...)

	for _, blocked := range a.Exceptions.Blocked {
		base = base.Remove(slot.TimeSlot{Start: blocked.Start, End: blocked.End, Location: loc})
	}

	if len(a.Exceptions.Available) > 0 {
		var adds []slot.TimeSlot
		for _, avail := range a.Exceptions.Available {
			s := timeutil.Clamp(avail.Start.In(loc), from, to)
			e := timeutil.Clamp(avail.End.In(loc), from, to)
			if e.After(s) {
				adds = append(adds, slot.TimeSlot{Start: s, End: e, Location: loc})
			}
		}
		base = base.Add(adds...)
	}

	base = base.Subtract(a.Bookings)
	return base.Merge()
}

func (a Availability) IsAvailable(t time.Time) bool {
	probe := slot.TimeSlot{Start: t, End: t.Add(time.Second), Location: a.locationOrUTC()}
	return len(a.GetSlots(t.Add(-24*time.Hour), t.Add(24*time.Hour)).FindOverlaps(probe)) > 0
}

func (a Availability) IsBooked(t time.Time) bool {
	probe := slot.TimeSlot{Start: t, End: t.Add(time.Second), Location: a.locationOrUTC()}
	return len(a.Bookings.FindOverlaps(probe)) > 0
}

func (a Availability) FindAvailableSlots(duration time.Duration, from, to time.Time) []slot.TimeSlot {
	if duration <= 0 {
		return nil
	}
	free := a.GetSlots(from, to)
	var out []slot.TimeSlot
	for _, s := range free.Slots() {
		for cur := s.Start; cur.Add(duration).Before(s.End) || cur.Add(duration).Equal(s.End); cur = cur.Add(duration) {
			out = append(out, slot.TimeSlot{Start: cur, End: cur.Add(duration), Location: s.Location})
		}
	}
	return out
}

func (a Availability) Validate() error {
	if err := a.Weekly.Validate(); err != nil {
		return err
	}
	if err := a.Exceptions.Validate(); err != nil {
		return err
	}
	for _, s := range a.Bookings.Slots() {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (a Availability) locationOrUTC() *time.Location {
	if a.Location != nil {
		return a.Location
	}
	if a.Weekly.Location != nil {
		return a.Weekly.Location
	}
	return time.UTC
}
