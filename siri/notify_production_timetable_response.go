package siri

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRINotifyProductionTimeTable struct {
	ProducerRef            string
	SubscriptionIdentifier string

	ResponseTimestamp time.Time
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string

	DatedTimetableVersionFrames []*SIRIDatedTimetableVersionFrame
}

func (notify *SIRINotifyProductionTimeTable) ErrorString() string {
	return fmt.Sprintf("%v: %v", notify.errorType(), notify.ErrorText)
}

func (notify *SIRINotifyProductionTimeTable) errorType() string {
	if notify.ErrorType == "OtherError" {
		return fmt.Sprintf("%v %v", notify.ErrorType, notify.ErrorNumber)
	}
	return notify.ErrorType
}

func (notify *SIRINotifyProductionTimeTable) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("production_timetable_notify%s.template", envType)

	// order StopPointRef lexicographically inside DatedCalls
	for _, dtvf := range notify.DatedTimetableVersionFrames {
		for _, dvj := range dtvf.DatedVehicleJourneys {
			sort.Slice(dvj.DatedCalls, func(i, j int) bool {
				return strings.ToLower(dvj.DatedCalls[i].StopPointRef) <
					strings.ToLower(dvj.DatedCalls[j].StopPointRef)
			})
		}
	}

	if err := templates.ExecuteTemplate(&buffer, templateName, notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}