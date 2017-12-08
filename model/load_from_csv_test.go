package model

import (
	"testing"
	"time"
)

func Test_LoadFromCSV(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSV("testdata/import.csv", "referential")

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
}
