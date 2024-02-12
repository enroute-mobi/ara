package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	originStopAreaCode model.Code
	partner            *Partner
	remoteCodeSpace    string

	stopMonitoringUpdateEvents *CollectUpdateEvents
}

func NewStopMonitoringUpdateEventBuilder(partner *Partner, originStopAreaCode model.Code) StopMonitoringUpdateEventBuilder {
	return StopMonitoringUpdateEventBuilder{
		originStopAreaCode:         originStopAreaCode,
		partner:                    partner,
		remoteCodeSpace:            partner.RemoteCodeSpace(),
		stopMonitoringUpdateEvents: NewCollectUpdateEvents(),
	}
}

func (builder *StopMonitoringUpdateEventBuilder) buildUpdateEvents(xmlStopVisitEvent *sxml.XMLMonitoredStopVisit) {
	origin := string(builder.partner.Slug())

	// StopAreas
	stopAreaCode := model.NewCode(builder.remoteCodeSpace, xmlStopVisitEvent.StopPointRef())

	_, ok := builder.stopMonitoringUpdateEvents.StopAreas[xmlStopVisitEvent.StopPointRef()]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin: origin,
			Code:   stopAreaCode,
			Name:   xmlStopVisitEvent.StopPointName(),
		}
		if builder.originStopAreaCode.Value() != "" && stopAreaCode.String() != builder.originStopAreaCode.String() {
			event.ParentCode = builder.originStopAreaCode
		}

		builder.stopMonitoringUpdateEvents.StopAreas[xmlStopVisitEvent.StopPointRef()] = event
		builder.stopMonitoringUpdateEvents.MonitoringRefs[xmlStopVisitEvent.StopPointRef()] = struct{}{}
	}

	// Lines
	lineCode := model.NewCode(builder.remoteCodeSpace, xmlStopVisitEvent.LineRef())

	_, ok = builder.stopMonitoringUpdateEvents.Lines[xmlStopVisitEvent.LineRef()]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin: origin,
			Code:   lineCode,
			Name:   xmlStopVisitEvent.PublishedLineName(),
		}

		builder.stopMonitoringUpdateEvents.Lines[xmlStopVisitEvent.LineRef()] = lineEvent
		builder.stopMonitoringUpdateEvents.LineRefs[xmlStopVisitEvent.LineRef()] = struct{}{}
	}

	// VehicleJourneys
	vjCode := model.NewCode(builder.remoteCodeSpace, xmlStopVisitEvent.DatedVehicleJourneyRef())

	_, ok = builder.stopMonitoringUpdateEvents.VehicleJourneys[xmlStopVisitEvent.DatedVehicleJourneyRef()]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:          origin,
			Code:            vjCode,
			LineCode:        lineCode,
			OriginRef:       xmlStopVisitEvent.OriginRef(),
			OriginName:      xmlStopVisitEvent.OriginName(),
			DirectionType:   builder.directionRef(xmlStopVisitEvent.DirectionRef()),
			DestinationRef:  xmlStopVisitEvent.DestinationRef(),
			DestinationName: xmlStopVisitEvent.DestinationName(),
			Monitored:       xmlStopVisitEvent.Monitored(),
			Occupancy:       model.NormalizedOccupancyName(xmlStopVisitEvent.Occupancy()),

			CodeSpace: builder.remoteCodeSpace,
			SiriXML:   &xmlStopVisitEvent.XMLMonitoredVehicleJourney,
		}

		builder.stopMonitoringUpdateEvents.VehicleJourneys[xmlStopVisitEvent.DatedVehicleJourneyRef()] = vjEvent
		builder.stopMonitoringUpdateEvents.VehicleJourneyRefs[xmlStopVisitEvent.DatedVehicleJourneyRef()] = struct{}{}
	}

	// StopVisits
	stopVisitCode := model.NewCode(builder.remoteCodeSpace, xmlStopVisitEvent.ItemIdentifier())

	_, ok = builder.stopMonitoringUpdateEvents.StopVisits[xmlStopVisitEvent.StopPointRef()][xmlStopVisitEvent.ItemIdentifier()]
	if !ok {
		svEvent := &model.StopVisitUpdateEvent{
			Origin:             origin,
			Code:               stopVisitCode,
			StopAreaCode:       stopAreaCode,
			VehicleJourneyCode: vjCode,
			DataFrameRef:       xmlStopVisitEvent.DataFrameRef(),
			PassageOrder:       xmlStopVisitEvent.Order(),
			Monitored:          xmlStopVisitEvent.Monitored(),
			VehicleAtStop:      xmlStopVisitEvent.VehicleAtStop(),
			ArrivalStatus:      model.SetStopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
			DepartureStatus:    model.SetStopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
			RecordedAt:         xmlStopVisitEvent.RecordedAt(),
			Schedules:          model.NewStopVisitSchedules(),

			CodeSpace: builder.remoteCodeSpace,
			SiriXML:   xmlStopVisitEvent,
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

func (builder *StopMonitoringUpdateEventBuilder) directionRef(direction string) (dir string) {
	in, out, err := builder.partner.PartnerSettings.SIRIDirectionType()
	if err {
		return direction
	}

	switch direction {
	case in:
		dir = model.VEHICLE_DIRECTION_INBOUND
	case out:
		dir = model.VEHICLE_DIRECTION_OUTBOUND
	default:
		dir = direction
	}

	return dir
}

func (builder *StopMonitoringUpdateEventBuilder) SetUpdateEvents(stopVisits []*sxml.XMLMonitoredStopVisit) {
	for _, xmlStopVisitEvent := range stopVisits {
		builder.buildUpdateEvents(xmlStopVisitEvent)
	}
}

// Used only in StopMonitoringSubscriptionCollector
func (builder *StopMonitoringUpdateEventBuilder) SetStopVisitCancellationEvents(delivery *sxml.XMLNotifyStopMonitoringDelivery) {
	for _, xmlStopVisitCancellationEvent := range delivery.XMLMonitoredStopVisitCancellations() {
		builder.stopMonitoringUpdateEvents.MonitoringRefs[xmlStopVisitCancellationEvent.MonitoringRef()] = struct{}{}

		code := model.NewCode(builder.remoteCodeSpace, xmlStopVisitCancellationEvent.ItemRef())

		var recordedAt time.Time
		if t := xmlStopVisitCancellationEvent.RecordedAt(); !t.IsZero() {
			recordedAt = t
		} else {
			recordedAt = delivery.ResponseTimestamp()
		}
		builder.stopMonitoringUpdateEvents.Cancellations = append(builder.stopMonitoringUpdateEvents.Cancellations, model.NewNotCollectedUpdateEvent(code, recordedAt))
	}
}

func (builder *StopMonitoringUpdateEventBuilder) UpdateEvents() CollectUpdateEvents {
	return *builder.stopMonitoringUpdateEvents
}
