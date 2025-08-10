package ical

import (
	"fmt"
	"strings"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/slot"
)

func ExportSlots(slots []slot.TimeSlot, calName string) ([]byte, error) {
	if calName == "" {
		calName = "TimeSlot Export"
	}
	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//timeslot//EN\r\n")
	b.WriteString("X-WR-CALNAME:" + escape(calName) + "\r\n")
	for i, s := range slots {
		if err := s.Validate(); err != nil {
			return nil, err
		}
		b.WriteString("BEGIN:VEVENT\r\n")
		b.WriteString(fmt.Sprintf("UID:timeslot-%d@local\r\n", i+1))
		b.WriteString("DTSTART:" + s.Start.UTC().Format("20060102T150405Z") + "\r\n")
		b.WriteString("DTEND:" + s.End.UTC().Format("20060102T150405Z") + "\r\n")
		b.WriteString("SUMMARY:Available Slot\r\n")
		b.WriteString("END:VEVENT\r\n")
	}
	b.WriteString("END:VCALENDAR\r\n")
	return []byte(b.String()), nil
}

func ExportAvailability(a availability.Availability, from, to time.Time) ([]byte, error) {
	slots := a.GetSlots(from, to).Slots()
	return ExportSlots(slots, "Availability")
}

func escape(v string) string {
	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, ",", "\\,")
	v = strings.ReplaceAll(v, ";", "\\;")
	v = strings.ReplaceAll(v, "\n", "\\n")
	return v
}
