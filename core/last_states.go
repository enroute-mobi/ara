package core

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type lastState interface {
	SetSubscription(*Subscription)
}

type stopMonitoringLastChange struct {
	optionParser
	schedulesHandler

	subscription *Subscription
	schedules    model.StopVisitSchedules

	departureStatus model.StopVisitDepartureStatus
	arrivalStatuts  model.StopVisitArrivalStatus
}

func (smlc *stopMonitoringLastChange) SetSubscription(sub *Subscription) {
	smlc.subscription = sub
}

func (smlc *stopMonitoringLastChange) InitState(sv *model.StopVisit, sub *Subscription) {
	smlc.subscription = sub
	sv.Schedules.Copy()
}

func (smlc *stopMonitoringLastChange) Haschanged(stopVisit model.StopVisit) bool {

	option, ok := smlc.subscription.subscriptionOptions["ChangeBeforeUpdates"]
	if !ok {
		return true
	}

	if smlc.handleArrivalStatus(stopVisit.ArrivalStatus, smlc.arrivalStatuts) || smlc.handleDepartureStatus(stopVisit.DepartureStatus, smlc.departureStatus) {
		return true
	}

	duration := smlc.getOptionDuration(option)
	for key, _ := range stopVisit.Schedules {
		ok = smlc.handleArrivalTime(stopVisit.Schedules[key], smlc.schedules[key], duration)
		ok = ok || smlc.handleDepartedTime(stopVisit.Schedules[key], smlc.schedules[key], duration)
		if ok {
			return true
		}
	}

	return false
}

func (smlc *stopMonitoringLastChange) UpdateState(stopVisit model.StopVisit) bool {
	smlc.schedules = stopVisit.Schedules.Copy()
	smlc.arrivalStatuts = stopVisit.ArrivalStatus
	smlc.departureStatus = stopVisit.DepartureStatus

	return true
}

type estimatedTimeTableLastChange struct {
	optionParser
	schedulesHandler

	subscription    *Subscription
	schedules       model.StopVisitSchedules
	departureStatus model.StopVisitDepartureStatus
	arrivalStatuts  model.StopVisitArrivalStatus
	vehicleAtStop   bool
}

func (ettlc *estimatedTimeTableLastChange) SetSubscription(sub *Subscription) {
	ettlc.subscription = sub
}

func (ettlc *estimatedTimeTableLastChange) InitState(sv *model.StopVisit, sub *Subscription) {
	ettlc.subscription = sub
	ettlc.schedules = sv.Schedules.Copy()
	ettlc.vehicleAtStop = sv.VehicleAtStop
	ettlc.departureStatus = sv.DepartureStatus
}

func (ettlc *estimatedTimeTableLastChange) Haschanged(stopVisit *model.StopVisit) bool {
	if ettlc.departureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED {
		return false
	}

	if stopVisit.VehicleAtStop == true {
		return true
	}

	if ettlc.handleArrivalStatus(stopVisit.ArrivalStatus, ettlc.arrivalStatuts) || ettlc.handleDepartureStatus(stopVisit.DepartureStatus, ettlc.departureStatus) {
		return true
	}

	option, ok := ettlc.subscription.subscriptionOptions["ChangeBeforeUpdates"]
	if !ok {
		return true
	}

	duration := ettlc.getOptionDuration(option)
	for key, _ := range stopVisit.Schedules {
		ok = ettlc.handleArrivalTime(stopVisit.Schedules[key], ettlc.schedules[key], duration)
		ok = ok || ettlc.handleDepartedTime(stopVisit.Schedules[key], ettlc.schedules[key], duration)
		if ok {
			return true
		}
	}
	return false
}

func (ettlc *estimatedTimeTableLastChange) UpdateState(sv model.StopVisit) {
	ettlc.vehicleAtStop = sv.VehicleAtStop
	ettlc.schedules = sv.Schedules.Copy()
	ettlc.departureStatus = sv.DepartureStatus
	ettlc.arrivalStatuts = sv.ArrivalStatus
}

type generalMessageLastChange struct {
	optionParser

	lastChange siri.SIRIGeneralMessage

	subscription *Subscription
}

func (sglc *generalMessageLastChange) Haschanged(situation model.Situation) bool {
	return true
}

func (sglc *generalMessageLastChange) UpdateState(sm siri.SIRIGeneralMessage) bool {
	sglc.lastChange = sm
	return true
}

func (sglc *generalMessageLastChange) SetSubscription(sub *Subscription) {
	sglc.subscription = sub
}

type schedulesHandler struct{}

func (sh *schedulesHandler) handleArrivalTime(sc, lssc *model.StopVisitSchedule, duration time.Duration) bool {
	return sh.handleTime(sc.ArrivalTime(), lssc.ArrivalTime(), duration)
}

func (sh *schedulesHandler) handleDepartedTime(sc, lssc *model.StopVisitSchedule, duration time.Duration) bool {
	return sh.handleTime(sc.DepartureTime(), lssc.DepartureTime(), duration)
}

func (sh *schedulesHandler) handleTime(scTime, lsscTime time.Time, duration time.Duration) bool {
	if scTime.IsZero() || !scTime.After(lsscTime.Add(2*time.Minute)) {
		return false
	}
	return true
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
