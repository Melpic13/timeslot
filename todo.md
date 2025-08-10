TimeSlot Go Library

## Project Identity

**Name:** TimeSlot
**Tagline:** Production-grade availability and scheduling logic for Go
**Repository:** github.com/[username]/timeslot
**License:** MIT
**Go Version:** 1.22+

---

## Mission Statement

Build a zero-dependency, production-ready Go library that solves complex availability management, time slot optimization, and booking conflict resolution. This library should be the go-to solution for any Go developer building scheduling features, replacing the need to rebuild this logic for every project.

---

## Core Design Principles

1. **Zero external dependencies** — Only use Go standard library
2. **Immutable by default** — All operations return new structs, never mutate
3. **Timezone-first** — Every time operation is timezone-aware from day one
4. **Composable API** — Builder pattern with method chaining
5. **Comprehensive error handling** — Typed errors with context
6. **100% test coverage** — Table-driven tests for all edge cases
7. **Performance conscious** — Benchmark critical paths, document Big-O complexity

---

## Project Structure

Create the following directory structure:

```
timeslot/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml                 # Test, lint, coverage on PR/push
│   │   ├── release.yml            # Automated releases with goreleaser
│   │   └── codeql.yml             # Security scanning
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── FUNDING.yml
├── availability/
│   ├── availability.go            # Core availability struct and methods
│   ├── availability_test.go
│   ├── weekly.go                  # Weekly recurring availability
│   ├── weekly_test.go
│   ├── exceptions.go              # Date exceptions (holidays, blocks)
│   └── exceptions_test.go
├── slot/
│   ├── slot.go                    # TimeSlot struct and operations
│   ├── slot_test.go
│   ├── collection.go              # SlotCollection with set operations
│   └── collection_test.go
├── provider/
│   ├── provider.go                # Provider entity with availability
│   ├── provider_test.go
│   ├── options.go                 # Functional options pattern
│   └── multi.go                   # Multi-provider availability
├── query/
│   ├── query.go                   # Query builder for finding slots
│   ├── query_test.go
│   ├── constraints.go             # Constraint definitions
│   └── optimizer.go               # Slot optimization algorithms
├── conflict/
│   ├── detector.go                # Conflict detection engine
│   ├── detector_test.go
│   ├── resolution.go              # Conflict resolution strategies
│   └── buffer.go                  # Buffer time handling
├── recurrence/
│   ├── rule.go                    # Recurrence rule engine (RFC 5545 subset)
│   ├── rule_test.go
│   ├── parser.go                  # Parse recurrence strings
│   └── generator.go               # Generate occurrences
├── ical/
│   ├── parser.go                  # iCal/ICS file parsing
│   ├── parser_test.go
│   ├── exporter.go                # Export to iCal format
│   └── exporter_test.go
├── timezone/
│   ├── timezone.go                # Timezone utilities
│   ├── timezone_test.go
│   └── dst.go                     # DST transition handling
├── examples/
│   ├── basic/
│   │   └── main.go
│   ├── multi-provider/
│   │   └── main.go
│   ├── recurring-availability/
│   │   └── main.go
│   └── booking-system/
│   │   └── main.go
├── internal/
│   ├── timeutil/
│   │   ├── timeutil.go            # Internal time helpers
│   │   └── timeutil_test.go
│   └── validate/
│       └── validate.go            # Internal validation helpers
├── testdata/
│   ├── calendars/                 # Sample iCal files for testing
│   └── fixtures/                  # Test fixtures
├── .gitignore
├── .golangci.yml                  # Linter configuration
├── .goreleaser.yml                # Release configuration
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── LICENSE
├── Makefile
├── README.md
├── SECURITY.md
├── doc.go                         # Package documentation
├── timeslot.go                    # Main entry point, re-exports
├── timeslot_test.go
├── errors.go                      # Typed errors
├── version.go                     # Version information
└── go.mod
```

---

## Core Types and Interfaces

### 1. TimeSlot (slot/slot.go)

