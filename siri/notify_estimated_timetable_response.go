package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRINotifyEstimatedTimeTable struct {
	Address                   string
	RequestMessageRef         string
	ProducerRef               string
	ResponseMessageIdentifier string
	SubscriberRef             string
	SubscriptionIdentifier    string

	ResponseTimestamp time.Time
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string

	EstimatedJourneyVersionFrames []*SIRIEstimatedJourneyVersionFrame
}

func (notify *SIRINotifyEstimatedTimeTable) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_timetable_notify.template", notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
