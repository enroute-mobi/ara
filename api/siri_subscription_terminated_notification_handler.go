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

type SIRISubscriptionTerminatedNotificationHandler struct {
	xmlRequest *siri.XMLSubscriptionTerminatedNotification
}

func (handler *SIRISubscriptionTerminatedNotificationHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRISubscriptionTerminatedNotificationHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRISubscriptionTerminatedNotificationHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("SubscriptionTerminatedNotification to cancel subscription: %s", handler.xmlRequest.SubscriptionRef())

	t := clock.DefaultClock().Now()

	connector.(core.SubscriptionRequestDispatcher).HandleSubscriptionTerminatedNotification(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)

	message.Type = "SubscriptionTerminatedNotification"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ProcessingTime = time.Since(t).Seconds()
	message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	message.SubscriptionIdentifiers = []string{handler.xmlRequest.SubscriptionRef()}
	audit.CurrentBigQuery().WriteEvent(message)
}