```go
// TimeSlot represents a specific time window with timezone awareness
type TimeSlot struct {
    Start    time.Time
    End      time.Time
    Location *time.Location
    Metadata map[string]any
}

// Methods to implement:
// - Duration() time.Duration
// - Contains(t time.Time) bool
// - Overlaps(other TimeSlot) bool
// - Intersection(other TimeSlot) (TimeSlot, bool)
// - Union(other TimeSlot) (TimeSlot, error)
// - Split(duration time.Duration) []TimeSlot
// - Shift(d time.Duration) TimeSlot
// - InTimezone(loc *time.Location) TimeSlot
// - IsZero() bool
// - Validate() error
// - Equal(other TimeSlot) bool
// - String() string
// - MarshalJSON() ([]byte, error)
// - UnmarshalJSON(data []byte) error
```

### 2. SlotCollection (slot/collection.go)

```go
// SlotCollection is an immutable, sorted collection of non-overlapping slots
type SlotCollection struct {
    slots    []TimeSlot
    location *time.Location
}

// Methods to implement:
// - Add(slots ...TimeSlot) SlotCollection
// - Remove(slots ...TimeSlot) SlotCollection
// - Merge() SlotCollection                    // Merge overlapping slots
// - Subtract(other SlotCollection) SlotCollection
// - Intersect(other SlotCollection) SlotCollection
// - Union(other SlotCollection) SlotCollection
// - Filter(fn func(TimeSlot) bool) SlotCollection
// - FindOverlaps(slot TimeSlot) []TimeSlot
// - TotalDuration() time.Duration
// - Gaps(within TimeSlot) SlotCollection
// - Slots() []TimeSlot                        // Returns copy
// - Len() int
// - IsEmpty() bool
// - First() (TimeSlot, bool)
// - Last() (TimeSlot, bool)
```

### 3. WeeklySchedule (availability/weekly.go)

```go
// WeeklySchedule defines recurring weekly availability
type WeeklySchedule struct {
    Monday    []TimeRange
    Tuesday   []TimeRange
    Wednesday []TimeRange
    Thursday  []TimeRange
    Friday    []TimeRange
    Saturday  []TimeRange
    Sunday    []TimeRange
    Location  *time.Location
}

// TimeRange represents a time window within a day (no date)
type TimeRange struct {
    Start TimeOfDay  // e.g., 09:00
    End   TimeOfDay  // e.g., 17:00
}

// TimeOfDay represents a time without a date
type TimeOfDay struct {
    Hour   int
    Minute int
    Second int
}

// Methods to implement:
// - SetDay(day time.Weekday, ranges ...TimeRange) WeeklySchedule
// - GetDay(day time.Weekday) []TimeRange
// - GenerateSlots(from, to time.Time) SlotCollection
// - IsAvailable(t time.Time) bool
// - NextAvailable(after time.Time) (time.Time, bool)
// - MergeWith(other WeeklySchedule) WeeklySchedule
// - Validate() error
```

### 4. Availability (availability/availability.go)

```go
// Availability combines weekly schedule with exceptions
type Availability struct {
    Weekly     WeeklySchedule
    Exceptions ExceptionSet
    Bookings   SlotCollection  // Already booked slots
    Location   *time.Location
}

// ExceptionSet handles date-based exceptions
type ExceptionSet struct {
    Blocked   []DateRange     // Completely unavailable
    Available []DateRange     // Override: available even if weekly says no
    Modified  []DateOverride  // Different hours on specific dates
}

// Methods to implement:
// - AddBlockedDates(dates ...time.Time) Availability
// - AddBlockedRange(start, end time.Time) Availability
// - AddAvailableOverride(date time.Time, ranges ...TimeRange) Availability
// - AddBooking(slot TimeSlot) Availability
// - RemoveBooking(slot TimeSlot) Availability
// - GetSlots(from, to time.Time) SlotCollection
// - IsAvailable(t time.Time) bool
// - IsBooked(t time.Time) bool
// - FindAvailableSlots(duration time.Duration, from, to time.Time) []TimeSlot
// - Validate() error
```

### 5. Provider (provider/provider.go)

