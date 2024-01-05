package model

import "time"

type StopVisitSelector func(*StopVisit) bool

func StopVisitSelectorByTime(startTime, endTime time.Time) StopVisitSelector {
	return func(stopVisit *StopVisit) bool {
		if stopVisit.ReferenceTime().IsZero() || stopVisit.ReferenceTime().Before(startTime) || stopVisit.ReferenceTime().After(endTime) {
			return false
		}
		return true
	}
}

func StopVisitSelectorAfterTime(startTime time.Time) StopVisitSelector {
	return func(stopVisit *StopVisit) bool {
		if stopVisit.ReferenceTime().IsZero() || stopVisit.ReferenceTime().Before(startTime) {
			return false
		}
		return true
	}
}

func StopVisitSelectorBeforeTime(endTime time.Time) StopVisitSelector {
	return func(stopVisit *StopVisit) bool {
		if stopVisit.ReferenceTime().IsZero() || stopVisit.ReferenceTime().After(endTime) {
			return false
		}
		return true
	}
}

func StopVisitSelectByStopAreaId(stopAreaId StopAreaId) StopVisitSelector {
	return func(stopVisit *StopVisit) bool {
		return stopVisit.StopAreaId == stopAreaId
	}
}

func StopVisitSelectorByLine(code Code) StopVisitSelector {
	return func(stopVisit *StopVisit) bool {
		vehicleJourney := stopVisit.VehicleJourney()
		if vehicleJourney == nil {
			return false
		}
		line := vehicleJourney.Line()
		if line == nil {
			return false
		}
		lineCode, ok := line.Code(code.CodeSpace())
		if ok {
			return lineCode.Value() == code.Value()
		}
		return false
	}
}

func CompositeStopVisitSelector(selectors []StopVisitSelector) StopVisitSelector {
	return func(stopVisit *StopVisit) bool {
		for _, selector := range selectors {
			if !selector(stopVisit) {
				return false
			}
		}
		return true
	}
}
