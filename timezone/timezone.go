package timezone

import (
	"fmt"
	"time"
)

func Load(name string) (*time.Location, error) {
	if name == "" {
		return nil, fmt.Errorf("timezone: empty location")
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("timezone: load %q: %w", name, err)
	}
	return loc, nil
}

func MustLoad(name string) *time.Location {
	loc, err := Load(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func Convert(t time.Time, loc *time.Location) time.Time {
	if loc == nil {
		return t
	}
	return t.In(loc)
}

func NowIn(loc *time.Location) time.Time {
	if loc == nil {
		return time.Now().UTC()
	}
	return time.Now().In(loc)
}

func EqualLocation(a, b *time.Location) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.String() == b.String()
}
