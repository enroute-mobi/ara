package core

import (
	"time"

	"github.com/af83/edwig/model"
)

type StopMonitoringSubscriptionCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	CancelStopVisitMonitoring(cancelledMap map[string][]string)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	Partner                   Partner
	StopAreaUpdateSubscribers []StopAreaUpdateSubscriber
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	siriStopMonitoringSubscriptionCollector := &SIRIStopMonitoringSubscriptionCollector{}
	siriStopMonitoringSubscriptionCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriStopMonitoringSubscriptionCollector.StopAreaUpdateSubscribers = manager.GetStopAreaUpdateSubscribers()

	return siriStopMonitoringSubscriptionCollector
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	for _, sr := range subscription.resourcesByObjectID {
		if sr.Reference.ObjectId.Value() == string(request.StopAreaId()) {
			sr.SubscribedUntil = sr.SubscribedUntil.Add(1 * time.Minute)
			return
		}
	}

	objId := model.NewObjectID("StopMonitoring", string(request.StopAreaId()))
	ref := model.Reference{
		ObjectId: &objId,
		Id:       string(request.StopAreaId()),
		Type:     "StopArea",
	}

	subscription.CreateAddNewResource(ref)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) SetStopAreaUpdateSubscriber(stopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	connector.StopAreaUpdateSubscribers = append(connector.StopAreaUpdateSubscribers, stopAreaUpdateSubscriber)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) broadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent) {
	for _, StopAreaUpdateSubscriber := range connector.StopAreaUpdateSubscribers {
		StopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) CancelStopVisitMonitoring(cancelledMap map[string][]string) {
	for key, stopVisitIds := range cancelledMap {
		stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID(), model.StopAreaId(key))
		for _, stopVisitId := range stopVisitIds {
			NotCollectedEvent := model.StopVisitNotCollectedEvent{StopVisitObjectId: model.NewObjectID("StopMonitoring", stopVisitId)}
			stopAreaUpdateEvent.StopVisitNotCollectedEvents = append(stopAreaUpdateEvent.StopVisitNotCollectedEvents, &NotCollectedEvent)
		}
		connector.broadcastStopAreaUpdateEvent(stopAreaUpdateEvent)
	}
}
