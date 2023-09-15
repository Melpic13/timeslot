package timezone

import "time"

func IsDST(t time.Time) bool {
	return t.IsDST()
}

func IsDSTTransitionDay(day time.Time) bool {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	_, offStart := start.Zone()
	_, offEnd := end.Zone()
	return offStart != offEnd
}
