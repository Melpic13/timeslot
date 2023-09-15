package timeutil

import "time"

func StartOfDay(t time.Time, loc *time.Location) time.Time {
	l := loc
	if l == nil {
		l = t.Location()
	}
	t = t.In(l)
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, l)
}

func EndOfDay(t time.Time, loc *time.Location) time.Time {
	return StartOfDay(t, loc).Add(24*time.Hour - time.Nanosecond)
}

func SameDay(a, b time.Time, loc *time.Location) bool {
	aa := a
	bb := b
	if loc != nil {
		aa = a.In(loc)
		bb = b.In(loc)
	}
	ay, am, ad := aa.Date()
	by, bm, bd := bb.Date()
	return ay == by && am == bm && ad == bd
}

func Clamp(t, min, max time.Time) time.Time {
	if t.Before(min) {
		return min
	}
	if t.After(max) {
		return max
	}
	return t
}

func Overlap(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
