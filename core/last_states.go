package core

import "github.com/af83/edwig/siri"

type lastState interface {
	Haschanged(interface{}) bool
}

type stopMonitoringLastChange struct {
	lastChange siri.SIRIMonitoredStopVisit
}

type vehicleJourneyLastChange struct {
}

type lineLastChange struct {
}

func (sve *stopMonitoringLastChange) Haschanged(svbeInterface interface{}) bool {
	// stopVisitBroadcastEvent := svbeInterface.(model.StopVisitBroadcastEvent)
	// stopVisitBroadcastEvent.Schedules == sve.Schedules
	return true
}
