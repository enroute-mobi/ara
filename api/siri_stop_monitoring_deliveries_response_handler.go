package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIStopMonitoringRequestDeliveriesResponseHandler struct {
	xmlRequest  *siri.XMLNotifyStopMonitoring
	referential *core.Referential
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIStopMonitoringRequestDeliveriesResponseHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifyStopMonitoring %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	params.connector.(core.StopMonitoringSubscriptionCollector).HandleNotifyStopMonitoring(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = "NotifyStopMonitoring"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	mRefs := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.StopMonitoringDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
		if !delivery.Status() {
			params.message.Status = "Error"
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
	params.message.SubscriptionIdentifiers = subs
	params.message.StopAreas = sas
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
