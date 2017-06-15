package model

import (
	"testing"
	"time"
)

func Test_CompositeStopVisitSelector_Empty(t *testing.T) {
	stopVisit := StopVisit{}

	selector := CompositeStopVisitSelector([]StopVisitSelector{})

	if !selector(stopVisit) {
		t.Errorf("Empty selector should return true, got false")
	}
}

func Test_StopVisitSelectorByTime(t *testing.T) {
	startTime := time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(2017, time.April, 1, 2, 0, 0, 0, time.UTC)

	selector := CompositeStopVisitSelector([]StopVisitSelector{StopVisitSelectorByTime(startTime, endTime)})

	stopVisit := StopVisit{
		Schedules: NewStopVisitSchedules(),
	}
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC), time.Time{})

	if !selector(stopVisit) {
		t.Errorf("Selector should return true, got false")
	}

	stopVisit2 := StopVisit{
		Schedules: NewStopVisitSchedules(),
	}
	stopVisit2.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 2, 1, 0, 0, 0, time.UTC), time.Time{})

	if selector(stopVisit2) {
		t.Errorf("Selector should return false, got true")
	}
}

func Test_StopVisitSelectorByLine(t *testing.T) {
	objectid := NewObjectID("kind", "value")
	selector := CompositeStopVisitSelector([]StopVisitSelector{StopVisitSelectorByLine(objectid)})

	model := NewMemoryModel()

	line := model.Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := model.VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Save()

	if !selector(stopVisit) {
		t.Errorf("Selector should return true, got false")
	}

	objectid2 := NewObjectID("kind", "value2")
	line2 := model.Lines().New()
	line2.SetObjectID(objectid2)
	line2.Save()

	vehicleJourney2 := model.VehicleJourneys().New()
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := model.StopVisits().New()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit2.Save()

	if selector(stopVisit2) {
		t.Errorf("Selector should return false, got true")
	}
}

func Test_CompositeStopVisitSelector(t *testing.T) {
	startTime := time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(2017, time.April, 1, 2, 0, 0, 0, time.UTC)
	objectid := NewObjectID("kind", "value")

	selector := CompositeStopVisitSelector([]StopVisitSelector{StopVisitSelectorByLine(objectid), StopVisitSelectorByTime(startTime, endTime)})

	model := NewMemoryModel()

	// Good VehicleJourney
	line := model.Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := model.VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC), time.Time{})
	stopVisit.Save()

	if !selector(stopVisit) {
		t.Errorf("Selector should return true, got false")
	}

	// Wrong Schedule
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 2, 1, 0, 0, 0, time.UTC), time.Time{})

	if selector(stopVisit) {
		t.Errorf("Selector should return false, got true")
	}

	// Wrong Line
	objectid2 := NewObjectID("kind", "value2")
	line2 := model.Lines().New()
	line2.SetObjectID(objectid2)
	line2.Save()

	vehicleJourney2 := model.VehicleJourneys().New()
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := model.StopVisits().New()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC), time.Time{})
	stopVisit2.Save()

	if selector(stopVisit2) {
		t.Errorf("Selector should return false, got true")
	}

	// All wrong
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 2, 1, 0, 0, 0, time.UTC), time.Time{})

	if selector(stopVisit2) {
		t.Errorf("Selector should return false, got true")
	}
}
