package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRISituationExchangeDeliveriesResponseHandler struct {
	xmlRequest  *sxml.XMLNotifySituationExchange
	referential *core.Referential
}

func (handler *SIRISituationExchangeDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRISituationExchangeDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRISituationExchangeDeliveriesResponseHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifySituationExchange %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	collectedRefs := params.connector.(core.SituationExchangeSubscriptionCollector).HandleNotifySituationExchange(handler.xmlRequest)
	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = audit.NOTIFY_SITUATION_EXCHANGE
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.SituationExchangesDeliveries() {
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
	params.message.Lines = collectedRefs.GetLines()
	params.message.StopAreas = collectedRefs.GetStopAreas()

	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
