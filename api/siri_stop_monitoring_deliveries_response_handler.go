package api

import (
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringDeliveriesResponseHandler struct {
	xmlRequest *siri.XMLNotifyStopMonitoring
	Partner    core.Partner
}

func (handler *SIRIStopMonitoringDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIStopMonitoringDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR
}

func (handler *SIRIStopMonitoringDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("NotifyStopMonitoring %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	connector.(core.StopMonitoringSubscriptionCollector).HandleNotifyStopMonitoring(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
