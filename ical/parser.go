package ical

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/recurrence"
	"github.com/Melpic13/timeslot/slot"
)

// Calendar represents a parsed iCal file.
type Calendar struct {
	Name     string
	Timezone *time.Location
	Events   []Event
}

type EventStatus int

const (
	EventStatusConfirmed EventStatus = iota
	EventStatusTentative
	EventStatusCancelled
)

type Event struct {
	UID        string
	Summary    string
	Start      time.Time
	End        time.Time
	Recurrence *recurrence.Rule
	Exceptions []time.Time
	Status     EventStatus
}

func Parse(r io.Reader) (*Calendar, error) {
	scanner := bufio.NewScanner(r)
	cal := &Calendar{Timezone: time.UTC}
	var current *Event
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch line {
		case "BEGIN:VEVENT":
			current = &Event{Status: EventStatusConfirmed}
			continue
		case "END:VEVENT":
			if current != nil {
				cal.Events = append(cal.Events, *current)
			}
			current = nil
			continue
		}
		key, params, value := parseICSLine(line)
		if current == nil {
			switch key {
			case "X-WR-CALNAME":
				cal.Name = value
			case "X-WR-TIMEZONE":
				loc, err := time.LoadLocation(value)
				if err == nil {
					cal.Timezone = loc
				}
			}
			continue
		}

		switch key {
		case "UID":
			current.UID = value
		case "SUMMARY":
			current.Summary = value
		case "DTSTART":
			loc := tzFromParams(params, cal.Timezone)
			t, err := parseDateTime(value, loc)
			if err != nil {
				return nil, err
			}
			current.Start = t
		case "DTEND":
			loc := tzFromParams(params, cal.Timezone)
			t, err := parseDateTime(value, loc)
			if err != nil {
				return nil, err
			}
			current.End = t
		case "RRULE":
			r, err := recurrence.ParseRule(value)
			if err != nil {
				return nil, err
			}
			current.Recurrence = r
		case "EXDATE":
			loc := tzFromParams(params, cal.Timezone)
			for _, v := range strings.Split(value, ",") {
				t, err := parseDateTime(strings.TrimSpace(v), loc)
				if err != nil {
					return nil, err
				}
				current.Exceptions = append(current.Exceptions, t)
			}
		case "STATUS":
			switch strings.ToUpper(value) {
			case "TENTATIVE":
				current.Status = EventStatusTentative
			case "CANCELLED":
				current.Status = EventStatusCancelled
			default:
				current.Status = EventStatusConfirmed
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cal, nil
}

func ParseFile(path string) (*Calendar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func (c *Calendar) ToAvailability() availability.Availability {
	loc := time.UTC
	if c != nil && c.Timezone != nil {
		loc = c.Timezone
	}
	a := availability.New(loc)
	for _, busy := range c.GetBusySlots(time.Time{}, time.Date(9999, 12, 31, 0, 0, 0, 0, loc)) {
		a = a.AddBooking(busy)
	}
	return a
}

func (c *Calendar) ToSlotCollection(from, to time.Time) slot.SlotCollection {
	return slot.NewCollection(c.GetBusySlots(from, to)...)
}

func (c *Calendar) GetBusySlots(from, to time.Time) []slot.TimeSlot {
	if c == nil {
		return nil
	}
	var out []slot.TimeSlot
	for _, e := range c.Events {
		if e.Status == EventStatusCancelled {
			continue
		}
		if e.Recurrence == nil {
			if overlapsTime(from, to, e.Start, e.End) {
				out = append(out, slot.TimeSlot{Start: maxTime(from, e.Start), End: minTime(to, e.End), Location: e.Start.Location()})
			}
			continue
		}
		occurrences := e.Recurrence.GenerateBetween(e.Start, from, to)
		d := e.End.Sub(e.Start)
		for _, occ := range occurrences {
			if isException(e.Exceptions, occ) {
				continue
			}
			end := occ.Add(d)
			if overlapsTime(from, to, occ, end) {
				out = append(out, slot.TimeSlot{Start: maxTime(from, occ), End: minTime(to, end), Location: occ.Location()})
			}
		}
	}
	return slot.NewCollection(out...).Slots()
}

func (c *Calendar) GetFreeSlots(from, to time.Time, within availability.WeeklySchedule) []slot.TimeSlot {
	all := within.GenerateSlots(from, to)
	busy := slot.NewCollection(c.GetBusySlots(from, to)...)
	return all.Subtract(busy).Slots()
}

func parseICSLine(line string) (key string, params map[string]string, value string) {
	params = map[string]string{}
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return strings.ToUpper(line), params, ""
	}
	head := parts[0]
	value = parts[1]
	headParts := strings.Split(head, ";")
	key = strings.ToUpper(strings.TrimSpace(headParts[0]))
	for _, p := range headParts[1:] {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			params[strings.ToUpper(strings.TrimSpace(kv[0]))] = strings.TrimSpace(kv[1])
		}
	}
	return key, params, strings.TrimSpace(value)
}

func tzFromParams(params map[string]string, fallback *time.Location) *time.Location {
	if v, ok := params["TZID"]; ok {
		if loc, err := time.LoadLocation(v); err == nil {
			return loc
		}
	}
	if fallback != nil {
		return fallback
	}
	return time.UTC
}

func parseDateTime(v string, loc *time.Location) (time.Time, error) {
	layouts := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	}
	for _, layout := range layouts {
		if layout == "20060102T150405Z" && strings.HasSuffix(v, "Z") {
			t, err := time.Parse(layout, v)
			if err == nil {
				return t, nil
			}
			continue
		}
		if layout == "20060102T150405" && strings.Contains(v, "T") && !strings.HasSuffix(v, "Z") {
			t, err := time.ParseInLocation(layout, v, loc)
			if err == nil {
				return t, nil
			}
			continue
		}
		if layout == "20060102" && !strings.Contains(v, "T") {
			t, err := time.ParseInLocation(layout, v, loc)
			if err == nil {
				return t, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("ical: unsupported datetime %q", v)
}

func overlapsTime(from, to, start, end time.Time) bool {
	if from.IsZero() && to.IsZero() {
		return true
	}
	if from.IsZero() {
		return start.Before(to)
	}
	if to.IsZero() {
		return end.After(from)
	}
	return start.Before(to) && from.Before(end)
}

func maxTime(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	if b.IsZero() {
		return a
	}
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	if b.IsZero() {
		return a
	}
	if a.Before(b) {
		return a
	}
	return b
}

func isException(ex []time.Time, t time.Time) bool {
	for _, e := range ex {
		if e.Equal(t) {
			return true
		}
	}
	return false
}
