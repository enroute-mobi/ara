package ls

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
)

type StopMonitoringLastChange struct {
	lastState
	optionParser
	schedulesHandler

	schedules       *schedules.StopVisitSchedules
	departureStatus model.StopVisitDepartureStatus
	arrivalStatuts  model.StopVisitArrivalStatus
}

func NewStopMonitoringLastChange(sv *model.StopVisit, sub subscription) *StopMonitoringLastChange {
	smlc := &StopMonitoringLastChange{}
	smlc.SetSubscription(sub)
	smlc.UpdateState(sv)
	return smlc
}

func (smlc *StopMonitoringLastChange) UpdateState(stopVisit *model.StopVisit) bool {
	smlc.schedules = stopVisit.Schedules.Copy()
	smlc.arrivalStatuts = stopVisit.ArrivalStatus
	smlc.departureStatus = stopVisit.DepartureStatus

	return true
}

func (smlc *StopMonitoringLastChange) Haschanged(stopVisit *model.StopVisit) bool {
	// Don't send info on cancelled or departed SV
	if smlc.departureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED || smlc.departureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED || smlc.arrivalStatuts == model.STOP_VISIT_ARRIVAL_CANCELLED {
		return false
	}

	// Check Arrival and Departure status
	if smlc.handleArrivalStatus(stopVisit.ArrivalStatus, smlc.arrivalStatuts) || smlc.handleDepartureStatus(stopVisit.DepartureStatus, smlc.departureStatus) {
		return true
	}

	option := smlc.subscription.SubscriptionOption("ChangeBeforeUpdates")
	if option == "" {
		return true
	}

	duration := smlc.getOptionDuration(option)
	if duration == 0 {
		duration = 1 * time.Minute
	}

	orderMap := []schedules.StopVisitScheduleType{schedules.Actual, schedules.Expected, schedules.Aimed}
	for _, kind := range orderMap {
		ok := smlc.handleArrivalTime(stopVisit.Schedules.Schedule(kind), smlc.schedules.Schedule(kind), duration)
		ok = ok || smlc.handleDepartedTime(stopVisit.Schedules.Schedule(kind), smlc.schedules.Schedule(kind), duration)
		if ok {
			return true
		}
	}

	return false
}
