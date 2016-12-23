package model

import (
	"testing"
	"time"
)

func Test_StopVisitUpdateManager_UpdateStopVisit_found(t *testing.T) {
	model := NewMemoryModel()
	stopVisit := model.StopVisits().New()
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
	model.StopVisits().Save(&stopVisit)
	manager := newStopVisitUpdateManager(model)

	event := &StopVisitUpdateEvent{
		StopVisitObjectid: objectid,
		DepartureStatus:   STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatuts:    STOP_VISIT_ARRIVAL_ONTIME,
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

func Test_StopVisitUpdateManager_UpdateStopVisit(t *testing.T) {
	model := NewMemoryModel()
	manager := newStopVisitUpdateManager(model)
	objectid := NewObjectID("kind", "value")

	event := &StopVisitUpdateEvent{
		StopVisitAttributes: &TestStopVisitUpdateAttributes{},
		StopVisitObjectid:   objectid,
	}

	manager.UpdateStopVisit(event)
	stopVisit, ok := model.StopVisits().FindByObjectId(objectid)
	if !ok {
		t.Errorf("StopVisit should be created by findOrCreateStopArea")
	}
	if len(model.Lines().FindAll()) != 1 {
		t.Errorf("Line should be created by findOrCreateStopArea")
	}
	if len(model.StopAreas().FindAll()) != 1 {
		t.Errorf("StopArea should be created by findOrCreateStopArea")
	}
	if len(model.VehicleJourneys().FindAll()) != 1 {
		t.Errorf("VehicleJourney should be created by findOrCreateStopArea")
	}

	if stopVisit.DepartureStatus() != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be cancelled, got: %v", stopVisit.DepartureStatus())
	}
	if stopVisit.ArrivalStatus() != STOP_VISIT_ARRIVAL_CANCELLED {
		t.Errorf("StopVisit ArrivalStatus should be cancelled, got: %v", stopVisit.ArrivalStatus())
	}
	if stopVisit.PassageOrder() != 1 {
		t.Errorf("StopVisit PassageOrder should be 1, got: %v", stopVisit.PassageOrder())
	}
}

func Test_StopVisitUpdateManager_findOrCreateStopArea_found(t *testing.T) {
	// Setup
	model := NewMemoryModel()
	stopArea := model.StopAreas().New()
	objectid := NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)
	model.StopAreas().Save(&stopArea)
	// Create manager and attributes
	stopVisitUpdateManager := newStopVisitUpdateManager(model)
	stopAreaAttributes := &StopAreaAttributes{ObjectId: &objectid}
	stopVisitUpdateManager.findOrCreateStopArea(stopAreaAttributes)

	if len(model.StopAreas().FindAll()) != 1 {
		t.Errorf("StopArea shouldn't be created by findOrCreateStopArea")
	}
}

func Test_StopVisitUpdateManager_findOrCreateStopArea(t *testing.T) {
	// Setup
	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	// Create manager and attributes
	stopVisitUpdateManager := newStopVisitUpdateManager(model)
	stopVisitUpdateManager.SetClock(NewFakeClock())
	stopAreaAttributes := &StopAreaAttributes{
		Name:     "stopArea",
		ObjectId: &objectid,
	}
	stopVisitUpdateManager.findOrCreateStopArea(stopAreaAttributes)

	stopArea, ok := model.StopAreas().FindByObjectId(objectid)
	if !ok {
		t.Errorf("StopArea should be created by findOrCreateStopArea")
	}
	if stopArea.Name != "stopArea" {
		t.Errorf("Wrong StopArea Name:\n expected: %v\n got: %v", "stopArea", stopArea.Name)
	}
	expected := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	if stopArea.UpdatedAt() != expected {
		t.Errorf("Wrong UpdatedAt:\n expected: %v\n got: %v", expected, stopArea.UpdatedAt())
	}
	if stopArea.RequestedAt() != expected {
		t.Errorf("Wrong RequestedAt:\n expected: %v\n got: %v", expected, stopArea.RequestedAt())
	}
}

func Test_StopVisitUpdateManager_findOrCreateLine_found(t *testing.T) {
	// Setup
	model := NewMemoryModel()
	line := model.Lines().New()
	objectid := NewObjectID("kind", "value")
	line.SetObjectID(objectid)
	model.Lines().Save(&line)
	// Create manager and attributes
	stopVisitUpdateManager := newStopVisitUpdateManager(model)
	lineAttributes := &LineAttributes{ObjectId: &objectid}
	stopVisitUpdateManager.findOrCreateLine(lineAttributes)

	if len(model.Lines().FindAll()) != 1 {
		t.Errorf("Line shouldn't be created by findOrCreateLine")
	}
}

func Test_StopVisitUpdateManager_findOrCreateLine(t *testing.T) {
	// Setup
	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	// Create manager and attributes
	stopVisitUpdateManager := newStopVisitUpdateManager(model)
	LineAttributes := &LineAttributes{
		Name:     "line",
		ObjectId: &objectid,
	}
	stopVisitUpdateManager.findOrCreateLine(LineAttributes)

	line, ok := model.Lines().FindByObjectId(objectid)
	if !ok {
		t.Errorf("Line should be created by findOrCreateStopArea")
	}
	if line.Name != "line" {
		t.Errorf("Wrong Line Name:\n expected: %v\n got: %v", "line", line.Name)
	}
}

func Test_StopVisitUpdateManager_findOrCreateVehicleJourney_found(t *testing.T) {
	// Setup
	model := NewMemoryModel()
	vehicleJourney := model.VehicleJourneys().New()
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)
	model.VehicleJourneys().Save(&vehicleJourney)
	// Create manager and attributes
	stopVisitUpdateManager := newStopVisitUpdateManager(model)
	vehicleJourneyAttributes := &VehicleJourneyAttributes{ObjectId: &objectid}
	stopVisitUpdateManager.findOrCreateVehicleJourney(vehicleJourneyAttributes)

	if len(model.VehicleJourneys().FindAll()) != 1 {
		t.Errorf("VehicleJourney shouldn't be created by findOrCreateVehicleJourney")
	}
}

func Test_StopVisitUpdateManager_findOrCreateVehicleJourney(t *testing.T) {
	// Setup
	model := NewMemoryModel()
	objectid := NewObjectID("kind", "value")
	line := model.Lines().New()
	line.SetObjectID(objectid)
	model.Lines().Save(&line)
	// Create manager and attributes
	stopVisitUpdateManager := newStopVisitUpdateManager(model)
	VehicleJourneyAttributes := &VehicleJourneyAttributes{
		ObjectId:     &objectid,
		LineObjectId: &objectid,
	}
	stopVisitUpdateManager.findOrCreateVehicleJourney(VehicleJourneyAttributes)

	vehicleJourney, ok := model.VehicleJourneys().FindByObjectId(objectid)
	if !ok {
		t.Errorf("VehicleJourney should be created by findOrCreateStopArea")
	}
	if vehicleJourney.lineId != line.Id() {
		t.Errorf("Wrong VehicleJourney LineId:\n expected: %v\n got: %v", line.Id(), vehicleJourney.lineId)
	}
}
