package api

import (
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringSubscriptionTerminatedResponseHandler struct {
	xmlRequest *siri.XMLStopMonitoringSubscriptionTerminatedResponse
	Partner    core.Partner
}

func (handler *SIRIStopMonitoringSubscriptionTerminatedResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIStopMonitoringSubscriptionTerminatedResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR
}

func (handler *SIRIStopMonitoringSubscriptionTerminatedResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("StopMonitoringSubscriptionTerminated from  %s \n", handler.xmlRequest.ProducerRef())

	connector.(core.StopMonitoringSubscriptionCollector).HandleTerminatedNotification(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
