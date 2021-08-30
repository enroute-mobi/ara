package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRISubscribeRequestHandler struct {
	xmlRequest  *siri.XMLSubscriptionRequest
	referential *core.Referential
}

func (handler *SIRISubscribeRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRISubscribeRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRISubscribeRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("SubscribeRequest %s\n", handler.xmlRequest.MessageIdentifier())
	t := handler.referential.Clock().Now()

	response, err := connector.(core.SubscriptionRequestDispatcher).Dispatch(handler.xmlRequest, message)
	if err != nil {
		siriErrorWithRequest("NotFound", err.Error(), handler.xmlRequest.RawXML(), string(handler.referential.Slug()), rw)
		return
	}

	xmlResponse, err := response.BuildXML()
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), string(handler.referential.Slug()), rw)
		return
	}

	// Wrap soap and send response
	soapEnvelope := remote.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	n, err := soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), string(handler.referential.Slug()), rw)
		return
	}

	message.ProcessingTime = handler.referential.Clock().Since(t).Seconds()
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ResponseRawMessage = xmlResponse
	message.ResponseSize = n
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)
}
