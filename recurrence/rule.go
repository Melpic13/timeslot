package recurrence

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Weekday mirrors time.Weekday.
type Weekday = time.Weekday

// Rule represents a recurrence rule (subset of RFC 5545 RRULE).
type Rule struct {
	Frequency  Frequency
	Interval   int
	Count      int
	Until      time.Time
	ByDay      []Weekday
	ByMonth    []time.Month
	ByMonthDay []int
	ByHour     []int
	ByMinute   []int
	WeekStart  time.Weekday
	Location   *time.Location
}

type Frequency int

const (
	Daily Frequency = iota
	Weekly
	Monthly
	Yearly
)

func ParseRule(rrule string) (*Rule, error) {
	return parseRule(rrule)
}

func (r *Rule) String() string {
	if r == nil {
		return ""
	}
	parts := []string{"FREQ=" + r.Frequency.String()}
	if r.Interval > 1 {
		parts = append(parts, "INTERVAL="+strconv.Itoa(r.Interval))
	}
	if r.Count > 0 {
		parts = append(parts, "COUNT="+strconv.Itoa(r.Count))
	}
	if !r.Until.IsZero() {
		parts = append(parts, "UNTIL="+r.Until.UTC().Format("20060102T150405Z"))
	}
	if len(r.ByDay) > 0 {
		day := make([]string, 0, len(r.ByDay))
		for _, d := range r.ByDay {
			day = append(day, weekdayToToken(d))
		}
		parts = append(parts, "BYDAY="+strings.Join(day, ","))
	}
	if len(r.ByMonth) > 0 {
		months := make([]string, 0, len(r.ByMonth))
		for _, m := range r.ByMonth {
			months = append(months, strconv.Itoa(int(m)))
		}
		parts = append(parts, "BYMONTH="+strings.Join(months, ","))
	}
	if len(r.ByMonthDay) > 0 {
		days := make([]string, 0, len(r.ByMonthDay))
		for _, d := range r.ByMonthDay {
			days = append(days, strconv.Itoa(d))
		}
		parts = append(parts, "BYMONTHDAY="+strings.Join(days, ","))
	}
	if len(r.ByHour) > 0 {
		h := make([]string, 0, len(r.ByHour))
		for _, n := range r.ByHour {
			h = append(h, strconv.Itoa(n))
		}
		parts = append(parts, "BYHOUR="+strings.Join(h, ","))
	}
	if len(r.ByMinute) > 0 {
		m := make([]string, 0, len(r.ByMinute))
		for _, n := range r.ByMinute {
			m = append(m, strconv.Itoa(n))
		}
		parts = append(parts, "BYMINUTE="+strings.Join(m, ","))
	}
	return strings.Join(parts, ";")
}

func (r *Rule) Generate(start time.Time, limit int) []time.Time {
	if r == nil {
		return nil
	}
	if limit <= 0 {
		limit = 1000
	}
	if r.Count > 0 && r.Count < limit {
		limit = r.Count
	}
	out := make([]time.Time, 0, limit)
	cur := start
	for len(out) < limit {
		if !r.Until.IsZero() && cur.After(r.Until) {
			break
		}
		if r.matches(cur) {
			out = append(out, cur)
		}
		next := r.step(cur)
		if !next.After(cur) {
			break
		}
		cur = next
	}
	return out
}

func (r *Rule) GenerateBetween(start, from, to time.Time) []time.Time {
	if !to.After(from) {
		return nil
	}
	all := r.Generate(start, 100000)
	out := make([]time.Time, 0)
	for _, t := range all {
		if (t.Equal(from) || t.After(from)) && t.Before(to) {
			out = append(out, t)
		}
		if t.After(to) {
			break
		}
	}
	return out
}

func (r *Rule) Next(after time.Time) (time.Time, bool) {
	probe := after
	for i := 0; i < 100000; i++ {
		probe = r.step(probe)
		if !r.Until.IsZero() && probe.After(r.Until) {
			return time.Time{}, false
		}
		if r.matches(probe) {
			return probe, true
		}
	}
	return time.Time{}, false
}

func (r *Rule) Contains(t time.Time) bool {
	if r == nil {
		return false
	}
	if !r.Until.IsZero() && t.After(r.Until) {
		return false
	}
	return r.matches(t)
}

func (r *Rule) Validate() error {
	if r == nil {
		return fmt.Errorf("recurrence: nil rule")
	}
	if r.Interval < 0 {
		return fmt.Errorf("recurrence: interval must be >= 0")
	}
	return nil
}

func (r *Rule) step(t time.Time) time.Time {
	interval := r.Interval
	if interval <= 0 {
		interval = 1
	}
	switch r.Frequency {
	case Daily:
		return t.AddDate(0, 0, interval)
	case Weekly:
		return t.AddDate(0, 0, 7*interval)
	case Monthly:
		return t.AddDate(0, interval, 0)
	case Yearly:
		return t.AddDate(interval, 0, 0)
	default:
		return t.AddDate(0, 0, interval)
	}
}

