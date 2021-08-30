package model

import "testing"

func Test_UpdateManager_CreateStopVisit(t *testing.T) {
	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	sa := model.StopAreas().New()
	sa.SetObjectID(objectid)
	sa.Save()

	l := model.Lines().New()
	l.SetObjectID(objectid)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetObjectID(objectid)
	vj.LineId = l.Id()
	vj.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		ObjectId:               objectid,
		StopAreaObjectId:       objectid,
		VehicleJourneyObjectId: objectid,
		DepartureStatus:        STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:          STOP_VISIT_ARRIVAL_ONTIME,
	}

	manager.Update(event)
	updatedStopVisit, ok := model.StopVisits().FindByObjectId(objectid)
	if !ok {
		t.Fatalf("StopVisit should be created")
	}
	if updatedStopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_ONTIME {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
	if !updatedStopVisit.IsCollected() {
		t.Errorf("StopVisit ArrivalStatus should be collected")
	}
	updatedStopArea, _ := model.StopAreas().Find(sa.Id())
	if !updatedStopArea.LineIds.Contains(l.Id()) {
		t.Errorf("StopArea LineIds should be updated")
	}
}

func Test_UpdateManager_UpdateStopVisit(t *testing.T) {
	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	sa := model.StopAreas().New()
	sa.SetObjectID(objectid)
	sa.Save()

	l := model.Lines().New()
	l.SetObjectID(objectid)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetObjectID(objectid)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetObjectID(objectid)
	stopVisit.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		ObjectId:               objectid,
		StopAreaObjectId:       objectid,
		VehicleJourneyObjectId: objectid,
		DepartureStatus:        STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:          STOP_VISIT_ARRIVAL_ONTIME,
	}

	manager.Update(event)
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
	updatedStopArea, _ := model.StopAreas().Find(sa.Id())
	if !updatedStopArea.LineIds.Contains(l.Id()) {
		t.Errorf("StopArea LineIds should be updated")
	}
}

func Test_UpdateManager_CreateStopVisit_NoStopAreaId(t *testing.T) {
	emptyObjectid := NewObjectID("kind", "")

	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	sa := model.StopAreas().New()
	sa.SetObjectID(objectid)
	sa.Save()

	l := model.Lines().New()
	l.SetObjectID(objectid)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetObjectID(objectid)
	vj.LineId = l.Id()
	vj.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		ObjectId:               objectid,
		StopAreaObjectId:       emptyObjectid,
		VehicleJourneyObjectId: objectid,
		DepartureStatus:        STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:          STOP_VISIT_ARRIVAL_ONTIME,
	}

	manager.Update(event)
	_, ok := model.StopVisits().FindByObjectId(objectid)
	if ok {
		t.Fatalf("StopVisit should not be created")
	}
}

func Test_UpdateManager_UpdateStopVisit_NoStopAreaId(t *testing.T) {
	emptyObjectid := NewObjectID("kind", "")

	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	sa := model.StopAreas().New()
	sa.SetObjectID(objectid)
	sa.Save()

	l := model.Lines().New()
	l.SetObjectID(objectid)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetObjectID(objectid)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetObjectID(objectid)
	stopVisit.StopAreaId = sa.Id()
	stopVisit.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		ObjectId:               objectid,
		StopAreaObjectId:       emptyObjectid,
		VehicleJourneyObjectId: objectid,
		DepartureStatus:        STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:          STOP_VISIT_ARRIVAL_ONTIME,
	}

	manager.Update(event)
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
	updatedStopArea, _ := model.StopAreas().Find(sa.Id())
	if !updatedStopArea.LineIds.Contains(l.Id()) {
		t.Errorf("StopArea LineIds should be updated")
	}
}

func Test_UpdateManager_UpdateStatus(t *testing.T) {
	model := NewMemoryModel()
	manager := newUpdateManager(model)

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

	event := NewStatusUpdateEvent(sa3.Id(), "test_origin", true)
	manager.Update(event)

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

func Test_UpdateManager_UpdateNotCollected(t *testing.T) {
	model := NewMemoryModel()
	manager := newUpdateManager(model)

	objectid := NewObjectID("kind", "value")
	stopVisit := model.StopVisits().New()
	stopVisit.SetObjectID(objectid)
	stopVisit.collected = true
	stopVisit.Save()

	manager.Update(NewNotCollectedUpdateEvent(objectid))
	updatedStopVisit, _ := model.StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_DEPARTED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_CANCELLED {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
	if updatedStopVisit.collected {
		t.Errorf("StopVisit Collected should be updated")
	}
}
