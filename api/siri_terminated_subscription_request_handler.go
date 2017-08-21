package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRITerminatedSubscriptionRequestHandler struct {
	xmlRequest *siri.XMLTerminatedSubscriptionRequest
	Partner    core.Partner
}

func (handler *SIRITerminatedSubscriptionRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRITerminatedSubscriptionRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRITerminatedSubscriptionRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("TerminatedSubscription %s cancel subscription: %s\n", handler.xmlRequest.MessageIdentifier(), handler.xmlRequest.SubscriptionRef())

	response := connector.(core.SubscriptionRequestDispatcher).CancelSubscription(handler.xmlRequest)

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
