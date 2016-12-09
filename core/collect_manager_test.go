package core

import (
	"testing"

	"github.com/af83/edwig/model"
)

func Test_CollectManager_Partners(t *testing.T) {
	partners := createTestPartnerManager()
	collectManager := NewCollectManager(partners)

	if collectManager.Partners() != partners {
		t.Errorf("CollectManager Partners() should return correct value")
	}
}

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

// WIP
func Test_CollectManager_UpdateStopArea(t *testing.T) {
	// func (manager *CollectManager) UpdateStopArea(request *StopAreaUpdateRequest)
	partners := createTestPartnerManager()
	collectManager := NewCollectManager(partners)
	partner := partners.New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partners.Save(partner)

	request := &StopAreaUpdateRequest{}
	collectManager.UpdateStopArea(request)

	// Check Events
	if len(collectManager.Events()) != 1 {
		t.Error("CollectManager UpdateStopArea should generate a stopAreaUpdateEvent")
	}
}
