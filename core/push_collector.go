package core

import (
	"strconv"
	"time"

	external_models "github.com/af83/ara-external-models"
	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type PushCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	BaseConnector

	subscriber UpdateSubscriber
}

type PushCollectorFactory struct{}

func (factory *PushCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewPushCollector(partner)
}

func (factory *PushCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting(PUSH_TOKEN)
	return ok
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

func (pc *PushCollector) HandlePushNotification(model *external_models.ExternalCompleteModel) {
	t := time.Now()
	pc.handleStopAreas(model.GetStopAreas())
	pc.handleLines(model.GetLines())

	total := len(model.GetStopAreas()) + len(model.GetLines()) + len(model.GetVehicleJourneys()) + len(model.GetStopVisits())
	logger.Log.Debugf("PushCollector handled %v models in %v", total, time.Since(t))

	logStashEvent := pc.newLogStashEvent()
	pc.logPushNotification(logStashEvent, model)
	audit.CurrentLogStash().WriteEvent(logStashEvent)
}

func (pc *PushCollector) handleStopAreas(sas []*external_models.ExternalStopArea) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().Setting("remote_objectid_kind")

	for i := range sas {
		sa := sas[i]
		event := model.NewStopAreaUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, sa.GetObjectid())
		event.Name = sa.GetName()
		event.CollectedAlways = true
		event.Longitude = sa.GetLongitude()
		event.Latitude = sa.GetLatitude()

		pc.broadcastUpdateEvent(event)
	}
}

func (pc *PushCollector) handleLines(lines []*external_models.ExternalLine) {
	partner := string(pc.Partner().Slug())
	id_kind := pc.Partner().Setting("remote_objectid_kind")

	for i := range lines {
		l := lines[i]
		event := model.NewLineUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(id_kind, l.GetObjectid())
		event.Name = l.GetName()

		pc.broadcastUpdateEvent(event)
	}
}

func (pc *PushCollector) newLogStashEvent() audit.LogStashEvent {
	event := pc.partner.NewLogStashEvent()
	event["connector"] = "PushCollector"
	return event
}

func (pc *PushCollector) logPushNotification(logStashEvent audit.LogStashEvent, model *external_models.ExternalCompleteModel) {
	logStashEvent["type"] = "PushNotification"
	logStashEvent["stopAreas"] = strconv.Itoa(len(model.GetStopAreas()))
	logStashEvent["lines"] = strconv.Itoa(len(model.GetLines()))
}
