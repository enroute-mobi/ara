package core

import (
	"time"

	em "bitbucket.org/enroute-mobi/ara-external-models"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
)

type PushCollector struct {
	connector

	vjEvents        map[string]model.UpdateEvent
	svEvents        []model.UpdateEvent
	vjOfIgnoredSv   map[string]struct{}
	vjWithStopVisit map[string]struct{}
	persistence     time.Duration
	subscriber      UpdateSubscriber
}

type PushCollectorFactory struct{}

func (factory *PushCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewPushCollector(partner)
}

func (factory *PushCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewPushCollector(partner *Partner) *PushCollector {
	connector := &PushCollector{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace()
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
	pc.persistence = pc.partner.Referential().ModelPersistenceDuration()
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
		event.Code = model.NewCode(pc.remoteCodeSpace, sa.GetObjectid())
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
		event.Code = model.NewCode(pc.remoteCodeSpace, l.GetObjectid())
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
		event.Code = model.NewCode(pc.remoteCodeSpace, vj.GetObjectid())
		event.LineCode = model.NewCode(pc.remoteCodeSpace, vj.GetLineRef())
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
		if event.Schedules.ReferenceTime().Before(pc.Clock().Now().Add(pc.persistence)) {
			pc.vjOfIgnoredSv[sv.GetVehicleJourneyRef()] = struct{}{} // Save vehicle journeys for which we didn't save at least 1 stop visit
			continue
		}

		event.Origin = partner
		event.Code = model.NewCode(pc.remoteCodeSpace, sv.GetObjectid())
		event.StopAreaCode = model.NewCode(pc.remoteCodeSpace, sv.GetStopAreaRef())
		event.VehicleJourneyCode = model.NewCode(pc.remoteCodeSpace, sv.GetVehicleJourneyRef())
		event.Monitored = sv.GetMonitored()
		event.PassageOrder = int(sv.GetPassageOrder())
		event.ArrivalStatus = model.SetStopVisitArrivalStatus(sv.GetArrivalStatus())
		event.DepartureStatus = model.SetStopVisitDepartureStatus(sv.GetDepartureStatus())

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
		event.Code = model.NewCode(pc.remoteCodeSpace, v.GetObjectid())
		event.VehicleJourneyCode = model.NewCode(pc.remoteCodeSpace, v.GetVehicleJourneyRef())
		event.Longitude = v.GetLongitude()
		event.Latitude = v.GetLatitude()
		event.Bearing = v.GetBearing()
		event.RecordedAt = v.GetRecordedAt().AsTime()

		vehicles = append(vehicles, v.GetObjectid())

		pc.broadcastUpdateEvent(event)
	}
	return
}

func handleSchedules(sc *schedules.StopVisitSchedules, protoDeparture, protoArrival *em.ExternalStopVisit_Times) {
	sc.SetSchedule(schedules.Aimed, protoDeparture.GetAimed().AsTime(), protoArrival.GetAimed().AsTime())
	sc.SetSchedule(schedules.Actual, protoDeparture.GetActual().AsTime(), protoArrival.GetActual().AsTime())
	sc.SetSchedule(schedules.Expected, protoDeparture.GetExpected().AsTime(), protoArrival.GetExpected().AsTime())
}
