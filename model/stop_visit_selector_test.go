package model

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/model/schedules"
)

func Test_CompositeStopVisitSelector_Empty(t *testing.T) {
	stopVisit := &StopVisit{}

	selector := CompositeStopVisitSelector([]StopVisitSelector{})

	if !selector(stopVisit) {
		t.Errorf("Empty selector should return true, got false")
	}
}

func Test_StopVisitSelectorByTime(t *testing.T) {
	startTime := time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(2017, time.April, 1, 2, 0, 0, 0, time.UTC)

	selector := CompositeStopVisitSelector([]StopVisitSelector{StopVisitSelectorByTime(startTime, endTime)})

	stopVisit := &StopVisit{
		Schedules: schedules.NewStopVisitSchedules(),
	}
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC), time.Time{})

	if !selector(stopVisit) {
		t.Errorf("Selector should return true, got false")
	}

	stopVisit2 := &StopVisit{
		Schedules: schedules.NewStopVisitSchedules(),
	}
	stopVisit2.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 2, 1, 0, 0, 0, time.UTC), time.Time{})

	if selector(stopVisit2) {
		t.Errorf("Selector should return false, got true")
	}
}

func Test_StopVisitSelectorByLine(t *testing.T) {
	model := NewTestMemoryModel()

	line := model.Lines().New()
	code := NewCode("codeSpace", "value")
	line.SetCode(code)
	line.Save()

	vehicleJourney := model.VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Save()

	selector := CompositeStopVisitSelector([]StopVisitSelector{StopVisitSelectorByLines([]LineId{line.id})})

	if !selector(stopVisit) {
		t.Errorf("Selector should return true, got false")
	}

	code2 := NewCode("codeSpace", "value2")
	line2 := model.Lines().New()
	line2.SetCode(code2)
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

	model := NewTestMemoryModel()

	// Good VehicleJourney
	line := model.Lines().New()
	code := NewCode("codeSpace", "value")
	line.SetCode(code)
	line.Save()

	vehicleJourney := model.VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 1, 1, 0, 0, 0, time.UTC), time.Time{})
	stopVisit.Save()

	selector := CompositeStopVisitSelector([]StopVisitSelector{StopVisitSelectorByLines([]LineId{line.id}), StopVisitSelectorByTime(startTime, endTime)})

	if !selector(stopVisit) {
		t.Errorf("Selector should return true, got false")
	}

	// Wrong Schedule
	stopVisit.Schedules.SetSchedule("aimed", time.Date(2017, time.April, 2, 1, 0, 0, 0, time.UTC), time.Time{})

	if selector(stopVisit) {
		t.Errorf("Selector should return false, got true")
	}

	// Wrong Line
	code2 := NewCode("codeSpace", "value2")
	line2 := model.Lines().New()
	line2.SetCode(code2)
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
