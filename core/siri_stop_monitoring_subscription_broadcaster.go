package core

import (
	"sync"

	"github.com/af83/edwig/model"
)

type StopMonitoringSubscriptionBroadcaster interface {
	model.Stopable
	model.Startable

	handleStopMonitoringBroadcastEvent(*model.StopMonitoringBroadcastEvent)
}

type SIRIStopMonitoringSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	toBroadcast               map[SubscriptionId][]model.StopVisitId
	mutex                     *sync.Mutex //protect the map
}

type SIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

type TestSIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

type TestStopMonitoringSubscriptionBroadcaster struct {
	model.UUIDConsumer

	events                    []*model.StopMonitoringBroadcastEvent
	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
}

func NewTestStopMonitoringSubscriptionBroadcaster() *TestStopMonitoringSubscriptionBroadcaster {
	connector := &TestStopMonitoringSubscriptionBroadcaster{}
	return connector
}

func (connector *TestStopMonitoringSubscriptionBroadcaster) handleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

// END OF TEST

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringSubscriptionBroadcaster()
}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return newSIRIStopMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func newSIRIStopMonitoringSubscriptionBroadcaster(partner *Partner) *SIRIStopMonitoringSubscriptionBroadcaster {
	siriStopMonitoringSubscriptionBroadcaster := &SIRIStopMonitoringSubscriptionBroadcaster{}
	siriStopMonitoringSubscriptionBroadcaster.partner = partner
	siriStopMonitoringSubscriptionBroadcaster.mutex = &sync.Mutex{}
	siriStopMonitoringSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	siriStopMonitoringSubscriptionBroadcaster.stopMonitoringBroadcaster = NewSIRIStopMonitoringBroadcaster(siriStopMonitoringSubscriptionBroadcaster)
	siriStopMonitoringSubscriptionBroadcaster.stopMonitoringBroadcaster.Run()

	return siriStopMonitoringSubscriptionBroadcaster
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Stop() {
	if connector.stopMonitoringBroadcaster != nil {
		connector.stopMonitoringBroadcaster.Stop()
	}
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Start() {
	if connector.stopMonitoringBroadcaster == nil {
		connector.stopMonitoringBroadcaster = NewSIRIStopMonitoringBroadcaster(connector)
	}
	connector.stopMonitoringBroadcaster.Run()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) handleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	switch event.ModelType {
	case "StopVisit":
		sv, ok := connector.Partner().Model().StopVisits().Find(model.StopVisitId(event.ModelId))
		subId, ok := connector.checkEvent(sv)
		if ok {
			connector.addStopVisit(subId, sv.Id())
		}
	case "VehicleJourney":
		for _, sv := range connector.Partner().Model().StopVisits().FindFollowingByVehicleJourneyId(model.VehicleJourneyId(event.ModelId)) {
			subId, ok := connector.checkEvent(sv)
			if ok {
				connector.addStopVisit(subId, sv.Id())
			}
		}
	default:
		return
	}
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) checkEvent(sv model.StopVisit) (SubscriptionId, bool) {
	subId := SubscriptionId(0) //just to return a correct type for errors

	stopArea, ok := connector.Partner().Model().StopAreas().Find(sv.StopAreaId)
	if !ok {
		return subId, false
	}

	obj, ok := stopArea.ObjectID(connector.partner.Setting("remote_objectid_kind"))
	if !ok {
		return subId, false
	}

	sub, ok := connector.partner.Subscriptions().FindByRessourceId(obj.String())
	if !ok {
		return subId, false
	}

	subId = sub.Id()

	ressources := sub.ResourcesByObjectID()

	ressource, ok := ressources[obj.String()]

	if !ok {
		return subId, false
	}

	lastState, ok := ressource.LastStates[string(sv.Id())]

	if ok == true && !lastState.Haschanged(sv) {
		return subId, false
	}
	return subId, true
}
