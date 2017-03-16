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
	partner.operationnalStatus = OPERATIONNAL_STATUS_UP
	partners.Save(partner)

	foundPartner := collectManager.(*CollectManager).bestPartner(NewStopAreaUpdateRequest(model.StopAreaId("id")))

	if foundPartner != partner {
		t.Errorf("collectManager.bestPartner should return correct partner:\n got: %v\n want: %v", foundPartner, partner)
	}
}

func Test_CollectManager_UpdateStopArea(t *testing.T) {
	partners := createTestPartnerManager()
	collectManager := &CollectManager{
		partners:                  partners,
		StopAreaUpdateSubscribers: make([]StopAreaUpdateSubscriber, 0),
	}
	partner := partners.New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partners.Save(partner)

	testManager := &TestCollectManager{}
	collectManager.HandleStopVisitUpdateEvent(testManager.TestStopAreaUpdateSubscriber)

	if len(collectManager.StopAreaUpdateSubscribers) != 1 {
		t.Error("CollectManager should have a subscriber after HandleStopVisitUpdateEvent call")
	}

	request := &StopAreaUpdateRequest{}
	collectManager.UpdateStopArea(request)

	if len(testManager.StopVisitEvents) != 1 {
		t.Errorf("Subscriber should be called by CollectManager UpdateStopArea %v", len(testManager.Events))
	}
}

func Test_CollectManager_StopVisitUpdate(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)
	referentials.Save(referential)

	stopVisit := referential.Model().StopVisits().New()
	objectid := model.NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
	stopVisit.Save()

	stopVisitUpdateEvent := &model.StopVisitUpdateEvent{
		StopVisitObjectid: objectid,
		DepartureStatus:   model.STOP_VISIT_DEPARTURE_ONTIME,
		ArrivalStatuts:    model.STOP_VISIT_ARRIVAL_ARRIVED,
	}
	referential.collectManager.(*CollectManager).broadcastStopVisitUpdateEvent(stopVisitUpdateEvent)

	updatedStopVisit, _ := referential.Model().StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.ArrivalStatus != model.STOP_VISIT_ARRIVAL_ARRIVED {
		t.Errorf("Wrong ArrivalStatus stopVisit should have been updated\n expected: %v\n got: %v", model.STOP_VISIT_ARRIVAL_ARRIVED, updatedStopVisit.ArrivalStatus)
	}
}
