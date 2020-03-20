package core

import (
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
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
		smbEventChan: make(chan model.StopMonitoringBroadcastEvent, 2000),
		gmbEventChan: make(chan model.GeneralMessageBroadcastEvent, 2000),
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
			manager.smsbEvent_handler(event)
			manager.ettsbEvent_handler(event)
		case event := <-manager.gmbEventChan:
			manager.gmsbEvent_handler(event)
		case <-manager.stop:
			logger.Log.Debugf("BroadcastManager Stop")
			return
		}
	}
}

func (manager *BroadcastManager) smsbEvent_handler(event model.StopMonitoringBroadcastEvent) {
	connectorTypes := []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER, TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.GetPartnersWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(&event)
			continue
		}

		// TEST
		connector, ok = partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) ettsbEvent_handler(event model.StopMonitoringBroadcastEvent) {
	connectorTypes := []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER, TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.GetPartnersWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).HandleBroadcastEvent(&event)
			continue
		}

		connector, ok = partner.Connector(TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestETTSubscriptionBroadcaster).HandleBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) gmsbEvent_handler(event model.GeneralMessageBroadcastEvent) {
	connectorTypes := []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER, TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.GetPartnersWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIGeneralMessageSubscriptionBroadcaster).HandleGeneralMessageBroadcastEvent(&event)
			continue
		}

		// TEST
		connector, ok = partner.Connector(TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestGeneralMessageSubscriptionBroadcaster).HandleGeneralMessageBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) Stop() {
	if manager.stop != nil {
		close(manager.stop)
	}
}
