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
	if !ok {
		return
	}
	stopVisit.schedules = event.Schedules
	stopVisit.departureStatus = event.DepartureStatus
	stopVisit.arrivalStatus = event.ArrivalStatuts

	tx.Model().StopVisits().Save(&stopVisit)
	tx.Commit()
}
