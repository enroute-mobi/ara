package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
)

func Test_CollectManager_StopVisitUpdate(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_REQUEST_COLLECTOR}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)
	referentials.Save(referential)

	stopArea := referential.Model().StopAreas().New()
	saCode := model.NewCode("codeSpace", "saValue")
	stopArea.SetCode(saCode)
	stopArea.Save()

	vj := referential.Model().VehicleJourneys().New()
	vjCode := model.NewCode("codeSpace", "vjValue")
	vj.SetCode(vjCode)
	vj.Save()

	stopVisit := referential.Model().StopVisits().New()
	code := model.NewCode("codeSpace", "value")
	stopVisit.SetCode(code)
	stopVisit.Save()

	event := &model.StopVisitUpdateEvent{
		Code:               code,
		StopAreaCode:       saCode,
		VehicleJourneyCode: vjCode,
		DepartureStatus:    model.STOP_VISIT_DEPARTURE_ONTIME,
		ArrivalStatus:      model.STOP_VISIT_ARRIVAL_ARRIVED,
		Schedules:          schedules.NewStopVisitSchedules(),
	}
	referential.collectManager.BroadcastUpdateEvent(event)

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
	stopArea.Origins.NewOrigin("partner")
	stopArea.Monitored = true
	stopArea.SetCode(model.NewCode("test", "value"))
	stopArea.Save()

	referential.CollectManager().HandlePartnerStatusChange("partner", false)

	updatedStopArea, _ := referential.Model().StopAreas().Find(stopArea.Id())
	if updatedStopArea.Monitored {
		t.Errorf("StopArea Monitored should be false after CollectManager UpdateStopArea")
	}
	if status, ok := updatedStopArea.Origins.Origin("partner"); !ok || status {
		t.Errorf("StopArea should have an Origin partner:false, got: %v", updatedStopArea.Origins.AllOrigin())
	}

	referential.CollectManager().HandlePartnerStatusChange("partner", true)

	updatedStopArea, _ = referential.Model().StopAreas().Find(stopArea.Id())
	if !updatedStopArea.Monitored {
		t.Errorf("StopArea Monitored should be true after CollectManager UpdateStopArea")
	}
	if status, ok := updatedStopArea.Origins.Origin("partner"); !ok || !status {
		t.Errorf("StopArea should have an Origin partner:false, got: %v", updatedStopArea.Origins.AllOrigin())
	}
}

func Test_CollectManager_StopAreaMonitoredWithReferent(t *testing.T) {
	// logger.Log.Debug = true
	// defer func() { logger.Log.Debug = false }()

	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	referentials.Save(referential)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Origins.NewOrigin("partner")
	stopArea.Monitored = true
	stopArea.SetCode(model.NewCode("test", "value"))
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.Origins.NewOrigin("partner2")
	stopArea2.ReferentId = stopArea.Id()
	stopArea2.Monitored = true
	stopArea2.SetCode(model.NewCode("test", "value"))
	stopArea2.Save()

	referential.CollectManager().HandlePartnerStatusChange("partner2", false)

	updatedStopArea, _ := referential.Model().StopAreas().Find(stopArea.Id())
	if updatedStopArea.Monitored {
		t.Errorf("StopArea Monitored should be false after CollectManager UpdateStopArea")
	}
	if status, ok := updatedStopArea.Origins.Origin("partner2"); !ok || status {
		t.Errorf("StopArea should have an Origin partner:false, got: %v", updatedStopArea.Origins.AllOrigin())
	}

	referential.CollectManager().HandlePartnerStatusChange("partner2", true)

	updatedStopArea, _ = referential.Model().StopAreas().Find(stopArea.Id())
	if !updatedStopArea.Monitored {
		t.Errorf("StopArea Monitored should be true after CollectManager UpdateStopArea")
	}
	if status, ok := updatedStopArea.Origins.Origin("partner2"); !ok || !status {
		t.Errorf("StopArea should have an Origin partner:false, got: %v", updatedStopArea.Origins.AllOrigin())
	}
}
