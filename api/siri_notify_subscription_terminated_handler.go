package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRINotifySubscriptionTerminatedHandler struct {
	xmlRequest  *sxml.XMLNotifySubscriptionTerminated
	referential *core.Referential
}

func (handler *SIRINotifySubscriptionTerminatedHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRINotifySubscriptionTerminatedHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRINotifySubscriptionTerminatedHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifySubscriptionTerminated %s to cancel subscription: %s", handler.xmlRequest.ResponseMessageIdentifier(), handler.xmlRequest.SubscriptionRef())

	t := clock.DefaultClock().Now()

	params.connector.(core.SubscriptionRequestDispatcher).HandleNotifySubscriptionTerminated(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = "NotifySubscriptionTerminated"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()
	params.message.SubscriptionIdentifiers = []string{handler.xmlRequest.SubscriptionRef()}
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
