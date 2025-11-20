package recurrence

import (
	"strings"
	"testing"
	"time"
)

func TestRuleStringAndFrequencyString(t *testing.T) {
	r := &Rule{
		Frequency:  Weekly,
		Interval:   2,
		Count:      3,
		Until:      time.Date(2025, 12, 31, 23, 0, 0, 0, time.UTC),
		ByDay:      []Weekday{time.Monday, time.Wednesday},
		ByMonth:    []time.Month{time.January, time.February},
		ByMonthDay: []int{1, 15},
		ByHour:     []int{9},
		ByMinute:   []int{30},
	}
	s := r.String()
	for _, want := range []string{"FREQ=WEEKLY", "INTERVAL=2", "COUNT=3", "BYDAY=MO,WE", "BYMONTH=1,2", "BYMONTHDAY=1,15", "BYHOUR=9", "BYMINUTE=30"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in %q", want, s)
		}
	}
	if (*Rule)(nil).String() != "" {
		t.Fatalf("nil rule string should be empty")
	}
	if Frequency(99).String() != "DAILY" {
		t.Fatalf("unknown frequency should default")
	}
}

func TestRuleGenerateNextContainsValidateAndStep(t *testing.T) {
	start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)

	if out := (*Rule)(nil).Generate(start, 1); out != nil {
		t.Fatalf("nil generate should return nil")
	}

	r := &Rule{Frequency: Daily, Interval: 1, Count: 3}
	got := r.Generate(start, 0)
	if len(got) != 3 {
		t.Fatalf("count should cap generate")
	}

	rUntil := &Rule{Frequency: Daily, Interval: 1, Until: start.Add(48 * time.Hour)}
	if out := rUntil.Generate(start, 10); len(out) == 0 || len(out) > 3 {
		t.Fatalf("until capping failed: %d", len(out))
	}

	if out := r.GenerateBetween(start, start.Add(24*time.Hour), start); out != nil {
		t.Fatalf("invalid range should be nil")
	}
	between := r.GenerateBetween(start, start, start.Add(4*24*time.Hour))
	if len(between) != 3 {
		t.Fatalf("generate between mismatch")
	}

	next, ok := r.Next(start)
	if !ok || !next.Equal(start.Add(24*time.Hour)) {
		t.Fatalf("next failed")
	}

	rShort := &Rule{Frequency: Daily, Interval: 1, Until: start}
	if _, ok := rShort.Next(start); ok {
		t.Fatalf("next should fail past until")
	}

	if !r.Contains(start) {
		t.Fatalf("contains should be true")
	}
	if (&Rule{Until: start}).Contains(start.Add(time.Hour)) {
		t.Fatalf("contains should fail past until")
	}
	if (*Rule)(nil).Contains(start) {
		t.Fatalf("nil contains should be false")
	}

	if err := (*Rule)(nil).Validate(); err == nil {
		t.Fatalf("nil validate should fail")
	}
	if err := (&Rule{Interval: -1}).Validate(); err == nil {
		t.Fatalf("negative interval should fail")
	}

	stepCases := []struct {
		freq Frequency
		want time.Time
	}{
		{Daily, start.AddDate(0, 0, 1)},
		{Weekly, start.AddDate(0, 0, 7)},
		{Monthly, start.AddDate(0, 1, 0)},
		{Yearly, start.AddDate(1, 0, 0)},
		{Frequency(99), start.AddDate(0, 0, 1)},
	}
	for _, tc := range stepCases {
		rule := &Rule{Frequency: tc.freq, Interval: 0}
		if got := rule.step(start); !got.Equal(tc.want) {
			t.Fatalf("step mismatch for %v: got %v want %v", tc.freq, got, tc.want)
		}
	}
}

