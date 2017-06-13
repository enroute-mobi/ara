package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopVisitUpdateEventBuilder struct {
	model.ClockConsumer
	model.UUIDConsumer

	partner *Partner
}

func newStopVisitUpdateEventBuilder(partner *Partner) StopVisitUpdateEventBuilder {
	return StopVisitUpdateEventBuilder{partner: partner}
}

func (builder *StopVisitUpdateEventBuilder) buildStopVisitUpdateEvent(event *model.StopAreaUpdateEvent, xmlStopVisitEvent *siri.XMLMonitoredStopVisit) {
	stopVisitEvent := &model.StopVisitUpdateEvent{
		Id:                     builder.NewUUID(),
		Created_at:             builder.Clock().Now(),
		RecordedAt:             xmlStopVisitEvent.RecordedAt(),
		VehicleAtStop:          xmlStopVisitEvent.VehicleAtStop(),
		StopVisitObjectid:      model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.ItemIdentifier()),
		StopAreaObjectId:       model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef()),
		Schedules:              make(model.StopVisitSchedules),
		DepartureStatus:        model.StopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
		ArrivalStatuts:         model.StopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
		DatedVehicleJourneyRef: xmlStopVisitEvent.DatedVehicleJourneyRef(),
		DestinationRef:         xmlStopVisitEvent.DestinationRef(),
		OriginRef:              xmlStopVisitEvent.OriginRef(),
		DestinationName:        xmlStopVisitEvent.DestinationName(),
		OriginName:             xmlStopVisitEvent.OriginName(),
		Attributes:             NewSIRIStopVisitUpdateAttributes(xmlStopVisitEvent, builder.partner.Setting("remote_objectid_kind")),
	}
	stopVisitEvent.Schedules = model.NewStopVisitSchedules()
	if !xmlStopVisitEvent.AimedDepartureTime().IsZero() || !xmlStopVisitEvent.AimedArrivalTime().IsZero() {
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, xmlStopVisitEvent.AimedDepartureTime(), xmlStopVisitEvent.AimedArrivalTime())
	}
	if !xmlStopVisitEvent.ExpectedDepartureTime().IsZero() || !xmlStopVisitEvent.ExpectedArrivalTime().IsZero() {
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, xmlStopVisitEvent.ExpectedDepartureTime(), xmlStopVisitEvent.ExpectedArrivalTime())
	}
	if !xmlStopVisitEvent.ActualDepartureTime().IsZero() || !xmlStopVisitEvent.ActualArrivalTime().IsZero() {
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, xmlStopVisitEvent.ActualDepartureTime(), xmlStopVisitEvent.ActualArrivalTime())
	}
	event.StopVisitUpdateEvents = append(event.StopVisitUpdateEvents, stopVisitEvent)
}

func (builder *StopVisitUpdateEventBuilder) setStopVisitUpdateEvents(event *model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringResponse) {
	xmlStopVisitEvents := xmlResponse.XMLMonitoredStopVisits()
	if len(xmlStopVisitEvents) == 0 {
		return
	}

	for _, xmlStopVisitEvent := range xmlStopVisitEvents {
		builder.buildStopVisitUpdateEvent(event, xmlStopVisitEvent)
	}
}
