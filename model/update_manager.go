package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type UpdateManager struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	transactionProvider TransactionProvider
}

func NewUpdateManager(transactionProvider TransactionProvider) func(UpdateEvent) {
	manager := &UpdateManager{transactionProvider: transactionProvider}
	return manager.Update
}

// Test method
func newUpdateManager(transactionProvider TransactionProvider) *UpdateManager {
	return &UpdateManager{transactionProvider: transactionProvider}
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

	tx := manager.transactionProvider.NewTransaction()

	stopArea, found := tx.Model().StopAreas().FindByObjectId(event.ObjectId)
	if !found {
		stopArea = tx.Model().StopAreas().New()

		stopArea.SetObjectID(event.ObjectId)
		stopArea.CollectGeneralMessages = true

		stopArea.Name = event.Name
		stopArea.CollectedAlways = event.CollectedAlways
		stopArea.Longitude = event.Longitude
		stopArea.Latitude = event.Latitude
	}

	if stopArea.ParentId == "" && event.ParentObjectId.Value() != "" {
		parentSA, _ := tx.Model().StopAreas().FindByObjectId(event.ParentObjectId)
		stopArea.ParentId = parentSA.Id()
	}

	stopArea.Updated(manager.Clock().Now())

	tx.Model().StopAreas().Save(&stopArea)
	tx.Commit()
	tx.Close()

	if event.Origin != "" {
		status, ok := stopArea.Origins.Origin(event.Origin)
		if !status || !ok {
			manager.updateStatus(NewStatusUpdateEvent(stopArea.Id(), event.Origin, true))
		}
	}
}

