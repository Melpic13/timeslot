package conflict

import (
	"time"

	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

// ConflictDetector checks for scheduling conflicts.
type ConflictDetector struct {
	providers []*provider.Provider
	options   DetectorOptions
}

type DetectorOptions struct {
	IncludeBuffers     bool
	TravelTimeFunc     func(from, to *provider.Provider) time.Duration
	AllowDoubleBooking bool
}

type Conflict struct {
	Type       ConflictType
	Slot       slot.TimeSlot
	Providers  []*provider.Provider
	Resolution []ResolutionOption
}

type ConflictType int

const (
	ConflictOverlap ConflictType = iota
	ConflictBuffer
	ConflictTravelTime
	ConflictDoubleBooking
)

func NewDetector(providers ...*provider.Provider) *ConflictDetector {
	return &ConflictDetector{providers: append([]*provider.Provider(nil), providers...)}
}

func (d *ConflictDetector) AddProvider(p *provider.Provider) *ConflictDetector {
	d.providers = append(d.providers, p)
	return d
}

func (d *ConflictDetector) Check(s slot.TimeSlot, p *provider.Provider) []Conflict {
	var out []Conflict
	if p == nil {
		return out
	}
	if !p.IsAvailable(s) {
		out = append(out, Conflict{Type: ConflictOverlap, Slot: s, Providers: []*provider.Provider{p}, Resolution: defaultResolutionOptions(s)})
	}
	if d.options.IncludeBuffers {
		expanded := p.EffectiveAvailability(s)
		if !expanded.Equal(s) && !p.IsAvailable(expanded) {
			out = append(out, Conflict{Type: ConflictBuffer, Slot: expanded, Providers: []*provider.Provider{p}, Resolution: defaultResolutionOptions(expanded)})
		}
	}
	if !d.options.AllowDoubleBooking {
		for _, existing := range p.GetBookings(slotWindowAround(s, 24*time.Hour).Start, slotWindowAround(s, 24*time.Hour).End) {
			if existing.Overlaps(s) {
				out = append(out, Conflict{Type: ConflictDoubleBooking, Slot: s, Providers: []*provider.Provider{p}, Resolution: defaultResolutionOptions(s)})
				break
			}
		}
	}
	return out
}

func (d *ConflictDetector) CheckAll(s slot.TimeSlot) []Conflict {
	var out []Conflict
	for _, p := range d.providers {
		out = append(out, d.Check(s, p)...)
	}
	return out
}

func (d *ConflictDetector) FindConflictFree(q query.Query) []slot.TimeSlot {
	return d.FindCommonAvailability(q)
}

func (d *ConflictDetector) Resolve(c Conflict, strategy ResolutionStrategy) (*slot.TimeSlot, error) {
	return resolveSlot(c, strategy)
}

func (d *ConflictDetector) FindAvailableSlots(q query.Query) map[string][]slot.TimeSlot {
	return provider.FindAny(d.providers, q)
}

func (d *ConflictDetector) FindCommonAvailability(q query.Query) []slot.TimeSlot {
	return provider.FindCommon(d.providers, q)
}
