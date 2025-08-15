package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIFacilityMonitoringRequestDeliveriesResponseHandler struct {
	xmlRequest  *sxml.XMLNotifyFacilityMonitoring
	referential *core.Referential
}

func (handler *SIRIFacilityMonitoringRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIFacilityMonitoringRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_FACILITY_MONITORING_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIFacilityMonitoringRequestDeliveriesResponseHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifyFacilityMonitoring %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	params.connector.(core.FacilityMonitoringSubscriptionCollector).HandleNotifyFacilityMonitoring(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = audit.NOTIFY_FACILITY_MONITORING
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()
	if !handler.xmlRequest.Status() {
		params.message.Status = "Error"
	}

	subIds := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.FacilityMonitoringDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
	}
	subs := make([]string, 0, len(subIds))
	for k := range subIds {
		subs = append(subs, k)
	}
	params.message.SubscriptionIdentifiers = subs
	// params.message.Facilitys = collectedRefs.GetFacilitys()

	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
