package slot

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

var (
	ErrInvalidTimeRange = errors.New("slot: end time must be after start time")
	ErrUnionDisjoint    = errors.New("slot: cannot union disjoint slots")
)

// TimeSlot represents a specific time window with timezone awareness.
type TimeSlot struct {
	Start    time.Time      `json:"start"`
	End      time.Time      `json:"end"`
	Location *time.Location `json:"-"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func New(start, end time.Time) (TimeSlot, error) {
	ts := TimeSlot{Start: start, End: end, Location: start.Location()}
	return ts, ts.Validate()
}

func (s TimeSlot) Duration() time.Duration {
	return s.End.Sub(s.Start)
}

func (s TimeSlot) Contains(t time.Time) bool {
	return (t.Equal(s.Start) || t.After(s.Start)) && t.Before(s.End)
}

func (s TimeSlot) Overlaps(other TimeSlot) bool {
	return s.Start.Before(other.End) && other.Start.Before(s.End)
}

func (s TimeSlot) Intersection(other TimeSlot) (TimeSlot, bool) {
	if !s.Overlaps(other) {
		return TimeSlot{}, false
	}
	start := maxTime(s.Start, other.Start)
	end := minTime(s.End, other.End)
	loc := s.locationOrUTC()
	if other.Location != nil {
		loc = other.Location
	}
	return TimeSlot{Start: start.In(loc), End: end.In(loc), Location: loc}, true
}

func (s TimeSlot) Union(other TimeSlot) (TimeSlot, error) {
	if !s.Overlaps(other) && !s.End.Equal(other.Start) && !other.End.Equal(s.Start) {
		return TimeSlot{}, ErrUnionDisjoint
	}
	start := minTime(s.Start, other.Start)
	end := maxTime(s.End, other.End)
	loc := s.locationOrUTC()
	return TimeSlot{Start: start.In(loc), End: end.In(loc), Location: loc}, nil
}

func (s TimeSlot) Split(duration time.Duration) []TimeSlot {
	if duration <= 0 || s.IsZero() || !s.End.After(s.Start) {
		return nil
	}
	var out []TimeSlot
	for cur := s.Start; cur.Before(s.End); {
		next := cur.Add(duration)
		if next.After(s.End) {
			next = s.End
		}
		out = append(out, TimeSlot{Start: cur, End: next, Location: s.locationOrUTC(), Metadata: cloneMetadata(s.Metadata)})
		if !next.After(cur) {
			break
		}
		cur = next
	}
	return out
}

func (s TimeSlot) Shift(d time.Duration) TimeSlot {
	return TimeSlot{Start: s.Start.Add(d), End: s.End.Add(d), Location: s.locationOrUTC(), Metadata: cloneMetadata(s.Metadata)}
}

func (s TimeSlot) InTimezone(loc *time.Location) TimeSlot {
	if loc == nil {
		loc = time.UTC
	}
	return TimeSlot{Start: s.Start.In(loc), End: s.End.In(loc), Location: loc, Metadata: cloneMetadata(s.Metadata)}
}

func (s TimeSlot) IsZero() bool {
	return s.Start.IsZero() || s.End.IsZero()
}

func (s TimeSlot) Validate() error {
	if s.IsZero() {
		return ErrInvalidTimeRange
	}
	if !s.End.After(s.Start) {
		return ErrInvalidTimeRange
	}
	return nil
}

func (s TimeSlot) Equal(other TimeSlot) bool {
	if !s.Start.Equal(other.Start) || !s.End.Equal(other.End) {
		return false
	}
	if s.locationOrUTC().String() != other.locationOrUTC().String() {
		return false
	}
	return mapsEqual(s.Metadata, other.Metadata)
}

func (s TimeSlot) String() string {
	return fmt.Sprintf("%s-%s", s.Start.Format(time.RFC3339), s.End.Format(time.RFC3339))
}

func (s TimeSlot) MarshalJSON() ([]byte, error) {
	type alias struct {
		Start    time.Time      `json:"start"`
		End      time.Time      `json:"end"`
		Location string         `json:"location,omitempty"`
		Metadata map[string]any `json:"metadata,omitempty"`
	}
	return json.Marshal(alias{
		Start:    s.Start,
		End:      s.End,
		Location: s.locationOrUTC().String(),
		Metadata: s.Metadata,
	})
}

func (s *TimeSlot) UnmarshalJSON(data []byte) error {
	var raw struct {
		Start    time.Time      `json:"start"`
		End      time.Time      `json:"end"`
		Location string         `json:"location"`
		Metadata map[string]any `json:"metadata"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	loc := time.UTC
	if raw.Location != "" {
		parsed, err := time.LoadLocation(raw.Location)
		if err != nil {
			return err
		}
		loc = parsed
	}
	s.Start = raw.Start.In(loc)
	s.End = raw.End.In(loc)
	s.Location = loc
	s.Metadata = cloneMetadata(raw.Metadata)
	return s.Validate()
}

func Sort(slots []TimeSlot) []TimeSlot {
	out := append([]TimeSlot(nil), slots...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Start.Equal(out[j].Start) {
			return out[i].End.Before(out[j].End)
		}
		return out[i].Start.Before(out[j].Start)
	})
	return out
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func (s TimeSlot) locationOrUTC() *time.Location {
	if s.Location != nil {
		return s.Location
	}
	if s.Start.Location() != nil {
		return s.Start.Location()
	}
	return time.UTC
}

func cloneMetadata(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
