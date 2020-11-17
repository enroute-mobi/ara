package core

import (
	"math/rand"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type ModelGuardian struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	gmTimer     time.Time
	stop        chan struct{}
	referential *Referential
}

func NewModelGuardian(referential *Referential) *ModelGuardian {
	return &ModelGuardian{referential: referential}
}

func (guardian *ModelGuardian) Start() {
	logger.Log.Debugf("Start models guardian")

	rand.Seed(time.Now().UTC().UnixNano())
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
	guardian.gmTimer = guardian.Clock().Now()

	for {
		select {
		case <-guardian.stop:
			logger.Log.Debugf("Model guardian stop")
			return
		case <-c:
			logger.Log.Debugf("Model guardian visit")

			if guardian.checkReloadModel() {
				return
			}

			guardian.refreshStopAreas()
			guardian.refreshLines()
			guardian.simulateActualAttributes()
			guardian.requestSituations()

			c = guardian.Clock().After(10 * time.Second)
		}
	}
}

func (guardian *ModelGuardian) checkReloadModel() bool {
	if guardian.Clock().Now().After(guardian.referential.NextReloadAt()) {
		guardian.referential.ReloadModel()
		return true
	}
	return false
}

func (guardian *ModelGuardian) refreshStopAreas() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	now := guardian.Clock().Now()

	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		if stopArea.ParentId != "" {
			parent, ok := stopArea.Parent()
			if ok && !parent.CollectChildren {
				continue
			}
		}
		if !stopArea.CollectedAlways && !stopArea.CollectedUntil.After(now) {
			continue
		}

		if !stopArea.NextCollectAt().Before(now) {
			continue
		}

		stopAreaTx := guardian.referential.NewTransaction()

		transactionnalStopArea, _ := stopAreaTx.Model().StopAreas().Find(stopArea.Id())

		randNb := time.Duration(rand.Intn(20)+40) * time.Second

		transactionnalStopArea.NextCollect(now.Add(randNb))
		transactionnalStopArea.Save()
		stopAreaTx.Commit()
		stopAreaTx.Close()

		stopAreaUpdateRequest := &StopAreaUpdateRequest{
			id:         StopAreaUpdateRequestId(guardian.NewUUID()),
			stopAreaId: transactionnalStopArea.Id(),
			createdAt:  now,
		}
		guardian.referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)

		if stopArea.CollectGeneralMessages {
			situationUpdateRequest := NewSituationUpdateRequest(SITUATION_UPDATE_REQUEST_STOP_AREA, string(transactionnalStopArea.Id()))
			guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
		}
	}
}

func (guardian *ModelGuardian) refreshLines() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	now := guardian.Clock().Now()

	for _, line := range tx.Model().Lines().FindAll() {
		if !line.CollectGeneralMessages {
			continue
		}

		if !line.NextCollectAt().Before(now) {
			continue
		}

		lineTx := guardian.referential.NewTransaction()

		transactionnalLine, _ := lineTx.Model().Lines().Find(line.Id())

		randNb := time.Duration(rand.Intn(20)+40) * time.Second

		transactionnalLine.NextCollect(now.Add(randNb))
		transactionnalLine.Save()
		lineTx.Commit()
		lineTx.Close()

		situationUpdateRequest := NewSituationUpdateRequest(SITUATION_UPDATE_REQUEST_LINE, string(transactionnalLine.Id()))
		guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
	}
}

func (guardian *ModelGuardian) requestSituations() {
	if guardian.Clock().Now().Before(guardian.gmTimer.Add(1 * time.Minute)) {
		return
	}

	guardian.gmTimer = guardian.gmTimer.Add(1 * time.Minute)

	situationUpdateRequest := &SituationUpdateRequest{
		id:        SituationUpdateRequestId(guardian.NewUUID()),
		kind:      SITUATION_UPDATE_REQUEST_ALL,
		createdAt: guardian.Clock().Now(),
	}
	guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
}

func (guardian *ModelGuardian) simulateActualAttributes() {
	tx := guardian.referential.NewTransaction()
	defer tx.Close()

	for _, stopVisit := range tx.Model().StopVisits().FindAll() {
		if stopVisit.IsCollected() {
			continue
		}

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
	clock.ClockConsumer

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
	return simulator.stopVisit.Schedules.ArrivalTimeFromKind([]model.StopVisitScheduleType{"expected", "aimed"})
}

func (simulator *ActualAttributesSimulator) AfterArrivalTime() bool {
	return simulator.Clock().Now().After(simulator.ArrivalTime())
}

func (simulator *ActualAttributesSimulator) DepartureTime() time.Time {
	return simulator.stopVisit.Schedules.DepartureTimeFromKind([]model.StopVisitScheduleType{"expected", "aimed"})
}

func (simulator *ActualAttributesSimulator) AfterDepartureTime() bool {
	return simulator.Clock().Now().After(simulator.DepartureTime())
}

func (simulator *ActualAttributesSimulator) Simulate() bool {
	if simulator.stopVisit.IsCollected() {
		return false
	}

	return simulator.simulateArrival() || simulator.simulateDeparture()
}

func (simulator *ActualAttributesSimulator) simulateArrival() bool {
	if simulator.AfterArrivalTime() && simulator.CanArrive() {
		simulator.stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_CANCELLED
		simulator.stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, simulator.ArrivalTime())

		logger.Log.Printf("Set StopVisit %s ArrivalStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_ARRIVAL_CANCELLED)

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
		simulator.stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_CANCELLED

		simulator.stopVisit.Schedules.SetDepartureTime(model.STOP_VISIT_SCHEDULE_ACTUAL, simulator.DepartureTime())
		simulator.stopVisit.VehicleAtStop = false

		logger.Log.Printf("Set StopVisit %s DepartureStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_DEPARTURE_CANCELLED)

		return true
	}
	return false
}
