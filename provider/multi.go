package provider

import (
	"github.com/Melpic13/timeslot/query"
	"github.com/Melpic13/timeslot/slot"
)

func FindAny(providers []*Provider, q query.Query) map[string][]slot.TimeSlot {
	out := make(map[string][]slot.TimeSlot, len(providers))
	for _, p := range providers {
		slots, err := p.FindSlots(q)
		if err != nil {
			continue
		}
		out[p.ID] = slots
	}
	return out
}

func FindCommon(providers []*Provider, q query.Query) []slot.TimeSlot {
	if len(providers) == 0 {
		return nil
	}
	first, err := providers[0].FindSlots(q)
	if err != nil {
		return nil
	}
	common := slot.NewCollection(first...)
	for i := 1; i < len(providers); i++ {
		next, err := providers[i].FindSlots(q)
		if err != nil {
			return nil
		}
		common = common.Intersect(slot.NewCollection(next...))
	}
	return common.Slots()
}
