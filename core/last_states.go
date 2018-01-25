package core

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/model"
)

type lastState interface {
	SetSubscription(*Subscription)
}

type stopMonitoringLastChange struct {
	optionParser
	schedulesHandler

	subscription *Subscription

	schedules       model.StopVisitSchedules
	departureStatus model.StopVisitDepartureStatus
	arrivalStatuts  model.StopVisitArrivalStatus
}

func (smlc *stopMonitoringLastChange) InitState(sv *model.StopVisit, sub *Subscription) {
	smlc.SetSubscription(sub)
	smlc.UpdateState(sv)
}

func (smlc *stopMonitoringLastChange) SetSubscription(sub *Subscription) {
	smlc.subscription = sub
}

func (smlc *stopMonitoringLastChange) UpdateState(stopVisit *model.StopVisit) bool {
	smlc.schedules = stopVisit.Schedules.Copy()
	smlc.arrivalStatuts = stopVisit.ArrivalStatus
	smlc.departureStatus = stopVisit.DepartureStatus

	return true
}

func (smlc *stopMonitoringLastChange) Haschanged(stopVisit model.StopVisit) bool {
	option, ok := smlc.subscription.subscriptionOptions["ChangeBeforeUpdates"]
	if !ok {
		return true
	}

	if stopVisit.ArrivalStatus == model.STOP_VISIT_ARRIVAL_CANCELLED && smlc.arrivalStatuts == stopVisit.ArrivalStatus {
		return false
	}

	if stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED && smlc.departureStatus == stopVisit.DepartureStatus {
		return false
	}

	if smlc.handleArrivalStatus(stopVisit.ArrivalStatus, smlc.arrivalStatuts) || smlc.handleDepartureStatus(stopVisit.DepartureStatus, smlc.departureStatus) {
		return true
	}

	duration := smlc.getOptionDuration(option)
	if duration == 0 {
		duration = 1 * time.Minute
	}

	for kind, _ := range stopVisit.Schedules {
		ok = smlc.handleArrivalTime(stopVisit.Schedules.Schedule(kind), smlc.schedules.Schedule(kind), duration)
		ok = ok || smlc.handleDepartedTime(stopVisit.Schedules.Schedule(kind), smlc.schedules.Schedule(kind), duration)
		if ok {
			return true
		}
	}

	return false
}

type estimatedTimeTableLastChange struct {
	optionParser
	schedulesHandler

	subscription *Subscription

	schedules       model.StopVisitSchedules
	departureStatus model.StopVisitDepartureStatus
	arrivalStatuts  model.StopVisitArrivalStatus
	vehicleAtStop   bool
}

func (ettlc *estimatedTimeTableLastChange) InitState(sv *model.StopVisit, sub *Subscription) {
	ettlc.SetSubscription(sub)
	ettlc.UpdateState(sv)
}

func (ettlc *estimatedTimeTableLastChange) SetSubscription(sub *Subscription) {
	ettlc.subscription = sub
}

func (ettlc *estimatedTimeTableLastChange) UpdateState(sv *model.StopVisit) {
	ettlc.vehicleAtStop = sv.VehicleAtStop
	ettlc.schedules = sv.Schedules.Copy()
	ettlc.departureStatus = sv.DepartureStatus
	ettlc.arrivalStatuts = sv.ArrivalStatus
}

func (ettlc *estimatedTimeTableLastChange) Haschanged(stopVisit *model.StopVisit) bool {
	option, ok := ettlc.subscription.subscriptionOptions["ChangeBeforeUpdates"]
	if !ok {
		return true
	}

	if ettlc.departureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED {
		return false
	}

	if stopVisit.ArrivalStatus == model.STOP_VISIT_ARRIVAL_CANCELLED && smlc.arrivalStatuts == stopVisit.ArrivalStatus {
		return false
	}

	if stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED && smlc.departureStatus == stopVisit.DepartureStatus {
		return false
	}

	if ettlc.vehicleAtStop != stopVisit.VehicleAtStop {
		return true
	}

	if ettlc.handleArrivalStatus(stopVisit.ArrivalStatus, ettlc.arrivalStatuts) || ettlc.handleDepartureStatus(stopVisit.DepartureStatus, ettlc.departureStatus) {
		return true
	}

	duration := ettlc.getOptionDuration(option)
	if duration == 0 {
		duration = 1 * time.Minute
	}

	for kind, _ := range stopVisit.Schedules {
		ok = ettlc.handleArrivalTime(stopVisit.Schedules.Schedule(kind), ettlc.schedules.Schedule(kind), duration)
		ok = ok || ettlc.handleDepartedTime(stopVisit.Schedules.Schedule(kind), ettlc.schedules.Schedule(kind), duration)
		if ok {
			return true
		}
	}
	return false
}

type generalMessageLastChange struct {
	subscription *Subscription

	version int
}

func (sglc *generalMessageLastChange) InitState(situation *model.Situation, sub *Subscription) {
	sglc.SetSubscription(sub)
	sglc.UpdateState(situation)
}

func (sglc *generalMessageLastChange) SetSubscription(sub *Subscription) {
	sglc.subscription = sub
}

func (sglc *generalMessageLastChange) UpdateState(situation *model.Situation) bool {
	sglc.version = situation.Version
	return true
}

func (sglc *generalMessageLastChange) Haschanged(situation *model.Situation) bool {
	return !(situation.Version == sglc.version)
}

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

	if svDs == model.STOP_VISIT_DEPARTURE_NOREPORT || svDs == model.STOP_VISIT_DEPARTURE_CANCELLED {
		return true
	}

	return false
}

type optionParser struct{}

func (subscription *optionParser) getOptionDuration(option string) time.Duration {

	durationRegex := regexp.MustCompile(`P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?`)
	matches := durationRegex.FindStringSubmatch(strings.TrimSpace(option))

	if len(matches) == 0 {
		return 0
	}
	years := subscription.parseDuration(matches[1]) * 24 * 365 * time.Hour
	months := subscription.parseDuration(matches[2]) * 30 * 24 * time.Hour
	days := subscription.parseDuration(matches[3]) * 24 * time.Hour
	hours := subscription.parseDuration(matches[4]) * time.Hour
	minutes := subscription.parseDuration(matches[5]) * time.Minute
	seconds := subscription.parseDuration(matches[6]) * time.Second

	return time.Duration(years + months + days + hours + minutes + seconds)
}

func (subscription *optionParser) parseDuration(value string) time.Duration {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return time.Duration(parsed)
}
