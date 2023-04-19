package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"

	"golang.org/x/exp/maps"
)

type UpdateManager struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	model Model
}

func NewUpdateManager(model Model) func(UpdateEvent) {
	manager := newUpdateManager(model)
	return manager.Update
}

// Test method
func newUpdateManager(model Model) *UpdateManager {
	return &UpdateManager{model: model}
}

func (manager *UpdateManager) Update(event UpdateEvent) {
	switch event.EventKind() {
	case STOP_AREA_EVENT:
		manager.updateStopArea(event.(*StopAreaUpdateEvent))
	case LINE_EVENT:
		manager.updateLine(event.(*LineUpdateEvent))
	case VEHICLE_JOURNEY_EVENT:
		manager.updateVehicleJourney(event.(*VehicleJourneyUpdateEvent))
	case STOP_VISIT_EVENT:
		manager.updateStopVisit(event.(*StopVisitUpdateEvent))
	case VEHICLE_EVENT:
		manager.updateVehicle(event.(*VehicleUpdateEvent))
	case STATUS_EVENT:
		manager.updateStatus(event.(*StatusUpdateEvent))
	case NOT_COLLECTED_EVENT:
		manager.updateNotCollected(event.(*NotCollectedUpdateEvent))
	}
}

func (manager *UpdateManager) updateStopArea(event *StopAreaUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a StopArea with an empty objectid
		return
	}

	stopArea, found := manager.model.StopAreas().FindByObjectId(event.ObjectId)
	if !found {
		stopArea = manager.model.StopAreas().New()

		stopArea.SetObjectID(event.ObjectId)
		stopArea.CollectGeneralMessages = true

		stopArea.Name = event.Name
		stopArea.CollectedAlways = event.CollectedAlways
		stopArea.Longitude = event.Longitude
		stopArea.Latitude = event.Latitude
	}

	if stopArea.ParentId == "" && event.ParentObjectId.Value() != "" {
		parentSA, _ := manager.model.StopAreas().FindByObjectId(event.ParentObjectId)
		stopArea.ParentId = parentSA.Id()
	}

	stopArea.Updated(manager.Clock().Now())

	manager.model.StopAreas().Save(stopArea)
	if event.Origin != "" {
		status, ok := stopArea.Origins.Origin(event.Origin)
		if !status || !ok {
			manager.updateStatus(NewStatusUpdateEvent(stopArea.Id(), event.Origin, true))
		}
	}
}

func (manager *UpdateManager) updateMonitoredStopArea(stopAreaId StopAreaId, partner string, status bool) {
	ascendants := manager.model.StopAreas().FindAscendants(stopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.SetPartnerStatus(partner, status)
		manager.model.StopAreas().Save(stopArea)
	}
}

func (manager *UpdateManager) updateLine(event *LineUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a Line with an empty objectid
		return
	}

	line, found := manager.model.Lines().FindByObjectId(event.ObjectId)
	if !found {
		line = manager.model.Lines().New()

		line.SetObjectID(event.ObjectId)
		line.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

		line.CollectGeneralMessages = true

		line.Name = event.Name

		line.SetOrigin(event.Origin)
	}

	line.Updated(manager.Clock().Now())

	manager.model.Lines().Save(line)
}

func (manager *UpdateManager) updateVehicleJourney(event *VehicleJourneyUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a VehicleJourney with an empty objectid
		return
	}

	vj, found := manager.model.VehicleJourneys().FindByObjectId(event.ObjectId)
	if !found {
		// LineObjectId
		l, ok := manager.model.Lines().FindByObjectId(event.LineObjectId)
		if !ok {
			logger.Log.Debugf("VehicleJourney update event without corresponding line: %v", event.LineObjectId.String())
			return
		}

		vj = manager.model.VehicleJourneys().New()

		vj.SetObjectID(event.ObjectId)
		vj.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

		vj.Origin = event.Origin
		vj.Name = event.Attributes()["VehicleJourneyName"]
		vj.LineId = l.Id()
	}

	maps.Copy(vj.Attributes, event.Attributes())

	if vj.References.IsEmpty() {
		vj.References = event.References()
	}

	vj.References.SetObjectId("OriginRef", NewObjectID(event.ObjectId.Kind(), event.OriginRef))
	vj.OriginName = event.OriginName

	vj.References.SetObjectId("DestinationRef", NewObjectID(event.ObjectId.Kind(), event.DestinationRef))
	vj.DestinationName = event.DestinationName

	if event.Direction != "" { // Only used for Push collector
		vj.Attributes.Set("DirectionName", event.Direction)
	}

	vj.Occupancy = event.Occupancy
	vj.Monitored = event.Monitored
	if event.DirectionType != "" { // Do not override unknown DirectionType
		vj.DirectionType = event.DirectionType
	}

	manager.model.VehicleJourneys().Save(vj)
}

