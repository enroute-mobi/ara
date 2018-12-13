package model

import "github.com/af83/edwig/logger"

type StopAreaUpdateManager struct {
	ClockConsumer
	UUIDConsumer

	transactionProvider TransactionProvider
}

type StopVisitUpdater struct {
	ClockConsumer

	tx    *Transaction
	event *StopVisitUpdateEvent
}

func NewStopAreaUpdateManager(transactionProvider TransactionProvider) func(*StopAreaUpdateEvent) {
	manager := newStopAreaUpdateManager(transactionProvider)
	return manager.UpdateStopArea
}

func newStopAreaUpdateManager(transactionProvider TransactionProvider) *StopAreaUpdateManager {
	return &StopAreaUpdateManager{transactionProvider: transactionProvider}
}

func (manager *StopAreaUpdateManager) UpdateStopArea(event *StopAreaUpdateEvent) {
	if event.StopAreaMonitoredEvent != nil {
		manager.UpdateMonitoredStopArea(event)
		return
	}

	tx := manager.transactionProvider.NewTransaction()

	stopArea, found := tx.Model().StopAreas().Find(event.StopAreaId)
	if !found {
		logger.Log.Debugf("StopAreaUpdateEvent for unknown StopArea %v", event.StopAreaId)

		stopArea = tx.Model().StopAreas().New()

		stopArea.SetObjectID(event.StopAreaAttributes.ObjectId)
		stopArea.Name = event.StopAreaAttributes.Name
		stopArea.CollectedAlways = event.StopAreaAttributes.CollectedAlways
		stopArea.CollectGeneralMessages = true
		stopArea.Monitored = true
		stopArea.id = StopAreaId(manager.NewUUID())

		event.StopAreaId = stopArea.id
	}

	if stopArea.ParentId == "" && event.StopAreaAttributes.ParentObjectId.Value() != "" {
		parentSA, _ := tx.Model().StopAreas().FindByObjectId(event.StopAreaAttributes.ParentObjectId)
		stopArea.ParentId = parentSA.Id()
	}

	stopArea.Updated(manager.Clock().Now())
	tx.Model().StopAreas().Save(&stopArea)
	tx.Commit()
	tx.Close()

	if event.Origin != "" {
		status, ok := stopArea.Origins.Origin(event.Origin)
		if !status || !ok {
			manager.UpdateMonitoredStopArea(NewStopAreaMonitoredEvent(manager.NewUUID(), event.StopAreaId, event.Origin, true))
		}
	}

	for _, stopVisitUpdateEvent := range event.StopVisitUpdateEvents {
		manager.UpdateStopVisit(stopVisitUpdateEvent)
	}
	for _, stopVisitNotCollectedEvent := range event.StopVisitNotCollectedEvents {
		manager.UpdateNotCollectedStopVisit(stopVisitNotCollectedEvent)
	}
}

func (manager *StopAreaUpdateManager) UpdateMonitoredStopArea(event *StopAreaUpdateEvent) {
	// Should never happen, but don't want to ever have a go nil pointer exception
	if event.StopAreaMonitoredEvent == nil {
		return
	}

	tx := manager.transactionProvider.NewTransaction()

	ascendants := tx.Model().StopAreas().FindAscendants(event.StopAreaId)
	for i, _ := range ascendants {
		stopArea := ascendants[i]
		stopArea.Origins.SetPartnerStatus(event.StopAreaMonitoredEvent.Partner, event.StopAreaMonitoredEvent.Status)
		stopArea.Monitored = stopArea.Origins.Monitored()
		tx.Model().StopAreas().Save(&stopArea)
	}

	tx.Commit()
	tx.Close()
}

func (manager *StopAreaUpdateManager) UpdateStopVisit(event *StopVisitUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()
	defer tx.Close()

	NewStopVisitUpdater(tx, event).Update()

	tx.Commit()
}

func (manager *StopAreaUpdateManager) UpdateNotCollectedStopVisit(event *StopVisitNotCollectedEvent) {
	tx := manager.transactionProvider.NewTransaction()
	defer tx.Close()

	stopVisit, found := tx.Model().StopVisits().FindByObjectId(event.StopVisitObjectId)
	if !found {
		logger.Log.Debugf("StopVisitNotCollectedEvent on unknown StopVisit: %#v", event)
		return
	}

	stopVisit.NotCollected()
	stopVisit.ArrivalStatus = STOP_VISIT_ARRIVAL_CANCELLED
	stopVisit.DepartureStatus = STOP_VISIT_DEPARTURE_DEPARTED

	stopVisit.Save()

	logger.Log.Debugf("StopVisit not Collected: %s (%v)", stopVisit.Id(), event.StopVisitObjectId)

	tx.Commit()
}

