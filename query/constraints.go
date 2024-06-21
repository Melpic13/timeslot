package query

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

// TimeOfDay is a query-specific time without date.
type TimeOfDay struct {
	Hour   int
	Minute int
	Second int
}

func NewTimeOfDay(hour, minute, second int) TimeOfDay {
	return TimeOfDay{Hour: hour, Minute: minute, Second: second}
}

func (t TimeOfDay) Validate() error {
	if t.Hour < 0 || t.Hour > 23 {
		return fmt.Errorf("query: invalid hour")
	}
	if t.Minute < 0 || t.Minute > 59 {
		return fmt.Errorf("query: invalid minute")
	}
	if t.Second < 0 || t.Second > 59 {
		return fmt.Errorf("query: invalid second")
	}
	return nil
}

func (t TimeOfDay) toTime(day time.Time, loc *time.Location) time.Time {
	if loc == nil {
		loc = day.Location()
	}
	d := day.In(loc)
	y, m, dd := d.Date()
	return time.Date(y, m, dd, t.Hour, t.Minute, t.Second, 0, loc)
}

// Constraint interface for extensibility.
type Constraint interface {
	IsSatisfied(slot.TimeSlot) bool
	String() string
}

// Preference scores candidate slots.
type Preference interface {
	Score(slot.TimeSlot) int
	String() string
}

// WeekdayConstraint allows only specific weekdays.
type WeekdayConstraint struct {
	Days map[time.Weekday]struct{}
}

func NewWeekdayConstraint(days ...time.Weekday) WeekdayConstraint {
	m := make(map[time.Weekday]struct{}, len(days))
	for _, d := range days {
		m[d] = struct{}{}
	}
	return WeekdayConstraint{Days: m}
}

func (c WeekdayConstraint) IsSatisfied(s slot.TimeSlot) bool {
	_, ok := c.Days[s.Start.Weekday()]
	return ok
}

func (c WeekdayConstraint) String() string { return "weekday" }

// TimeOfDayConstraint constrains slot start between Start and End.
type TimeOfDayConstraint struct {
	Start *TimeOfDay
	End   *TimeOfDay
}

func (c TimeOfDayConstraint) IsSatisfied(s slot.TimeSlot) bool {
	loc := s.Start.Location()
	if c.Start != nil {
		min := c.Start.toTime(s.Start, loc)
		if s.Start.Before(min) {
			return false
		}
	}
	if c.End != nil {
		max := c.End.toTime(s.Start, loc)
		if s.Start.After(max) {
			return false
		}
	}
	return true
}

func (c TimeOfDayConstraint) String() string { return "time-of-day" }

type NotBeforeConstraint struct{ Time TimeOfDay }

func (c NotBeforeConstraint) IsSatisfied(s slot.TimeSlot) bool {
	limit := c.Time.toTime(s.Start, s.Start.Location())
	return !s.Start.Before(limit)
}

func (c NotBeforeConstraint) String() string { return "not-before" }

type NotAfterConstraint struct{ Time TimeOfDay }

func (c NotAfterConstraint) IsSatisfied(s slot.TimeSlot) bool {
	limit := c.Time.toTime(s.Start, s.Start.Location())
	return !s.Start.After(limit)
}

func (c NotAfterConstraint) String() string { return "not-after" }

type MinGapConstraint struct {
	Gap      time.Duration
	Bookings []slot.TimeSlot
}

func (c MinGapConstraint) IsSatisfied(candidate slot.TimeSlot) bool {
	for _, b := range c.Bookings {
		expanded := slot.TimeSlot{Start: b.Start.Add(-c.Gap), End: b.End.Add(c.Gap), Location: b.Location}
		if expanded.Overlaps(candidate) {
			return false
		}
	}
	return true
}

func (c MinGapConstraint) String() string { return "min-gap" }

type preferEarlier struct{}

func (p preferEarlier) Score(s slot.TimeSlot) int { return -int(s.Start.Unix()) }
func (p preferEarlier) String() string            { return "prefer-earlier" }

type preferLater struct{}

func (p preferLater) Score(s slot.TimeSlot) int { return int(s.Start.Unix()) }
func (p preferLater) String() string            { return "prefer-later" }

type preferTime struct{ Target TimeOfDay }

func (p preferTime) Score(s slot.TimeSlot) int {
	target := p.Target.toTime(s.Start, s.Start.Location())
	delta := s.Start.Sub(target)
	if delta < 0 {
		delta = -delta
	}
	return -int(delta / time.Second)
}

func (p preferTime) String() string { return "prefer-time" }
