package core

import (
	"fmt"
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
		fmt.Println("salut les gens \n")
		duration := smlc.getOptionDuration(option)
		refTime := stopVisit.ReferenceTime()
		fmt.Println(refTime, "  ", smlc.lastTimeUpdate.Add(duration), "   ", duration)
		if smlc.lastTimeUpdate.Add(duration).After(refTime) {
			return false
		}
	}
	return true
}

func (smlc *stopMonitoringLastChange) UpdateState(stopVisit model.StopVisit) bool {
	smlc.lastTimeUpdate = stopVisit.ReferenceTime()
	fmt.Println("The reference time is  %v", stopVisit.ReferenceTime())
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
