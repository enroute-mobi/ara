package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGetStopMonitoringRequest struct {
	SIRIStopMonitoringRequest

	RequestorRef string
}

type SIRIStopMonitoringRequest struct {
	MessageIdentifier string
	MonitoringRef     string
	StopVisitTypes    string

	RequestTimestamp time.Time
}

func NewSIRIGetStopMonitoringRequest(
	messageIdentifier,
	monitoringRef,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetStopMonitoringRequest {
	request := &SIRIGetStopMonitoringRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.MonitoringRef = monitoringRef
	request.RequestTimestamp = requestTimestamp
	return request
}

func (request *SIRIGetStopMonitoringRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_stop_monitoring_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIStopMonitoringRequest) BuildStopMonitoringRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "stop_monitoring_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
