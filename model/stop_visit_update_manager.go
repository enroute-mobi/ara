package model

import (
	"time"

	"github.com/af83/edwig/logger"
)

type StopAreaUpdateManager struct {
	ClockConsumer

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

	stopArea, found := tx.Model().StopAreas().Find(event.StopAreaId)
	if !found {
		logger.Log.Debugf("StopAreaUpdateEvent for unknown StopArea %v", event.StopAreaId)
		return
	}

	logger.Log.Debugf("Update StopArea %v", stopArea.Id())
	stopArea.Updated(manager.Clock().Now())
	stopArea.Save()

	tx.Commit()

	for _, stopVisitUpdateEvent := range event.StopVisitUpdateEvents {
		manager.UpdateStopVisit(stopVisitUpdateEvent)
	}
	for _, stopVisitNotCollectedEvent := range event.StopVisitNotCollectedEvents {
		manager.UpdateNotCollectedStopVisit(stopVisitNotCollectedEvent)
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
	existingStopVisit, ok := updater.tx.Model().StopVisits().FindByObjectId(updater.event.StopVisitObjectid)

	if ok {
		// too verbose
		// logger.Log.Debugf("Update StopVisit %v", existingStopVisit.Id())
		existingStopVisit.Schedules.Merge(updater.event.Schedules)
		existingStopVisit.DepartureStatus = updater.event.DepartureStatus
		existingStopVisit.ArrivalStatus = updater.event.ArrivalStatuts
		existingStopVisit.RecordedAt = updater.event.RecordedAt
		existingStopVisit.VehicleAtStop = updater.event.VehicleAtStop
		existingStopVisit.Collected(updater.Clock().Now())

		existingStopVisit.Save()

		stopArea, found := updater.tx.Model().StopAreas().FindByObjectId(updater.event.StopAreaObjectId)
		if found {
			stopArea.Updated(updater.Clock().Now())
			stopArea.Save()
		} else {
			logger.Log.Debugf("StopVisitUpdateEvent associated to unknown StopArea: %v", updater.event.StopAreaObjectId)
		}

		foundVehicleJourney, ok := updater.tx.Model().VehicleJourneys().FindByObjectId(NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.DatedVehicleJourneyRef))
		if ok {
			foundVehicleJourney.References.SetObjectId("DestinationRef", NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.DestinationRef), "")
			foundVehicleJourney.References.SetObjectId("DestinationName", NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.DestinationName), "")
			foundVehicleJourney.References.SetObjectId("OriginRef", NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.OriginRef), "")
			foundVehicleJourney.References.SetObjectId("OriginName", NewObjectID(updater.event.StopVisitObjectid.Kind(), updater.event.OriginName), "")
		} else {
			foundVehicleJourney = *updater.CreateVehicleJourney(updater.event.Attributes.VehicleJourneyAttributes())
		}
		updater.resolveVehiculeJourneyReferences(foundVehicleJourney)

		foundVehicleJourney.Save()

		return
	}

	foundStopArea := updater.findOrCreateStopArea(updater.event.Attributes.StopAreaAttributes())

	updater.findOrCreateLine(updater.event.Attributes.LineAttributes())

	foundVehicleJourney := updater.findOrCreateVehicleJourney(updater.event.Attributes.VehicleJourneyAttributes())

	stopVisitAttributes := updater.event.Attributes.StopVisitAttributes()

	stopVisit := updater.tx.Model().StopVisits().New()

	logger.Log.Debugf("Create new StopVisit, objectid: %v", stopVisit.Id())

	stopVisit.StopAreaId = foundStopArea.Id()
	stopVisit.VehicleJourneyId = foundVehicleJourney.Id()
	stopVisit.PassageOrder = stopVisitAttributes.PassageOrder
	stopVisit.VehicleAtStop = stopVisitAttributes.VehicleAtStop
	stopVisit.RecordedAt = stopVisitAttributes.RecordedAt
	stopVisit.Collected(updater.Clock().Now())

	stopVisit.SetObjectID(stopVisitAttributes.ObjectId)
	stopVisit.SetObjectID(NewObjectID("_default", stopVisitAttributes.ObjectId.HashValue()))
	stopVisit.Schedules = stopVisitAttributes.Schedules

	stopVisit.DepartureStatus = stopVisitAttributes.DepartureStatus
	stopVisit.ArrivalStatus = stopVisitAttributes.ArrivalStatus
	stopVisit.Attributes = stopVisitAttributes.Attributes
	stopVisit.References = stopVisitAttributes.References

	stopVisit.Save()
}

