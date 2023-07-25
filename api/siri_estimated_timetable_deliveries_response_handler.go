package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIEstimatedTimetableRequestDeliveriesResponseHandler struct {
	xmlRequest  *sxml.XMLNotifyEstimatedTimetable
	referential *core.Referential
}

func (handler *SIRIEstimatedTimetableRequestDeliveriesResponseHandler) RequestorRef() string {
	return handler.xmlRequest.ProducerRef()
}

func (handler *SIRIEstimatedTimetableRequestDeliveriesResponseHandler) ConnectorType() string {
	return core.SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_COLLECTOR
}

func (handler *SIRIEstimatedTimetableRequestDeliveriesResponseHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("NotifyEstimatedTimetable %s\n", handler.xmlRequest.ResponseMessageIdentifier())

	t := clock.DefaultClock().Now()

	updateEvents := params.connector.(core.EstimatedTimetableSubscriptionCollector).HandleNotifyEstimatedTimetable(handler.xmlRequest)
	params.rw.WriteHeader(http.StatusOK)

	params.message.Type = "NotifyEstimatedTimetable"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	params.message.RequestIdentifier = handler.xmlRequest.RequestMessageRef()
	params.message.ResponseIdentifier = handler.xmlRequest.ResponseMessageIdentifier()

	subIds := make(map[string]struct{})
	for _, delivery := range handler.xmlRequest.EstimatedTimetableDeliveries() {
		subIds[delivery.SubscriptionRef()] = struct{}{}
		if !delivery.Status() {
			params.message.Status = "Error"
		}
	}
	subs := make([]string, 0, len(subIds))
	for k := range subIds {
		subs = append(subs, k)
	}
	params.message.SubscriptionIdentifiers = subs
	params.message.Lines = updateEvents.GetLines()
	params.message.VehicleJourneys = updateEvents.GetVehicleJourneys()
	params.message.StopAreas = updateEvents.GetStopAreas()

	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
