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
			logger.Log.Debugf("Model guardian visit")

			guardian.refreshStopAreas()
			guardian.checkReloadModel()
			guardian.simulateActualAttributes()

			c = guardian.Clock().After(10 * time.Second)
		}
	}
}

func (guardian *ModelGuardian) checkReloadModel() {
	if guardian.Clock().Now().After(guardian.referential.NextReloadAt()) {
		guardian.referential.ReloadModel()
	}
}

func (guardian *ModelGuardian) refreshStopAreas() {
	// Open a new transaction
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

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

		logger.Log.Debugf("Simulate actual attributes on StopVisit %v", stopVisit.Id())

		stopVisitTx := guardian.referential.NewTransaction()
		defer stopVisitTx.Close()

		transactionnalStopVisit, _ := tx.Model().StopVisits().Find(stopVisit.Id())
		simulator := NewActualAttributesSimulator(&transactionnalStopVisit)
		simulator.SetClock(guardian.Clock())
		if simulator.Simulate() {
			transactionnalStopVisit.Save()
			stopVisitTx.Commit()
		} else {
			stopVisitTx.Rollback()
		}
	}
}

type ActualAttributesSimulator struct {
	model.ClockConsumer

	stopVisit *model.StopVisit
	now       time.Time
}

func NewActualAttributesSimulator(stopVisit *model.StopVisit) *ActualAttributesSimulator {
	return &ActualAttributesSimulator{stopVisit: stopVisit}
}

func (simulator *ActualAttributesSimulator) Now() time.Time {
	if simulator.now.IsZero() {
		simulator.now = simulator.Clock().Now()
	}
	return simulator.now
}

func (simulator *ActualAttributesSimulator) ArrivalTime() time.Time {
	return simulator.stopVisit.Schedules.ArrivalTimeFromKind([]model.StopVisitScheduleType{"aimed", "expected"})
}

func (simulator *ActualAttributesSimulator) AfterArrivalTime() bool {
	return simulator.Clock().Now().After(simulator.ArrivalTime())
}

func (simulator *ActualAttributesSimulator) DepartureTime() time.Time {
	return simulator.stopVisit.Schedules.DepartureTimeFromKind([]model.StopVisitScheduleType{"aimed", "expected"})
}

func (simulator *ActualAttributesSimulator) AfterDepartureTime() bool {
	return simulator.Clock().Now().After(simulator.DepartureTime())
}

func (simulator *ActualAttributesSimulator) Simulate() bool {
	if simulator.stopVisit.IsCollected() == true {
		return false
	}

	return simulator.simulateArrival() || simulator.simulateDeparture()
}

func (simulator *ActualAttributesSimulator) simulateArrival() bool {
	if simulator.AfterArrivalTime() && simulator.CanArrive() {
		simulator.stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_ARRIVED
		simulator.stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, simulator.Now())

		logger.Log.Printf("Set StopVisit %s ArrivalStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_ARRIVAL_ARRIVED)

		if !simulator.AfterDepartureTime() {
			simulator.stopVisit.VehicleAtStop = true
			logger.Log.Printf("Set StopVisit %s VehicleAtStop at true", simulator.stopVisit.Id())
		}

		return true
	}

	return false
}

func (simulator *ActualAttributesSimulator) CanArrive() bool {
	switch simulator.stopVisit.ArrivalStatus {
	case model.STOP_VISIT_ARRIVAL_ONTIME, model.STOP_VISIT_ARRIVAL_EARLY, model.STOP_VISIT_ARRIVAL_DELAYED:
		return true
	default:
		return false
	}
}

func (simulator *ActualAttributesSimulator) CanDepart() bool {
	switch simulator.stopVisit.DepartureStatus {
	case model.STOP_VISIT_DEPARTURE_ONTIME, model.STOP_VISIT_DEPARTURE_EARLY, model.STOP_VISIT_DEPARTURE_DELAYED:
		return true
	default:
		return false
	}
}

func (simulator *ActualAttributesSimulator) simulateDeparture() bool {
	if simulator.AfterDepartureTime() && simulator.CanDepart() {
		simulator.stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_DEPARTED

		simulator.stopVisit.Schedules.SetDepartureTime(model.STOP_VISIT_SCHEDULE_ACTUAL, simulator.Now())
		simulator.stopVisit.VehicleAtStop = false

		logger.Log.Printf("Set StopVisit %s DepartureStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_DEPARTURE_DEPARTED)

		return true
	}
	return false
}
