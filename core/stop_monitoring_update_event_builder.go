package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	originStopAreaObjectId model.ObjectID
	partner                *Partner
	remoteObjectidKind     string

	stopMonitoringUpdateEvents *CollectUpdateEvents
}

func NewStopMonitoringUpdateEventBuilder(partner *Partner, originStopAreaObjectId model.ObjectID) StopMonitoringUpdateEventBuilder {
	return StopMonitoringUpdateEventBuilder{
		originStopAreaObjectId:     originStopAreaObjectId,
		partner:                    partner,
		remoteObjectidKind:         partner.Setting(REMOTE_OBJECTID_KIND),
		stopMonitoringUpdateEvents: NewCollectUpdateEvents(),
	}
}

func (builder *StopMonitoringUpdateEventBuilder) buildUpdateEvents(xmlStopVisitEvent *siri.XMLMonitoredStopVisit) {
	origin := string(builder.partner.Slug())

	// StopAreas
	stopAreaObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlStopVisitEvent.StopPointRef())

	_, ok := builder.stopMonitoringUpdateEvents.StopAreas[xmlStopVisitEvent.StopPointRef()]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin:   origin,
			ObjectId: stopAreaObjectId,
			Name:     xmlStopVisitEvent.StopPointName(),
		}
		if builder.originStopAreaObjectId.Value() != "" && stopAreaObjectId.String() != builder.originStopAreaObjectId.String() {
			event.ParentObjectId = builder.originStopAreaObjectId
		}

		builder.stopMonitoringUpdateEvents.StopAreas[xmlStopVisitEvent.StopPointRef()] = event
		builder.stopMonitoringUpdateEvents.MonitoringRefs[xmlStopVisitEvent.StopPointRef()] = struct{}{}
	}

	// Lines
	lineObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlStopVisitEvent.LineRef())

	_, ok = builder.stopMonitoringUpdateEvents.Lines[xmlStopVisitEvent.LineRef()]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin:   origin,
			ObjectId: lineObjectId,
			Name:     xmlStopVisitEvent.PublishedLineName(),
		}

		builder.stopMonitoringUpdateEvents.Lines[xmlStopVisitEvent.LineRef()] = lineEvent
	}

	// VehicleJourneys
	vjObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlStopVisitEvent.DatedVehicleJourneyRef())

	_, ok = builder.stopMonitoringUpdateEvents.VehicleJourneys[xmlStopVisitEvent.DatedVehicleJourneyRef()]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:          origin,
			ObjectId:        vjObjectId,
			LineObjectId:    lineObjectId,
			OriginRef:       xmlStopVisitEvent.OriginRef(),
			OriginName:      xmlStopVisitEvent.OriginName(),
			DestinationRef:  xmlStopVisitEvent.DestinationRef(),
			DestinationName: xmlStopVisitEvent.DestinationName(),
			Monitored:       xmlStopVisitEvent.Monitored(),

			ObjectidKind: builder.remoteObjectidKind,
			SiriXML:      &xmlStopVisitEvent.XMLMonitoredVehicleJourney,
		}

		builder.stopMonitoringUpdateEvents.VehicleJourneys[xmlStopVisitEvent.DatedVehicleJourneyRef()] = vjEvent
	}

	// StopVisits
	stopVisitObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlStopVisitEvent.ItemIdentifier())

	_, ok = builder.stopMonitoringUpdateEvents.StopVisits[xmlStopVisitEvent.StopPointRef()][xmlStopVisitEvent.ItemIdentifier()]
	if !ok {
		svEvent := &model.StopVisitUpdateEvent{
			Origin:                 origin,
			ObjectId:               stopVisitObjectId,
			StopAreaObjectId:       stopAreaObjectId,
			VehicleJourneyObjectId: vjObjectId,
			DataFrameRef:           xmlStopVisitEvent.DataFrameRef(),
			PassageOrder:           xmlStopVisitEvent.Order(),
			Monitored:              xmlStopVisitEvent.Monitored(),
			VehicleAtStop:          xmlStopVisitEvent.VehicleAtStop(),
			ArrivalStatus:          model.StopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
			DepartureStatus:        model.StopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
			RecordedAt:             xmlStopVisitEvent.RecordedAt(),
			Schedules:              model.NewStopVisitSchedules(),

			ObjectidKind: builder.remoteObjectidKind,
			SiriXML:      xmlStopVisitEvent,
		}

		if !xmlStopVisitEvent.AimedDepartureTime().IsZero() || !xmlStopVisitEvent.AimedArrivalTime().IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, xmlStopVisitEvent.AimedDepartureTime(), xmlStopVisitEvent.AimedArrivalTime())
		}
		if !xmlStopVisitEvent.ExpectedDepartureTime().IsZero() || !xmlStopVisitEvent.ExpectedArrivalTime().IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, xmlStopVisitEvent.ExpectedDepartureTime(), xmlStopVisitEvent.ExpectedArrivalTime())
		}
		if !xmlStopVisitEvent.ActualDepartureTime().IsZero() || !xmlStopVisitEvent.ActualArrivalTime().IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, xmlStopVisitEvent.ActualDepartureTime(), xmlStopVisitEvent.ActualArrivalTime())
		}

		if builder.stopMonitoringUpdateEvents.StopVisits[xmlStopVisitEvent.StopPointRef()] == nil {
			builder.stopMonitoringUpdateEvents.StopVisits[xmlStopVisitEvent.StopPointRef()] = make(map[string]*model.StopVisitUpdateEvent)
		}
		builder.stopMonitoringUpdateEvents.StopVisits[xmlStopVisitEvent.StopPointRef()][xmlStopVisitEvent.ItemIdentifier()] = svEvent

	}
}

func (builder *StopMonitoringUpdateEventBuilder) SetUpdateEvents(stopVisits []*siri.XMLMonitoredStopVisit) {
	for _, xmlStopVisitEvent := range stopVisits {
		builder.buildUpdateEvents(xmlStopVisitEvent)
	}
}

// Used only in StopMonitoringSubscriptionCollector
func (builder *StopMonitoringUpdateEventBuilder) SetStopVisitCancellationEvents(delivery *siri.XMLNotifyStopMonitoringDelivery) {
	for _, xmlStopVisitCancellationEvent := range delivery.XMLMonitoredStopVisitCancellations() {
		builder.stopMonitoringUpdateEvents.MonitoringRefs[xmlStopVisitCancellationEvent.MonitoringRef()] = struct{}{}

		builder.stopMonitoringUpdateEvents.Cancellations = append(builder.stopMonitoringUpdateEvents.Cancellations, model.NewNotCollectedUpdateEvent(model.NewObjectID(builder.remoteObjectidKind, xmlStopVisitCancellationEvent.ItemRef())))
	}
}

func (builder *StopMonitoringUpdateEventBuilder) UpdateEvents() CollectUpdateEvents {
	return *builder.stopMonitoringUpdateEvents
}
