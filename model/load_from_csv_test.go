package model

import (
	"testing"
	"time"
)

func Test_LoadFromCSV(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSV("testdata/import.csv", "referential", false)

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
	_, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find VehicleJourney: %v", model.VehicleJourneys().FindAll())
	}
	_, ok = model.StopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find StopVisit: %v", model.StopVisits().FindAll())
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
	_, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find VehicleJourney: %v", model.VehicleJourneys().FindAll())
	}
	_, ok = model.StopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find StopVisit: %v", model.StopVisits().FindAll())
	}
	_, ok = model.Operators().Find("03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if !ok {
		t.Errorf("Can't find Operator: %v", model.Operators().FindAll())
	}
}

func Test_LoadFromCSV_Force(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSV("testdata/import.csv", "referential", false)

	forceBuilder := newLoader("", "referential", true)

	sa := []string{"stop_area", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "2017-01-01", "", "", "", "", "", "", "", ""}
	o := []string{"operator", "03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", ""}
	l := []string{"line", "f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", "", "", ""}
	vj := []string{"vehicle_journey", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "", ""}
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
