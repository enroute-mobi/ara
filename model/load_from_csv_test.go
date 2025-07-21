package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_LoadFromCSVFile(t *testing.T) {
	assert := assert.New(t)
	var vj *VehicleJourney
	var li *Line

	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSVFile("testdata/import.csv", "referential", false)

	// Fetch data from the db
	model := NewTestMemoryModel("referential")
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	model.Load()

	sag, ok := model.StopAreaGroups().Find("cf3e1970-7a7e-4379-ae67-a67abe1c7c1b")
	assert.True(ok, "Can't find StopAreaGroup: \"cf3e1970-7a7e-4379-ae67-a67abe1c7c1b\"")
	assert.Equal("Name", sag.Name)
	assert.Equal("ShortName", sag.ShortName)
	assert.Len(sag.StopAreaIds, 1)
	assert.ElementsMatch(sag.StopAreaIds, []StopAreaId{StopAreaId("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")})

	_, ok = model.StopAreas().Find("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find StopArea: \"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")

	li, ok = model.Lines().Find("f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find Line: \"f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")
	assert.Equal("L1", li.Number)

	lg, ok := model.LineGroups().Find("59208069-3cad-4108-968f-349c5d50a351")
	assert.True(ok, "Can't find LineGroup: \"59208069-3cad-4108-968f-349c5d50a35\"")
	assert.Equal("AgroupName", lg.Name)
	assert.Equal("groupShortName", lg.ShortName)
	assert.ElementsMatch(lg.LineIds, []LineId{LineId("f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")})

	vj, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find VehicleJourney: \"01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")
	assert.Equal("outbound", vj.DirectionType)
	assert.Equal("bus", vj.Attributes["VehicleMode"])

	_, ok = model.ScheduledStopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find StopVisit: \"02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")

	_, ok = model.Operators().Find("03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find Operator: \"03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")

	_, ok = model.Facilities().Find("05cb9be3-e78b-4f76-b644-459e23ba5f1c")
	assert.True(ok, "Can't find Facility: \"05cb9be3-e78b-4f76-b644-459e23ba5f1c\"")

	model = NewTestMemoryModel("referential")
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   2,
	}
	model.Load()

	_, ok = model.StopAreaGroups().Find("cf3e1970-7a7e-4379-ae67-a67abe1c7c1b")
	assert.False(ok, "No StopAreaGroup should exist for this model date")

	_, ok = model.StopAreas().Find("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find StopArea: \"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")

	li, ok = model.Lines().Find("f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find Line: \"f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")
	assert.Zero(li.Number)

	_, ok = model.lineGroups.Find("59208069-3cad-4108-968f-349c5d50a351")
	assert.True(ok, "Can't find LineGroup : \"59208069-3cad-4108-968f-349c5d50a351\"")

	vj, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find VehicleJourney: \"01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")
	assert.Equal("inbound", vj.DirectionType)
	assert.Equal("bus", vj.Attributes["VehicleMode"])

	_, ok = model.ScheduledStopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find StopVisit: \"02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")

	_, ok = model.Operators().Find("03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find Operator: \"03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")

	_, ok = model.Facilities().Find("05cb9be3-e78b-4f76-b644-459e23ba5f1c")
	assert.True(ok, "Can't find Facility: \"05cb9be3-e78b-4f76-b644-459e23ba5f1c\"")

	model = NewTestMemoryModel("referential")
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   3,
	}
	model.Load()

	vj, ok = model.VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok, "Can't find VehicleJourney: \"01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\"")
	assert.Zero(vj.DirectionType)
}

func Test_LoadFromCSVFile_Force(t *testing.T) {
	assert := assert.New(t)

	InitTestDb(t)
	defer CleanTestDb(t)

	// Fill DB
	LoadFromCSVFile("testdata/import.csv", "referential", false)

	forceBuilder := NewLoader("referential", true, true)

	sa := []string{"stop_area", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "2017-01-01", "", "", "", "", "", "", "", ""}
	sag := []string{"stop_area_group", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", ""}
	o := []string{"operator", "03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", ""}
	l := []string{"line", "f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", "", "", ""}
	lg := []string{"line_group", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", ""}
	vj := []string{"vehicle_journey", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "", "", ""}
	sv := []string{"stop_visit", "02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "2017-01-01", "", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "", "", "", ""}

	err := forceBuilder.handleStopArea(sa)
	assert.NoError(err, "Import StopArea with force should not return an error")

	err = forceBuilder.handleStopAreaGroup(sag)
	assert.NoError(err, "Import StopAreaGroup with force should not return an error")

	err = forceBuilder.handleOperator(o)
	assert.NoError(err, "Import Operator with force should not return an error")

	err = forceBuilder.handleLine(l)
	assert.NoError(err, "Import Line with force should not return an error")

	err = forceBuilder.handleLineGroup(lg)
	assert.NoError(err, "Import LineGroup with force should not return an error")

	err = forceBuilder.handleVehicleJourney(vj)
	assert.NoError(err, "Import VehicleJourney with force should not return an error")

	err = forceBuilder.handleStopVisit(sv)
	assert.NoError(err, "Import StopVisit with force should not return an error")
}
