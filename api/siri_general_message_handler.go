package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIGeneralMessageRequestHandler struct {
	xmlRequest *siri.XMLGeneralMessageRequest
}

func (handler *SIRIGeneralMessageRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIGeneralMessageRequestHandler) ConnectorType() string {
	return core.SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER
}

func (handler *SIRIGeneralMessageRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("General Message %s\n", handler.xmlRequest.MessageIdentifier())

	tmp := connector.(*core.SIRIGeneralMessageRequestBroadcaster)
	response, _ := tmp.Situations(handler.xmlRequest)
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
