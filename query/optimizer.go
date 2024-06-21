package query

import (
	"sort"

	"github.com/Melpic13/timeslot/slot"
)

func OptimizeSlots(slots []slot.TimeSlot, q Query) []slot.TimeSlot {
	out := append([]slot.TimeSlot(nil), slots...)
	sort.SliceStable(out, func(i, j int) bool {
		si := score(out[i], q)
		sj := score(out[j], q)
		if si == sj {
			return out[i].Start.Before(out[j].Start)
		}
		return si > sj
	})
	if q.Limit > 0 && len(out) > q.Limit {
		return out[:q.Limit]
	}
	return out
}

func score(s slot.TimeSlot, q Query) int {
	total := 0
	for _, p := range q.Preferences {
		total += p.Score(s)
	}
	return total
}
