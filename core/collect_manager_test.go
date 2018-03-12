package core

import (
	"testing"

	"github.com/af83/edwig/model"
)

func Test_CollectManager_BestPartner(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(model.NewObjectID("internal", "boarle"))
	stopArea.Save()

	partners := referential.Partners()
	collectManager := NewCollectManager(referential)
	partner := partners.New("partner")
	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partner.PartnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UP
	partner.Settings["collect.include_stop_areas"] = "boarle"
	partner.Settings["remote_objectid_kind"] = "internal"
	partners.Save(partner)

	foundPartner := collectManager.(*CollectManager).bestPartner(stopArea)

	if foundPartner != partner {
		t.Errorf("collectManager.bestPartner should return correct partner:\n got: %v\n want: %v", foundPartner, partner)
	}
}

// Already tested by siriStopMonitoringRequestCollectorTest
// func Test_CollectManager_UpdateStopArea(t *testing.T) {
// 	referentials := NewMemoryReferentials()
// 	referential := referentials.New("referential")
//
// 	partners := referential.Partners()
// 	collectManager := &CollectManager{
// 		referential:               referential,
// 		StopAreaUpdateSubscribers: make([]StopAreaUpdateSubscriber, 0),
// 	}
//
// 	stopArea := referential.Model().StopAreas().New()
// 	stopArea.SetObjectID(model.NewObjectID("internal", "boarle"))
// 	stopArea.Save()
//
// 	partner := partners.New("partner")
// 	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
// 	partner.RefreshConnectors()
// 	partner.PartnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UP
// 	partner.Settings["remote_objectid_kind"] = "internal"
//
// 	partners.Save(partner)
//
// 	testManager := &TestCollectManager{}
// 	collectManager.HandleStopAreaUpdateEvent(testManager.TestStopAreaUpdateSubscriber)
//
// 	if len(collectManager.StopAreaUpdateSubscribers) != 1 {
// 		t.Error("CollectManager should have a subscriber after HandleStopVisitUpdateEvent call")
// 	}
//
// 	request := NewStopAreaUpdateRequest(stopArea.Id())
// 	collectManager.UpdateStopArea(request)
//
// 	time.Sleep(50 * time.Millisecond)
// 	if len(testManager.StopVisitEvents) != 1 {
// 		t.Errorf("Subscriber should be called by CollectManager UpdateStopArea %v", len(testManager.StopVisitEvents))
// 	}
// }

func Test_CollectManager_StopVisitUpdate(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)
	referentials.Save(referential)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	stopVisit := referential.Model().StopVisits().New()
	objectid := model.NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
	stopVisit.Save()

	stopVisitUpdateEvent := &model.StopVisitUpdateEvent{
		StopVisitObjectid: objectid,
		DepartureStatus:   model.STOP_VISIT_DEPARTURE_ONTIME,
		ArrivalStatus:     model.STOP_VISIT_ARRIVAL_ARRIVED,
		Attributes:        &model.TestStopVisitUpdateAttributes{},
	}
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent("test", stopArea.Id())
	stopAreaUpdateEvent.StopVisitUpdateEvents = []*model.StopVisitUpdateEvent{stopVisitUpdateEvent}
	referential.collectManager.BroadcastStopAreaUpdateEvent(stopAreaUpdateEvent)

	updatedStopVisit, _ := referential.Model().StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.ArrivalStatus != model.STOP_VISIT_ARRIVAL_ARRIVED {
		t.Errorf("Wrong ArrivalStatus stopVisit should have been updated\n expected: %v\n got: %v", model.STOP_VISIT_ARRIVAL_ARRIVED, updatedStopVisit.ArrivalStatus)
	}
}

func Test_CollectManager_StopAreaMonitored(t *testing.T) {
	// logger.Log.Debug = true
	// defer func() { logger.Log.Debug = false }()

	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	referentials.Save(referential)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Monitored = true
	stopArea.SetObjectID(model.NewObjectID("test", "value"))
	stopArea.Save()

	stopAreaUpdateRequest := &StopAreaUpdateRequest{
		id:         StopAreaUpdateRequestId(model.DefaultUUIDGenerator().NewUUID()),
		stopAreaId: stopArea.Id(),
		createdAt:  referential.Clock().Now(),
	}
	referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)

	updatedStopArea, _ := referential.Model().StopAreas().Find(stopArea.Id())
	if updatedStopArea.Monitored {
		t.Errorf("StopArea Monitored should be false after CollectManager UpdateStopArea")
	}

	partner := referential.Partners().New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	partner.Settings["remote_objectid_kind"] = "test"
	partner.PartnerStatus = PartnerStatus{OperationnalStatus: OPERATIONNAL_STATUS_UP}
	referential.Partners().Save(partner)

	referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)
	updatedStopArea, _ = referential.Model().StopAreas().Find(stopArea.Id())
	if !updatedStopArea.Monitored {
		t.Errorf("StopArea Monitored should be true after CollectManager UpdateStopArea")
	}
}
