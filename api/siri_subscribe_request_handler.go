package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRISubscribeRequestHandler struct {
	xmlRequest *siri.XMLSubscriptionRequest
}

func (handler *SIRISubscribeRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRISubscribeRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRISubscribeRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("SubscribeRequest %s\n", handler.xmlRequest.MessageIdentifier())

	response, err := connector.(core.SubscriptionRequestDispatcher).Dispatch(handler.xmlRequest)
	if err != nil {
		siriErrorWithRequest("NotFound", err.Error(), handler.xmlRequest.RawXML(), rw)
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

	_, err = soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}
}
