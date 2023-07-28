package core

import (
	"bitbucket.org/enroute-mobi/ara/model"
)

type CollectUpdateEvents struct {
	StopAreas          map[string]*model.StopAreaUpdateEvent
	Lines              map[string]*model.LineUpdateEvent
	VehicleJourneys    map[string]*model.VehicleJourneyUpdateEvent
	StopVisits         map[string]map[string]*model.StopVisitUpdateEvent
	Vehicles           map[string]*model.VehicleUpdateEvent
	Cancellations      []*model.NotCollectedUpdateEvent
	MonitoringRefs     map[string]struct{}
	LineRefs           map[string]struct{}
	VehicleJourneyRefs map[string]struct{}
}

func NewCollectUpdateEvents() *CollectUpdateEvents {
	return &CollectUpdateEvents{
		StopAreas:          make(map[string]*model.StopAreaUpdateEvent),
		Lines:              make(map[string]*model.LineUpdateEvent),
		VehicleJourneys:    make(map[string]*model.VehicleJourneyUpdateEvent),
		StopVisits:         make(map[string]map[string]*model.StopVisitUpdateEvent),
		Vehicles:           make(map[string]*model.VehicleUpdateEvent),
		MonitoringRefs:     make(map[string]struct{}),
		LineRefs:           make(map[string]struct{}),
		VehicleJourneyRefs: make(map[string]struct{}),
	}
}

func (events CollectUpdateEvents) GetLines() []string {
	return GetModelReferenceSlice(events.LineRefs)
}

func (events CollectUpdateEvents) GetVehicleJourneys() []string {
	return GetModelReferenceSlice(events.VehicleJourneyRefs)
}

func (events CollectUpdateEvents) GetStopAreas() []string {
	return GetModelReferenceSlice(events.MonitoringRefs)
}
