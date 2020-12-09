package api

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
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

func (handler *SIRIStopMonitoringRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("StopMonitoring %s\n", handler.xmlRequest.MessageIdentifier())

	t := clock.DefaultClock().Now()

	response := connector.(core.StopMonitoringRequestBroadcaster).RequestStopArea(handler.xmlRequest, message)
	xmlResponse, err := response.BuildXML()
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	n, err := soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}

	message.Type = "StopMonitoringRequest"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ResponseRawMessage = xmlResponse
	message.ResponseSize = n
	message.ProcessingTime = time.Since(t).Seconds()
	audit.CurrentBigQuery().WriteEvent(message)
}
