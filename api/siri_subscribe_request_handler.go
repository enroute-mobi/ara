package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRISubscribeRequestHandler struct {
	xmlRequest *siri.XMLSubscriptionRequest
	Partner    core.Partner
}

func (handler *SIRISubscribeRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRISubscribeRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST
}

func (handler *SIRISubscribeRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("SubscribeRequest %s\n", handler.xmlRequest.MessageIdentifier())

	response := connector.(core.SubscriptionRequest).Dispatch(handler.xmlRequest)

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