func (manager *UpdateManager) updateMonitoredStopArea(stopAreaId StopAreaId, partner string, status bool) {
	tx := manager.transactionProvider.NewTransaction()

	ascendants := tx.Model().StopAreas().FindAscendants(stopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.SetPartnerStatus(partner, status)
		tx.Model().StopAreas().Save(&stopArea)
	}

	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateLine(event *LineUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a Line with an empty objectid
		return
	}

	tx := manager.transactionProvider.NewTransaction()

	line, found := tx.Model().Lines().FindByObjectId(event.ObjectId)
	if !found {
		line = tx.Model().Lines().New()

		line.SetObjectID(event.ObjectId)
		line.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

		line.CollectGeneralMessages = true

		line.Name = event.Name

		line.SetOrigin(event.Origin)
	}

	line.Updated(manager.Clock().Now())

	tx.Model().Lines().Save(&line)
	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateVehicleJourney(event *VehicleJourneyUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a VehicleJourney with an empty objectid
		return
	}

	tx := manager.transactionProvider.NewTransaction()

	vj, found := tx.Model().VehicleJourneys().FindByObjectId(event.ObjectId)
	if !found {
		// LineObjectId
		l, ok := tx.Model().Lines().FindByObjectId(event.LineObjectId)
		if !ok {
			logger.Log.Debugf("VehicleJourney update event without corresponding line: %v", event.LineObjectId.String())
			return
		}

		vj = tx.Model().VehicleJourneys().New()

		vj.SetObjectID(event.ObjectId)
		vj.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

		vj.Origin = event.Origin
		vj.Name = event.Attributes()["VehicleJourneyName"]
		vj.LineId = l.Id()

		vj.Attributes = event.Attributes()
		vj.References = event.References()
	}

	vj.References.SetObjectId("OriginRef", NewObjectID(event.ObjectId.Kind(), event.OriginRef))
	vj.OriginName = event.OriginName

	vj.References.SetObjectId("DestinationRef", NewObjectID(event.ObjectId.Kind(), event.DestinationRef))
	vj.DestinationName = event.DestinationName

	if event.Direction != "" {
		vj.Attributes.Set("DirectionName", event.Direction)
	}

	vj.Monitored = event.Monitored

	tx.Model().VehicleJourneys().Save(&vj)
	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateStopVisit(event *StopVisitUpdateEvent) {
	if event.ObjectId.Value() == "" { // Avoid creating a StopVisit with an empty objectid
		return
	}

	tx := manager.transactionProvider.NewTransaction()

	vj, ok := tx.Model().VehicleJourneys().FindByObjectId(event.VehicleJourneyObjectId)
	if !ok {
		logger.Log.Debugf("StopVisit update event without corresponding vehicle journey: %v", event.VehicleJourneyObjectId.String())
		return
	}

	var sa StopArea
	var sv StopVisit
	if event.StopAreaObjectId.Value() == "" {
		sv, ok = tx.Model().StopVisits().FindByObjectId(event.ObjectId)
		if !ok {
			logger.Log.Debugf("Can't find Stopvisit from update event without stop area id")
			return
		}
		sa = sv.StopArea()
	} else {
		sa, ok = tx.Model().StopAreas().FindByObjectId(event.StopAreaObjectId)
		if !ok {
			logger.Log.Debugf("StopVisit update event without corresponding stop area: %v", event.StopAreaObjectId.String())
			return
		}

		sv, ok = tx.Model().StopVisits().FindByObjectId(event.ObjectId)
		if !ok {
			sv = tx.Model().StopVisits().New()

			sv.SetObjectID(event.ObjectId)
			sv.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

			sv.StopAreaId = sa.Id()
			sv.VehicleJourneyId = vj.Id()

			sv.Origin = event.Origin
			sv.PassageOrder = event.PassageOrder
			sv.DataFrameRef = event.DataFrameRef

			sv.Attributes = event.Attributes()
			sv.References = event.References()
		}
	}

	// Update StopArea Lines
	l := vj.Line()
	if l != nil {
		sa.LineIds.Add(l.Id())
		referent, ok := tx.Model().StopAreas().Find(sa.ReferentId)
		if ok {
			referent.LineIds.Add(l.Id())
			tx.Model().StopAreas().Save(&referent)
		}
		tx.Model().StopAreas().Save(&sa)
	}

	if !event.RecordedAt.IsZero() {
		sv.RecordedAt = event.RecordedAt
	} else if !sv.Schedules.Include(&event.Schedules) {
		sv.RecordedAt = manager.Clock().Now()
	}

	sv.Schedules.Merge(&event.Schedules)
	sv.DepartureStatus = event.DepartureStatus
	sv.ArrivalStatus = event.ArrivalStatus
	sv.VehicleAtStop = event.VehicleAtStop
	sv.Collected(manager.Clock().Now())

	if event.Monitored != vj.Monitored {
		vj.Monitored = event.Monitored
		tx.Model().VehicleJourneys().Save(&vj)
	}

	if event.Origin != "" {
		status, ok := sa.Origins.Origin(event.Origin)
		if status != event.Monitored || !ok {
			manager.updateMonitoredStopArea(sa.Id(), event.Origin, event.Monitored)
		}
	}

	tx.Model().StopVisits().Save(&sv)
	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateVehicle(event *VehicleUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()

	vj, _ := tx.Model().VehicleJourneys().FindByObjectId(event.VehicleJourneyObjectId)
	line := vj.Line()

	vehicle, found := tx.Model().Vehicles().FindByObjectId(event.ObjectId)
	if !found {
		vehicle = tx.Model().Vehicles().New()

		vehicle.SetObjectID(event.ObjectId)
	}

	vehicle.VehicleJourneyId = vj.Id()
	vehicle.Longitude = event.Longitude
	vehicle.Latitude = event.Latitude
	vehicle.Bearing = event.Bearing
	vehicle.RecordedAtTime = manager.Clock().Now()

	if line != nil {
		vehicle.LineId = line.Id()
	}

	tx.Model().Vehicles().Save(&vehicle)
	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateStatus(event *StatusUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()

	ascendants := tx.Model().StopAreas().FindAscendants(event.StopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.SetPartnerStatus(event.Partner, event.Status)
		tx.Model().StopAreas().Save(&stopArea)
	}

	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateNotCollected(event *NotCollectedUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()

	stopVisit, found := tx.Model().StopVisits().FindByObjectId(event.ObjectId)
	if !found {
		logger.Log.Debugf("StopVisitNotCollectedEvent on unknown StopVisit: %#v", event)
		tx.Close()
		return
	}

	stopVisit.NotCollected()
	tx.Model().StopVisits().Save(&stopVisit)

	logger.Log.Debugf("StopVisit not Collected: %s (%v)", stopVisit.Id(), event.ObjectId)

	tx.Commit()
	tx.Close()
}
