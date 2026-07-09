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

func (es *CollectUpdateEvents) StopAreasUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.StopAreas {
		evs = append(evs, es.StopAreas[i])
	}

	return evs
}

func (es *CollectUpdateEvents) LinesUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.Lines {
		evs = append(evs, es.Lines[i])
	}

	return evs
}

func (es *CollectUpdateEvents) VehicleJourneysUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.VehicleJourneys {
		evs = append(evs, es.VehicleJourneys[i])
	}

	return evs
}

func (es *CollectUpdateEvents) StopVisitsUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.StopVisits {
		for j := range es.StopVisits[i] {
			evs = append(evs, es.StopVisits[i][j])
		}
	}

	return evs
}

func (es *CollectUpdateEvents) SituationsUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.Situations {
		evs = append(evs, es.Situations[i])
	}

	return evs
}

func (es *CollectUpdateEvents) VehiclesUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.Vehicles {
		evs = append(evs, es.Vehicles[i])
	}

	return evs
}

func (es *CollectUpdateEvents) FacilitiesUpdateEvents() []model.UpdateEvent {
	evs := []model.UpdateEvent{}

	for i := range es.Facilities {
		evs = append(evs, es.Facilities[i])
	}

	return evs
}

type CollectedRefs struct {
	MonitoringRefs     map[string]struct{}
	LineRefs           map[string]struct{}
	VehicleJourneyRefs map[string]struct{}
	VehicleRefs        map[string]struct{}
	FacilityRefs       map[string]struct{}
}

func NewCollectedRefs() *CollectedRefs {
	return &CollectedRefs{
		MonitoringRefs:     make(map[string]struct{}),
		LineRefs:           make(map[string]struct{}),
		VehicleJourneyRefs: make(map[string]struct{}),
		VehicleRefs:        make(map[string]struct{}),
		FacilityRefs:       make(map[string]struct{}),
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

func (events *CollectedRefs) GetFacilities() []string {
	return GetModelReferenceSlice(events.FacilityRefs)
}
