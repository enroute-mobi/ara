package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type BroadcastManagerInterface interface {
	model.Stopable

	Run()
	GetStopMonitoringBroadcastEventChan() chan model.StopMonitoringBroadcastEvent
}

type BroadcastManager struct {
	Referential  *Referential
	smbEventChan chan model.StopMonitoringBroadcastEvent
	stop         chan struct{}
}

func NewBroadcastManager(referential *Referential) *BroadcastManager {
	return &BroadcastManager{
		Referential:  referential,
		smbEventChan: make(chan model.StopMonitoringBroadcastEvent, 0),
	}
}

func (manager *BroadcastManager) GetStopMonitoringBroadcastEventChan() chan model.StopMonitoringBroadcastEvent {
	return manager.smbEventChan
}

func (manager *BroadcastManager) GetPartnersWithConnector(event *model.StopMonitoringBroadcastEvent) []*Partner {
	partners := []*Partner{}

	for _, partner := range manager.Referential.Partners().FindAll() {
		_, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		_, testConnectorPresent := partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

		if !ok && !testConnectorPresent {
			continue
		}
		partners = append(partners, partner)
	}
	return partners
}

func (manager *BroadcastManager) Run() {
	logger.Log.Debugf("BroadcastManager start")

	go manager.run()
}

func (manager *BroadcastManager) run() {
	for {
		select {
		case event := <-manager.smbEventChan:
			for _, partner := range manager.GetPartnersWithConnector(&event) {
				connector, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
				if ok {
					connector.(*SIRIStopMonitoringSubscriptionBroadcaster).handleStopMonitoringBroadcastEvent(&event)
					continue
				}

				// TEST
				connector, ok = partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
				if ok {
					connector.(*TestStopMonitoringSubscriptionBroadcaster).handleStopMonitoringBroadcastEvent(&event)
					continue
				}
			}
		case <-manager.stop:
			logger.Log.Debugf("BroadcastManager Stop")
			return
		}
	}
}

func (manager *BroadcastManager) Stop() {
	close(manager.stop)
}
