package timeslot

import (
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

type (
	TimeSlot       = slot.TimeSlot
	SlotCollection = slot.SlotCollection
	Availability   = availability.Availability
	WeeklySchedule = availability.WeeklySchedule
	TimeRange      = availability.TimeRange
	TimeOfDay      = availability.TimeOfDay
	Provider       = provider.Provider
	Query          = query.Query
)

func NewSlot(start, end time.Time) (TimeSlot, error) {
	return slot.New(start, end)
}

func NewCollection(slots ...TimeSlot) SlotCollection {
	return slot.NewCollection(slots...)
}

func NewWeeklySchedule(loc *time.Location) WeeklySchedule {
	return availability.NewWeeklySchedule(loc)
}

func NewAvailability(loc *time.Location) Availability {
	return availability.New(loc)
}

func NewProvider(id string, opts ...provider.ProviderOption) *Provider {
	return provider.NewProvider(id, opts...)
}

func NewQuery() *query.QueryBuilder {
	return query.NewQuery()
}
