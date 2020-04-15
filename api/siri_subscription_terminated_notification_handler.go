package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/edwig/core"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/siri"
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

func (handler *SIRISubscriptionTerminatedNotificationHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("SubscriptionTerminatedNotification to cancel subscription: %s", handler.xmlRequest.SubscriptionRef())

	connector.(core.SubscriptionRequestDispatcher).HandleSubscriptionTerminatedNotification(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
