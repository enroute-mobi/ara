package core

import (
	"time"

	em "bitbucket.org/enroute-mobi/ara-external-models"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type PushCollector struct {
	connector

	vjEvents        map[string]model.UpdateEvent
	svEvents        []model.UpdateEvent
	vjOfIgnoredSv   map[string]struct{}
	vjWithStopVisit map[string]struct{}
	persistence     persistence
	subscriber      UpdateSubscriber
}

type persistence struct {
	p  time.Duration
	ok bool
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
	connector.remoteObjectidKind = partner.RemoteObjectIDKind()
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

func (pc *PushCollector) refresh() {
	pc.persistence.p, pc.persistence.ok = pc.partner.Referential().ModelPersistenceDuration()
	pc.vjEvents = make(map[string]model.UpdateEvent)
	pc.svEvents = []model.UpdateEvent{}
	pc.vjOfIgnoredSv = make(map[string]struct{})
	pc.vjWithStopVisit = make(map[string]struct{})
}

func (pc *PushCollector) HandlePushNotification(model *em.ExternalCompleteModel, message *audit.BigQueryMessage) {
	t := clock.DefaultClock().Now()
	pc.refresh()

	pc.handleStopAreas(model.GetStopAreas())
	pc.handleLines(model.GetLines())
	pc.handleVehicleJourneysAndStopVisits(model.GetVehicleJourneys(), model.GetStopVisits())
	pc.handleVehicles(model.GetVehicles())
	processingTime := clock.DefaultClock().Since(t)

	total := len(model.GetStopAreas()) + len(model.GetLines()) + len(model.GetVehicleJourneys()) + len(model.GetStopVisits())
	logger.Log.Debugf("PushCollector handled %v models in %v", total, processingTime)

	pc.partner.Pushed()

	message.ProcessingTime = processingTime.Seconds()
}

func (pc *PushCollector) handleStopAreas(sas []*em.ExternalStopArea) (stopAreas []string) {
	partner := string(pc.Partner().Slug())

	for i := range sas {
		sa := sas[i]
		event := model.NewStopAreaUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(pc.remoteObjectidKind, sa.GetObjectid())
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

	for i := range lines {
		l := lines[i]
		event := model.NewLineUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(pc.remoteObjectidKind, l.GetObjectid())
		event.Name = l.GetName()

		lineIds = append(lineIds, l.GetObjectid())

		pc.broadcastUpdateEvent(event)
	}
	return
}

func (pc *PushCollector) handleVehicleJourneysAndStopVisits(vjs []*em.ExternalVehicleJourney, svs []*em.ExternalStopVisit) {
	pc.handleVehicleJourneys(vjs)
	pc.handleStopVisits(svs)

	// For each vehicle journey for which we didn't save a stopvisit, check if we did save at least one
	// If not, we don't save them at all
	for id := range pc.vjOfIgnoredSv {
		if _, ok := pc.vjWithStopVisit[id]; !ok {
			delete(pc.vjEvents, id)
		}
	}

	for k := range pc.vjEvents {
		pc.broadcastUpdateEvent(pc.vjEvents[k])
	}
	for i := range pc.svEvents {
		pc.broadcastUpdateEvent(pc.svEvents[i])
	}
}

func (pc *PushCollector) handleVehicleJourneys(vjs []*em.ExternalVehicleJourney) {
	partner := string(pc.Partner().Slug())

	for i := range vjs {
		vj := vjs[i]
		event := model.NewVehicleJourneyUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(pc.remoteObjectidKind, vj.GetObjectid())
		event.LineObjectId = model.NewObjectID(pc.remoteObjectidKind, vj.GetLineRef())
		event.OriginRef = vj.GetOriginRef()
		event.OriginName = vj.GetOriginName()
		event.DestinationRef = vj.GetDestinationRef()
		event.DestinationName = vj.GetDestinationName()
		event.Direction = vj.GetDirection()

		pc.vjEvents[vj.GetObjectid()] = event
	}
}

func (pc *PushCollector) handleStopVisits(svs []*em.ExternalStopVisit) {
	partner := string(pc.Partner().Slug())

	for i := range svs {
		sv := svs[i]
		event := model.NewStopVisitUpdateEvent()

		handleSchedules(event.Schedules, sv.GetDepartureTimes(), sv.GetArrivalTimes())
		if pc.persistence.ok && event.Schedules.ReferenceTime().Before(pc.Clock().Now().Add(pc.persistence.p)) {
			pc.vjOfIgnoredSv[sv.GetVehicleJourneyRef()] = struct{}{} // Save vehicle journeys for which we didn't save at least 1 stop visit
			continue
		}

		event.Origin = partner
		event.ObjectId = model.NewObjectID(pc.remoteObjectidKind, sv.GetObjectid())
		event.StopAreaObjectId = model.NewObjectID(pc.remoteObjectidKind, sv.GetStopAreaRef())
		event.VehicleJourneyObjectId = model.NewObjectID(pc.remoteObjectidKind, sv.GetVehicleJourneyRef())
		event.Monitored = sv.GetMonitored()
		event.PassageOrder = int(sv.GetPassageOrder())
		event.ArrivalStatus = model.StopVisitArrivalStatus(sv.GetArrivalStatus())
		event.DepartureStatus = model.StopVisitDepartureStatus(sv.GetDepartureStatus())

		pc.vjWithStopVisit[sv.GetVehicleJourneyRef()] = struct{}{} // Save vehicle journeys for which we did save at least 1 stop visit
		pc.svEvents = append(pc.svEvents, event)
	}
}

func (pc *PushCollector) handleVehicles(vs []*em.ExternalVehicle) (vehicles []string) {
	partner := string(pc.Partner().Slug())

	for i := range vs {
		v := vs[i]
		event := model.NewVehicleUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(pc.remoteObjectidKind, v.GetObjectid())
		event.VehicleJourneyObjectId = model.NewObjectID(pc.remoteObjectidKind, v.GetVehicleJourneyRef())
		event.Longitude = v.GetLongitude()
		event.Latitude = v.GetLatitude()
		event.Bearing = v.GetBearing()
		event.RecordedAt = v.GetRecordedAt().AsTime()

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
