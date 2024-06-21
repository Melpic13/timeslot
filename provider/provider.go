package provider

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

// Provider represents an entity with availability.
type Provider struct {
	ID           string
	Availability availability.Availability
	BufferBefore time.Duration
	BufferAfter  time.Duration
	MinNotice    time.Duration
	MaxAdvance   time.Duration
	Metadata     map[string]any
}

func NewProvider(id string, opts ...ProviderOption) *Provider {
	p := &Provider{
		ID:           id,
		Availability: availability.New(time.UTC),
		Metadata:     map[string]any{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Provider) WithBlockedDates(dates ...time.Time) *Provider {
	copy := p.clone()
	copy.Availability = copy.Availability.AddBlockedDates(dates...)
	return copy
}

func (p *Provider) FindSlots(q query.Query) ([]slot.TimeSlot, error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}
	free := p.Availability.GetSlots(q.From, q.To)
	var candidates []slot.TimeSlot
	for _, s := range free.Slots() {
		for cur := s.Start; cur.Add(q.Duration).Before(s.End) || cur.Add(q.Duration).Equal(s.End); cur = cur.Add(q.Duration) {
			candidate := slot.TimeSlot{Start: cur, End: cur.Add(q.Duration), Location: s.Location}
			if p.passesConstraints(candidate, q.Constraints) {
				candidates = append(candidates, candidate)
			}
		}
	}
	candidates = query.OptimizeSlots(candidates, q)
	if q.Limit > 0 && len(candidates) > q.Limit {
		candidates = candidates[:q.Limit]
	}
	return candidates, nil
}

func (p *Provider) IsAvailable(s slot.TimeSlot) bool {
	from := s.Start.Add(-24 * time.Hour)
	to := s.End.Add(24 * time.Hour)
	free := p.Availability.GetSlots(from, to)
	overlap := free.FindOverlaps(s)
	if len(overlap) == 0 {
		return false
	}
	for _, existing := range p.Availability.Bookings.Slots() {
		if p.EffectiveAvailability(existing).Overlaps(s) || existing.Overlaps(p.EffectiveAvailability(s)) {
			return false
		}
	}
	return true
}

func (p *Provider) Book(s slot.TimeSlot) (*Provider, error) {
	now := time.Now().In(p.locationOrUTC())
	if s.Start.Before(now) {
		return nil, fmt.Errorf("provider: cannot book in the past")
	}
	if p.MinNotice > 0 && s.Start.Before(now.Add(p.MinNotice)) {
		return nil, fmt.Errorf("provider: insufficient notice")
	}
	if p.MaxAdvance > 0 && s.Start.After(now.Add(p.MaxAdvance)) {
		return nil, fmt.Errorf("provider: booking too far in advance")
	}
	if !p.IsAvailable(s) {
		return nil, fmt.Errorf("provider: slot not available")
	}
	copy := p.clone()
	copy.Availability = copy.Availability.AddBooking(s)
	return copy, nil
}

func (p *Provider) CancelBooking(s slot.TimeSlot) (*Provider, error) {
	copy := p.clone()
	before := copy.Availability.Bookings.Len()
	copy.Availability = copy.Availability.RemoveBooking(s)
	if before == copy.Availability.Bookings.Len() {
		return nil, fmt.Errorf("provider: booking not found")
	}
	return copy, nil
}

func (p *Provider) GetBookings(from, to time.Time) []slot.TimeSlot {
	window := slot.NewCollection(slot.TimeSlot{Start: from, End: to, Location: p.locationOrUTC()})
	return p.Availability.Bookings.Intersect(window).Slots()
}

func (p *Provider) EffectiveAvailability(s slot.TimeSlot) slot.TimeSlot {
	return slot.TimeSlot{
		Start:    s.Start.Add(-p.BufferBefore),
		End:      s.End.Add(p.BufferAfter),
		Location: s.Location,
		Metadata: s.Metadata,
	}
}

func (p *Provider) passesConstraints(candidate slot.TimeSlot, constraints []query.Constraint) bool {
	for _, c := range constraints {
		if !c.IsSatisfied(candidate) {
			return false
		}
	}
	return true
}

func (p *Provider) clone() *Provider {
	meta := map[string]any{}
	for k, v := range p.Metadata {
		meta[k] = v
	}
	copy := *p
	copy.Metadata = meta
	return &copy
}

func (p *Provider) locationOrUTC() *time.Location {
	if p.Availability.Location != nil {
		return p.Availability.Location
	}
	return time.UTC
}
