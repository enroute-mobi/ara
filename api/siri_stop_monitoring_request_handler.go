package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/edwig/core"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type SIRIStopMonitoringRequestHandler struct {
	xmlRequest *siri.XMLGetStopMonitoring
}

func (handler *SIRIStopMonitoringRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIStopMonitoringRequestHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_REQUEST_BROADCASTER
}

func (handler *SIRIStopMonitoringRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("StopMonitoring %s\n", handler.xmlRequest.MessageIdentifier())

	response := connector.(core.StopMonitoringRequestBroadcaster).RequestStopArea(handler.xmlRequest)
	xmlResponse, err := response.BuildXML()
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err = soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}
}
