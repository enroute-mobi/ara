package siri

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGetVehicleMonitoringRequest struct {
	SIRIVehicleMonitoringRequest

	RequestorRef string
}

type SIRIVehicleMonitoringRequest struct {
	MessageIdentifier string
	LineRef           string

	RequestTimestamp time.Time
}

func NewSIRIGetVehicleMonitoringRequest() *SIRIGetVehicleMonitoringRequest {
	return &SIRIGetVehicleMonitoringRequest{}
}

func (request *SIRIGetVehicleMonitoringRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "" && envelopeType[0] != "soap" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_vehicle_monitoring_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIVehicleMonitoringRequest) BuildVehicleMonitoringRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "vehicle_monitoring_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (request *SIRIVehicleMonitoringRequest) BuildVehicleMonitoringRequestXMLRaw() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "vehicle_monitoring_request_raw.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
