package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/edwig/core"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type SIRIStopMonitoringSubscriptionTerminatedResponseHandler struct {
	xmlRequest *siri.XMLStopMonitoringSubscriptionTerminatedResponse
}

func (handler *SIRIStopMonitoringSubscriptionTerminatedResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIStopMonitoringSubscriptionTerminatedResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIStopMonitoringSubscriptionTerminatedResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("StopMonitoringSubscriptionTerminated from  %s \n", handler.xmlRequest.ProducerRef())

	connector.(core.StopMonitoringSubscriptionCollector).HandleTerminatedNotification(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
