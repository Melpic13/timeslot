package availability

import "testing"

func TestParseTimeRange(t *testing.T) {
	r, err := ParseTimeRange("09:00-17:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Start.Hour != 9 || r.End.Hour != 17 {
		t.Fatalf("unexpected parsed range: %+v", r)
	}
}

func FuzzParseTimeRange(f *testing.F) {
	f.Add("09:00-17:00")
	f.Add("00:00-23:59")
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = ParseTimeRange(input)
	})
}
