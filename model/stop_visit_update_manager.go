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
	tx := manager.transactionProvider.NewTransaction()
	defer tx.Close()

	if event.StopAreaMonitoredEvent != nil {
		manager.UpdateMonitoredStopArea(event, tx)
		tx.Commit()
		return
	}

	stopArea, found := tx.Model().StopAreas().Find(event.StopAreaId)
	if !found {
		logger.Log.Debugf("StopAreaUpdateEvent for unknown StopArea %v", event.StopAreaId)

		stopArea = tx.Model().StopAreas().New()
		parentSA, _ := tx.Model().StopAreas().FindByObjectId(event.StopAreaAttributes.ParentObjectId)

		stopArea.SetObjectID(event.StopAreaAttributes.ObjectId)
		stopArea.ParentId = parentSA.Id()
		stopArea.Name = event.StopAreaAttributes.Name
		stopArea.CollectedAlways = event.StopAreaAttributes.CollectedAlways
		stopArea.CollectGeneralMessages = true
		stopArea.Monitored = true
		stopArea.Save()

		event.StopAreaId = stopArea.Id()
	}

	logger.Log.Debugf("Update StopArea %v", stopArea.Id())
	stopArea.Updated(manager.Clock().Now())
	stopArea.Save()
	if event.Origin != "" {
		status, ok := stopArea.Origins.Origin(event.Origin)
		if !status || !ok {
			manager.UpdateMonitoredStopArea(NewStopAreaMonitoredEvent(manager.NewUUID(), event.StopAreaId, event.Origin, true), tx)
		}
	}

	tx.Commit()

	for _, stopVisitUpdateEvent := range event.StopVisitUpdateEvents {
		manager.UpdateStopVisit(stopVisitUpdateEvent)
	}
	for _, stopVisitNotCollectedEvent := range event.StopVisitNotCollectedEvents {
		manager.UpdateNotCollectedStopVisit(stopVisitNotCollectedEvent)
	}
}

func (manager *StopAreaUpdateManager) UpdateMonitoredStopArea(event *StopAreaUpdateEvent, tx *Transaction) {
	// Should never happen, but don't want to ever have a go nil pointer exception
	if event.StopAreaMonitoredEvent == nil {
		return
	}

	for _, stopArea := range tx.Model().StopAreas().FindAscendants(event.StopAreaId) {
		stopArea.Origins.SetPartnerStatus(event.StopAreaMonitoredEvent.Partner, event.StopAreaMonitoredEvent.Status)
		stopArea.Monitored = stopArea.Origins.Monitored()
		stopArea.Save()
	}
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

	logger.Log.Printf("StopVisit not Collected: %s (%v)", stopVisit.Id(), event.StopVisitObjectId)

	tx.Commit()
}

func NewStopVisitUpdater(tx *Transaction, event *StopVisitUpdateEvent) *StopVisitUpdater {
	return &StopVisitUpdater{tx: tx, event: event}
}

func (updater *StopVisitUpdater) Update() {
	stopArea, ok := updater.tx.Model().StopAreas().FindByObjectId(updater.event.StopAreaObjectId)
	if !ok { // Should never happen
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

	line.Save()

	return &line
}
