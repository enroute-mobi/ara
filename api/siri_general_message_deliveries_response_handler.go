package api

import (
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIGeneralMessageRequestDeliveriesResponseHandler struct {
	xmlRequest *siri.XMLNotifyGeneralMessage
	Partner    core.Partner
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_GENERAL_MESSAGE_DELIVERIES_RESPONSE_COLLECTOR
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("NotifyGeneralMessage %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	connector.(core.GeneralMessageSubscriptionCollector).HandleNotifyGeneralMessage(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
