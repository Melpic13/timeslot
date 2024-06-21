package provider

import (
	"time"

	"github.com/Melpic13/timeslot/availability"
)

// ProviderOption configures a Provider.
type ProviderOption func(*Provider)

func WithWeeklySchedule(ws availability.WeeklySchedule) ProviderOption {
	return func(p *Provider) {
		p.Availability.Weekly = ws
		if p.Availability.Location == nil {
			p.Availability.Location = ws.Location
		}
	}
}

func WithBufferBefore(d time.Duration) ProviderOption {
	return func(p *Provider) { p.BufferBefore = d }
}

func WithBufferAfter(d time.Duration) ProviderOption {
	return func(p *Provider) { p.BufferAfter = d }
}

func WithBuffer(d time.Duration) ProviderOption {
	return func(p *Provider) {
		p.BufferBefore = d
		p.BufferAfter = d
	}
}

func WithMinNotice(d time.Duration) ProviderOption {
	return func(p *Provider) { p.MinNotice = d }
}

func WithMaxAdvance(d time.Duration) ProviderOption {
	return func(p *Provider) { p.MaxAdvance = d }
}

func WithTimezone(loc *time.Location) ProviderOption {
	return func(p *Provider) {
		if loc != nil {
			p.Availability.Location = loc
			p.Availability.Weekly.Location = loc
		}
	}
}

func WithMetadata(key string, value any) ProviderOption {
	return func(p *Provider) {
		if p.Metadata == nil {
			p.Metadata = map[string]any{}
		}
		p.Metadata[key] = value
	}
}

func WithBlockedDates(dates ...time.Time) ProviderOption {
	return func(p *Provider) {
		p.Availability = p.Availability.AddBlockedDates(dates...)
	}
}
