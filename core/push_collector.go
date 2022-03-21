package core

import (
	"strconv"

	em "bitbucket.org/enroute-mobi/ara-external-models"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type PushCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	subscriber UpdateSubscriber
}

type PushCollectorFactory struct{}

func (factory *PushCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewPushCollector(partner)
}

func (factory *PushCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewPushCollector(partner *Partner) *PushCollector {
	connector := &PushCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.subscriber = manager.BroadcastUpdateEvent

	return connector
}

func (pc *PushCollector) SetSubscriber(subscriber UpdateSubscriber) {
	pc.subscriber = subscriber
}

func (pc *PushCollector) broadcastUpdateEvent(event model.UpdateEvent) {
	if pc.subscriber != nil {
		pc.subscriber(event)
	}
}

func (pc *PushCollector) HandlePushNotification(model *em.ExternalCompleteModel, message *audit.BigQueryMessage) {
	t := clock.DefaultClock().Now()

	pc.handleStopAreas(model.GetStopAreas())
	pc.handleLines(model.GetLines())
	pc.handleVehicleJourneys(model.GetVehicleJourneys())
	pc.handleStopVisits(model.GetStopVisits())
	pc.handleVehicles(model.GetVehicles())

	processingTime := clock.DefaultClock().Since(t)

	total := len(model.GetStopAreas()) + len(model.GetLines()) + len(model.GetVehicleJourneys()) + len(model.GetStopVisits())
	logger.Log.Debugf("PushCollector handled %v models in %v", total, processingTime)

	pc.partner.Pushed()

	logStashEvent := pc.newLogStashEvent()
	pc.logPushNotification(logStashEvent, model)
	audit.CurrentLogStash().WriteEvent(logStashEvent)

	message.ProcessingTime = processingTime.Seconds()
}

func (pc *PushCollector) handleStopAreas(sas []*em.ExternalStopArea) (stopAreas []string) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().RemoteObjectIDKind()

	for i := range sas {
		sa := sas[i]
		event := model.NewStopAreaUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, sa.GetObjectid())
		event.Name = sa.GetName()
		event.CollectedAlways = true
		event.Longitude = sa.GetLongitude()
		event.Latitude = sa.GetLatitude()

		stopAreas = append(stopAreas, sa.GetObjectid())

		pc.broadcastUpdateEvent(event)
	}
	return
}

func (pc *PushCollector) handleLines(lines []*em.ExternalLine) (lineIds []string) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().RemoteObjectIDKind()

	for i := range lines {
		l := lines[i]
		event := model.NewLineUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, l.GetObjectid())
		event.Name = l.GetName()

		lineIds = append(lineIds, l.GetObjectid())

		pc.broadcastUpdateEvent(event)
	}
	return
}

func (pc *PushCollector) handleVehicleJourneys(vjs []*em.ExternalVehicleJourney) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().RemoteObjectIDKind()

	for i := range vjs {
		vj := vjs[i]
		event := model.NewVehicleJourneyUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, vj.GetObjectid())
		event.LineObjectId = model.NewObjectID(id_kind, vj.GetLineRef())
		event.OriginRef = vj.GetOriginRef()
		event.OriginName = vj.GetOriginName()
		event.DestinationRef = vj.GetDestinationRef()
		event.DestinationName = vj.GetDestinationName()
		event.Direction = vj.GetDirection()

		pc.broadcastUpdateEvent(event)
	}
}

func (pc *PushCollector) handleStopVisits(svs []*em.ExternalStopVisit) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().RemoteObjectIDKind()

	for i := range svs {
		sv := svs[i]
		event := model.NewStopVisitUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, sv.GetObjectid())
		event.StopAreaObjectId = model.NewObjectID(id_kind, sv.GetStopAreaRef())
		event.VehicleJourneyObjectId = model.NewObjectID(id_kind, sv.GetVehicleJourneyRef())
		event.Monitored = sv.GetMonitored()
		event.PassageOrder = int(sv.GetPassageOrder())
		event.ArrivalStatus = model.StopVisitArrivalStatus(sv.GetArrivalStatus())
		event.DepartureStatus = model.StopVisitDepartureStatus(sv.GetDepartureStatus())

		handleSchedules(&event.Schedules, sv.GetDepartureTimes(), sv.GetArrivalTimes())

		pc.broadcastUpdateEvent(event)
	}
}

func (pc *PushCollector) handleVehicles(vs []*em.ExternalVehicle) (vehicles []string) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().RemoteObjectIDKind()

	for i := range vs {
		v := vs[i]
		event := model.NewVehicleUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, v.GetObjectid())
		event.VehicleJourneyObjectId = model.NewObjectID(id_kind, v.GetVehicleJourneyRef())
		event.Longitude = v.GetLongitude()
		event.Latitude = v.GetLatitude()
		event.Bearing = v.GetBearing()

		vehicles = append(vehicles, v.GetObjectid())

		pc.broadcastUpdateEvent(event)
	}
	return
}

func handleSchedules(sc *model.StopVisitSchedules, protoDeparture, protoArrival *em.ExternalStopVisit_Times) {
	sc.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, protoDeparture.GetAimed().AsTime(), protoArrival.GetAimed().AsTime())
	sc.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, protoDeparture.GetActual().AsTime(), protoArrival.GetActual().AsTime())
	sc.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, protoDeparture.GetExpected().AsTime(), protoArrival.GetExpected().AsTime())
}

func (pc *PushCollector) newLogStashEvent() audit.LogStashEvent {
	event := pc.partner.NewLogStashEvent()
	event["connector"] = "PushCollector"
	return event
}

func (pc *PushCollector) logPushNotification(logStashEvent audit.LogStashEvent, model *em.ExternalCompleteModel) {
	logStashEvent["type"] = "PushNotification"
	logStashEvent["stopAreas"] = strconv.Itoa(len(model.GetStopAreas()))
	logStashEvent["lines"] = strconv.Itoa(len(model.GetLines()))
	logStashEvent["vehicleJourneys"] = strconv.Itoa(len(model.GetVehicleJourneys()))
	logStashEvent["StopVisits"] = strconv.Itoa(len(model.GetStopVisits()))
}