```go
// Provider represents an entity with availability (person, room, resource)
type Provider struct {
    ID           string
    Availability Availability
    BufferBefore time.Duration
    BufferAfter  time.Duration
    MinNotice    time.Duration    // Minimum booking notice
    MaxAdvance   time.Duration    // Maximum advance booking
    Metadata     map[string]any
}

// Builder pattern with functional options
func NewProvider(id string, opts ...ProviderOption) *Provider

// ProviderOption type
type ProviderOption func(*Provider)

// Option functions:
// - WithWeeklySchedule(ws WeeklySchedule) ProviderOption
// - WithBufferBefore(d time.Duration) ProviderOption
// - WithBufferAfter(d time.Duration) ProviderOption
// - WithBuffer(d time.Duration) ProviderOption  // Sets both
// - WithMinNotice(d time.Duration) ProviderOption
// - WithMaxAdvance(d time.Duration) ProviderOption
// - WithTimezone(loc *time.Location) ProviderOption
// - WithMetadata(key string, value any) ProviderOption
// - WithBlockedDates(dates ...time.Time) ProviderOption

// Methods to implement:
// - FindSlots(query Query) ([]TimeSlot, error)
// - IsAvailable(slot TimeSlot) bool
// - Book(slot TimeSlot) (*Provider, error)    // Returns new Provider
// - CancelBooking(slot TimeSlot) (*Provider, error)
// - GetBookings(from, to time.Time) []TimeSlot
// - EffectiveAvailability(slot TimeSlot) TimeSlot  // Accounts for buffers
```

### 6. Query (query/query.go)

```go
// Query defines search criteria for finding available slots
type Query struct {
    Duration    time.Duration
    From        time.Time
    To          time.Time
    Constraints []Constraint
    Preferences []Preference
    Limit       int
    Location    *time.Location
}

// QueryBuilder provides fluent API
type QueryBuilder struct {
    query Query
}

func NewQuery() *QueryBuilder

// Methods:
// - Duration(d time.Duration) *QueryBuilder
// - Between(from, to time.Time) *QueryBuilder
// - From(t time.Time) *QueryBuilder
// - To(t time.Time) *QueryBuilder
// - InNext(d time.Duration) *QueryBuilder      // Convenience: next 7 days, etc.
// - OnWeekdays(days ...time.Weekday) *QueryBuilder
// - OnlyMornings() *QueryBuilder               // Before 12:00
// - OnlyAfternoons() *QueryBuilder             // 12:00-17:00
// - OnlyEvenings() *QueryBuilder               // After 17:00
// - NotBefore(tod TimeOfDay) *QueryBuilder
// - NotAfter(tod TimeOfDay) *QueryBuilder
// - WithConstraint(c Constraint) *QueryBuilder
// - PreferEarlier() *QueryBuilder
// - PreferLater() *QueryBuilder
// - PreferTime(tod TimeOfDay) *QueryBuilder    // Prefer slots near this time
// - Limit(n int) *QueryBuilder
// - InTimezone(loc *time.Location) *QueryBuilder
// - Build() Query
// - Validate() error

// Constraint interface for extensibility
type Constraint interface {
    IsSatisfied(slot TimeSlot) bool
    String() string
}

// Built-in constraints:
// - WeekdayConstraint
// - TimeOfDayConstraint
// - NotBeforeConstraint
// - NotAfterConstraint
// - MinGapConstraint (minimum gap from other bookings)
```

### 7. Conflict Detection (conflict/detector.go)

```go
// ConflictDetector checks for scheduling conflicts
type ConflictDetector struct {
    providers []*Provider
    options   DetectorOptions
}

type DetectorOptions struct {
    IncludeBuffers     bool
    TravelTimeFunc     func(from, to *Provider) time.Duration
    AllowDoubleBooking bool
}

type Conflict struct {
    Type       ConflictType
    Slot       TimeSlot
    Providers  []*Provider
    Resolution []ResolutionOption
}

type ConflictType int

const (
    ConflictOverlap ConflictType = iota
    ConflictBuffer
    ConflictTravelTime
    ConflictDoubleBooking
)

// Methods:
// - AddProvider(p *Provider) *ConflictDetector
// - Check(slot TimeSlot, provider *Provider) []Conflict
// - CheckAll(slot TimeSlot) []Conflict
// - FindConflictFree(query Query) []TimeSlot
// - Resolve(conflict Conflict, strategy ResolutionStrategy) (*TimeSlot, error)
```

### 8. Recurrence Rules (recurrence/rule.go)

