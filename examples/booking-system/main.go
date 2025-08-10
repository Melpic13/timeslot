package main

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/availability"
	"github.com/Melpic13/timeslot/provider"
	"github.com/Melpic13/timeslot/query"
)

func main() {
	weekly := availability.NewWeeklySchedule(time.UTC).
		SetDay(time.Monday, availability.TimeRange{Start: availability.NewTimeOfDay(9, 0, 0), End: availability.NewTimeOfDay(17, 0, 0)})
	p := provider.NewProvider("room-1", provider.WithWeeklySchedule(weekly), provider.WithBuffer(15*time.Minute))

	from := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	q := query.NewQuery().Duration(time.Hour).Between(from, to).Limit(3).Build()
	slots, err := p.FindSlots(q)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d candidate slots\n", len(slots))
	if len(slots) == 0 {
		return
	}

	booked, err := p.Book(slots[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Booked: %s - %s\n", slots[0].Start.Format(time.RFC3339), slots[0].End.Format(time.RFC3339))

	remaining, _ := booked.FindSlots(q)
	fmt.Printf("Remaining slots: %d\n", len(remaining))
}
