package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIVehicleMonitoringRequestDeliveriesResponseHandler struct {
	xmlRequest  *sxml.XMLNotifyVehicleMonitoring
	referential *core.Referential
}

func (handler *SIRIVehicleMonitoringRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIVehicleMonitoringRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_VEHICLE_MONITORING_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIVehicleMonitoringRequestDeliveriesResponseHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifyVehicleMonitoring %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	collectedRefs := params.connector.(core.VehicleMonitoringSubscriptionCollector).HandleNotifyVehicleMonitoring(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = audit.NOTIFY_VEHICLE_MONITORING
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()
	if !handler.xmlRequest.Status() {
		params.message.Status = "Error"
	}

	subIds := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.VehicleMonitoringDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
	}
	subs := make([]string, 0, len(subIds))
	for k := range subIds {
		subs = append(subs, k)
	}
	params.message.SubscriptionIdentifiers = subs
	params.message.Vehicles = collectedRefs.GetVehicles()
	params.message.Lines = collectedRefs.GetLines()
	params.message.VehicleJourneys = collectedRefs.GetVehicleJourneys()
	params.message.StopAreas = collectedRefs.GetStopAreas()

	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
