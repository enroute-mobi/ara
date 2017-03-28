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
		model: model.NewMemoryModel(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	fakeClock := model.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules = make(model.StopVisitSchedules)

	stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_ONTIME
	stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_ONTIME

	stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_AIMED, fakeClock.Now().Add(1*time.Minute))
	stopVisit.Schedules.SetDepartureTime(model.STOP_VISIT_SCHEDULE_AIMED, fakeClock.Now().Add(10*time.Minute))

	stopVisit.Save()

	fakeClock.Advance(1*time.Minute + 1*time.Second)
	referential.modelGuardian.simulateActualAttributes()

	stopVisit, _ = referential.Model().StopVisits().Find(stopVisit.Id())
	if expected := model.STOP_VISIT_ARRIVAL_CANCELLED; stopVisit.ArrivalStatus != expected {
		t.Errorf("Wrong StopVisit ArrivalStatus at %s\n want: %#v\n got: %#v", fakeClock.Now(), expected, stopVisit.ArrivalStatus)
	}
	if !stopVisit.VehicleAtStop {
		t.Errorf("Wrong StopVisit VehicleAtStop at %s\n want: %#v\n got: %#v", fakeClock.Now(), true, stopVisit.VehicleAtStop)
	}

	fakeClock.Advance(10 * time.Minute)
	referential.modelGuardian.simulateActualAttributes()

	stopVisit, _ = referential.Model().StopVisits().Find(stopVisit.Id())
	if expected := model.STOP_VISIT_ARRIVAL_CANCELLED; stopVisit.ArrivalStatus != expected {
		t.Errorf("Wrong StopVisit ArrivalStatus at %s\n want: %#v\n got: %#v", fakeClock.Now(), expected, stopVisit.ArrivalStatus)
	}
	if expected := model.STOP_VISIT_DEPARTURE_CANCELLED; stopVisit.DepartureStatus != expected {
		t.Errorf("Wrong StopVisit DepartureStatus at %s\n want: %#v\n got: %#v", fakeClock.Now(), expected, stopVisit.DepartureStatus)
	}
	if stopVisit.VehicleAtStop {
		t.Errorf("Wrong StopVisit VehicleAtStop at %s\n want: %#v\n got: %#v", fakeClock.Now(), false, stopVisit.VehicleAtStop)
	}
}
