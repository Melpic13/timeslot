package availability

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseTimeOfDay parses HH:MM or HH:MM:SS.
func ParseTimeOfDay(input string) (TimeOfDay, error) {
	parts := strings.Split(strings.TrimSpace(input), ":")
	if len(parts) < 2 || len(parts) > 3 {
		return TimeOfDay{}, fmt.Errorf("availability: invalid time of day %q", input)
	}
	vals := make([]int, 3)
	for i := 0; i < len(parts); i++ {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return TimeOfDay{}, fmt.Errorf("availability: invalid time component %q", parts[i])
		}
		vals[i] = n
	}
	tod := NewTimeOfDay(vals[0], vals[1], vals[2])
	if err := tod.Validate(); err != nil {
		return TimeOfDay{}, err
	}
	return tod, nil
}

// ParseTimeRange parses "HH:MM-HH:MM" or "HH:MM:SS-HH:MM:SS".
func ParseTimeRange(input string) (TimeRange, error) {
	parts := strings.Split(strings.TrimSpace(input), "-")
	if len(parts) != 2 {
		return TimeRange{}, fmt.Errorf("availability: invalid range %q", input)
	}
	start, err := ParseTimeOfDay(parts[0])
	if err != nil {
		return TimeRange{}, err
	}
	end, err := ParseTimeOfDay(parts[1])
	if err != nil {
		return TimeRange{}, err
	}
	r := TimeRange{Start: start, End: end}
	if err := r.Validate(); err != nil {
		return TimeRange{}, err
	}
	return r, nil
}
