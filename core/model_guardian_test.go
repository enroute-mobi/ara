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

func Test_ModelGuardian_Run_On_MonitoredUntil(t *testing.T) {
	referential := &Referential{
		model:          model.NewMemoryModel(),
		collectManager: NewTestCollectManager(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	requested := time.Time{}
	fakeClock := model.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	stopArea.MonitoredAlways = false
	stopArea.MonitoredUntil = fakeClock.Now().Add(15 * time.Minute)
	referential.Model().StopAreas().Save(&stopArea)
	stopAreaId := stopArea.Id()

	referential.ModelGuardian().Start()
	defer referential.ModelGuardian().Stop()

	// Wait for the guardian to launch Run
	fakeClock.BlockUntil(1)
	// Advance time
	fakeClock.Advance(11 * time.Second)

	select {
	case <-referential.CollectManager().(*TestCollectManager).Done:
		updatedStopArea, ok := referential.Model().StopAreas().Find(stopAreaId)
		if !ok {
			t.Error("StopArea should still be found after guardian work")
		} else if requested = updatedStopArea.RequestedAt(); requested != fakeClock.Now() {
			t.Errorf("StopArea should have RequestedAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.RequestedAt())
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Guardian CheckPartnerStatus with TestCheckStatusClient timed out")
	}

	fakeClock.Advance(16 * time.Minute)

	time.Sleep(100 * time.Millisecond)
	updatedStopArea, ok := referential.Model().StopAreas().Find(stopAreaId)
	if !ok {
		t.Error("StopArea should still be found after guardian work")
	}
	if updatedStopArea.RequestedAt() != requested {
		t.Errorf("StopArea should have RequestedAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.RequestedAt())
	}
}
