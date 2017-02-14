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

func (guardian *ModelGuardian) Reload() {
	guardian.deleteAllVehiculeJourney()
	guardian.deleteAllStopVisits()
	guardian.deleteLines()
}

func (guardian *ModelGuardian) Run() {
	guardian.referential.model.SetDate(guardian.referential.Setting("reload_at"))
	c := guardian.Clock().After(10 * time.Second)
	d := guardian.Clock().After(10 * time.Minute)

	for {
		select {
		case <-guardian.stop:
			return
		case <-c:
			guardian.refreshStopAreas()
			c = guardian.Clock().After(10 * time.Second)
		case <-d:
			if time.Now().After(guardian.referential.NextReloadAt()) == true {
				guardian.referential.ReloadModel()
				guardian.referential.model.SetDate(guardian.referential.Setting("reload_at"))
			}
			d = guardian.Clock().After(10 * time.Minute)
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

func (guardian *ModelGuardian) deleteAllVehiculeJourney() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	for _, vehicleJourney := range tx.Model().VehicleJourneys().FindAll() {

		transactionnalVehicleJourney, _ := tx.Model().VehicleJourneys().Find(vehicleJourney.Id())
		tx.Model().VehicleJourneys().Delete(&transactionnalVehicleJourney)
		tx.Commit()
	}
}

func (guardian *ModelGuardian) deleteAllStopVisits() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	for _, stopVisit := range tx.Model().StopVisits().FindAll() {
		transactionnalStopVisit, _ := tx.Model().StopVisits().Find(stopVisit.Id())
		guardian.referential.model.StopVisits().Delete(&transactionnalStopVisit)
		tx.Commit()
	}
}

func (guardian *ModelGuardian) deleteLines() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	for _, line := range tx.Model().Lines().FindAll() {
		transactionnalLine, _ := tx.Model().Lines().Find(line.Id())
		guardian.referential.model.Lines().Delete(&transactionnalLine)
		tx.Commit()
	}
}
