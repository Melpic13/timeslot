package ical

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkICalParse_LargeCalendar(b *testing.B) {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\n")
	for i := 0; i < 10000; i++ {
		sb.WriteString("BEGIN:VEVENT\n")
		sb.WriteString(fmt.Sprintf("UID:%d\n", i))
		sb.WriteString("DTSTART:20250101T090000Z\n")
		sb.WriteString("DTEND:20250101T093000Z\n")
		sb.WriteString("SUMMARY:Bench\n")
		sb.WriteString("END:VEVENT\n")
	}
	sb.WriteString("END:VCALENDAR\n")
	payload := sb.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(strings.NewReader(payload))
	}
}
