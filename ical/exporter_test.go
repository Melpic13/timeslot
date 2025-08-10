package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/Melpic13/timeslot/slot"
)

func TestExportSlots(t *testing.T) {
	s := slot.TimeSlot{Start: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), Location: time.UTC}
	b, err := ExportSlots([]slot.TimeSlot{s}, "Test")
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if !strings.Contains(string(b), "BEGIN:VCALENDAR") {
		t.Fatalf("unexpected content")
	}
}
