package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRISubscriptionTerminatedNotificationHandler struct {
	xmlRequest  *siri.XMLSubscriptionTerminatedNotification
	referential *core.Referential
}

func (handler *SIRISubscriptionTerminatedNotificationHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRISubscriptionTerminatedNotificationHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRISubscriptionTerminatedNotificationHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("SubscriptionTerminatedNotification to cancel subscription: %s", handler.xmlRequest.SubscriptionRef())

	t := clock.DefaultClock().Now()

	params.connector.(core.SubscriptionRequestDispatcher).HandleSubscriptionTerminatedNotification(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = "SubscriptionTerminatedNotification"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.SubscriptionIdentifiers = []string{handler.xmlRequest.SubscriptionRef()}
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
