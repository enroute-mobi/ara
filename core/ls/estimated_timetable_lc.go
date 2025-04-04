package ls

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
)

type EstimatedTimetableLastChange struct {
	lastState
	optionParser
	schedulesHandler

	schedules       *schedules.StopVisitSchedules
	departureStatus model.StopVisitDepartureStatus
	arrivalStatuts  model.StopVisitArrivalStatus
	vehicleAtStop   bool
}

func NewEstimatedTimetableLastChange(sv *model.StopVisit, sub subscription) *EstimatedTimetableLastChange {
	ettlc := &EstimatedTimetableLastChange{}
	ettlc.SetSubscription(sub)
	ettlc.UpdateState(sv)
	return ettlc
}

func (ettlc *EstimatedTimetableLastChange) UpdateState(sv *model.StopVisit) {
	ettlc.vehicleAtStop = sv.VehicleAtStop
	ettlc.schedules = sv.Schedules.Copy()
	ettlc.departureStatus = sv.DepartureStatus
	ettlc.arrivalStatuts = sv.ArrivalStatus
}

func (ettlc *EstimatedTimetableLastChange) Haschanged(stopVisit *model.StopVisit) bool {
	// Don't send info on cancelled or departed SV
	if ettlc.departureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED || ettlc.departureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED || ettlc.arrivalStatuts == model.STOP_VISIT_ARRIVAL_ARRIVED || ettlc.arrivalStatuts == model.STOP_VISIT_ARRIVAL_CANCELLED {
		return false
	}

	// Check Departure Status
	if ettlc.handleArrivalStatus(stopVisit.ArrivalStatus, ettlc.arrivalStatuts) || ettlc.handleDepartureStatus(stopVisit.DepartureStatus, ettlc.departureStatus) {
		return true
	}

	// Check VehicleAtStop
	if ettlc.vehicleAtStop != stopVisit.VehicleAtStop && stopVisit.VehicleAtStop {
		return true
	}

	// Check Schedules
	option := ettlc.subscription.SubscriptionOption("ChangeBeforeUpdates")
	if option == "" {
		return true
	}

	duration := ettlc.getOptionDuration(option)
	if duration == 0 {
		duration = 1 * time.Minute
	}

	orderMap := []schedules.StopVisitScheduleType{schedules.Actual, schedules.Expected, schedules.Aimed}
	for _, kind := range orderMap {
		ok := ettlc.handleArrivalTime(stopVisit.Schedules.Schedule(kind), ettlc.schedules.Schedule(kind), duration)
		ok = ok || ettlc.handleDepartedTime(stopVisit.Schedules.Schedule(kind), ettlc.schedules.Schedule(kind), duration)
		if ok {
			return true
		}
	}
	return false
}
