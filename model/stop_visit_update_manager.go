package model

import "github.com/af83/edwig/logger"

type StopVisitUpdateManager struct {
	ClockConsumer

	model Model
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

	stopVisit, ok := tx.Model().StopVisits().FindByObjectId(event.StopVisitObjectid)
	if ok {
		logger.Log.Debugf("Update StopVisit %v", stopVisit.Id())
		stopArea, _ := tx.Model().StopAreas().FindByObjectId(event.StopAreaObjectId)
		stopArea.Updated(manager.Clock().Now())
		stopVisit.schedules = event.Schedules
		stopVisit.departureStatus = event.DepartureStatus
		stopVisit.arrivalStatus = event.ArrivalStatuts

		tx.Model().StopVisits().Save(&stopVisit)
		tx.Commit()
		return
	}

	manager.findOrCreateStopArea(event.StopVisitAttributes.StopAreaAttributes())
	manager.findOrCreateLine(event.StopVisitAttributes.LineAttributes())
	manager.findOrCreateVehicleJourney(event.StopVisitAttributes.VehicleJourneyAttributes())

	stopVisitAttributes := event.StopVisitAttributes.StopVisitAttributes()

	logger.Log.Debugf("Create new StopVisit, objectid: %v", stopVisitAttributes.ObjectId)

	stopVisit = tx.Model().StopVisits().New()
	foundStopArea, _ := tx.Model().StopAreas().FindByObjectId(*stopVisitAttributes.StopAreaObjectId)
	stopVisit.stopAreaId = foundStopArea.Id()
	foundVehicleJourney, _ := tx.Model().VehicleJourneys().FindByObjectId(*stopVisitAttributes.VehicleJourneyObjectId)
	stopVisit.vehicleJourneyId = foundVehicleJourney.Id()
	stopVisit.passageOrder = stopVisitAttributes.PassageOrder
	stopVisit.SetObjectID(*stopVisitAttributes.ObjectId)
	stopVisit.schedules = stopVisitAttributes.Schedules
	stopVisit.departureStatus = stopVisitAttributes.DepartureStatus
	stopVisit.arrivalStatus = stopVisitAttributes.ArrivalStatus

	tx.Model().StopVisits().Save(&stopVisit)
	tx.Commit()
}

func (manager *StopVisitUpdateManager) findOrCreateStopArea(stopAreaAttributes *StopAreaAttributes) {
	stopArea, ok := manager.model.StopAreas().FindByObjectId(*stopAreaAttributes.ObjectId)
	if ok {
		return
	}
	tx := NewTransaction(manager.model)
	defer tx.Close()

	logger.Log.Debugf("Create new StopArea %v, objectid: %v", stopAreaAttributes.Name, *stopAreaAttributes.ObjectId)

	stopArea = tx.Model().StopAreas().New()
	stopArea.SetObjectID(*stopAreaAttributes.ObjectId)
	stopArea.Name = stopAreaAttributes.Name
	stopArea.Requested(manager.Clock().Now())
	stopArea.Updated(manager.Clock().Now())
	tx.Model().StopAreas().Save(&stopArea)
	tx.Commit()
}

func (manager *StopVisitUpdateManager) findOrCreateLine(lineAttributes *LineAttributes) {
	line, ok := manager.model.Lines().FindByObjectId(*lineAttributes.ObjectId)
	if ok {
		return
	}
	tx := NewTransaction(manager.model)
	defer tx.Close()

	logger.Log.Debugf("Create new Line, objectid: %v", *lineAttributes.ObjectId)

	line = tx.Model().Lines().New()
	line.SetObjectID(*lineAttributes.ObjectId)
	line.Name = lineAttributes.Name
	tx.Model().Lines().Save(&line)
	tx.Commit()
}

func (manager *StopVisitUpdateManager) findOrCreateVehicleJourney(vehicleJourneyAttributes *VehicleJourneyAttributes) {
	vehicleJourney, ok := manager.model.VehicleJourneys().FindByObjectId(*vehicleJourneyAttributes.ObjectId)
	if ok {
		return
	}
	tx := NewTransaction(manager.model)
	defer tx.Close()

	logger.Log.Debugf("Create new VehicleJourney, objectid: %v", *vehicleJourneyAttributes.ObjectId)

	vehicleJourney = tx.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(*vehicleJourneyAttributes.ObjectId)
	foundLine, _ := tx.Model().Lines().FindByObjectId(*vehicleJourneyAttributes.LineObjectId)
	vehicleJourney.lineId = foundLine.Id()
	tx.Model().VehicleJourneys().Save(&vehicleJourney)
	tx.Commit()
}
