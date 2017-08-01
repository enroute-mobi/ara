package core

import "github.com/af83/edwig/model"

type lastState interface {
	Ischanged(interface{}) bool
}

type stopVisitLastChange struct {
	Schedules model.StopVisitSchedules
}

type vehicleJourneyLastChange struct {
}

type lineLastChange struct {
}

func (sve *stopVisitLastChange) Ischanged(svbeInterface interface{}) bool {
	// stopVisitBroadcastEvent := svbeInterface.(model.StopVisitBroadcastEvent)
	// stopVisitBroadcastEvent.Schedules == sve.Schedules
	return false
}
