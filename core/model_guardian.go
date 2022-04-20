package core

import (
	"math/rand"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/monitoring"
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
	defer monitoring.HandlePanic()

	now := guardian.Clock().Now()

	sas := guardian.referential.Model().StopAreas().FindAll()
	for i := range sas {
		if sas[i].ParentId != "" {
			parent, ok := sas[i].Parent()
			if ok && !parent.CollectChildren {
				continue
			}
		}
		if !sas[i].CollectedAlways && !sas[i].CollectedUntil.After(now) {
			continue
		}

		if !sas[i].NextCollectAt().Before(now) {
			continue
		}

		stopArea, _ := guardian.referential.Model().StopAreas().Find(sas[i].Id())

		randNb := time.Duration(rand.Intn(20)+40) * time.Second

		stopArea.NextCollect(now.Add(randNb))
		stopArea.Save()

		stopAreaUpdateRequest := &StopAreaUpdateRequest{
			id:         StopAreaUpdateRequestId(guardian.NewUUID()),
			stopAreaId: stopArea.Id(),
			createdAt:  now,
		}
		guardian.referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)

		if sas[i].CollectGeneralMessages {
			situationUpdateRequest := NewSituationUpdateRequest(SITUATION_UPDATE_REQUEST_STOP_AREA, string(stopArea.Id()))
			guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
		}
	}
}

func (guardian *ModelGuardian) refreshLines() {
	defer monitoring.HandlePanic()

	now := guardian.Clock().Now()

	lines := guardian.referential.Model().Lines().FindAll()
	for i := range lines {
		if !lines[i].NextCollectAt().Before(now) {
			continue
		}

		line, _ := guardian.referential.Model().Lines().Find(lines[i].Id())

		randNb := time.Duration(rand.Intn(20)+40) * time.Second

		line.NextCollect(now.Add(randNb))
		line.Save()

		if lines[i].CollectGeneralMessages {
			situationUpdateRequest := NewSituationUpdateRequest(SITUATION_UPDATE_REQUEST_LINE, string(line.Id()))
			guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
		}

		vehicleUpdateRequest := NewVehicleUpdateRequest(line.Id())
		guardian.referential.CollectManager().UpdateVehicle(vehicleUpdateRequest)
	}
}

func (guardian *ModelGuardian) requestSituations() {
	defer monitoring.HandlePanic()

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
	defer monitoring.HandlePanic()

	svs := guardian.referential.Model().StopVisits().FindAll()
	for i := range svs {
		if svs[i].IsCollected() {
			continue
		}

		stopVisit, _ := guardian.referential.Model().StopVisits().Find(svs[i].Id())
		simulator := NewActualAttributesSimulator(&stopVisit)
		simulator.SetClock(guardian.Clock())
		if simulator.Simulate() {
			stopVisit.Save()
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
		simulator.stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_ARRIVED
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
		simulator.stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_DEPARTED

		simulator.stopVisit.Schedules.SetDepartureTime(model.STOP_VISIT_SCHEDULE_ACTUAL, simulator.DepartureTime())
		simulator.stopVisit.VehicleAtStop = false

		logger.Log.Printf("Set StopVisit %s DepartureStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_DEPARTURE_CANCELLED)

		return true
	}
	return false
}
