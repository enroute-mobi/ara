package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type BrocasterManagerInterface interface {
	model.Stopable

	Run()
	StopVisitBroadcastEvent() chan model.StopVisitBroadcastEvent
}

type BrocasterManager struct {
	Referential             *Referential
	stopVisitBroadcastEvent chan model.StopVisitBroadcastEvent
	stop                    chan struct{}
}

func NewBroadcasterManager(referential *Referential) *BrocasterManager {
	return &BrocasterManager{
		Referential:             referential,
		stopVisitBroadcastEvent: make(chan model.StopVisitBroadcastEvent, 0),
	}
}

func (manager *BrocasterManager) StopVisitBroadcastEvent() chan model.StopVisitBroadcastEvent {
	return manager.stopVisitBroadcastEvent
}

func (manager *BrocasterManager) GetPartnersInterrestedByStopVisitBroadcastEvent(event *model.StopVisitBroadcastEvent) []*Partner {
	partners := []*Partner{}

	for _, partner := range manager.Referential.Partners().FindAll() {
		_, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		_, testConnectorPresent := partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

		if !ok && !testConnectorPresent {
			continue
		}

		stopArea, ok := manager.Referential.Model().StopAreas().Find(event.StopAreaId)
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

		lastState, ok := ressource.LastStates[string(event.Id)].(*stopVisitLastChange)

		if ok == true && !lastState.Ischanged(event) {
			continue
		}
		partners = append(partners, partner)
	}
	return partners
}

func (manager *BrocasterManager) Run() {
	logger.Log.Debugf("BroadcasterManager start")

	go manager.run()
}

func (manager *BrocasterManager) run() {
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

func (manager *BrocasterManager) Stop() {
	close(manager.stop)
}
