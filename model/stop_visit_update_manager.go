package model

import "github.com/af83/edwig/logger"

type StopAreaUpdateManager struct {
	model Model
}

type StopVisitUpdater struct {
	ClockConsumer

	tx    *Transaction
	event *StopVisitUpdateEvent
}

func NewStopAreaUpdateManager(model Model) func(*StopVisitUpdateEvent) {
	manager := newStopAreaUpdateManager(model)
	return manager.UpdateStopArea
}

func newStopAreaUpdateManager(model Model) *StopAreaUpdateManager {
	return &StopAreaUpdateManager{model: model}
}

func (manager *StopAreaUpdateManager) UpdateStopArea(event *StopVisitUpdateEvent) {
	tx := NewTransaction(manager.model)
	defer tx.Close()

	NewStopVisitUpdater(tx, event).Update()

	tx.Commit()
}

func NewStopVisitUpdater(tx *Transaction, event *StopVisitUpdateEvent) *StopVisitUpdater {
	return &StopVisitUpdater{tx: tx, event: event}
}

func (updater *StopVisitUpdater) Update() {
	existingStopVisit, ok := updater.tx.Model().StopVisits().FindByObjectId(updater.event.StopVisitObjectid)

	if ok {
		logger.Log.Debugf("Update StopVisit %v", existingStopVisit.Id())
		stopArea, _ := updater.tx.Model().StopAreas().FindByObjectId(updater.event.StopAreaObjectId)
		stopArea.Updated(updater.Clock().Now())
		existingStopVisit.Schedules.Merge(updater.event.Schedules)
		existingStopVisit.DepartureStatus = updater.event.DepartureStatus
		existingStopVisit.ArrivalStatus = updater.event.ArrivalStatuts
		existingStopVisit.RecordedAt = updater.event.RecordedAt
		existingStopVisit.VehicleAtStop = updater.event.VehicleAtStop

		updater.tx.Model().StopVisits().Save(&existingStopVisit)
		updater.tx.Model().StopAreas().Save(&stopArea)
		return
	}

	foundStopArea := updater.findOrCreateStopArea(updater.event.Attributes.StopAreaAttributes())

	updater.findOrCreateLine(updater.event.Attributes.LineAttributes())
	foundVehicleJourney := updater.findOrCreateVehicleJourney(updater.event.Attributes.VehicleJourneyAttributes())

	stopVisitAttributes := updater.event.Attributes.StopVisitAttributes()

	logger.Log.Debugf("Create new StopVisit, objectid: %v", stopVisitAttributes.ObjectId)

	stopVisit := updater.tx.Model().StopVisits().New()

	stopVisit.StopAreaId = foundStopArea.Id()
	stopVisit.VehicleJourneyId = foundVehicleJourney.Id()
	stopVisit.PassageOrder = stopVisitAttributes.PassageOrder
	stopVisit.VehicleAtStop = stopVisitAttributes.VehicleAtStop
	stopVisit.RecordedAt = stopVisitAttributes.RecordedAt

	stopVisit.SetObjectID(stopVisitAttributes.ObjectId)
	stopVisit.SetObjectID(NewObjectID("_default", stopVisitAttributes.ObjectId.HashValue()))
	stopVisit.Schedules = stopVisitAttributes.Schedules

	stopVisit.DepartureStatus = stopVisitAttributes.DepartureStatus
	stopVisit.ArrivalStatus = stopVisitAttributes.ArrivalStatus
	stopVisit.Attributes = stopVisitAttributes.Attributes
	stopVisit.References = stopVisitAttributes.References
	updater.tx.Model().StopVisits().Save(&stopVisit)
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
	stopArea.Requested(updater.Clock().Now())
	stopArea.Updated(updater.Clock().Now())
	updater.tx.Model().StopAreas().Save(&stopArea)

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
	updater.tx.Model().Lines().Save(&line)

	return &line
}

func (updater *StopVisitUpdater) findOrCreateVehicleJourney(vehicleJourneyAttributes *VehicleJourneyAttributes) *VehicleJourney {
	vehicleJourney, ok := updater.tx.Model().VehicleJourneys().FindByObjectId(vehicleJourneyAttributes.ObjectId)
	if ok {
		updater.resolveVehiculeJourneyReferences(vehicleJourney)
		return &vehicleJourney
	}

	logger.Log.Debugf("Create new VehicleJourney, objectid: %v", vehicleJourneyAttributes.ObjectId)

	vehicleJourney = updater.tx.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(vehicleJourneyAttributes.ObjectId)
	vehicleJourney.SetObjectID(NewObjectID("_default", vehicleJourneyAttributes.ObjectId.HashValue()))
	foundLine, _ := updater.tx.Model().Lines().FindByObjectId(vehicleJourneyAttributes.LineObjectId)
	vehicleJourney.LineId = foundLine.Id()
	vehicleJourney.Attributes = updater.event.Attributes.VehicleJourneyAttributes().Attributes
	vehicleJourney.References = updater.event.Attributes.VehicleJourneyAttributes().References
	vehicleJourney.Name = vehicleJourney.Attributes["VehicleJourneyName"]
	updater.resolveVehiculeJourneyReferences(vehicleJourney)
	updater.tx.Model().VehicleJourneys().Save(&vehicleJourney)

	return &vehicleJourney
}