```go
// Rule represents a recurrence rule (subset of RFC 5545 RRULE)
type Rule struct {
    Frequency  Frequency
    Interval   int
    Count      int           // 0 = unlimited
    Until      time.Time     // Zero = unlimited
    ByDay      []Weekday
    ByMonth    []time.Month
    ByMonthDay []int
    ByHour     []int
    ByMinute   []int
    WeekStart  time.Weekday
    Location   *time.Location
}

type Frequency int

const (
    Daily Frequency = iota
    Weekly
    Monthly
    Yearly
)

// Parse RRULE string
func ParseRule(rrule string) (*Rule, error)

// Methods:
// - String() string                            // To RRULE format
// - Generate(start time.Time, limit int) []time.Time
// - GenerateBetween(start, from, to time.Time) []time.Time
// - Next(after time.Time) (time.Time, bool)
// - Contains(t time.Time) bool
// - Validate() error
```

### 9. iCal Support (ical/parser.go)

```go
// Calendar represents a parsed iCal file
type Calendar struct {
    Name       string
    Timezone   *time.Location
    Events     []Event
}

type Event struct {
    UID         string
    Summary     string
    Start       time.Time
    End         time.Time
    Recurrence  *Rule
    Exceptions  []time.Time  // EXDATE
    Status      EventStatus
}

// Functions:
func Parse(r io.Reader) (*Calendar, error)
func ParseFile(path string) (*Calendar, error)

// Methods on Calendar:
// - ToAvailability() Availability              // Busy times become blocked
// - ToSlotCollection(from, to time.Time) SlotCollection
// - GetBusySlots(from, to time.Time) []TimeSlot
// - GetFreeSlots(from, to time.Time, within WeeklySchedule) []TimeSlot

// Export functions (ical/exporter.go):
func ExportSlots(slots []TimeSlot, calName string) ([]byte, error)
func ExportAvailability(a Availability, from, to time.Time) ([]byte, error)
```

---

## Error Handling

Create typed errors in `errors.go`:

```go
var (
    ErrInvalidTimeRange   = errors.New("timeslot: end time must be after start time")
    ErrInvalidDuration    = errors.New("timeslot: duration must be positive")
    ErrSlotOverlap        = errors.New("timeslot: slots overlap")
    ErrNoAvailability     = errors.New("timeslot: no availability found")
    ErrConflict           = errors.New("timeslot: booking conflict detected")
    ErrInvalidTimezone    = errors.New("timeslot: invalid timezone")
    ErrPastTime           = errors.New("timeslot: cannot book in the past")
    ErrInsufficientNotice = errors.New("timeslot: insufficient booking notice")
    ErrTooFarAdvance      = errors.New("timeslot: booking too far in advance")
    ErrInvalidRecurrence  = errors.New("timeslot: invalid recurrence rule")
    ErrInvalidICS         = errors.New("timeslot: invalid iCal format")
)

// Wrap errors with context
type SlotError struct {
    Op   string    // Operation that failed
    Slot TimeSlot  // Slot involved
    Err  error     // Underlying error
}

func (e *SlotError) Error() string
func (e *SlotError) Unwrap() error
```

---

## Testing Requirements

### Test Coverage
- Minimum 95% code coverage
- 100% coverage on core types (TimeSlot, SlotCollection, Availability)

### Test Categories

1. **Unit Tests** — Every public method
2. **Table-Driven Tests** — Use for all edge cases
3. **Property-Based Tests** — Use `testing/quick` for:
   - Slot operations (intersection, union)
   - Collection operations maintain invariants
   - Recurrence generation
4. **Timezone Tests** — Explicit tests for:
   - DST transitions
   - Cross-timezone operations
   - Timezone conversions
5. **Benchmark Tests** — For:
   - Slot finding with large availability
   - Collection operations
   - Conflict detection with many providers
6. **Fuzz Tests** — For:
   - iCal parsing
   - Recurrence rule parsing
   - Time range parsing
7. **Integration Tests** — End-to-end scenarios in `examples/`

### Test Fixtures

Create test fixtures in `testdata/`:
- Sample iCal files (valid and invalid)
- Complex availability scenarios
- Edge case time ranges

### Example Test Structure

