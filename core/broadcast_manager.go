package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type BroadcastManagerInterface interface {
	model.Startable
	model.Stopable

	GetStopMonitoringBroadcastEventChan() chan model.StopMonitoringBroadcastEvent
	GetGeneralMessageBroadcastEventChan() chan model.GeneralMessageBroadcastEvent
}

type BroadcastManager struct {
	Referential *Referential

	smbEventChan chan model.StopMonitoringBroadcastEvent
	gmbEventChan chan model.GeneralMessageBroadcastEvent
	stop         chan struct{}
}

func NewBroadcastManager(referential *Referential) *BroadcastManager {
	return &BroadcastManager{
		Referential:  referential,
		smbEventChan: make(chan model.StopMonitoringBroadcastEvent, 0),
		gmbEventChan: make(chan model.GeneralMessageBroadcastEvent, 0),
	}
}

func (manager *BroadcastManager) GetStopMonitoringBroadcastEventChan() chan model.StopMonitoringBroadcastEvent {
	return manager.smbEventChan
}

func (manager *BroadcastManager) GetGeneralMessageBroadcastEventChan() chan model.GeneralMessageBroadcastEvent {
	return manager.gmbEventChan
}

func (manager *BroadcastManager) GetPartnersWithConnector(connectorTypes []string) []*Partner {
	partners := []*Partner{}

	for _, partner := range manager.Referential.Partners().FindAll() {
		ok := false
		for _, connectorType := range connectorTypes {
			if _, present := partner.Connector(connectorType); !present {
				continue
			}
			ok = true
		}
		if !ok {
			continue
		}
		partners = append(partners, partner)
	}
	return partners
}

func (manager *BroadcastManager) Start() {
	logger.Log.Debugf("BroadcastManager start")

	manager.stop = make(chan struct{})

	go manager.run()
}

func (manager *BroadcastManager) run() {
	for {
		select {
		case event := <-manager.smbEventChan:
			connectorTypes := []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER, TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
			for _, partner := range manager.GetPartnersWithConnector(connectorTypes) {
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
		case event := <-manager.gmbEventChan:
			connectorTypes := []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER, TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
			for _, partner := range manager.GetPartnersWithConnector(connectorTypes) {
				connector, ok := partner.Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
				if ok {
					connector.(*SIRIGeneralMessageSubscriptionBroadcaster).handleGeneralMessageBroadcastEvent(&event)
					continue
				}

				// TEST
				connector, ok = partner.Connector(TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
				if ok {
					connector.(*TestGeneralMessageSubscriptionBroadcaster).handleGeneralMessageBroadcastEvent(&event)
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
	if manager.stop != nil {
		close(manager.stop)
	}
}
