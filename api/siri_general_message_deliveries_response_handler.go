package api

import (
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIGeneralMessageRequestDeliveriesResponseHandler struct {
	xmlRequest  *siri.XMLNotifyGeneralMessage
	referential *core.Referential
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_GENERAL_MESSAGE_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("NotifyGeneralMessage: %s", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	connector.(core.GeneralMessageSubscriptionCollector).HandleNotifyGeneralMessage(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)

	message.Type = "NotifyGeneralMessage"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ProcessingTime = time.Since(t).Seconds()
	message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.GeneralMessagesDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
		if !delivery.Status() {
			message.Status = "Error"
		}
	}
	subs := make([]string, 0, len(subIds))
	for k := range subIds {
		subs = append(subs, k)
	}
	message.SubscriptionIdentifiers = subs
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)
}
