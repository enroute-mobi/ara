package core

import "github.com/af83/edwig/model"

type CollectManagerInterface interface {
	UpdateStopArea(request *StopAreaUpdateRequest)
	Partners() Partners
	Events() []*model.StopAreaUpdateEvent
}

type CollectManager struct {
	partners Partners

	events []*model.StopAreaUpdateEvent
}

type TestCollectManager struct {
	Done   chan bool
	events []*model.StopAreaUpdateEvent
}

func NewTestCollectManager() CollectManagerInterface {
	return &TestCollectManager{
		Done: make(chan bool, 1),
	}
}

func (manager *TestCollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	event := &model.StopAreaUpdateEvent{}
	manager.events = append(manager.events, event)

	manager.Done <- true
}

func (manager *TestCollectManager) Partners() Partners {
	return nil
}

func (manager *TestCollectManager) Events() []*model.StopAreaUpdateEvent {
	return manager.events
}

func NewCollectManager(partners Partners) CollectManagerInterface {
	return &CollectManager{partners: partners}
}

func (manager *CollectManager) Events() []*model.StopAreaUpdateEvent {
	return manager.events
}

func (manager *CollectManager) Partners() Partners {
	return manager.partners
}

func (manager *CollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	partner := manager.bestPartner(request)
	if partner == nil {
		// WIP
		return
	}

	event, err := manager.requestStopAreaUpdate(partner, request)
	if err != nil {
		// WIP: Handle error
		return
	}

	manager.events = append(manager.events, event)
}

func (manager *CollectManager) bestPartner(request *StopAreaUpdateRequest) *Partner {
	var testPartner *Partner
	for _, partner := range manager.partners.FindAll() {
		if partner.isConnectorDefined(SIRI_STOP_MONITORING_REQUEST_COLLECTOR) {
			return partner
		} else if partner.isConnectorDefined(TEST_STOP_MONITORING_REQUEST_COLLECTOR) {
			testPartner = partner
		}
	}
	if testPartner != nil {
		return testPartner
	}
	return nil
}

func (manager *CollectManager) requestStopAreaUpdate(partner *Partner, request *StopAreaUpdateRequest) (*model.StopAreaUpdateEvent, error) {
	event, err := partner.StopMonitoringRequestCollector().RequestStopAreaUpdate(request)
	if err != nil {
		return nil, err
	}
	return event, nil
}
