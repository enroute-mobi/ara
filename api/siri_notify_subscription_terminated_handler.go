package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/edwig/core"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type SIRINotifySubscriptionTerminatedHandler struct {
	xmlRequest *siri.XMLNotifySubscriptionTerminated
}

func (handler *SIRINotifySubscriptionTerminatedHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRINotifySubscriptionTerminatedHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRINotifySubscriptionTerminatedHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("NotifySubscriptionTerminated %s to cancel subscription: %s", handler.xmlRequest.ResponseMessageIdentifier(), handler.xmlRequest.SubscriptionRef())

	connector.(core.SubscriptionRequestDispatcher).HandleNotifySubscriptionTerminated(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
