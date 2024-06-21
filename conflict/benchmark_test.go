package conflict

import (
	"fmt"
	"testing"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
)

func BenchmarkConflictDetector_100Providers(b *testing.B) {
	ws := availability.NewWeeklySchedule(time.UTC).SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(17, 0, 0)})
	providers := make([]*provider.Provider, 0, 100)
	for i := 0; i < 100; i++ {
		providers = append(providers, provider.NewProvider(fmt.Sprintf("p-%d", i), provider.WithWeeklySchedule(ws)))
	}
	d := NewDetector(providers...)
	from := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	q := query.NewQuery().Duration(time.Hour).Between(from, to).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.FindCommonAvailability(q)
	}
}
