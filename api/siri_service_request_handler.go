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

type SIRIServiceRequestHandler struct {
	xmlRequest *siri.XMLSiriServiceRequest
}

func (handler *SIRIServiceRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIServiceRequestHandler) ConnectorType() string {
	return "siri-service-request-broadcaster"
}

func (handler *SIRIServiceRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("SiriService %s\n", handler.xmlRequest.MessageIdentifier())

	t := clock.DefaultClock().Now()

	response := connector.(core.ServiceRequestBroadcaster).HandleRequests(handler.xmlRequest, message)
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

	message.Type = "SiriServiceRequest"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ResponseRawMessage = xmlResponse
	message.ResponseSize = n
	message.ProcessingTime = time.Since(t).Seconds()
	audit.CurrentBigQuery().WriteEvent(message)
}
