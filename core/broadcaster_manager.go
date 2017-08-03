package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type BroadcasterManagerInterface interface {
	model.Stopable

	Run()
	StopVisitBroadcastEvent() chan model.StopVisitBroadcastEvent
}

type BroadcasterManager struct {
	Referential             *Referential
	stopVisitBroadcastEvent chan model.StopVisitBroadcastEvent
	stop                    chan struct{}
}

func NewBroadcasterManager(referential *Referential) *BroadcasterManager {
	return &BroadcasterManager{
		Referential:             referential,
		stopVisitBroadcastEvent: make(chan model.StopVisitBroadcastEvent, 0),
	}
}

func (manager *BroadcasterManager) StopVisitBroadcastEvent() chan model.StopVisitBroadcastEvent {
	return manager.stopVisitBroadcastEvent
}

func (manager *BroadcasterManager) GetPartnersInterrestedByStopVisitBroadcastEvent(event *model.StopVisitBroadcastEvent) []*Partner {
	partners := []*Partner{}

	for _, partner := range manager.Referential.Partners().FindAll() {
		_, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		_, testConnectorPresent := partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

		if !ok && !testConnectorPresent {
			continue
		}

		stopArea, ok := manager.GetStopAreaFromEvent(event)
		if !ok {
			continue
		}

		obj, ok := stopArea.ObjectID(partner.Setting("remote_objectid_kind"))
		if !ok {
			continue
		}

		subs, ok := partner.Subscriptions().FindByRessourceId(obj.String())
		if !ok {
			continue
		}

		event.SubscriptionId = string(subs.Id())

		ressources := subs.ResourcesByObjectID()

		ressource, ok := ressources[obj.String()]

		if !ok {
			continue
		}

		lastState, ok := ressource.LastStates[string(event.ModelId)]

		if ok == true && !lastState.Ischanged(event) {
			continue
		}
		partners = append(partners, partner)
	}
	return partners
}

func (manager *BroadcasterManager) FindStopAreaFromStopVisitId(svId model.StopVisitId) (*model.StopArea, bool) {
	sv, ok := manager.Referential.Model().StopVisits().Find(svId)
	if !ok {
		return nil, false
	}

	sa, ok := manager.Referential.Model().StopAreas().Find(sv.StopAreaId)
	if !ok {
		return nil, false
	}

	return &sa, true
}

func (manager *BroadcasterManager) GetStopAreaFromEvent(event *model.StopVisitBroadcastEvent) (*model.StopArea, bool) {

	switch event.ModelType {
	case "StopVisit":
		return manager.FindStopAreaFromStopVisitId(model.StopVisitId(event.ModelId))
	default:
		return nil, false
	}
}

func (manager *BroadcasterManager) Run() {
	logger.Log.Debugf("BroadcasterManager start")

	go manager.run()
}

func (manager *BroadcasterManager) run() {
	for {
		event := <-manager.stopVisitBroadcastEvent
		for _, partner := range manager.GetPartnersInterrestedByStopVisitBroadcastEvent(&event) {
			connector, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
			if ok {
				connector.(*SIRIStopMonitoringSubscriptionBroadcaster).handleStopVisitBroadcastEvent(&event)
				continue
			}

			// TEST
			connector, ok = partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
			if ok {
				connector.(*TestStopMonitoringSubscriptionBroadcaster).handleStopVisitBroadcastEvent(&event)
				continue
			}
		}
	}
}

func (manager *BroadcasterManager) Stop() {
	close(manager.stop)
}
