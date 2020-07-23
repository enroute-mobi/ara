package model

import "bitbucket.org/enroute-mobi/ara/logger"

type UpdateManager struct {
	ClockConsumer
	UUIDConsumer

	transactionProvider TransactionProvider
}

func NewUpdateManager(transactionProvider TransactionProvider) func(UpdateEvent) {
	manager := &UpdateManager{transactionProvider: transactionProvider}
	return manager.Update
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
	}
}

func (manager *UpdateManager) updateStopArea(event *StopAreaUpdateEvent) {
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
}

func (manager *UpdateManager) updateMonitoredStopArea(stopAreaId StopAreaId, partner string, status bool) {
	tx := manager.transactionProvider.NewTransaction()

	ascendants := tx.Model().StopAreas().FindAscendants(stopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.Origins.SetPartnerStatus(partner, status)
		stopArea.Monitored = stopArea.Origins.Monitored()
		tx.Model().StopAreas().Save(&stopArea)
	}

	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateLine(event *LineUpdateEvent) {
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
	tx := manager.transactionProvider.NewTransaction()

	vj, found := tx.Model().VehicleJourneys().FindByObjectId(event.ObjectId)
	if !found {
		vj = tx.Model().VehicleJourneys().New()

		vj.SetObjectID(event.ObjectId)
		vj.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

		// LineObjectId
		l, ok := tx.Model().Lines().FindByObjectId(event.LineObjectId)
		if !ok {
			logger.Log.Debugf("VehicleJourney update event without corresponding line: %v", event.LineObjectId.String())
			return
		}
		vj.LineId = l.Id()

		vj.Origin = event.Origin
	}

	vj.References.SetObjectId("OriginRef", NewObjectID(event.ObjectId.Kind(), event.OriginRef))
	vj.OriginName = event.OriginName

	vj.References.SetObjectId("DestinationRef", NewObjectID(event.ObjectId.Kind(), event.DestinationRef))
	vj.DestinationName = event.DestinationName

	vj.Attributes.Set("DirectionName", event.Direction)

	tx.Model().VehicleJourneys().Save(&vj)
	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateStopVisit(event *StopVisitUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()

	sa, ok := tx.Model().StopAreas().FindByObjectId(event.StopAreaObjectId)
	if !ok {
		logger.Log.Debugf("StopVisit update event without corresponding stop area: %v", event.StopAreaObjectId.String())
		return
	}

	vj, ok := tx.Model().VehicleJourneys().FindByObjectId(event.VehicleJourneyObjectId)
	if !ok {
		logger.Log.Debugf("StopVisit update event without corresponding vehicle journey: %v", event.VehicleJourneyObjectId.String())
		return
	}

	sv, found := tx.Model().StopVisits().FindByObjectId(event.ObjectId)
	if !found {
		sv = tx.Model().StopVisits().New()

		sv.SetObjectID(event.ObjectId)
		sv.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))

		// StopAreaObjectId
		sv.StopAreaId = sa.Id()

		// VehicleJourneyObjectId
		sv.VehicleJourneyId = vj.Id()

		sv.Origin = event.Origin

		sv.PassageOrder = event.PassageOrder
	}

	if sv.Schedules.Eq(&event.Schedules) {
		sv.RecordedAt = manager.Clock().Now()
	}
	sv.Schedules.Merge(&event.Schedules)
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
