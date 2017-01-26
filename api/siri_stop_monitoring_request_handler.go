package api

import "github.com/af83/edwig/core"

type SIRIStopMonitoringRequestHandler struct{}

func (handler *SIRIStopMonitoringRequestHandler) RequestorRef() string {
	return ""
}

func (handler *SIRIStopMonitoringRequestHandler) ConnectorType() string {
	return "siri-stop-monitoring-request-collector"
}

func (handler *SIRIStopMonitoringRequestHandler) XMLResponse(core.Connector) string {
	return ""
}
