package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopVisitUpdateEventBuilder struct {
	model.ClockConsumer
	model.UUIDConsumer

	originStopAreaObjectId model.ObjectID
	partner                *Partner
}

func newStopVisitUpdateEventBuilder(partner *Partner, originStopAreaObjectId model.ObjectID) StopVisitUpdateEventBuilder {
	return StopVisitUpdateEventBuilder{
		partner:                partner,
		originStopAreaObjectId: originStopAreaObjectId,
	}
}

func (builder *StopVisitUpdateEventBuilder) buildStopVisitUpdateEvent(events map[string]*model.StopAreaUpdateEvent, xmlStopVisitEvent *siri.XMLMonitoredStopVisit) {
	stopVisitEvent := &model.StopVisitUpdateEvent{
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
		event = &model.StopAreaUpdateEvent{}
		event.Origin = string(builder.partner.Slug())
		event.StopAreaAttributes.Name = xmlStopVisitEvent.StopPointName()
		event.StopAreaAttributes.ObjectId = model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef())
		event.StopAreaAttributes.CollectedAlways = false
		events[stopVisitEvent.StopAreaObjectId.String()] = event
		if builder.originStopAreaObjectId.Value() != "" && stopAreaObjectidString != builder.originStopAreaObjectId.String() {
			event.StopAreaAttributes.ParentObjectId = builder.originStopAreaObjectId
		}
	}
	event.StopVisitUpdateEvents = append(event.StopVisitUpdateEvents, stopVisitEvent)
}

func (builder *StopVisitUpdateEventBuilder) setStopVisitUpdateEvents(events map[string]*model.StopAreaUpdateEvent, stopVisits []*siri.XMLMonitoredStopVisit) {
	for _, xmlStopVisitEvent := range stopVisits {
		builder.buildStopVisitUpdateEvent(events, xmlStopVisitEvent)
	}
}
