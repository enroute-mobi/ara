package core

import (
	"sync"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type EstimatedTimeTableSubscriptionBroadcaster interface {
	HandleStopMonitoringBroadcastEvent(*model.StopMonitoringBroadcastEvent)
	HandleSubscriptionRequest([]*siri.XMLStopMonitoringSubscriptionRequestEntry)
}

type SIRIEstimatedTimeTableSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	//stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	toBroadcast map[SubscriptionId][]model.StopVisitId
	mutex       *sync.Mutex //protect the map
}

type SIRIEstimatedTimetableSubscriptionBroadcasterFactory struct{}

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIStopMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func newSIRIEstimatedTimeTableSubscriptionBroadcaster(partner *Partner) *SIRIEstimatedTimeTableSubscriptionBroadcaster {
	siriEstimatedTimeTableSubscriptionBroadcaster := &SIRIEstimatedTimeTableSubscriptionBroadcaster{}
	siriEstimatedTimeTableSubscriptionBroadcaster.partner = partner
	siriEstimatedTimeTableSubscriptionBroadcaster.mutex = &sync.Mutex{}
	siriEstimatedTimeTableSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	return siriEstimatedTimeTableSubscriptionBroadcaster
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkEvent(sv model.StopVisit, tx *model.Transaction) (SubscriptionId, bool) {
	return SubscriptionId(0), true
}