```go
func TestTimeSlot_Overlaps(t *testing.T) {
    loc := time.UTC
    
    tests := []struct {
        name     string
        slot1    TimeSlot
        slot2    TimeSlot
        expected bool
    }{
        {
            name:     "completely before",
            slot1:    NewSlot(time.Date(2024, 1, 1, 9, 0, 0, 0, loc), time.Date(2024, 1, 1, 10, 0, 0, 0, loc)),
            slot2:    NewSlot(time.Date(2024, 1, 1, 11, 0, 0, 0, loc), time.Date(2024, 1, 1, 12, 0, 0, 0, loc)),
            expected: false,
        },
        {
            name:     "adjacent (no overlap)",
            slot1:    NewSlot(time.Date(2024, 1, 1, 9, 0, 0, 0, loc), time.Date(2024, 1, 1, 10, 0, 0, 0, loc)),
            slot2:    NewSlot(time.Date(2024, 1, 1, 10, 0, 0, 0, loc), time.Date(2024, 1, 1, 11, 0, 0, 0, loc)),
            expected: false,
        },
        {
            name:     "partial overlap",
            slot1:    NewSlot(time.Date(2024, 1, 1, 9, 0, 0, 0, loc), time.Date(2024, 1, 1, 11, 0, 0, 0, loc)),
            slot2:    NewSlot(time.Date(2024, 1, 1, 10, 0, 0, 0, loc), time.Date(2024, 1, 1, 12, 0, 0, 0, loc)),
            expected: true,
        },
        {
            name:     "slot1 contains slot2",
            slot1:    NewSlot(time.Date(2024, 1, 1, 9, 0, 0, 0, loc), time.Date(2024, 1, 1, 17, 0, 0, 0, loc)),
            slot2:    NewSlot(time.Date(2024, 1, 1, 10, 0, 0, 0, loc), time.Date(2024, 1, 1, 12, 0, 0, 0, loc)),
            expected: true,
        },
        {
            name:     "identical slots",
            slot1:    NewSlot(time.Date(2024, 1, 1, 9, 0, 0, 0, loc), time.Date(2024, 1, 1, 10, 0, 0, 0, loc)),
            slot2:    NewSlot(time.Date(2024, 1, 1, 9, 0, 0, 0, loc), time.Date(2024, 1, 1, 10, 0, 0, 0, loc)),
            expected: true,
        },
        // Add DST transition cases
        // Add cross-timezone cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.slot1.Overlaps(tt.slot2)
            if got != tt.expected {
                t.Errorf("Overlaps() = %v, want %v", got, tt.expected)
            }
            // Test symmetry
            gotReverse := tt.slot2.Overlaps(tt.slot1)
            if gotReverse != tt.expected {
                t.Errorf("Overlaps() (reverse) = %v, want %v", gotReverse, tt.expected)
            }
        })
    }
}
```

---

## Documentation Requirements

### README.md

Structure:
1. **Header** — Logo placeholder, badges (Go version, CI status, coverage, Go Report Card, GoDoc)
2. **One-line description**
3. **Features** — Bullet list of capabilities
4. **Installation** — `go get` command
5. **Quick Start** — Simple, copy-paste example
6. **Core Concepts** — Brief explanation of main types
7. **Examples** — Links to examples directory
8. **API Reference** — Link to pkg.go.dev
9. **Performance** — Benchmark results
10. **Contributing** — Link to CONTRIBUTING.md
11. **License**

### GoDoc

- Package-level documentation in `doc.go`
- Every exported type, function, and method must have documentation
- Include code examples in doc comments
- Use `Example` functions for complex usage

### CONTRIBUTING.md

Include:
- Development setup
- Code style (follow Go conventions)
- Testing requirements
- PR process
- Issue guidelines

### CHANGELOG.md

Follow Keep a Changelog format with:
- Unreleased section
- Semantic versioning

---

## CI/CD Configuration

### .github/workflows/ci.yml

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.22', '1.23']
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      
      - name: Verify dependencies
        run: go mod verify
      
      - name: Run go vet
        run: go vet ./...
      
      - name: Install staticcheck
        run: go install honnef.co/go/tools/cmd/staticcheck@latest
      
      - name: Run staticcheck
        run: staticcheck ./...
      
      - name: Install golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
          fail_ci_if_error: true
      
      - name: Run benchmarks
        run: go test -bench=. -benchmem ./... | tee benchmark.txt
      
      - name: Store benchmark result
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: benchmark.txt
          fail-on-alert: true

  fuzz:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Run fuzz tests
        run: |
          go test -fuzz=FuzzParseTimeRange -fuzztime=30s ./...
          go test -fuzz=FuzzParseRule -fuzztime=30s ./recurrence/...
          go test -fuzz=FuzzParseICS -fuzztime=30s ./ical/...
