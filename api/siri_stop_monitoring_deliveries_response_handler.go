package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/edwig/core"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type SIRIStopMonitoringRequestDeliveriesResponseHandler struct {
	xmlRequest *siri.XMLNotifyStopMonitoring
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("NotifyStopMonitoring %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	connector.(core.StopMonitoringSubscriptionCollector).HandleNotifyStopMonitoring(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)
}
