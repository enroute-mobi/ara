package model

import (
	"testing"
	"time"
)

func Test_LoadFromCSVFile(t *testing.T) {
	var vj *VehicleJourney
	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSVFile("testdata/import.csv", "referential", false)

	// Fetch data from the db
	model := NewMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	model.Load("referential")

	_, ok := model.StopAreas().Find("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find StopArea: %v", model.StopAreas().FindAll())
	}
	_, ok = model.Lines().Find("f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find Line: %v", model.Lines().FindAll())
	}
	vj, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find VehicleJourney: %v", model.VehicleJourneys().FindAll())
	}
	if vj.DirectionType != "outbound" {
		t.Errorf("Wrong direction_type for VehicleJourney: expected \"outbound\", got: %v", vj.DirectionType)
	}
	if vj.Attributes["VehicleMode"] != "bus" {
		t.Errorf("Wrong Attributes for VehicleJourney: expected \"bus\", got: %v", vj.Attributes["VehicleMode"])
	}
	_, ok = model.ScheduledStopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find StopVisit: %v", model.ScheduledStopVisits().FindAll())
	}
	_, ok = model.Operators().Find("03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find Operator: %v", model.Operators().FindAll())
	}

	model = NewMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   2,
	}
	model.Load("referential")

	_, ok = model.StopAreas().Find("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find StopArea: %v", model.StopAreas().FindAll())
	}
	_, ok = model.Lines().Find("f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find Line: %v", model.Lines().FindAll())
	}
	vj, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find VehicleJourney: %v", model.VehicleJourneys().FindAll())
	}
	if vj.DirectionType != "inbound" {
		t.Errorf("Wrong direction_type for VehicleJourney: expected \"inbound\", got: %v", vj.DirectionType)
	}
	_, ok = model.ScheduledStopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find StopVisit: %v", model.ScheduledStopVisits().FindAll())
	}
	_, ok = model.Operators().Find("03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find Operator: %v", model.Operators().FindAll())
	}

	model = NewMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   3,
	}
	model.Load("referential")

	vj, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find VehicleJourney: %v", model.VehicleJourneys().FindAll())
	}
	if vj.DirectionType != "" {
		t.Errorf("Wrong direction_type for VehicleJourney: expected \"\", got: %v", vj.DirectionType)
	}
}

func Test_LoadFromCSVFile_Force(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSVFile("testdata/import.csv", "referential", false)

	forceBuilder := NewLoader("referential", true, true)

	sa := []string{"stop_area", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "2017-01-01", "", "", "", "", "", "", "", ""}
	o := []string{"operator", "03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", ""}
	l := []string{"line", "f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", "", "", ""}
	vj := []string{"vehicle_journey", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "", "", ""}
	sv := []string{"stop_visit", "02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "", ""}

	err := forceBuilder.handleStopArea(sa)
	if err != nil {
		t.Errorf("Import StopArea with force should not return an error, got: %v", err)
	}
	err = forceBuilder.handleOperator(o)
	if err != nil {
		t.Errorf("Import Operator with force should not return an error, got: %v", err)
	}
	err = forceBuilder.handleLine(l)
	if err != nil {
		t.Errorf("Import Line with force should not return an error, got: %v", err)
	}
	err = forceBuilder.handleVehicleJourney(vj)
	if err != nil {
		t.Errorf("Import VehicleJourney with force should not return an error, got: %v", err)
	}
	err = forceBuilder.handleStopVisit(sv)
	if err != nil {
		t.Errorf("Import StopVisit with force should not return an error, got: %v", err)
	}
}
