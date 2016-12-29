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
	for {
		select {
		case <-guardian.stop:
			return
		case <-guardian.Clock().After(10 * time.Second):
			guardian.refreshStopAreas()
		}
	}
}

func (guardian *ModelGuardian) refreshStopAreas() {
	// Open a new transaction
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Check StopAreas status")

	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		now := guardian.Clock().Now()
		outdated := now.Add(-1 * time.Minute)

		if stopArea.RequestedAt().Before(outdated) && stopArea.UpdatedAt().Before(outdated) {
			stopAreaUpdateRequest := &StopAreaUpdateRequest{
				id:         StopAreaUpdateRequestId(guardian.NewUUID()),
				stopAreaId: stopArea.Id(),
				createdAt:  guardian.Clock().Now(),
			}
			guardian.referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)

			stopArea.Requested(now)
			tx.Model().StopAreas().Save(&stopArea)
		}
	}
	tx.Commit()
}
