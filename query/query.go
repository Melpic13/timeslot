package query

import (
	"fmt"
	"time"
)

// Query defines search criteria for finding available slots.
type Query struct {
	Duration    time.Duration
	From        time.Time
	To          time.Time
	Constraints []Constraint
	Preferences []Preference
	Limit       int
	Location    *time.Location
}

// QueryBuilder provides fluent API.
type QueryBuilder struct {
	query Query
}

func NewQuery() *QueryBuilder {
	return &QueryBuilder{query: Query{Location: time.UTC}}
}

func (b *QueryBuilder) Duration(d time.Duration) *QueryBuilder {
	b.query.Duration = d
	return b
}

func (b *QueryBuilder) Between(from, to time.Time) *QueryBuilder {
	b.query.From = from
	b.query.To = to
	return b
}

func (b *QueryBuilder) From(t time.Time) *QueryBuilder {
	b.query.From = t
	return b
}

func (b *QueryBuilder) To(t time.Time) *QueryBuilder {
	b.query.To = t
	return b
}

func (b *QueryBuilder) InNext(d time.Duration) *QueryBuilder {
	now := time.Now().In(b.locationOrUTC())
	b.query.From = now
	b.query.To = now.Add(d)
	return b
}

func (b *QueryBuilder) OnWeekdays(days ...time.Weekday) *QueryBuilder {
	b.query.Constraints = append(b.query.Constraints, NewWeekdayConstraint(days...))
	return b
}

func (b *QueryBuilder) OnlyMornings() *QueryBuilder {
	start := NewTimeOfDay(0, 0, 0)
	end := NewTimeOfDay(11, 59, 59)
	b.query.Constraints = append(b.query.Constraints, TimeOfDayConstraint{Start: &start, End: &end})
	return b
}

func (b *QueryBuilder) OnlyAfternoons() *QueryBuilder {
	start := NewTimeOfDay(12, 0, 0)
	end := NewTimeOfDay(17, 0, 0)
	b.query.Constraints = append(b.query.Constraints, TimeOfDayConstraint{Start: &start, End: &end})
	return b
}

func (b *QueryBuilder) OnlyEvenings() *QueryBuilder {
	start := NewTimeOfDay(17, 0, 0)
	end := NewTimeOfDay(23, 59, 59)
	b.query.Constraints = append(b.query.Constraints, TimeOfDayConstraint{Start: &start, End: &end})
	return b
}

func (b *QueryBuilder) NotBefore(tod TimeOfDay) *QueryBuilder {
	b.query.Constraints = append(b.query.Constraints, NotBeforeConstraint{Time: tod})
	return b
}

func (b *QueryBuilder) NotAfter(tod TimeOfDay) *QueryBuilder {
	b.query.Constraints = append(b.query.Constraints, NotAfterConstraint{Time: tod})
	return b
}

func (b *QueryBuilder) WithConstraint(c Constraint) *QueryBuilder {
	if c != nil {
		b.query.Constraints = append(b.query.Constraints, c)
	}
	return b
}

func (b *QueryBuilder) PreferEarlier() *QueryBuilder {
	b.query.Preferences = append(b.query.Preferences, preferEarlier{})
	return b
}

func (b *QueryBuilder) PreferLater() *QueryBuilder {
	b.query.Preferences = append(b.query.Preferences, preferLater{})
	return b
}

func (b *QueryBuilder) PreferTime(tod TimeOfDay) *QueryBuilder {
	b.query.Preferences = append(b.query.Preferences, preferTime{Target: tod})
	return b
}

func (b *QueryBuilder) Limit(n int) *QueryBuilder {
	b.query.Limit = n
	return b
}

func (b *QueryBuilder) InTimezone(loc *time.Location) *QueryBuilder {
	if loc != nil {
		b.query.Location = loc
	}
	return b
}

func (b *QueryBuilder) Build() Query {
	q := b.query
	if q.Location == nil {
		q.Location = time.UTC
	}
	if q.From.IsZero() {
		q.From = time.Now().In(q.Location)
	}
	if q.To.IsZero() && q.Duration > 0 {
		q.To = q.From.Add(7 * 24 * time.Hour)
	}
	return q
}

func (b *QueryBuilder) Validate() error {
	return b.Build().Validate()
}

func (q Query) Validate() error {
	if q.Duration <= 0 {
		return fmt.Errorf("query: duration must be positive")
	}
	if !q.To.After(q.From) {
		return fmt.Errorf("query: to must be after from")
	}
	return nil
}

func (b *QueryBuilder) locationOrUTC() *time.Location {
	if b.query.Location == nil {
		return time.UTC
	}
	return b.query.Location
}
