package core

import "fmt"

type CollectManager struct {
	partners Partners

	Events []*StopAreaUpdateEvent
}

func NewCollectManager(partners Partners) *CollectManager {
	return &CollectManager{partners: partners}
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

	manager.Events = append(manager.Events, event)
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

func (manager *CollectManager) requestStopAreaUpdate(partner *Partner, request *StopAreaUpdateRequest) (*StopAreaUpdateEvent, error) {
	fmt.Println("***************************************")
	event, err := partner.StopMonitoringRequestCollector().RequestStopAreaUpdate(request)
	if err != nil {
		return nil, err
	}
	return event, nil
}