package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
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

func (handler *SIRICheckStatusRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("CheckStatus %s\n", handler.xmlRequest.MessageIdentifier())

	response, err := connector.(core.CheckStatusServer).CheckStatus(handler.xmlRequest)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
	}
	xmlResponse, err := response.BuildXML()
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
	}

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err = soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
	}
}
