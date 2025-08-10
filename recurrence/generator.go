package recurrence

import "time"

func GenerateOccurrences(rule *Rule, start time.Time, limit int) []time.Time {
	if rule == nil {
		return nil
	}
	return rule.Generate(start, limit)
}
