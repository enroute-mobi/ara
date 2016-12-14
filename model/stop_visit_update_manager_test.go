package model

import "testing"

func Test_StopVisitUpdateManager_UpdateStopVisit(t *testing.T) {
	model := NewMemoryModel()
	stopVisit := model.StopVisits().New()
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
	model.StopVisits().Save(&stopVisit)
	manager := newStopVisitUpdateManager(model)

	event := &StopVisitUpdateEvent{
		Stop_visit_objectid: objectid,
		DepartureStatus:     STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatuts:      STOP_VISIT_ARRIVAL_ONTIME,
	}

	manager.UpdateStopVisit(event)
	updatedStopVisit, _ := model.StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.DepartureStatus() != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus() != STOP_VISIT_ARRIVAL_ONTIME {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
}
