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

	subscription   *Subscription
	lastTimeUpdate time.Time
}

func (smlc *stopMonitoringLastChange) SetSubscription(sub *Subscription) {
	smlc.subscription = sub
}

func (smlc *stopMonitoringLastChange) Haschanged(stopVisit model.StopVisit) bool {
	option, ok := smlc.subscription.subscriptionOptions["ChangeBeforeUpdates"]
	if ok {
		duration := smlc.getOptionDuration(option)
		refTime := stopVisit.ReferenceTime()
		if smlc.lastTimeUpdate.Add(duration).After(refTime) {
			return false
		}
	}
	return true
}

func (smlc *stopMonitoringLastChange) UpdateState(stopVisit model.StopVisit) bool {
	smlc.lastTimeUpdate = stopVisit.ReferenceTime()
	return true
}

type estimatedTimeTable struct {
	optionParser

	subscription    *Subscription
	lastTimeUpdate  time.Time
	departureStatus model.StopVisitDepartureStatus
	vehicleAtStop   bool
}

func (ettlc *estimatedTimeTable) SetSubscription(sub *Subscription) {
	ettlc.subscription = sub
}

func (ettlc *estimatedTimeTable) InitState(sv *model.StopVisit, sub *Subscription) {
	ettlc.subscription = sub
	ettlc.lastTimeUpdate = sv.ReferenceTime()
	ettlc.vehicleAtStop = sv.VehicleAtStop
}

func (ettlc *estimatedTimeTable) Haschanged(stopVisit *model.StopVisit) bool {
	if ettlc.departureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED {
		return false
	}

	if stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED {
		ettlc.departureStatus = model.STOP_VISIT_DEPARTURE_DEPARTED
		return true
	}

	if stopVisit.VehicleAtStop == true {
		ettlc.vehicleAtStop = true
		return true
	}

	option, ok := ettlc.subscription.subscriptionOptions["ChangeBeforeUpdates"]
	if ok {
		duration := ettlc.getOptionDuration(option)
		refTime := stopVisit.ReferenceTime()
		if ettlc.lastTimeUpdate.Add(duration).After(refTime) {
			return false
		}
	}
	return true
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
