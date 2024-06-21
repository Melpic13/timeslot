package availability

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/internal/timeutil"
)

// DateRange is an inclusive start, exclusive end date-time range.
type DateRange struct {
	Start time.Time
	End   time.Time
}

func (r DateRange) Validate() error {
	if !r.End.After(r.Start) {
		return fmt.Errorf("availability: invalid date range")
	}
	return nil
}

func (r DateRange) Contains(t time.Time) bool {
	return (t.Equal(r.Start) || t.After(r.Start)) && t.Before(r.End)
}

// DateOverride sets custom ranges for a specific date.
type DateOverride struct {
	Date   time.Time
	Ranges []TimeRange
}

// ExceptionSet handles date-based exceptions.
type ExceptionSet struct {
	Blocked   []DateRange
	Available []DateRange
	Modified  []DateOverride
}

func (e ExceptionSet) AddBlockedDates(dates ...time.Time) ExceptionSet {
	out := e
	for _, d := range dates {
		start := timeutil.StartOfDay(d, d.Location())
		out.Blocked = append(out.Blocked, DateRange{Start: start, End: start.Add(24 * time.Hour)})
	}
	return out
}

func (e ExceptionSet) AddBlockedRange(start, end time.Time) ExceptionSet {
	out := e
	out.Blocked = append(out.Blocked, DateRange{Start: start, End: end})
	return out
}

func (e ExceptionSet) AddAvailableOverride(date time.Time, ranges ...TimeRange) ExceptionSet {
	out := e
	if len(ranges) == 0 {
		start := timeutil.StartOfDay(date, date.Location())
		out.Available = append(out.Available, DateRange{Start: start, End: start.Add(24 * time.Hour)})
		return out
	}
	out.Modified = append(out.Modified, DateOverride{Date: date, Ranges: normalizeRanges(ranges)})
	return out
}

func (e ExceptionSet) IsBlocked(t time.Time) bool {
	for _, r := range e.Blocked {
		if r.Contains(t) {
			return true
		}
	}
	return false
}

func (e ExceptionSet) HasAvailableOverride(t time.Time) bool {
	for _, r := range e.Available {
		if r.Contains(t) {
			return true
		}
	}
	return false
}

func (e ExceptionSet) ModifiedForDate(day time.Time) ([]TimeRange, bool) {
	for _, m := range e.Modified {
		if timeutil.SameDay(m.Date, day, day.Location()) {
			return append([]TimeRange(nil), m.Ranges...), true
		}
	}
	return nil, false
}

func (e ExceptionSet) Validate() error {
	for _, r := range e.Blocked {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	for _, r := range e.Available {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	for _, m := range e.Modified {
		for _, r := range m.Ranges {
			if err := r.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}
