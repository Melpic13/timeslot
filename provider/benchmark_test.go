package provider

import (
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/query"
)

func BenchmarkProvider_FindSlots_LargeRange(b *testing.B) {
	ws := availability.NewWeeklySchedule(time.UTC)
	for _, day := range []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday} {
		ws = ws.SetDay(day, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(17, 0, 0)})
	}
	p := NewProvider("bench", WithWeeklySchedule(ws))
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(1, 0, 0)
	q := query.NewQuery().Duration(time.Hour).Between(from, to).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.FindSlots(q)
	}
}
