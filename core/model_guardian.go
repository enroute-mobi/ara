package core

import (
	"fmt"
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
	guardian.simulateActualAttributes()
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Check StopAreas status")

	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		now := guardian.Clock().Now()
		if !stopArea.CollectedAlways && !stopArea.CollectedUntil.After(now) {
			continue
		}

		stopAreaTx := guardian.referential.NewTransaction()
		defer stopAreaTx.Close()
		transactionnalStopArea, _ := stopAreaTx.Model().StopAreas().Find(stopArea.Id())

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

func (guardian *ModelGuardian) simulateActualAttributes() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()
	for _, stopVisit := range tx.Model().StopVisits().FindAll() {
		if stopVisit.IsCollected() == true {
			continue
		}
		arrivalTime := stopVisit.Schedules.ArrivalTimeFromKind([]model.StopVisitScheduleType{"aimed", "expected"})
		departureTime := stopVisit.Schedules.DepartureTimeFromKind([]model.StopVisitScheduleType{"aimed", "expected"})

		now := guardian.Clock().Now()
		if now.After(arrivalTime) {
			stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_ARRIVED
			stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, now)
		}
		fmt.Println(arrivalTime, departureTime)
		if guardian.Clock().Now().After(arrivalTime) && departureTime.After(guardian.Clock().Now()) {
			stopVisit.VehicleAtStop = true
		}
		if guardian.Clock().Now().After(departureTime) {
			stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_DEPARTED

			stopVisit.Schedules.SetDepartureTime(model.STOP_VISIT_SCHEDULE_ACTUAL, now)
			stopVisit.VehicleAtStop = false
		}
		stopVisit.Save()
	}
}
