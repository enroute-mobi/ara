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

type SIRINotifySubscriptionTerminatedHandler struct {
	xmlRequest  *siri.XMLNotifySubscriptionTerminated
	referential *core.Referential
}

func (handler *SIRINotifySubscriptionTerminatedHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRINotifySubscriptionTerminatedHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRINotifySubscriptionTerminatedHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("NotifySubscriptionTerminated %s to cancel subscription: %s", handler.xmlRequest.ResponseMessageIdentifier(), handler.xmlRequest.SubscriptionRef())

	t := clock.DefaultClock().Now()

	connector.(core.SubscriptionRequestDispatcher).HandleNotifySubscriptionTerminated(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)

	message.Type = "NotifySubscriptionTerminated"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ProcessingTime = time.Since(t).Seconds()
	message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()
	message.SubscriptionIdentifiers = []string{handler.xmlRequest.SubscriptionRef()}
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)
}