func (manager *UpdateManager) updateStopVisit(event *StopVisitUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a StopVisit with an empty objectid
		return
	}

	vj, ok := manager.model.VehicleJourneys().FindByObjectId(event.VehicleJourneyObjectId)
	if !ok {
		logger.Log.Debugf("StopVisit update event without corresponding vehicle journey: %v", event.VehicleJourneyObjectId.String())
		return
	}

	var sa *StopArea
	var sv *StopVisit
	if event.StopAreaObjectId.Value() == "" {
		sv, ok = manager.model.StopVisits().FindByObjectId(event.ObjectId)
		if !ok {
			logger.Log.Debugf("Can't find Stopvisit from update event without stop area id")
			return
		}
		sa = sv.StopArea()
		if sa == nil {
			logger.Log.Printf("StopVisit in memory without a StopArea: %v", sv.Id())
			return
		}
	} else {
		sa, ok = manager.model.StopAreas().FindByObjectId(event.StopAreaObjectId)
		if !ok {
			logger.Log.Debugf("StopVisit update event without corresponding stop area: %v", event.StopAreaObjectId.String())
			return
		}

		sv, ok = manager.model.StopVisits().FindByObjectId(event.ObjectId)
		if !ok {
			sv = manager.model.StopVisits().New()

			sv.SetObjectID(event.ObjectId)
			sv.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

			sv.StopAreaId = sa.Id()
			sv.VehicleJourneyId = vj.Id()

			sv.Origin = event.Origin
			sv.PassageOrder = event.PassageOrder
			sv.DataFrameRef = event.DataFrameRef
		}
	}

	if sv.Attributes.IsEmpty() {
		sv.Attributes = event.Attributes()
	}
	if sv.References.IsEmpty() {
		sv.References = event.References()
	}

	// Update StopArea Lines
	l := vj.Line()
	if l != nil {
		sa.LineIds.Add(l.Id())
		referent, ok := manager.model.StopAreas().Find(sa.ReferentId)
		if ok {
			referent.LineIds.Add(l.Id())
			manager.model.StopAreas().Save(referent)
		}
		manager.model.StopAreas().Save(sa)
	}

	if !event.RecordedAt.IsZero() {
		sv.RecordedAt = event.RecordedAt
	} else if !sv.Schedules.Include(event.Schedules) {
		sv.RecordedAt = manager.Clock().Now()
	}

	sv.Schedules.Merge(event.Schedules)
	if event.DepartureStatus != "" {
		sv.DepartureStatus = event.DepartureStatus
	}
	if event.ArrivalStatus != "" {
		sv.ArrivalStatus = event.ArrivalStatus
	}
	sv.VehicleAtStop = event.VehicleAtStop
	sv.Collected(manager.Clock().Now())

	if event.Monitored != vj.Monitored {
		vj.Monitored = event.Monitored
		manager.model.VehicleJourneys().Save(vj)
	}

	if event.Origin != "" {
		status, ok := sa.Origins.Origin(event.Origin)
		if status != event.Monitored || !ok {
			manager.updateMonitoredStopArea(sa.Id(), event.Origin, event.Monitored)
		}
	}

	manager.model.StopVisits().Save(sv)

	// VehicleJourney stop sequence
	if !vj.HasCompleteStopSequence {
		completeStopSequence := vj.model.ScheduledStopVisits().StopVisitsLenByVehicleJourney(vj.Id()) == vj.model.StopVisits().StopVisitsLenByVehicleJourney(vj.Id())
		if completeStopSequence {
			vj.HasCompleteStopSequence = true
			manager.model.VehicleJourneys().Save(vj)
		}
	}

	// long term historisation
	if sv.IsArchivable() {
		sva := &StopVisitArchiver{
			Model:     manager.model,
			StopVisit: sv,
		}
		sva.Archive()
	}

}

func (manager *UpdateManager) updateVehicle(event *VehicleUpdateEvent) {
	sa, _ := manager.model.StopAreas().FindByObjectId(event.StopAreaObjectId)
	vj, _ := manager.model.VehicleJourneys().FindByObjectId(event.VehicleJourneyObjectId)
	line := vj.Line()

	vehicle, found := manager.model.Vehicles().FindByObjectId(event.ObjectId)
	if !found {
		vehicle = manager.model.Vehicles().New()

		vehicle.SetObjectID(event.ObjectId)
	}

	vehicle.StopAreaId = sa.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.DriverRef = event.DriverRef
	vehicle.Longitude = event.Longitude
	vehicle.Latitude = event.Latitude
	vehicle.Bearing = event.Bearing
	vehicle.LinkDistance = event.LinkDistance
	vehicle.Percentage = event.Percentage
	vehicle.ValidUntilTime = event.ValidUntilTime
	if event.RecordedAt.IsZero() {
		vehicle.RecordedAtTime = manager.Clock().Now()
	} else {
		vehicle.RecordedAtTime = event.RecordedAt
	}
	vehicle.Occupancy = event.Occupancy

	if line != nil {
		vehicle.LineId = line.Id()
	}

	manager.model.Vehicles().Save(vehicle)
}

func (manager *UpdateManager) updateStatus(event *StatusUpdateEvent) {
	ascendants := manager.model.StopAreas().FindAscendants(event.StopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.SetPartnerStatus(event.Partner, event.Status)
		manager.model.StopAreas().Save(stopArea)
	}
}

func (manager *UpdateManager) updateNotCollected(event *NotCollectedUpdateEvent) {
	stopVisit, found := manager.model.StopVisits().FindByObjectId(event.ObjectId)
	if !found {
		logger.Log.Debugf("StopVisitNotCollectedEvent on unknown StopVisit: %#v", event)
		return
	}

	stopVisit.NotCollected()
	manager.model.StopVisits().Save(stopVisit)

	logger.Log.Debugf("StopVisit not Collected: %s (%v)", stopVisit.Id(), event.ObjectId)
}
