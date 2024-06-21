package availability

import (
	"fmt"
	"sort"
	"time"

	"github.com/Melpic13/timeslot/internal/timeutil"
	"github.com/Melpic13/timeslot/slot"
)

// TimeOfDay represents a time without a date.
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
		return fmt.Errorf("availability: invalid hour %d", t.Hour)
	}
	if t.Minute < 0 || t.Minute > 59 {
		return fmt.Errorf("availability: invalid minute %d", t.Minute)
	}
	if t.Second < 0 || t.Second > 59 {
		return fmt.Errorf("availability: invalid second %d", t.Second)
	}
	return nil
}

func (t TimeOfDay) ToTime(day time.Time, loc *time.Location) time.Time {
	l := loc
	if l == nil {
		l = day.Location()
	}
	d := day.In(l)
	y, m, dd := d.Date()
	return time.Date(y, m, dd, t.Hour, t.Minute, t.Second, 0, l)
}

// TimeRange represents a time window within a day (no date).
type TimeRange struct {
	Start TimeOfDay
	End   TimeOfDay
}

func (r TimeRange) Validate() error {
	if err := r.Start.Validate(); err != nil {
		return err
	}
	if err := r.End.Validate(); err != nil {
		return err
	}
	if !r.End.after(r.Start) {
		return fmt.Errorf("availability: invalid range %v-%v", r.Start, r.End)
	}
	return nil
}

func (t TimeOfDay) after(other TimeOfDay) bool {
	if t.Hour != other.Hour {
		return t.Hour > other.Hour
	}
	if t.Minute != other.Minute {
		return t.Minute > other.Minute
	}
	return t.Second > other.Second
}

// WeeklySchedule defines recurring weekly availability.
type WeeklySchedule struct {
	Monday    []TimeRange
	Tuesday   []TimeRange
	Wednesday []TimeRange
	Thursday  []TimeRange
	Friday    []TimeRange
	Saturday  []TimeRange
	Sunday    []TimeRange
	Location  *time.Location
}

func NewWeeklySchedule(loc *time.Location) WeeklySchedule {
	if loc == nil {
		loc = time.UTC
	}
	return WeeklySchedule{Location: loc}
}

func (w WeeklySchedule) SetDay(day time.Weekday, ranges ...TimeRange) WeeklySchedule {
	out := w
	clean := normalizeRanges(ranges)
	switch day {
	case time.Monday:
		out.Monday = clean
	case time.Tuesday:
		out.Tuesday = clean
	case time.Wednesday:
		out.Wednesday = clean
	case time.Thursday:
		out.Thursday = clean
	case time.Friday:
		out.Friday = clean
	case time.Saturday:
		out.Saturday = clean
	case time.Sunday:
		out.Sunday = clean
	}
	return out
}

func (w WeeklySchedule) GetDay(day time.Weekday) []TimeRange {
	switch day {
	case time.Monday:
		return append([]TimeRange(nil), w.Monday...)
	case time.Tuesday:
		return append([]TimeRange(nil), w.Tuesday...)
	case time.Wednesday:
		return append([]TimeRange(nil), w.Wednesday...)
	case time.Thursday:
		return append([]TimeRange(nil), w.Thursday...)
	case time.Friday:
		return append([]TimeRange(nil), w.Friday...)
	case time.Saturday:
		return append([]TimeRange(nil), w.Saturday...)
	case time.Sunday:
		return append([]TimeRange(nil), w.Sunday...)
	default:
		return nil
	}
}

func (w WeeklySchedule) GenerateSlots(from, to time.Time) slot.SlotCollection {
	if !to.After(from) {
		return slot.NewCollection()
	}
	loc := w.locationOrUTC()
	startDay := timeutil.StartOfDay(from, loc)
	endDay := timeutil.StartOfDay(to, loc)
	if endDay.Before(startDay) {
		return slot.NewCollection()
	}
	var out []slot.TimeSlot
	for d := startDay; !d.After(endDay); d = d.Add(24 * time.Hour) {
		ranges := w.GetDay(d.In(loc).Weekday())
		for _, r := range ranges {
			rs := r.Start.ToTime(d, loc)
			re := r.End.ToTime(d, loc)
			if !re.After(rs) {
				continue
			}
			s := rs
			if s.Before(from) {
				s = from
			}
			e := re
			if e.After(to) {
				e = to
			}
			if e.After(s) {
				out = append(out, slot.TimeSlot{Start: s, End: e, Location: loc})
			}
		}
	}
	return slot.NewCollection(out...)
}

func (w WeeklySchedule) IsAvailable(t time.Time) bool {
	loc := w.locationOrUTC()
	t = t.In(loc)
	for _, r := range w.GetDay(t.Weekday()) {
		start := r.Start.ToTime(t, loc)
		end := r.End.ToTime(t, loc)
		if (t.Equal(start) || t.After(start)) && t.Before(end) {
			return true
		}
	}
	return false
}

func (w WeeklySchedule) NextAvailable(after time.Time) (time.Time, bool) {
	loc := w.locationOrUTC()
	probe := after.In(loc)
	for i := 0; i < 14; i++ {
		day := timeutil.StartOfDay(probe, loc).Add(time.Duration(i) * 24 * time.Hour)
		ranges := w.GetDay(day.Weekday())
		for _, r := range ranges {
			start := r.Start.ToTime(day, loc)
			if start.After(after) {
				return start, true
			}
			end := r.End.ToTime(day, loc)
			if (after.Equal(start) || after.After(start)) && after.Before(end) {
				return after, true
			}
		}
	}
	return time.Time{}, false
}

func (w WeeklySchedule) MergeWith(other WeeklySchedule) WeeklySchedule {
	out := w
	for _, day := range []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	} {
		merged := append(w.GetDay(day), other.GetDay(day)...)
		out = out.SetDay(day, merged...)
	}
	if out.Location == nil {
		out.Location = other.Location
	}
	return out
}

func (w WeeklySchedule) Validate() error {
	for _, day := range []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	} {
		for _, r := range w.GetDay(day) {
			if err := r.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w WeeklySchedule) locationOrUTC() *time.Location {
	if w.Location == nil {
		return time.UTC
	}
	return w.Location
}

func normalizeRanges(ranges []TimeRange) []TimeRange {
	valid := make([]TimeRange, 0, len(ranges))
	for _, r := range ranges {
		if err := r.Validate(); err == nil {
			valid = append(valid, r)
		}
	}
	sort.Slice(valid, func(i, j int) bool {
		a := valid[i]
		b := valid[j]
		if a.Start.Hour != b.Start.Hour {
			return a.Start.Hour < b.Start.Hour
		}
		if a.Start.Minute != b.Start.Minute {
			return a.Start.Minute < b.Start.Minute
		}
		return a.Start.Second < b.Start.Second
	})
	if len(valid) == 0 {
		return nil
	}
	out := []TimeRange{valid[0]}
	for i := 1; i < len(valid); i++ {
		curr := out[len(out)-1]
		next := valid[i]
		if !next.Start.after(curr.End) {
			if next.End.after(curr.End) {
				curr.End = next.End
				out[len(out)-1] = curr
			}
			continue
		}
		out = append(out, next)
	}
	return out
}
