package core

import (
	"bitbucket.org/enroute-mobi/ara/model"
)

type CollectUpdateEvents struct {
	StopAreas       map[string]*model.StopAreaUpdateEvent
	Lines           map[string]*model.LineUpdateEvent
	VehicleJourneys map[string]*model.VehicleJourneyUpdateEvent
	StopVisits      map[string]map[string]*model.StopVisitUpdateEvent
	Vehicles        map[string]*model.VehicleUpdateEvent
	Situations      []*model.SituationUpdateEvent
	Facilities      map[string]*model.FacilityUpdateEvent
	Cancellations   []*model.NotCollectedUpdateEvent
	*CollectedRefs
}

type CollectedRefs struct {
	MonitoringRefs     map[string]struct{}
	LineRefs           map[string]struct{}
	VehicleJourneyRefs map[string]struct{}
	VehicleRefs        map[string]struct{}
}

func NewCollectedRefs() *CollectedRefs {
	return &CollectedRefs{
		MonitoringRefs:     make(map[string]struct{}),
		LineRefs:           make(map[string]struct{}),
		VehicleJourneyRefs: make(map[string]struct{}),
		VehicleRefs:        make(map[string]struct{}),
	}
}
func NewCollectUpdateEvents() *CollectUpdateEvents {
	collectedUpdateEvents := &CollectUpdateEvents{
		StopAreas:       make(map[string]*model.StopAreaUpdateEvent),
		Lines:           make(map[string]*model.LineUpdateEvent),
		VehicleJourneys: make(map[string]*model.VehicleJourneyUpdateEvent),
		StopVisits:      make(map[string]map[string]*model.StopVisitUpdateEvent),
		Vehicles:        make(map[string]*model.VehicleUpdateEvent),
		Facilities:      make(map[string]*model.FacilityUpdateEvent),
	}
	collectedUpdateEvents.CollectedRefs = NewCollectedRefs()
	return collectedUpdateEvents
}

func (events *CollectedRefs) GetLines() []string {
	return GetModelReferenceSlice(events.LineRefs)
}

func (events *CollectedRefs) GetVehicleJourneys() []string {
	return GetModelReferenceSlice(events.VehicleJourneyRefs)
}

func (events *CollectedRefs) GetStopAreas() []string {
	return GetModelReferenceSlice(events.MonitoringRefs)
}

func (events *CollectedRefs) GetVehicles() []string {
	return GetModelReferenceSlice(events.VehicleRefs)
}
