package api

import (
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringDeliveriesResponseHandler struct {
	xmlRequest *siri.XMLStopMonitoringResponse
}

func (handler *SIRIStopMonitoringDeliveriesResponseHandler) RequestorRef() string {
	return "None"
}

func (handler *SIRIStopMonitoringDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR
}

func (handler *SIRIStopMonitoringDeliveriesResponseHandler) Respond(connector core.Connector, rw http.ResponseWriter) {

	stopvisits := handler.xmlRequest.XMLMonitoredStopVisits()
	for _, stopvisit := range stopvisits {
		if stopvisit.StopPointRef() == "" {
			continue
		}
		stopAreaUpdateRequest := core.NewStopAreaUpdateRequest(model.StopAreaId(stopvisit.StopPointRef()))
		connector.(core.StopMonitoringSubscriptionCollector).RequestStopAreaUpdate(stopAreaUpdateRequest)
	}
}
