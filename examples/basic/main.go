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
		SetDay(time.Tuesday, availability.TimeRange{
			Start: availability.NewTimeOfDay(9, 0, 0),
			End:   availability.NewTimeOfDay(17, 0, 0),
		}).
		SetDay(time.Thursday, availability.TimeRange{
			Start: availability.NewTimeOfDay(10, 0, 0),
			End:   availability.NewTimeOfDay(14, 0, 0),
		})

	p := provider.NewProvider("stylist-1",
		provider.WithWeeklySchedule(weekly),
		provider.WithBuffer(15*time.Minute),
		provider.WithMinNotice(2*time.Hour),
	)

	holidays := []time.Time{
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
	}
	p = p.WithBlockedDates(holidays...)

	now := time.Now().UTC()
	q := query.NewQuery().
		Duration(time.Hour).
		Between(now, now.Add(14*24*time.Hour)).
		OnlyMornings().
		PreferEarlier().
		Limit(5).
		Build()

	slots, err := p.FindSlots(q)
	if err != nil {
		panic(err)
	}

	fmt.Println("Available slots:")
	for _, s := range slots {
		fmt.Printf("  %s - %s\n", s.Start.Format("Mon Jan 2 15:04"), s.End.Format("15:04"))
	}
}
