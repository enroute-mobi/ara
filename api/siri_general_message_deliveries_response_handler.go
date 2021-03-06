package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIGeneralMessageRequestDeliveriesResponseHandler struct {
	xmlRequest *siri.XMLNotifyGeneralMessage
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_GENERAL_MESSAGE_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("NotifyGeneralMessage: %s", handler.xmlRequest.ResponseMessageIdentifier())

	connector.(core.GeneralMessageSubscriptionCollector).HandleNotifyGeneralMessage(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
