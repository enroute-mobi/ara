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

	params.connector.(core.VehicleMonitoringSubscriptionCollector).HandleNotifyVehicleMonitoring(handler.xmlRequest)

	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = "NotifyVehicleMonitoring"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	vehicleRefs := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.VehicleMonitoringDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
		if !delivery.Status() {
			params.message.Status = "Error"
		}
		for _, vehicleActivity := range delivery.VehicleActivities() {
			vehicleRefs[vehicleActivity.XMLMonitoredVehicleJourney.VehicleRef()] = struct{}{}
		}
	}
	subs := make([]string, 0, len(subIds))
	vehicles := make([]string, 0, len(vehicleRefs))
	for k := range subIds {
		subs = append(subs, k)
	}
	for j := range vehicleRefs {
		vehicles = append(vehicles, j)
	}

	params.message.Vehicles = vehicles
	params.message.SubscriptionIdentifiers = subs
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
