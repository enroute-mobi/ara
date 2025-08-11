package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type SIRIFacilityMonitoringResponse struct {
	SIRIFacilityMonitoringDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIFacilityMonitoringDelivery struct {
	SIRIFacilityCondition

	RequestMessageRef string
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string
	ResponseTimestamp time.Time
}

type SIRIFacilityCondition struct {
	FacilityRef    string
	FacilityStatus string
}

func (response *SIRIFacilityMonitoringResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("facility_monitoring_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIFacilityMonitoringDelivery) ErrorString() string {
	return fmt.Sprintf("%v: %v", delivery.errorType(), delivery.ErrorText)
}

func (delivery *SIRIFacilityMonitoringDelivery) errorType() string {
	if delivery.ErrorType == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", delivery.ErrorType, delivery.ErrorNumber)
	}
	return delivery.ErrorType
}

func (delivery *SIRIFacilityMonitoringDelivery) BuildFacilityMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "facility_monitoring_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIFacilityMonitoringDelivery) BuildFacilityMonitoringDeliveryXMLRaw() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "facility_monitoring_delivery_raw.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIFacilityCondition) BuildFacilityConditionXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "facility_condition.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIFacilityCondition) BuildFacilityConditionXMLRaw() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "facility_condition_raw.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
