package validate

import (
	"errors"
	"time"
)

var (
	ErrInvalidTimeRange = errors.New("end time must be after start time")
	ErrInvalidDuration  = errors.New("duration must be positive")
)

func TimeRange(start, end time.Time) error {
	if !end.After(start) {
		return ErrInvalidTimeRange
	}
	return nil
}

func PositiveDuration(d time.Duration) error {
	if d <= 0 {
		return ErrInvalidDuration
	}
	return nil
}

func LocationOrUTC(loc *time.Location) *time.Location {
	if loc == nil {
		return time.UTC
	}
	return loc
}
