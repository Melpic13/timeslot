package ical

import (
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	ics := `BEGIN:VCALENDAR
VERSION:2.0
X-WR-CALNAME:Test
BEGIN:VEVENT
UID:1
DTSTART:20240101T090000Z
DTEND:20240101T100000Z
SUMMARY:Busy
END:VEVENT
END:VCALENDAR`
	cal, err := Parse(strings.NewReader(ics))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(cal.Events) != 1 {
		t.Fatalf("expected 1 event")
	}
	busy := cal.GetBusySlots(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))
	if len(busy) != 1 {
		t.Fatalf("expected 1 busy slot")
	}
}

func FuzzParseICS(f *testing.F) {
	f.Add("BEGIN:VCALENDAR\nEND:VCALENDAR")
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = Parse(strings.NewReader(input))
	})
}
