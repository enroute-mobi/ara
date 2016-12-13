package model

func NewStopVisitUpdateManager(model Model) func(*StopVisitUpdateEvent) {
	return func(event *StopVisitUpdateEvent) {
		tx := NewTransaction(model)
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
}
