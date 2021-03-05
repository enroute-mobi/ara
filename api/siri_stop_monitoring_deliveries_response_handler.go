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

type SIRIStopMonitoringRequestDeliveriesResponseHandler struct {
	xmlRequest *siri.XMLNotifyStopMonitoring
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter, message *audit.BigQueryMessage) {
	logger.Log.Debugf("NotifyStopMonitoring %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	connector.(core.StopMonitoringSubscriptionCollector).HandleNotifyStopMonitoring(handler.xmlRequest)

	rw.WriteHeader(http.StatusOK)

	message.Type = "NotifyStopMonitoring"
	message.RequestRawMessage = handler.xmlRequest.RawXML()
	message.ProcessingTime = time.Since(t).Seconds()
	message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	mRefs := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.StopMonitoringDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
		if !delivery.Status() {
			message.Status = "Error"
		}
		mRefs[delivery.MonitoringRef()] = struct{}{}
	}
	subs := make([]string, 0, len(subIds))
	sas := make([]string, 0, len(mRefs))
	for k := range subIds {
		subs = append(subs, k)
	}
	for k := range mRefs {
		sas = append(sas, k)
	}
	message.SubscriptionIdentifiers = subs
	message.StopAreas = sas
	audit.CurrentBigQuery().WriteEvent(message)
}
