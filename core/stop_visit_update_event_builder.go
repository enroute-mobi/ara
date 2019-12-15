package core

import (
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type LegacyStopVisitUpdateEventBuilder struct {
	model.ClockConsumer
	model.UUIDConsumer

	originStopAreaObjectId model.ObjectID
	partner                *Partner
}

func newLegacyStopVisitUpdateEventBuilder(partner *Partner, originStopAreaObjectId model.ObjectID) LegacyStopVisitUpdateEventBuilder {
	return LegacyStopVisitUpdateEventBuilder{
		partner:                partner,
		originStopAreaObjectId: originStopAreaObjectId,
	}
}

func (builder *LegacyStopVisitUpdateEventBuilder) buildLegacyStopVisitUpdateEvent(events map[string]*model.LegacyStopAreaUpdateEvent, xmlStopVisitEvent *siri.XMLMonitoredStopVisit) {
	stopVisitEvent := &model.LegacyStopVisitUpdateEvent{
		Id:                     builder.NewUUID(),
		Origin:                 string(builder.partner.Slug()),
		DataFrameRef:           xmlStopVisitEvent.DataFrameRef(),
		Created_at:             builder.Clock().Now(),
		RecordedAt:             xmlStopVisitEvent.RecordedAt(),
		VehicleAtStop:          xmlStopVisitEvent.VehicleAtStop(),
		StopVisitObjectid:      model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.ItemIdentifier()),
		StopAreaObjectId:       model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef()),
		Schedules:              model.NewStopVisitSchedules(),
		DepartureStatus:        model.StopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
		ArrivalStatus:          model.StopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
		DatedVehicleJourneyRef: xmlStopVisitEvent.DatedVehicleJourneyRef(),
		DestinationRef:         xmlStopVisitEvent.DestinationRef(),
		OriginRef:              xmlStopVisitEvent.OriginRef(),
		DestinationName:        xmlStopVisitEvent.DestinationName(),
		OriginName:             xmlStopVisitEvent.OriginName(),
		Monitored:              xmlStopVisitEvent.Monitored(),
		Attributes:             NewSIRIStopVisitUpdateAttributes(xmlStopVisitEvent, builder.partner.Setting("remote_objectid_kind")),
	}

	if !xmlStopVisitEvent.AimedDepartureTime().IsZero() || !xmlStopVisitEvent.AimedArrivalTime().IsZero() {
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, xmlStopVisitEvent.AimedDepartureTime(), xmlStopVisitEvent.AimedArrivalTime())
	}
	if !xmlStopVisitEvent.ExpectedDepartureTime().IsZero() || !xmlStopVisitEvent.ExpectedArrivalTime().IsZero() {
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, xmlStopVisitEvent.ExpectedDepartureTime(), xmlStopVisitEvent.ExpectedArrivalTime())
	}
	if !xmlStopVisitEvent.ActualDepartureTime().IsZero() || !xmlStopVisitEvent.ActualArrivalTime().IsZero() {
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, xmlStopVisitEvent.ActualDepartureTime(), xmlStopVisitEvent.ActualArrivalTime())
	}

	stopAreaObjectidString := stopVisitEvent.StopAreaObjectId.String()
	event, ok := events[stopAreaObjectidString]
	if !ok {
		event = &model.LegacyStopAreaUpdateEvent{}
		event.Origin = string(builder.partner.Slug())
		event.StopAreaAttributes.Name = xmlStopVisitEvent.StopPointName()
		event.StopAreaAttributes.ObjectId = model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef())
		event.StopAreaAttributes.CollectedAlways = false
		events[stopVisitEvent.StopAreaObjectId.String()] = event
		if builder.originStopAreaObjectId.Value() != "" && stopAreaObjectidString != builder.originStopAreaObjectId.String() {
			event.StopAreaAttributes.ParentObjectId = builder.originStopAreaObjectId
		}
	}
	event.LegacyStopVisitUpdateEvents = append(event.LegacyStopVisitUpdateEvents, stopVisitEvent)
}

func (builder *LegacyStopVisitUpdateEventBuilder) setLegacyStopVisitUpdateEvents(events map[string]*model.LegacyStopAreaUpdateEvent, stopVisits []*siri.XMLMonitoredStopVisit) {
	for _, xmlStopVisitEvent := range stopVisits {
		builder.buildLegacyStopVisitUpdateEvent(events, xmlStopVisitEvent)
	}
}