func NewStopVisitUpdater(tx *Transaction, event *StopVisitUpdateEvent) *StopVisitUpdater {
	return &StopVisitUpdater{tx: tx, event: event}
}

func (updater *StopVisitUpdater) Update() {
	stopArea, ok := updater.tx.Model().StopAreas().FindByObjectId(updater.event.StopAreaObjectId)
	if !ok { // Should never happen
		logger.Log.Debugf("can't find SA: %v", updater.event.StopAreaObjectId)
		return
	}

	// Find the Line
	line := updater.findOrCreateLine()

	// Update StopArea
	stopArea.LineIds.Add(line.Id())
	referent, ok := updater.tx.Model().StopAreas().Find(stopArea.ReferentId)
	if ok {
		referent.LineIds.Add(line.Id())
		referent.Save()
	}
	stopArea.Save()

	// Create or update VJ
	vehicleJourney, ok := updater.tx.Model().VehicleJourneys().FindByObjectId(NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.DatedVehicleJourneyRef))
	if !ok {
		vehicleJourneyAttributes := updater.event.Attributes.VehicleJourneyAttributes()

		vehicleJourney = updater.tx.Model().VehicleJourneys().New()
		vehicleJourney.SetObjectID(vehicleJourneyAttributes.ObjectId)
		vehicleJourney.SetObjectID(NewObjectID("_default", vehicleJourneyAttributes.ObjectId.HashValue()))

		vehicleJourney.LineId = line.Id()
		vehicleJourney.Name = vehicleJourneyAttributes.Attributes["VehicleJourneyName"]

		vehicleJourney.Attributes = vehicleJourneyAttributes.Attributes
		vehicleJourney.References = vehicleJourneyAttributes.References
	} else {
		vehicleJourney.References.SetObjectId("OriginRef", NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.OriginRef))
		vehicleJourney.References.SetObjectId("DestinationRef", NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.DestinationRef))
	}

	vehicleJourney.Monitored = updater.event.Monitored
	vehicleJourney.OriginName = updater.event.OriginName
	vehicleJourney.DestinationName = updater.event.DestinationName
	vehicleJourney.Save()

	// Create or update SV
	stopVisit, ok := updater.tx.Model().StopVisits().FindByObjectId(updater.event.StopVisitObjectid)
	if !ok {
		stopVisitAttributes := updater.event.Attributes.StopVisitAttributes()

		stopVisit = updater.tx.Model().StopVisits().New()
		stopVisit.Origin = updater.event.Origin
		stopVisit.SetObjectID(stopVisitAttributes.ObjectId)
		stopVisit.SetObjectID(NewObjectID("_default", stopVisitAttributes.ObjectId.HashValue()))

		stopVisit.StopAreaId = stopArea.Id()
		stopVisit.VehicleJourneyId = vehicleJourney.Id()
		stopVisit.PassageOrder = stopVisitAttributes.PassageOrder

		stopVisit.Attributes = stopVisitAttributes.Attributes
		stopVisit.References = stopVisitAttributes.References
	}

	stopVisit.Schedules.Merge(&updater.event.Schedules)
	stopVisit.DepartureStatus = updater.event.DepartureStatus
	stopVisit.ArrivalStatus = updater.event.ArrivalStatus
	stopVisit.RecordedAt = updater.event.RecordedAt
	stopVisit.VehicleAtStop = updater.event.VehicleAtStop
	stopVisit.Collected(updater.Clock().Now())

	stopVisit.Save()
}

func (updater *StopVisitUpdater) findOrCreateLine() *Line {
	lineAttributes := updater.event.Attributes.LineAttributes()

	line, ok := updater.tx.Model().Lines().FindByObjectId(lineAttributes.ObjectId)
	if ok {
		return &line
	}

	// logger.Log.Debugf("Create new Line, objectid: %v", lineAttributes.ObjectId)
	line = updater.tx.Model().Lines().New()
	line.SetObjectID(lineAttributes.ObjectId)
	line.SetObjectID(NewObjectID("_default", lineAttributes.ObjectId.HashValue()))
	line.Name = lineAttributes.Name
	line.CollectGeneralMessages = true
	line.SetOrigin(updater.event.Origin)

	line.Save()

	return &line
}
