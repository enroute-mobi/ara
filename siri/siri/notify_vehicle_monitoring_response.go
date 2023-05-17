package siri

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRINotifyVehicleMonitoring struct {
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

	VehicleActivities []*SIRIVehicleActivity

	SortPayloadForTest bool
}

func (notify *SIRINotifyVehicleMonitoring) ErrorString() string {
	return fmt.Sprintf("%v: %v", notify.errorType(), notify.ErrorText)
}

func (notify *SIRINotifyVehicleMonitoring) errorType() string {
	if notify.ErrorType == "OtherError" {
		return fmt.Sprintf("%v %v", notify.ErrorType, notify.ErrorNumber)
	}
	return notify.ErrorType
}

func (notify *SIRINotifyVehicleMonitoring) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("vehicle_monitoring_notify%s.template", envType)

	if notify.SortPayloadForTest {
		sort.Sort(SortByVehicleMonitoringRef{notify.VehicleActivities})
	}

	if err := templates.ExecuteTemplate(&buffer, templateName, notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
