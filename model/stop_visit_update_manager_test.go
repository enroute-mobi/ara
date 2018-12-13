package model

import "testing"

func Test_StopVisitUpdateManager_UpdateStopVisit_found(t *testing.T) {
	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	sa := model.StopAreas().New()
	sa.SetObjectID(objectid)
	sa.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetObjectID(objectid)
	model.StopVisits().Save(&stopVisit)

	manager := newStopAreaUpdateManager(model)

	event := &StopVisitUpdateEvent{
		StopAreaObjectId:  objectid,
		StopVisitObjectid: objectid,
		DepartureStatus:   STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:     STOP_VISIT_ARRIVAL_ONTIME,
		Attributes:        &TestStopVisitUpdateAttributes{},
	}

	manager.UpdateStopVisit(event)
	updatedStopVisit, _ := model.StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_ONTIME {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
	if !updatedStopVisit.IsCollected() {
		t.Errorf("StopVisit ArrivalStatus should be collected")
	}
}

func Test_StopVisitUpdateManager_UpdateStopVisit(t *testing.T) {
	model := NewMemoryModel()
	manager := newStopAreaUpdateManager(model)
	objectid := NewObjectID("kind", "value")

	sa := model.StopAreas().New()
	sa.SetObjectID(objectid)
	sa.Save()

	event := &StopVisitUpdateEvent{
		Attributes:        &TestStopVisitUpdateAttributes{},
		StopAreaObjectId:  objectid,
		StopVisitObjectid: objectid,
		DepartureStatus:   STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:     STOP_VISIT_ARRIVAL_CANCELLED,
	}

	manager.UpdateStopVisit(event)
	stopVisit, ok := model.StopVisits().FindByObjectId(objectid)
	if !ok {
		t.Fatalf("StopVisit should be created by findOrCreateStopVisit")
	}
	if len(model.Lines().FindAll()) != 1 {
		t.Fatalf("Line should be created by findOrCreateLine")
	}
	if len(model.StopAreas().FindAll()) != 1 {
		t.Fatalf("StopArea should be created by findOrCreateStopArea")
	}
	if len(model.VehicleJourneys().FindAll()) != 1 {
		t.Fatalf("VehicleJourney should be created by findOrCreateVehicleJourney")
	}

	if stopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be cancelled, got: %v", stopVisit.DepartureStatus)
	}
	if stopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_CANCELLED {
		t.Errorf("StopVisit ArrivalStatus should be cancelled, got: %v", stopVisit.ArrivalStatus)
	}
	if stopVisit.PassageOrder != 1 {
		t.Errorf("StopVisit PassageOrder should be 1, got: %v", stopVisit.PassageOrder)
	}

	// Check if the SA have the Line in LineIds
	lineId := model.Lines().FindAll()[0].Id()
	lineIds := model.StopAreas().FindAll()[0].LineIds
	if len(lineIds) != 1 && lineIds[0] != lineId {
		t.Errorf("StopArea should have %v in LineIds, got: %v", lineId, lineIds)
	}
}

func Test_StopVisitUpdateManager_UpdateMonitoredStopArea(t *testing.T) {
	model := NewMemoryModel()
	manager := newStopAreaUpdateManager(model)

	sa := model.StopAreas().New()
	sa.Name = "Parent"
	sa.Save()

	sa2 := model.StopAreas().New()
	sa2.Name = "Son"
	sa2.ParentId = sa.id
	sa2.Save()

	sa3 := model.StopAreas().New()
	sa3.Name = "Grandson"
	sa3.ParentId = sa2.id
	sa3.Save()

	event := NewStopAreaMonitoredEvent(manager.NewUUID(), sa3.Id(), "test_origin", true)
	manager.UpdateMonitoredStopArea(event)

	stopArea, _ := model.StopAreas().Find(sa.Id())
	if status, ok := stopArea.Origins.Origin("test_origin"); !ok || !status {
		t.Errorf("Parent StopArea status should have been updated, got found origin: %v and status: %v", ok, status)
	}

	stopArea2, _ := model.StopAreas().Find(sa2.Id())
	if status, ok := stopArea2.Origins.Origin("test_origin"); !ok || !status {
		t.Errorf("StopArea status should have been updated, got found origin: %v and status: %v", ok, status)
	}

	stopArea3, _ := model.StopAreas().Find(sa3.Id())
	if status, ok := stopArea3.Origins.Origin("test_origin"); !ok || !status {
		t.Errorf("StopArea status should have been updated, got found origin: %v and status: %v", ok, status)
	}
}