func (r *Rule) matches(t time.Time) bool {
	if len(r.ByDay) > 0 {
		ok := false
		for _, d := range r.ByDay {
			if t.Weekday() == d {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if len(r.ByMonth) > 0 {
		ok := false
		for _, m := range r.ByMonth {
			if t.Month() == m {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if len(r.ByMonthDay) > 0 {
		ok := false
		for _, d := range r.ByMonthDay {
			if t.Day() == d {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if len(r.ByHour) > 0 {
		ok := false
		for _, h := range r.ByHour {
			if t.Hour() == h {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if len(r.ByMinute) > 0 {
		ok := false
		for _, m := range r.ByMinute {
			if t.Minute() == m {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (f Frequency) String() string {
	switch f {
	case Daily:
		return "DAILY"
	case Weekly:
		return "WEEKLY"
	case Monthly:
		return "MONTHLY"
	case Yearly:
		return "YEARLY"
	default:
		return "DAILY"
	}
}

func parseRule(input string) (*Rule, error) {
	r := &Rule{Frequency: Daily, Interval: 1, WeekStart: time.Monday, Location: time.UTC}
	for _, part := range strings.Split(strings.TrimSpace(input), ";") {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("recurrence: malformed token %q", part)
		}
		key := strings.ToUpper(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])
		switch key {
		case "FREQ":
			switch strings.ToUpper(val) {
			case "DAILY":
				r.Frequency = Daily
			case "WEEKLY":
				r.Frequency = Weekly
			case "MONTHLY":
				r.Frequency = Monthly
			case "YEARLY":
				r.Frequency = Yearly
			default:
				return nil, fmt.Errorf("recurrence: unsupported frequency %q", val)
			}
		case "INTERVAL":
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			r.Interval = n
		case "COUNT":
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			r.Count = n
		case "UNTIL":
			t, err := parseDateTime(val, time.UTC)
			if err != nil {
				return nil, err
			}
			r.Until = t
		case "BYDAY":
			r.ByDay = nil
			for _, tok := range strings.Split(val, ",") {
				d, err := tokenToWeekday(tok)
				if err != nil {
					return nil, err
				}
				r.ByDay = append(r.ByDay, d)
			}
		case "BYMONTH":
			vals, err := parseInts(val)
			if err != nil {
				return nil, err
			}
			r.ByMonth = nil
			for _, n := range vals {
				r.ByMonth = append(r.ByMonth, time.Month(n))
			}
		case "BYMONTHDAY":
			vals, err := parseInts(val)
			if err != nil {
				return nil, err
			}
			r.ByMonthDay = vals
		case "BYHOUR":
			vals, err := parseInts(val)
			if err != nil {
				return nil, err
			}
			r.ByHour = vals
		case "BYMINUTE":
			vals, err := parseInts(val)
			if err != nil {
				return nil, err
			}
			r.ByMinute = vals
		}
	}
	sort.Slice(r.ByMonthDay, func(i, j int) bool { return r.ByMonthDay[i] < r.ByMonthDay[j] })
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return r, nil
}

func parseInts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

func parseDateTime(v string, loc *time.Location) (time.Time, error) {
	layouts := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	}
	for _, layout := range layouts {
		switch layout {
		case "20060102T150405Z":
			if strings.HasSuffix(v, "Z") {
				if t, err := time.Parse(layout, v); err == nil {
					return t, nil
				}
			}
		case "20060102T150405":
			if strings.Contains(v, "T") && !strings.HasSuffix(v, "Z") {
				if t, err := time.ParseInLocation(layout, v, loc); err == nil {
					return t, nil
				}
			}
		case "20060102":
			if !strings.Contains(v, "T") {
				if t, err := time.ParseInLocation(layout, v, loc); err == nil {
					return t, nil
				}
			}
		}
	}
	return time.Time{}, fmt.Errorf("recurrence: invalid datetime %q", v)
}

func tokenToWeekday(tok string) (Weekday, error) {
	switch strings.ToUpper(strings.TrimSpace(tok)) {
	case "MO":
		return time.Monday, nil
	case "TU":
		return time.Tuesday, nil
	case "WE":
		return time.Wednesday, nil
	case "TH":
		return time.Thursday, nil
	case "FR":
		return time.Friday, nil
	case "SA":
		return time.Saturday, nil
	case "SU":
		return time.Sunday, nil
	default:
		return time.Sunday, fmt.Errorf("recurrence: invalid weekday %q", tok)
	}
}

func weekdayToToken(day Weekday) string {
	switch day {
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	case time.Sunday:
		return "SU"
	default:
		return "MO"
	}
}
