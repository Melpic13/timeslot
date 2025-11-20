package ical

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/recurrence"
)

func TestParseStatusesAndHelpers(t *testing.T) {
	input := `BEGIN:VCALENDAR
X-WR-CALNAME:Team Calendar
X-WR-TIMEZONE:UTC
BEGIN:VEVENT
UID:evt1
DTSTART:20250101T090000Z
DTEND:20250101T100000Z
STATUS:TENTATIVE
END:VEVENT
BEGIN:VEVENT
UID:evt2
DTSTART;TZID=UTC:20250102T090000
DTEND;TZID=UTC:20250102T100000
STATUS:CANCELLED
END:VEVENT
BEGIN:VEVENT
UID:evt3
DTSTART:20250103
DTEND:20250104
RRULE:FREQ=DAILY;COUNT=3
EXDATE:20250104
STATUS:CONFIRMED
END:VEVENT
END:VCALENDAR`

	cal, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if cal.Name != "Team Calendar" || len(cal.Events) != 3 {
		t.Fatalf("unexpected calendar parse result")
	}

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	busy := cal.GetBusySlots(from, to)
	if len(busy) == 0 {
		t.Fatalf("expected busy slots")
	}

	weekly := availability.NewWeeklySchedule(time.UTC).SetDay(time.Wednesday, availability.TimeRange{Start: availability.NewTimeOfDay(8, 0, 0), End: availability.NewTimeOfDay(12, 0, 0)})
	free := cal.GetFreeSlots(from, to, weekly)
	if free == nil {
		t.Fatalf("expected free slots slice")
	}

	if k, _, _ := parseICSLine("X-FOO"); k == "" {
		t.Fatalf("parseICSLine no-colon path failed")
	}
	k, p, v := parseICSLine("DTSTART;TZID=UTC:20250101T090000")
	if k != "DTSTART" || p["TZID"] != "UTC" || v == "" {
		t.Fatalf("parseICSLine with params failed")
	}

	if tz := tzFromParams(map[string]string{"TZID": "UTC"}, nil); tz.String() != "UTC" {
		t.Fatalf("tzFromParams valid TZID failed")
	}
	if tz := tzFromParams(map[string]string{"TZID": "Invalid/Zone"}, time.UTC); tz.String() != "UTC" {
		t.Fatalf("tzFromParams fallback failed")
	}
	if tz := tzFromParams(map[string]string{"TZID": "Invalid/Zone"}, nil); tz.String() != "UTC" {
		t.Fatalf("tzFromParams default UTC failed")
	}

	if _, err := parseDateTime("20250101T090000Z", time.UTC); err != nil {
		t.Fatalf("parseDateTime Z failed")
	}
	if _, err := parseDateTime("20250101T090000", time.UTC); err != nil {
		t.Fatalf("parseDateTime local failed")
	}
	if _, err := parseDateTime("20250101", time.UTC); err != nil {
		t.Fatalf("parseDateTime date failed")
	}
	if _, err := parseDateTime("oops", time.UTC); err == nil {
		t.Fatalf("expected parseDateTime error")
	}

	start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	if !overlapsTime(time.Time{}, time.Time{}, start, end) {
		t.Fatalf("zero bounds overlap should be true")
	}
	if !overlapsTime(time.Time{}, end.Add(time.Hour), start, end) {
		t.Fatalf("from zero overlap failed")
	}
	if !overlapsTime(start.Add(-time.Hour), time.Time{}, start, end) {
		t.Fatalf("to zero overlap failed")
	}
	if overlapsTime(end, end.Add(time.Hour), start, end) {
		t.Fatalf("non-overlap expected")
	}

	if got := maxTime(time.Time{}, end); !got.Equal(end) {
		t.Fatalf("maxTime zero case failed")
	}
	if got := maxTime(start, time.Time{}); !got.Equal(start) {
		t.Fatalf("maxTime zero case 2 failed")
	}
	if got := minTime(time.Time{}, end); !got.Equal(end) {
		t.Fatalf("minTime zero case failed")
	}
	if got := minTime(start, time.Time{}); !got.Equal(start) {
		t.Fatalf("minTime zero case 2 failed")
	}

	if !isException([]time.Time{start}, start) {
		t.Fatalf("expected exception hit")
	}
	if isException([]time.Time{start}, end) {
		t.Fatalf("unexpected exception hit")
	}
}

func TestCalendarConversionsAndParseFile(t *testing.T) {
	rec, _ := recurrence.ParseRule("FREQ=DAILY;COUNT=2")
	cal := &Calendar{
		Name:     "Synthetic",
		Timezone: time.UTC,
		Events: []Event{
			{UID: "1", Start: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC), Status: EventStatusConfirmed},
			{UID: "2", Start: time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC), Status: EventStatusCancelled},
			{UID: "3", Start: time.Date(2025, 1, 3, 9, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 3, 10, 0, 0, 0, time.UTC), Status: EventStatusConfirmed, Recurrence: rec, Exceptions: []time.Time{time.Date(2025, 1, 4, 9, 0, 0, 0, time.UTC)}},
		},
	}

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC)
	busy := cal.GetBusySlots(from, to)
	if len(busy) == 0 {
		t.Fatalf("expected busy slots")
	}

	if c := (*Calendar)(nil); c.GetBusySlots(from, to) != nil {
		t.Fatalf("nil calendar busy should be nil")
	}

	availability := cal.ToAvailability()
	if availability.Bookings.Len() == 0 {
		t.Fatalf("expected bookings in availability conversion")
	}
	if c := (*Calendar)(nil); c.ToAvailability().Location == nil {
		t.Fatalf("nil calendar availability should still have location")
	}

	slots := cal.ToSlotCollection(from, to)
	if slots.Len() == 0 {
		t.Fatalf("expected slot collection")
	}

	samplePath := filepath.Join("..", "testdata", "calendars", "sample.ics")
	parsed, err := ParseFile(samplePath)
	if err != nil {
		t.Fatalf("parse file failed: %v", err)
	}
	if len(parsed.Events) == 0 {
		t.Fatalf("expected events from sample file")
	}
	if _, err := ParseFile(filepath.Join("..", "testdata", "calendars", "missing.ics")); err == nil {
		t.Fatalf("expected parse file error for missing file")
	}

	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.ics")
	if err := os.WriteFile(badFile, []byte("BEGIN:VCALENDAR\nBEGIN:VEVENT\nDTSTART:bad\nEND:VEVENT\nEND:VCALENDAR"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if _, err := ParseFile(badFile); err == nil {
		t.Fatalf("expected parse error for bad file")
	}
}
