package main

import (
	"fmt"
	"time"

	"github.com/Melpic13/timeslot/recurrence"
)

func main() {
	rule, err := recurrence.ParseRule("FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,WE,FR;COUNT=10")
	if err != nil {
		panic(err)
	}
	start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	occurrences := rule.Generate(start, 10)
	for _, o := range occurrences {
		fmt.Println(o.Format(time.RFC3339))
	}
}
