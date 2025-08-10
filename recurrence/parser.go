package recurrence

func Parse(rrule string) (*Rule, error) {
	return ParseRule(rrule)
}
