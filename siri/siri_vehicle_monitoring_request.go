package siri

import (
	"bytes"
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

func (request *SIRIGetVehicleMonitoringRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "get_vehicle_monitoring_request.template", request); err != nil {
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
	return buffer.String(), nil
}
