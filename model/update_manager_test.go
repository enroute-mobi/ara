package model

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

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
		Schedules:              NewStopVisitSchedules(),
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
		Schedules:              NewStopVisitSchedules(),
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
		Schedules:              NewStopVisitSchedules(),
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
		Schedules:              NewStopVisitSchedules(),
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

func Test_UpdateManager_UpdateFreshVehicleJourney(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Insert Data in the test db
	databaseVehicleJourney := DatabaseVehicleJourney{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "vehicleJourney",
		ObjectIDs:       `{"internal":"value"}`,
		LineId:          "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Attributes:      "{}",
		References:      `{}`,
	}

	Database.AddTableWithName(databaseVehicleJourney, "vehicle_journeys")
	err := Database.Insert(&databaseVehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	model := NewMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	vehicleJourneys := model.VehicleJourneys().(*MemoryVehicleJourneys)
	err = vehicleJourneys.Load("referential")
	if err != nil {
		t.Fatal(err)
	}

	vehicleJourneyId := VehicleJourneyId(databaseVehicleJourney.Id)
	_, ok := vehicleJourneys.Find(vehicleJourneyId)
	if !ok {
		t.Fatal("Loaded VehicleJourneys should be found")
	}

	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLStopMonitoringResponseFromContent(content)

	manager := newUpdateManager(model)

	objectid := NewObjectID("internal", "value")
	event := &VehicleJourneyUpdateEvent{
		ObjectidKind: "internal",
		ObjectId:     objectid,
		SiriXML:      &response.StopMonitoringDeliveries()[0].XMLMonitoredStopVisits()[0].XMLMonitoredVehicleJourney,
	}

	manager.Update(event)

	updatedVehicleJourney, _ := vehicleJourneys.Find(vehicleJourneyId)
	if updatedVehicleJourney.Attributes.IsEmpty() {
		t.Fatal("Attributes shouldn't be empty after update")
	}

	if updatedVehicleJourney.References.IsEmpty() {
		t.Fatalf("References shouldn't be empty after update: %v", updatedVehicleJourney.References)
	}
}
