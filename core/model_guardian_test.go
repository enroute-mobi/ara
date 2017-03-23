package core

import (
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_ModelGuardian_Run(t *testing.T) {
	referential := &Referential{
		model:          model.NewMemoryModel(),
		collectManager: NewTestCollectManager(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	fakeClock := model.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	referential.Model().StopAreas().Save(&stopArea)
	stopAreaId := stopArea.Id()

	referential.ModelGuardian().Start()
	defer referential.ModelGuardian().Stop()

	// Wait for the guardian to launch Run
	fakeClock.BlockUntil(1)
	// Advance time
	fakeClock.Advance(11 * time.Second)
	// Wait for the Test CollectManager to finish Status()
	select {
	case <-referential.CollectManager().(*TestCollectManager).Done:
		updatedStopArea, ok := referential.Model().StopAreas().Find(stopAreaId)
		if !ok {
			t.Error("StopArea should still be found after guardian work")
		} else if updatedStopArea.RequestedAt() != fakeClock.Now() {
			t.Errorf("StopArea should have RequestedAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.RequestedAt())
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Guardian CheckPartnerStatus with TestCheckStatusClient timed out")
	}
}

func Test_ModelGuardian_RefreshStopAreas_CollectedUntil(t *testing.T) {
	referential := &Referential{
		model:          model.NewMemoryModel(),
		collectManager: NewTestCollectManager(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	fakeClock := model.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	stopArea.CollectedAlways = false
	stopArea.CollectedUntil = fakeClock.Now().Add(15 * time.Minute)
	referential.Model().StopAreas().Save(&stopArea)

	referential.modelGuardian.refreshStopAreas()

	updatedStopArea, ok := referential.Model().StopAreas().Find(stopArea.Id())
	if !ok {
		t.Fatal("StopArea not found after guardian work")
	}

	if updatedStopArea.RequestedAt() != fakeClock.Now() {
		t.Errorf("StopArea should have RequestedAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.RequestedAt())
	}

	requestedAt := updatedStopArea.RequestedAt()

	fakeClock.Advance(15*time.Minute + time.Second)

	referential.modelGuardian.refreshStopAreas()

	updatedStopArea, ok = referential.Model().StopAreas().Find(stopArea.Id())
	if !ok {
		t.Error("StopArea should still be found after guardian work")
	}
	if updatedStopArea.RequestedAt() != requestedAt {
		t.Errorf("StopArea should have RequestedAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.RequestedAt())
	}
}

func Test_ModelGuardian_Run_simulateActualAttributes(t *testing.T) {

	referential := &Referential{
		model:          model.NewMemoryModel(),
		collectManager: NewTestCollectManager(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	fakeClock := model.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	stopArea.CollectedAlways = true
	referential.Model().StopAreas().Save(&stopArea)

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules = make(model.StopVisitSchedules)

	stopVisit.Schedules[model.STOP_VISIT_SCHEDULE_AIMED] = &model.StopVisitSchedule{}
	stopVisit.Schedules[model.STOP_VISIT_SCHEDULE_AIMED].SetArrivalTime(referential.ModelGuardian().Clock().Now())
	stopVisit.Schedules[model.STOP_VISIT_SCHEDULE_AIMED].SetDepartureTime(referential.ModelGuardian().Clock().Now().Add(10 * time.Minute))
	stopVisit.Schedules[model.STOP_VISIT_SCHEDULE_ACTUAL] = &model.StopVisitSchedule{}

	referential.Model().StopVisits().Save(&stopVisit)
	stopVisit.Save()
	referential.ModelGuardian().Start()
	defer referential.ModelGuardian().Stop()

	// Wait for the guardian to launch Run
	fakeClock.BlockUntil(1)
	// Advance time
	fakeClock.Advance(11 * time.Second)
	if referential.Model().StopVisits().FindByStopAreaId(stopArea.Id())[0].VehicleAtStop != false {
		t.Errorf("VehicleAtStop should be set at false")
	}
	fakeClock.Advance(5 * time.Minute)
	time.Sleep(100 * time.Millisecond)

	if referential.Model().StopVisits().FindByStopAreaId(stopArea.Id())[0].VehicleAtStop != true {
		t.Errorf("VehicleAtStop should be set at true")
	}

	fakeClock.Advance(10 * time.Minute)
	time.Sleep(100 * time.Millisecond)

	if referential.Model().StopVisits().FindByStopAreaId(stopArea.Id())[0].VehicleAtStop != false {
		t.Errorf("VehicleAtStop should be set at false")
	}
}
