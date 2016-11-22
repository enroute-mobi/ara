package siri

type SIRIStopMonitoringRequest struct {
	MonitoringRef string
}

func NewSIRIStopMonitoringRequest(monitoringRef string) *SIRIStopMonitoringRequest {
	return &SIRIStopMonitoringRequest{MonitoringRef: monitoringRef}
}

func (request *SIRIStopMonitoringRequest) BuildXML() string {
	return ""
}
