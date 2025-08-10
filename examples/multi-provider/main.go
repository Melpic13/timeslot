package main

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/conflict"
	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
)

func main() {
	providers := []*provider.Provider{
		createProvider("alice", time.Tuesday, time.Wednesday, time.Thursday),
		createProvider("bob", time.Monday, time.Tuesday, time.Friday),
		createProvider("carol", time.Wednesday, time.Thursday, time.Friday),
	}

	detector := conflict.NewDetector(providers...)

	now := time.Now().UTC()
	q := query.NewQuery().
		Duration(30*time.Minute).
		Between(now, now.Add(7*24*time.Hour)).
		Build()

	results := detector.FindAvailableSlots(q)
	for id, slots := range results {
		fmt.Printf("%s has %d available slots\n", id, len(slots))
	}

	common := detector.FindCommonAvailability(q)
	fmt.Printf("\nTimes when everyone is available: %d slots\n", len(common))
}

func createProvider(id string, days ...time.Weekday) *provider.Provider {
	weekly := availability.NewWeeklySchedule(time.UTC)
	for _, day := range days {
		weekly = weekly.SetDay(day, availability.TimeRange{
			Start: availability.NewTimeOfDay(9, 0, 0),
			End:   availability.NewTimeOfDay(17, 0, 0),
		})
	}
	return provider.NewProvider(id, provider.WithWeeklySchedule(weekly))
}
