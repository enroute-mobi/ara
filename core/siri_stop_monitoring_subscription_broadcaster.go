package core

import (
	"sync"

	"github.com/af83/edwig/model"
)

type StopMonitoringSubscriptionBroadcaster interface {
	model.Stopable
	model.Startable

	handleStopVisitBroadcastEvent(*model.StopVisitBroadcastEvent)
}

type SIRIStopMonitoringSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	events                    map[SubscriptionId][]*model.StopVisitBroadcastEvent
	mutex                     *sync.Mutex //protect the map
}

type SIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

type TestSIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

type TestStopMonitoringSubscriptionBroadcaster struct {
	model.UUIDConsumer

	events                    []*model.StopVisitBroadcastEvent
	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
}

func NewTestStopMonitoringSubscriptionBroadcaster() *TestStopMonitoringSubscriptionBroadcaster {
	connector := &TestStopMonitoringSubscriptionBroadcaster{}
	return connector
}

func (connector *TestStopMonitoringSubscriptionBroadcaster) handleStopVisitBroadcastEvent(event *model.StopVisitBroadcastEvent) {
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
	siriStopMonitoringSubscriptionBroadcaster.events = make(map[SubscriptionId][]*model.StopVisitBroadcastEvent)

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

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) handleStopVisitBroadcastEvent(event *model.StopVisitBroadcastEvent) {
	connector.mutex.Lock()
	connector.events[SubscriptionId(event.SubscriptionId)] = append(connector.events[SubscriptionId(event.SubscriptionId)], event)
	connector.mutex.Unlock()
}
