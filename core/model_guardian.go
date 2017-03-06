package core

import (
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type ModelGuardian struct {
	model.ClockConsumer
	model.UUIDConsumer

	stop        chan struct{}
	referential *Referential
}

func NewModelGuardian(referential *Referential) *ModelGuardian {
	return &ModelGuardian{referential: referential}
}

func (guardian *ModelGuardian) Start() {
	logger.Log.Debugf("Start models guardian")

	guardian.stop = make(chan struct{})
	go guardian.Run()
}

func (guardian *ModelGuardian) Stop() {
	if guardian.stop != nil {
		close(guardian.stop)
	}
}

func (guardian *ModelGuardian) Run() {
	c := guardian.Clock().After(10 * time.Second)

	for {
		select {
		case <-guardian.stop:
			return
		case <-c:
			guardian.refreshStopAreas()

			if guardian.Clock().Now().After(guardian.referential.NextReloadAt()) {
				guardian.referential.ReloadModel()
			}

			c = guardian.Clock().After(10 * time.Second)
		}
	}
}

func (guardian *ModelGuardian) refreshStopAreas() {
	// Open a new transaction
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Check StopAreas status")

	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		stopAreaTx := guardian.referential.NewTransaction()
		defer stopAreaTx.Close()
		transactionnalStopArea, _ := stopAreaTx.Model().StopAreas().Find(stopArea.Id())

		now := guardian.Clock().Now()
		outdated := now.Add(-1 * time.Minute)

		if transactionnalStopArea.RequestedAt().Before(outdated) && transactionnalStopArea.UpdatedAt().Before(outdated) {
			transactionnalStopArea.Requested(now)
			stopAreaTx.Model().StopAreas().Save(&transactionnalStopArea)
			stopAreaTx.Commit()

			stopAreaUpdateRequest := &StopAreaUpdateRequest{
				id:         StopAreaUpdateRequestId(guardian.NewUUID()),
				stopAreaId: transactionnalStopArea.Id(),
				createdAt:  guardian.Clock().Now(),
			}
			guardian.referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)
		}
	}
}
