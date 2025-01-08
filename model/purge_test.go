package model

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
)

func fill_purifier_test_db(t *testing.T) {
	databaseStopArea := DatabaseStopArea{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "stopArea",
		Codes:           "{}",
		LineIds:         "[]",
		Attributes:      "{}",
		References:      "{}",
	}

	Database.AddTableWithName(databaseStopArea, "stop_areas")
	err := Database.Insert(&databaseStopArea)
	if err != nil {
		t.Fatal(err)
	}

	databaseStopArea.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseStopArea.ModelName = "2017-01-02"
	err = Database.Insert(&databaseStopArea)
	if err != nil {
		t.Fatal(err)
	}

	databaseStopAreaGroup := DatabaseStopAreaGroup{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "stopAreaGroup",
		ShortName:       "stop_area_group_short_name",
		StopAreaIds:     `["d0eebc99-9c0b","e0eebc99-9c0b"]`,
	}

	Database.AddTableWithName(databaseStopAreaGroup, "stop_area_groups")
	err = Database.Insert(&databaseStopAreaGroup)
	if err != nil {
		t.Fatal(err)
	}

	databaseStopAreaGroup.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseStopAreaGroup.ModelName = "2017-01-02"
	err = Database.Insert(&databaseStopAreaGroup)
	if err != nil {
		t.Fatal(err)
	}

	// Insert Data in the test db
	databaseLine := DatabaseLine{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "line",
		Codes:           "{}",
		Attributes:      "{}",
		References:      "{}",
	}

	Database.AddTableWithName(databaseLine, "lines")
	err = Database.Insert(&databaseLine)
	if err != nil {
		t.Fatal(err)
	}

	databaseLine.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseLine.ModelName = "2017-01-02"
	err = Database.Insert(&databaseLine)
	if err != nil {
		t.Fatal(err)
	}

	databaseLineGroup := DatabaseLineGroup{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "lineGroup",
		ShortName:       "line_group_short_name",
		LineIds:         `["d0eebc99-9c0b","e0eebc99-9c0b"]`,
	}

	Database.AddTableWithName(databaseLineGroup, "line_groups")
	err = Database.Insert(&databaseLineGroup)
	if err != nil {
		t.Fatal(err)
	}

	databaseLineGroup.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseLineGroup.ModelName = "2017-01-02"
	err = Database.Insert(&databaseLineGroup)
	if err != nil {
		t.Fatal(err)
	}

	databaseVehicleJourney := DatabaseVehicleJourney{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "vehicleJourney",
		LineId:          "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		DirectionType:   "",
		Codes:           "{}",
		Attributes:      "{}",
		References:      "{}",
	}

	Database.AddTableWithName(databaseVehicleJourney, "vehicle_journeys")
	err = Database.Insert(&databaseVehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	databaseVehicleJourney.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseVehicleJourney.ModelName = "2017-01-02"
	err = Database.Insert(&databaseVehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	databaseStopVisit := DatabaseStopVisit{
		Id:               "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug:  "referential",
		ModelName:        "2017-01-01",
		StopAreaId:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		VehicleJourneyId: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Codes:            "{}",
		Schedules:        "[]",
		Attributes:       "{}",
		References:       "{}",
	}

	Database.AddTableWithName(databaseStopVisit, "stop_visits")
	err = Database.Insert(&databaseStopVisit)
	if err != nil {
		t.Fatal(err)
	}

	databaseStopVisit.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseStopVisit.ModelName = "2017-01-02"
	err = Database.Insert(&databaseStopVisit)
	if err != nil {
		t.Fatal(err)
	}

	databaseOperator := DatabaseOperator{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		Name:            "operator",
		Codes:           "{}",
		ModelName:       "2017-01-01",
	}

	Database.AddTableWithName(databaseOperator, "operators")
	err = Database.Insert(&databaseOperator)
	if err != nil {
		t.Fatal(err)
	}

	databaseOperator.Id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"
	databaseOperator.ModelName = "2017-01-02"
	err = Database.Insert(&databaseOperator)
	if err != nil {
		t.Fatal(err)
	}
}

func check_record_count(expected_count int, t *testing.T) {
	table_names := []string{
		"stop_areas",
		"stop_area_groups",
		"lines",
		"line_groups",
		"vehicle_journeys",
		"stop_visits",
		"operators",
	}

	for i := range table_names {
		r, err := Database.Exec(fmt.Sprintf("select * from %v;", table_names[i]))
		if err != nil {
			t.Errorf("Error while executing query: %v", err)
			continue
		}
		i, err := r.RowsAffected()
		if err != nil {
			t.Errorf("Error while executing query: %v", err)
			continue
		}
		if int(i) != expected_count {
			t.Errorf("Wrong number of records in db, expected: %v got: %v", expected_count, i)
		}
	}
}

func test_Purifier_Purge(date time.Time, expected_count int, t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	fill_purifier_test_db(t)

	clock.SetDefaultClock(clock.NewFakeClockAt(date))

	purifier := NewPurifier(0)
	err := purifier.Purge()
	if err != nil {
		t.Fatal(err)
	}

	check_record_count(expected_count, t)
}

func Test_Purifier_Purge_0day(t *testing.T) {
	test_Purifier_Purge(time.Date(2017, time.January, 01, 0, 0, 0, 0, time.UTC), 2, t)
}

func Test_Purifier_Purge_1day(t *testing.T) {
	test_Purifier_Purge(time.Date(2017, time.January, 02, 0, 0, 0, 0, time.UTC), 1, t)
}

func Test_Purifier_Purge_2days(t *testing.T) {
	test_Purifier_Purge(time.Date(2017, time.January, 04, 0, 0, 0, 0, time.UTC), 0, t)
}
