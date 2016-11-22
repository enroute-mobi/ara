package core

import (
	"testing"

	"github.com/af83/edwig/model"
)

func Test_CollectManager_Partners(t *testing.T) {
	partners := NewPartnerManager(model.NewMemoryModel())
	collectManager := NewCollectManager(partners)

	if collectManager.Partners() != partners {
		t.Errorf("CollectManager Partners() should return correct value")
	}
}

func Test_CollectManager_BestPartner(t *testing.T) {
	partners := NewPartnerManager(model.NewMemoryModel())
	collectManager := NewCollectManager(partners)
	partner := partners.New("partner")
	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partners.Save(partner)

	foundPartner := collectManager.bestPartner(NewStopAreaUpdateRequest(model.StopAreaId("id")))

	if foundPartner != partner {
		t.Errorf("collectManager.bestPartner should return correct partner:\n got: %v\n want: %v", foundPartner, partner)
	}
}
