package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIDeleteSubscriptionRequestHandler struct {
	xmlRequest *siri.XMLDeleteSubscriptionRequest
	Partner    core.Partner
}

func (handler *SIRIDeleteSubscriptionRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIDeleteSubscriptionRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRIDeleteSubscriptionRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("DeleteSubscription %s cancel subscription: %s\n", handler.xmlRequest.MessageIdentifier(), handler.xmlRequest.SubscriptionRef())

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
