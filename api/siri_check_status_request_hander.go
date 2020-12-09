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

type SIRICheckStatusRequestHandler struct {
	xmlRequest *siri.XMLCheckStatusRequest
}

func (handler *SIRICheckStatusRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRICheckStatusRequestHandler) ConnectorType() string {
	return "siri-check-status-server"
}

func (handler *SIRICheckStatusRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("CheckStatus %s", handler.xmlRequest.MessageIdentifier())

	t := clock.DefaultClock().Now()

	response, err := connector.(core.CheckStatusServer).CheckStatus(handler.xmlRequest, message)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}
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

	message.Type = "CheckStatusRequest"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ResponseRawMessage = xmlResponse
	message.ResponseSize = n
	message.ProcessingTime = time.Since(t).Seconds()
	audit.CurrentBigQuery().WriteEvent(message)
}
