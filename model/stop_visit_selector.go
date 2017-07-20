package model

import "time"

type StopVisitSelector func(StopVisit) bool

func FindStopVisitBy(collection map[StopVisitId]*StopVisit, selectors ...StopVisitSelector) (stopVisits []StopVisit) {
	if len(collection) == 0 {
		return []StopVisit{}
	}
	for _, stopVisit := range collection {
		if checkSelector(*stopVisit, selectors...) {
			stopVisits = append(stopVisits, *stopVisit)
		}
	}
	return
}

func checkSelector(stopVisit StopVisit, selectors ...StopVisitSelector) bool {
	for _, selector := range selectors {
		if !selector(stopVisit) {
			return false
		}
	}
	return true
}

func CompositeStopVisitSelector(selectors []StopVisitSelector) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		for _, selector := range selectors {
			if !selector(stopVisit) {
				return false
			}
		}
		return true
	}
}

func StopVisitSelectorByTime(startTime, endTime time.Time) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		if stopVisit.ReferenceTime().IsZero() || stopVisit.ReferenceTime().Before(startTime) || stopVisit.ReferenceTime().After(endTime) {
			return false
		}
		return true
	}
}

func StopVisitSelectorFollowing(referenceTime time.Time) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		if stopVisit.ReferenceTime().After(referenceTime) {
			return true
		}
		return false
	}
}

func StopVisitSelectorByLine(objectid ObjectID) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		vehicleJourney := stopVisit.VehicleJourney()
		if vehicleJourney == nil {
			return false
		}
		line := vehicleJourney.Line()
		if line == nil {
			return false
		}
		lineObjectid, ok := line.ObjectID(objectid.Kind())
		if ok {
			return lineObjectid.Value() == objectid.Value()
		}
		return false
	}
}

func StopVisitSelectorByVehicleJourneyId(id VehicleJourneyId) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		return stopVisit.VehicleJourneyId == id
	}
}

func StopVisitSelectorByStopAreaId(id StopAreaId) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		return stopVisit.StopAreaId == id
	}
}

func StopVisitSelectorByStopAreaIds(ids []StopAreaId) StopVisitSelector {
	return func(stopVisit StopVisit) bool {
		for _, id := range ids {
			if stopVisit.StopAreaId == id {
				return true
			}
		}
		return false
	}
}
