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
