# TimeSlot

Production-grade availability and scheduling logic for Go.

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue)](#)
[![CI](https://img.shields.io/badge/ci-github_actions-brightgreen)](#)
[![Coverage](https://img.shields.io/badge/coverage-92.3%25-success)](#)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- Zero-dependency scheduling primitives
- Immutable-style slot and collection operations
- Weekly schedules with exception handling
- Provider-level booking, buffers, and constraints
- Multi-provider conflict detection
- RRULE subset recurrence support
- iCal import/export helpers

## Installation

```bash
go get github.com/Melpic13/timeslot
```

## Quick Start

```go
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
        })

    p := provider.NewProvider("stylist-1", provider.WithWeeklySchedule(weekly))

    now := time.Now().UTC()
    q := query.NewQuery().
        Duration(time.Hour).
        Between(now, now.Add(7*24*time.Hour)).
        Limit(5).
        Build()

    slots, _ := p.FindSlots(q)
    fmt.Println("slots:", len(slots))
}
```

## Quality Gates

- `go test -race ./...` must pass
- Coverage must stay at or above `90%` (`make coverage-check`)
- CI enforces `go mod tidy`, `go vet`, lint, and gosec

## Core Concepts

- `slot.TimeSlot`: one concrete time window
- `slot.SlotCollection`: immutable collection operations
- `availability.WeeklySchedule`: recurring weekly windows
- `availability.Availability`: weekly schedule + exceptions + bookings
- `provider.Provider`: bookable resource with options
- `query.Query`: fluent API for slot searches

## Examples

- `examples/basic/main.go`
- `examples/multi-provider/main.go`
- `examples/recurring-availability/main.go`
- `examples/booking-system/main.go`

## API Reference

Use `pkg.go.dev/github.com/Melpic13/timeslot`.

## Performance

Benchmark stubs are included in test files. Run with:

```bash
go test -bench=. -benchmem ./...
```

## Contributing

See `CONTRIBUTING.md`.

## License

MIT