func TestRuleMatchesAndParseHelpers(t *testing.T) {
	candidate := time.Date(2025, 1, 6, 9, 30, 0, 0, time.UTC) // Monday
	r := &Rule{
		ByDay:      []Weekday{time.Monday},
		ByMonth:    []time.Month{time.January},
		ByMonthDay: []int{6},
		ByHour:     []int{9},
		ByMinute:   []int{30},
	}
	if !r.matches(candidate) {
		t.Fatalf("expected full match")
	}
	if (&Rule{ByDay: []Weekday{time.Tuesday}}).matches(candidate) {
		t.Fatalf("day mismatch should fail")
	}
	if (&Rule{ByMonth: []time.Month{time.February}}).matches(candidate) {
		t.Fatalf("month mismatch should fail")
	}
	if (&Rule{ByMonthDay: []int{7}}).matches(candidate) {
		t.Fatalf("month day mismatch should fail")
	}
	if (&Rule{ByHour: []int{10}}).matches(candidate) {
		t.Fatalf("hour mismatch should fail")
	}
	if (&Rule{ByMinute: []int{0}}).matches(candidate) {
		t.Fatalf("minute mismatch should fail")
	}

	parsed, err := parseRule("FREQ=MONTHLY;INTERVAL=1;COUNT=2;BYDAY=MO,WE;BYMONTH=1,2;BYMONTHDAY=10,1;BYHOUR=9;BYMINUTE=30")
	if err != nil {
		t.Fatalf("parse rule failed: %v", err)
	}
	if len(parsed.ByMonthDay) != 2 || parsed.ByMonthDay[0] != 1 {
		t.Fatalf("expected sorted BYMONTHDAY")
	}

	if _, err := parseRule("BADTOKEN"); err == nil {
		t.Fatalf("expected malformed token error")
	}
	if _, err := parseRule("FREQ=NOPE"); err == nil {
		t.Fatalf("expected unsupported freq error")
	}
	if _, err := parseRule("INTERVAL=x"); err == nil {
		t.Fatalf("expected invalid interval parse error")
	}
	if _, err := parseRule("BYDAY=XX"); err == nil {
		t.Fatalf("expected invalid weekday")
	}
	if _, err := parseRule("INTERVAL=-1"); err == nil {
		t.Fatalf("expected invalid interval validation")
	}

	if _, err := parseInts("1, 2,3"); err != nil {
		t.Fatalf("parse ints failed: %v", err)
	}
	if _, err := parseInts("1,a"); err == nil {
		t.Fatalf("expected parse ints error")
	}

	if _, err := parseDateTime("20250101T090000Z", time.UTC); err != nil {
		t.Fatalf("parse UTC datetime failed: %v", err)
	}
	if _, err := parseDateTime("20250101T090000", time.UTC); err != nil {
		t.Fatalf("parse local datetime failed: %v", err)
	}
	if _, err := parseDateTime("20250101", time.UTC); err != nil {
		t.Fatalf("parse date failed: %v", err)
	}
	if _, err := parseDateTime("not-a-date", time.UTC); err == nil {
		t.Fatalf("expected invalid datetime error")
	}

	for tok, day := range map[string]time.Weekday{"MO": time.Monday, "TU": time.Tuesday, "WE": time.Wednesday, "TH": time.Thursday, "FR": time.Friday, "SA": time.Saturday, "SU": time.Sunday} {
		got, err := tokenToWeekday(tok)
		if err != nil || got != day {
			t.Fatalf("tokenToWeekday %s failed", tok)
		}
	}
	if _, err := tokenToWeekday("XX"); err == nil {
		t.Fatalf("expected invalid token")
	}
	for _, day := range []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday} {
		if weekdayToToken(day) == "" {
			t.Fatalf("weekday token should not be empty")
		}
	}
	if weekdayToToken(time.Weekday(99)) != "MO" {
		t.Fatalf("default weekday token mismatch")
	}
}

func TestRecurrenceWrappers(t *testing.T) {
	r, err := Parse("FREQ=DAILY;COUNT=2")
	if err != nil {
		t.Fatalf("parse wrapper failed: %v", err)
	}
	if _, err := ParseRule("FREQ=DAILY;COUNT=2"); err != nil {
		t.Fatalf("parse rule failed: %v", err)
	}
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if len(GenerateOccurrences(nil, start, 1)) != 0 {
		t.Fatalf("nil rule wrapper should return nil slice")
	}
	if len(GenerateOccurrences(r, start, 2)) != 2 {
		t.Fatalf("generate occurrences wrapper failed")
	}
}
