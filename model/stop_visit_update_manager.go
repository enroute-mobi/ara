package model

import "github.com/af83/edwig/logger"

type StopVisitUpdateManager struct {
	model Model
}

type StopVisitUpdater struct {
	ClockConsumer

	tx    *Transaction
	event *StopVisitUpdateEvent
}

func NewStopVisitUpdateManager(model Model) func(*StopVisitUpdateEvent) {
	manager := newStopVisitUpdateManager(model)
	return manager.UpdateStopVisit
}

func newStopVisitUpdateManager(model Model) *StopVisitUpdateManager {
	return &StopVisitUpdateManager{model: model}
}

func (manager *StopVisitUpdateManager) UpdateStopVisit(event *StopVisitUpdateEvent) {
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
		existingStopVisit.schedules = updater.event.Schedules
		existingStopVisit.departureStatus = updater.event.DepartureStatus
		existingStopVisit.arrivalStatus = updater.event.ArrivalStatuts
		existingStopVisit.recordedAt = updater.event.RecordedAt

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
	stopVisit.stopAreaId = foundStopArea.Id()
	stopVisit.vehicleJourneyId = foundVehicleJourney.Id()
	stopVisit.passageOrder = stopVisitAttributes.PassageOrder
	stopVisit.recordedAt = stopVisitAttributes.RecordedAt
	stopVisit.SetObjectID(stopVisitAttributes.ObjectId)
	stopVisit.schedules = stopVisitAttributes.Schedules
	stopVisit.departureStatus = stopVisitAttributes.DepartureStatus
	stopVisit.arrivalStatus = stopVisitAttributes.ArrivalStatus

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
		return &vehicleJourney
	}

	logger.Log.Debugf("Create new VehicleJourney, objectid: %v", vehicleJourneyAttributes.ObjectId)

	vehicleJourney = updater.tx.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(vehicleJourneyAttributes.ObjectId)
	foundLine, _ := updater.tx.Model().Lines().FindByObjectId(vehicleJourneyAttributes.LineObjectId)
	vehicleJourney.lineId = foundLine.Id()
	updater.tx.Model().VehicleJourneys().Save(&vehicleJourney)

	return &vehicleJourney
}
