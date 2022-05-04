package ls

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
)

type schedulesHandler struct{}

func (sh *schedulesHandler) handleArrivalTime(sc, lssc *model.StopVisitSchedule, duration time.Duration) bool {
	if sc.ArrivalTime().IsZero() {
		return false
	}
	if lssc.ArrivalTime().IsZero() {
		return true
	}
	return !(sc.ArrivalTime().Before(lssc.ArrivalTime().Add(duration)) && sc.ArrivalTime().After(lssc.ArrivalTime().Add(-duration)))
}

func (sh *schedulesHandler) handleDepartedTime(sc, lssc *model.StopVisitSchedule, duration time.Duration) bool {
	if sc.DepartureTime().IsZero() {
		return false
	}
	if lssc.DepartureTime().IsZero() {
		return true
	}
	return !(sc.DepartureTime().Before(lssc.DepartureTime().Add(duration)) && sc.DepartureTime().After(lssc.DepartureTime().Add(-duration)))
}

func (sh *schedulesHandler) handleArrivalStatus(svAs model.StopVisitArrivalStatus, ettlcAs model.StopVisitArrivalStatus) bool {
	if svAs == ettlcAs {
		return false
	}

	if svAs == model.STOP_VISIT_ARRIVAL_MISSED || svAs == model.STOP_VISIT_ARRIVAL_NOT_EXPECTED || svAs == model.STOP_VISIT_ARRIVAL_CANCELLED || svAs == model.STOP_VISIT_ARRIVAL_NOREPORT {
		return true
	}

	return false
}

func (sh *schedulesHandler) handleDepartureStatus(svDs model.StopVisitDepartureStatus, ettlcDs model.StopVisitDepartureStatus) bool {
	if svDs == ettlcDs {
		return false
	}

	if svDs == model.STOP_VISIT_DEPARTURE_NOREPORT || svDs == model.STOP_VISIT_DEPARTURE_CANCELLED || svDs == model.STOP_VISIT_DEPARTURE_DEPARTED {
		return true
	}

	return false
}
