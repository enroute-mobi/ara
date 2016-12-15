package model

type StopVisitUpdateManager struct {
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

	stopVisit, ok := tx.Model().StopVisits().FindByObjectId(event.Stop_visit_objectid)
	if ok {
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

	stopArea = tx.Model().StopAreas().New()
	stopArea.SetObjectID(*stopAreaAttributes.ObjectId)
	stopArea.Name = stopAreaAttributes.Name
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

	line = tx.Model().Lines().New()
	line.SetObjectID(*lineAttributes.ObjectId)
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

	vehicleJourney = tx.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(*vehicleJourneyAttributes.ObjectId)
	foundLine, _ := tx.Model().Lines().FindByObjectId(*vehicleJourneyAttributes.LineObjectId)
	vehicleJourney.lineId = foundLine.Id()
	tx.Model().VehicleJourneys().Save(&vehicleJourney)
	tx.Commit()
}
