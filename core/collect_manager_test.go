package core

import (
	"testing"

	"github.com/af83/edwig/model"
)

func Test_CollectManager_BestPartner(t *testing.T) {
	partners := createTestPartnerManager()
	collectManager := NewCollectManager(partners)
	partner := partners.New("partner")
	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partners.Save(partner)

	foundPartner := collectManager.(*CollectManager).bestPartner(NewStopAreaUpdateRequest(model.StopAreaId("id")))

	if foundPartner != partner {
		t.Errorf("collectManager.bestPartner should return correct partner:\n got: %v\n want: %v", foundPartner, partner)
	}
}

func Test_CollectManager_UpdateStopArea(t *testing.T) {
	partners := createTestPartnerManager()
	collectManager := &CollectManager{
		partners:                   partners,
		stopVisitUpdateSubscribers: make([]StopVisitUpdateSubscriber, 0),
	}
	partner := partners.New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partners.Save(partner)

	testManager := &TestCollectManager{}
	collectManager.HandleStopVisitUpdateEvent(testManager.TestStopVisitUpdateSubscriber)

	if len(collectManager.stopVisitUpdateSubscribers) != 1 {
		t.Error("CollectManager should have a subscriber after HandleStopVisitUpdateEvent call")
	}

	request := &StopAreaUpdateRequest{}
	collectManager.UpdateStopArea(request)

	if len(testManager.StopVisitEvents) != 1 {
		t.Errorf("Subscriber should be called by CollectManager UpdateStopArea %v", len(testManager.Events))
	}
}
