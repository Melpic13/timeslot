package slot

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkSlotCollection_FindOverlaps(b *testing.B) {
	loc := time.UTC
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, loc)
	slots := make([]TimeSlot, 0, 1000)
	for i := 0; i < 1000; i++ {
		start := base.Add(time.Duration(i*30) * time.Minute)
		slots = append(slots, TimeSlot{
			Start:    start,
			End:      start.Add(15 * time.Minute),
			Location: loc,
			Metadata: map[string]any{"i": strconv.Itoa(i)},
		})
	}
	c := NewCollection(slots...)
	probe := TimeSlot{Start: base.Add(500 * 30 * time.Minute), End: base.Add(500*30*time.Minute + time.Hour), Location: loc}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.FindOverlaps(probe)
	}
}