```

### .github/workflows/release.yml

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### .golangci.yml

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - goconst
    - gocyclo
    - gosec
    - prealloc
    - exportloopref
    - noctx
    - bodyclose
    - exhaustive

linters-settings:
  gocyclo:
    min-complexity: 15
  goconst:
    min-len: 3
    min-occurrences: 3
  misspell:
    locale: US
  exhaustive:
    default-signifies-exhaustive: true

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - gocyclo
```

---

## Makefile

```makefile
.PHONY: all build test test-race test-cover bench lint fmt vet clean help

all: lint test build

build:
	go build ./...

test:
	go test -v ./...

test-race:
	go test -v -race ./...

test-cover:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	goimports -w .

vet:
	go vet ./...

fuzz:
	go test -fuzz=FuzzParseTimeRange -fuzztime=1m ./...

clean:
	rm -f coverage.out coverage.html
	go clean -testcache

help:
	@echo "Available targets:"
	@echo "  build      - Build the library"
	@echo "  test       - Run tests"
	@echo "  test-race  - Run tests with race detector"
	@echo "  test-cover - Run tests with coverage"
	@echo "  bench      - Run benchmarks"
	@echo "  lint       - Run linters"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  fuzz       - Run fuzz tests"
	@echo "  clean      - Clean build artifacts"
```

---

## Example Implementations

### examples/basic/main.go

```go
package main

import (
	"fmt"
	"time"

	"github.com/[username]/timeslot"
	"github.com/[username]/timeslot/availability"
	"github.com/[username]/timeslot/provider"
	"github.com/[username]/timeslot/query"
)

func main() {
	// Create a weekly schedule: available Tue & Thu, 9am-5pm
	weekly := availability.NewWeeklySchedule(time.UTC).
		SetDay(time.Tuesday, availability.TimeRange{
			Start: availability.NewTimeOfDay(9, 0, 0),
			End:   availability.NewTimeOfDay(17, 0, 0),
		}).
		SetDay(time.Thursday, availability.TimeRange{
			Start: availability.NewTimeOfDay(10, 0, 0),
			End:   availability.NewTimeOfDay(14, 0, 0),
		})

	// Create a provider with the schedule
	p := provider.NewProvider("stylist-1",
		provider.WithWeeklySchedule(weekly),
		provider.WithBuffer(15*time.Minute),
		provider.WithMinNotice(2*time.Hour),
	)

	// Block out some dates
	holidays := []time.Time{
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
	}
	p = p.WithBlockedDates(holidays...)

	// Find available 1-hour slots in the next 2 weeks
	q := query.NewQuery().
		Duration(1 * time.Hour).
		InNext(14 * 24 * time.Hour).
		OnlyMornings().
		PreferEarlier().
		Limit(5).
		Build()

	slots, err := p.FindSlots(q)
	if err != nil {
		panic(err)
	}

	fmt.Println("Available slots:")
	for _, slot := range slots {
		fmt.Printf("  %s - %s\n",
			slot.Start.Format("Mon Jan 2 15:04"),
			slot.End.Format("15:04"))
	}
}
```

### examples/multi-provider/main.go

```go
package main

import (
	"fmt"
	"time"

	"github.com/[username]/timeslot/availability"
	"github.com/[username]/timeslot/conflict"
	"github.com/[username]/timeslot/provider"
	"github.com/[username]/timeslot/query"
)

func main() {
	// Create multiple providers
	providers := []*provider.Provider{
		createProvider("alice", time.Tuesday, time.Wednesday, time.Thursday),
		createProvider("bob", time.Monday, time.Tuesday, time.Friday),
		createProvider("carol", time.Wednesday, time.Thursday, time.Friday),
	}

	// Find times when ANY provider is available
	detector := conflict.NewDetector(providers...)
	
	q := query.NewQuery().
		Duration(30 * time.Minute).
		InNext(7 * 24 * time.Hour).
		Build()

	// Get slots grouped by provider
	results := detector.FindAvailableSlots(q)
	
	for providerID, slots := range results {
		fmt.Printf("%s has %d available slots\n", providerID, len(slots))
	}

	// Find times when ALL providers are available (for a meeting)
	commonSlots := detector.FindCommonAvailability(q)
	fmt.Printf("\nTimes when everyone is available: %d slots\n", len(commonSlots))
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
```