func (updater *StopVisitUpdater) findOrCreateStopArea(stopAreaAttributes *StopAreaAttributes) *StopArea {
	stopArea, ok := updater.tx.Model().StopAreas().FindByObjectId(stopAreaAttributes.ObjectId)
	if ok {
		return &stopArea
	}

	logger.Log.Debugf("Create new StopArea %v, objectid: %v", stopAreaAttributes.Name, stopAreaAttributes.ObjectId)
	stopArea = updater.tx.Model().StopAreas().New()
	stopArea.SetObjectID(stopAreaAttributes.ObjectId)
	stopArea.Name = stopAreaAttributes.Name
	stopArea.Updated(updater.Clock().Now())
	stopArea.CollectedAlways = false
	stopArea.NextCollectAt = updater.Clock().Now().Add(1 * time.Minute)
	stopArea.Save()
	return &stopArea
}

func (updater *StopVisitUpdater) resolveVehiculeJourneyReferences(foundVehicleJourney VehicleJourney) error {
	toResolve := []string{"PlaceRef", "OriginRef", "DestinationRef"}

	for _, ref := range toResolve {
		if foundVehicleJourney.References[ref] != (Reference{}) {
			foundStopArea, ok := updater.tx.Model().StopAreas().FindByObjectId(*(foundVehicleJourney.References[ref].ObjectId))
			if ok {
				reference := foundVehicleJourney.References[ref]
				reference.Id = string(foundStopArea.Id())
				foundVehicleJourney.References[ref] = reference
			}
		}
	}
	return nil
}

func (updater *StopVisitUpdater) findOrCreateLine(lineAttributes *LineAttributes) *Line {
	line, ok := updater.tx.Model().Lines().FindByObjectId(lineAttributes.ObjectId)
	if ok {
		return &line
	}

	logger.Log.Debugf("Create new Line, objectid: %v", lineAttributes.ObjectId)

	line = updater.tx.Model().Lines().New()
	line.SetObjectID(lineAttributes.ObjectId)
	line.Name = lineAttributes.Name

	line.Save()

	return &line
}

func (updater *StopVisitUpdater) CreateVehicleJourney(vehicleJourneyAttributes *VehicleJourneyAttributes) *VehicleJourney {
	logger.Log.Debugf("Create new VehicleJourney, objectid: %v", vehicleJourneyAttributes.ObjectId)

	vehicleJourney := updater.tx.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(vehicleJourneyAttributes.ObjectId)
	vehicleJourney.SetObjectID(NewObjectID("_default", vehicleJourneyAttributes.ObjectId.HashValue()))
	foundLine, _ := updater.tx.Model().Lines().FindByObjectId(vehicleJourneyAttributes.LineObjectId)
	vehicleJourney.LineId = foundLine.Id()

	vehicleJourney.Attributes = vehicleJourneyAttributes.Attributes
	vehicleJourney.References = vehicleJourneyAttributes.References
	vehicleJourney.Name = vehicleJourney.Attributes["VehicleJourneyName"]

	updater.resolveVehiculeJourneyReferences(vehicleJourney)

	vehicleJourney.Save()

	return &vehicleJourney
}

func (updater *StopVisitUpdater) findOrCreateVehicleJourney(vehicleJourneyAttributes *VehicleJourneyAttributes) *VehicleJourney {
	vehicleJourney, ok := updater.tx.Model().VehicleJourneys().FindByObjectId(vehicleJourneyAttributes.ObjectId)
	if ok {
		updater.resolveVehiculeJourneyReferences(vehicleJourney)
		return &vehicleJourney
	}

	return updater.CreateVehicleJourney(vehicleJourneyAttributes)
}
