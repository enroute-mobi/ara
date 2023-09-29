package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/remote"
)

type SIRIGeneralMessageRequestDeliveriesResponseHandler struct {
	xmlRequest  *sxml.XMLNotifyGeneralMessage
	referential *core.Referential
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_GENERAL_MESSAGE_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIGeneralMessageRequestDeliveriesResponseHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifyGeneralMessage: %s", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	params.connector.(core.GeneralMessageSubscriptionCollector).HandleNotifyGeneralMessage(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	emptySOAPResponse := params.connector.Partner().PartnerSettings.SiriSoapEmptyResponseOnNotification()
	if params.envelopeType == remote.SOAP_SIRI_ENVELOPE && emptySOAPResponse {
		buffer := remote.NewSIRIBuffer(params.envelopeType)
		buffer.WriteXML("")
		buffer.WriteTo(params.rw)
	}

	params.message.Type = "NotifyGeneralMessage"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.GeneralMessagesDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
		if !delivery.Status() {
			params.message.Status = "Error"
		}
	}
	subs := make([]string, 0, len(subIds))
	for k := range subIds {
		subs = append(subs, k)
	}
	params.message.SubscriptionIdentifiers = subs
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