---

## Implementation Order

Build the library in this order to ensure dependencies are met:

1. **Phase 1: Foundation**
   - `internal/timeutil` — Time helper functions
   - `internal/validate` — Validation helpers  
   - `errors.go` — Error types
   - `timezone/timezone.go` — Timezone utilities

2. **Phase 2: Core Types**
   - `slot/slot.go` — TimeSlot struct
   - `slot/collection.go` — SlotCollection
   - `availability/weekly.go` — WeeklySchedule
   - `availability/exceptions.go` — Exception handling
   - `availability/availability.go` — Availability combining all

3. **Phase 3: Provider & Query**
   - `provider/options.go` — Functional options
   - `provider/provider.go` — Provider entity
   - `query/constraints.go` — Constraint types
   - `query/query.go` — Query builder
   - `query/optimizer.go` — Slot optimization

4. **Phase 4: Advanced Features**
   - `conflict/buffer.go` — Buffer time handling
   - `conflict/detector.go` — Conflict detection
   - `conflict/resolution.go` — Resolution strategies
   - `provider/multi.go` — Multi-provider operations
   - `recurrence/rule.go` — Recurrence rules
   - `recurrence/parser.go` — RRULE parsing
   - `recurrence/generator.go` — Occurrence generation

5. **Phase 5: Integrations**
   - `ical/parser.go` — iCal parsing
   - `ical/exporter.go` — iCal export

6. **Phase 6: Polish**
   - `timeslot.go` — Main entry point with re-exports
   - `doc.go` — Package documentation
   - `version.go` — Version info
   - All documentation files
   - All examples
   - CI/CD configuration

---

## Performance Requirements

### Benchmarks to Include

```go
func BenchmarkSlotCollection_FindOverlaps(b *testing.B) {
    // Test with 1000 slots
}

func BenchmarkProvider_FindSlots_LargeRange(b *testing.B) {
    // Find slots in a 1-year range
}

func BenchmarkConflictDetector_100Providers(b *testing.B) {
    // Detect conflicts with 100 providers
}

func BenchmarkICalParse_LargeCalendar(b *testing.B) {
    // Parse a calendar with 10,000 events
}
```

### Performance Targets

- `SlotCollection.FindOverlaps`: O(log n) using binary search
- `Provider.FindSlots` for 1 year: < 1ms
- `ConflictDetector` with 100 providers: < 10ms
- iCal parsing 10k events: < 100ms

---

## Security Considerations

- No unsafe pointer operations
- Validate all external input (iCal files, time strings)
- Prevent integer overflow in duration calculations
- Handle timezone database loading failures gracefully
- No goroutines that can leak
- Document thread-safety guarantees for each type

---

## Final Checklist Before v1.0.0

- [ ] All tests passing with 95%+ coverage
- [ ] All benchmarks documented with results
- [ ] golangci-lint passing with zero issues
- [ ] All public APIs documented with examples
- [ ] README complete with badges
- [ ] CHANGELOG has all changes
- [ ] Examples compile and run
- [ ] CI/CD pipelines working
- [ ] Security scan passing (gosec)
- [ ] Fuzz tests running without crashes
- [ ] go mod tidy produces no changes
- [ ] Version tagged as v1.0.0

---

## Commands to Execute

Start by running these commands in order:

```bash
# Initialize the project
mkdir -p timeslot && cd timeslot
go mod init github.com/[username]/timeslot

# Create directory structure
mkdir -p .github/workflows .github/ISSUE_TEMPLATE
mkdir -p availability slot provider query conflict recurrence ical timezone
mkdir -p internal/timeutil internal/validate
mkdir -p examples/basic examples/multi-provider examples/recurring-availability examples/booking-system
mkdir -p testdata/calendars testdata/fixtures

# Start implementing in order specified above
# Run tests frequently
go test -v ./...

# Check coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

**Remember:** This is a zero-dependency library. Every feature must be implemented using only the Go standard library. No external packages allowed except for development tooling (linters, etc.).
