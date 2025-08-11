package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGetFacilityMonitoringRequest struct {
	SIRIFacilityMonitoringRequest

	RequestorRef string
}

type SIRIFacilityMonitoringRequest struct {
	MessageIdentifier string
	FacilityRef       string
	StopVisitTypes    string

	RequestTimestamp time.Time
}

func NewSIRIGetFacilityMonitoringRequest(
	messageIdentifier,
	facilityRef,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetFacilityMonitoringRequest {

	request := &SIRIGetFacilityMonitoringRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.FacilityRef = facilityRef
	request.RequestTimestamp = requestTimestamp

	return request
}

func (request *SIRIGetFacilityMonitoringRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "" && envelopeType[0] != "soap" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_facility_monitoring_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIFacilityMonitoringRequest) BuildFacilityMonitoringRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "facility_monitoring_request.template", request); err != nil {
		logger.Log.Debugf("Errorw hile executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIFacilityMonitoringRequest) BuildFacilityMonitoringRequestXMLRaw() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "facility_monitoring_request_raw.template", request); err != nil {
		logger.Log.Debugf("Errorw hile executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
