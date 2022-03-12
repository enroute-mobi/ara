package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIDeleteSubscriptionRequestHandler struct {
	xmlRequest  *siri.XMLDeleteSubscriptionRequest
	referential *core.Referential
}

func (handler *SIRIDeleteSubscriptionRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIDeleteSubscriptionRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRIDeleteSubscriptionRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("DeleteSubscription %s cancel subscription: %s", handler.xmlRequest.MessageIdentifier(), handler.xmlRequest.SubscriptionRef())

	t := clock.DefaultClock().Now()

	response := connector.(core.SubscriptionRequestDispatcher).CancelSubscription(handler.xmlRequest, message)

	xmlResponse, err := response.BuildXML()
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), string(handler.referential.Slug()), rw)
		return
	}

	// Wrap soap and send response
	soapEnvelope := remote.NewSIRIBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	n, err := soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), string(handler.referential.Slug()), rw)
		return
	}

	message.Type = "DeleteSubscriptionRequest"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ResponseRawMessage = xmlResponse
	message.ResponseSize = n
	message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)
}
