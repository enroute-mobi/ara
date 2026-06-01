package core

import (
	"context"
	"math/rand"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/monitoring"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
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

	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
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

			guardian.routineWork()

			c = guardian.Clock().After(10 * time.Second)
		}
	}
}

func (guardian *ModelGuardian) routineWork() {
	ctx := context.Background()

	span, spanContext := tracer.StartSpanFromContext(ctx, "model_guardian.routine")
	defer span.Finish()
	span.SetTag("referential", guardian.referential.Slug())

	guardian.refreshStopAreas(spanContext)
	guardian.refreshLines(spanContext)
	guardian.cleanOrUpdateStopVisits(spanContext)
	guardian.requestSituations(spanContext)
	guardian.refreshFacilities(spanContext)
}

func (guardian *ModelGuardian) checkReloadModel() bool {
	if guardian.Clock().Now().After(guardian.referential.NextReloadAt()) {
		guardian.referential.ReloadModel()
		return true
	}
	return false
}

func (guardian *ModelGuardian) refreshStopAreas(ctx context.Context) {
	child, _ := tracer.StartSpanFromContext(ctx, "refresh_stop_areas")
	defer child.Finish()

	defer monitoring.HandlePanic()

	now := guardian.Clock().Now()
	refresh := guardian.referential.ModelRefreshTime()

	sas := guardian.referential.Model().StopAreas().CollectableStopAreas(now)
	child.SetTag("collectable_stop_areas_count", len(sas))
	for i := range sas {
		sas[i].NextCollect(now.Add(guardian.randDuration(refresh)))
		sas[i].Save()

		stopAreaUpdateRequest := &StopAreaUpdateRequest{
			stopAreaId: sas[i].Id(),
			createdAt:  now,
		}
		guardian.referential.CollectManager().UpdateStopArea(stopAreaUpdateRequest)

		if sas[i].CollectSituations {
			situationUpdateRequest := NewSituationUpdateRequest(SITUATION_UPDATE_REQUEST_STOP_AREA, string(sas[i].Id()))
			guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
		}
	}
}

func (guardian *ModelGuardian) refreshFacilities(ctx context.Context) {
	child, childContext := tracer.StartSpanFromContext(ctx, "refresh_facilities")
	defer child.Finish()

	defer monitoring.HandlePanic()

	now := guardian.Clock().Now()
	refresh := guardian.referential.ModelRefreshTime()

	facilities := guardian.referential.Model().Facilities().CollectableFacilities(now)
	child.SetTag("collectable_facilities_count", len(facilities))

	for i := range facilities {
		facilities[i].NextCollect(now.Add(guardian.randDuration(refresh)))
		facilities[i].Save()

		facilityUpdateRequest := NewFacilityUpdateRequest(facilities[i].Id())
		guardian.referential.CollectManager().UpdateFacility(childContext, facilityUpdateRequest)
	}
}

func (guardian *ModelGuardian) refreshLines(ctx context.Context) {
	child, childContext := tracer.StartSpanFromContext(ctx, "refresh_lines")
	defer child.Finish()

	defer monitoring.HandlePanic()

	now := guardian.Clock().Now()
	refresh := guardian.referential.ModelRefreshTime()

	lines := guardian.referential.Model().Lines().CollectableLines(now)
	child.SetTag("collectable_lines_count", len(lines))
	for i := range lines {
		lines[i].NextCollect(now.Add(guardian.randDuration(refresh)))
		lines[i].Save()

		if lines[i].CollectSituations {
			situationUpdateRequest := NewSituationUpdateRequest(SITUATION_UPDATE_REQUEST_LINE, string(lines[i].Id()))
			guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
		}

		lineUpdateRequest := NewLineUpdateRequest(lines[i].Id())
		guardian.referential.CollectManager().UpdateLine(childContext, lineUpdateRequest)

		vehicleUpdateRequest := NewVehicleUpdateRequest(lines[i].Id())
		guardian.referential.CollectManager().UpdateVehicle(childContext, vehicleUpdateRequest)
	}
}

func (guardian *ModelGuardian) randDuration(refresh time.Duration) time.Duration {
	return time.Duration(rand.Intn(20)-10)*time.Second + refresh
}

func (guardian *ModelGuardian) requestSituations(ctx context.Context) {
	child, _ := tracer.StartSpanFromContext(ctx, "request_situations")
	defer child.Finish()
	defer monitoring.HandlePanic()

	if guardian.Clock().Now().Before(guardian.gmTimer.Add(1 * time.Minute)) {
		return
	}

	guardian.gmTimer = guardian.gmTimer.Add(1 * time.Minute)

	situationUpdateRequest := &SituationUpdateRequest{
		kind:      SITUATION_UPDATE_REQUEST_ALL,
		createdAt: guardian.Clock().Now(),
	}
	guardian.referential.CollectManager().UpdateSituation(situationUpdateRequest)
}

func (guardian *ModelGuardian) cleanOrUpdateStopVisits(ctx context.Context) {
	child, _ := tracer.StartSpanFromContext(ctx, "clean_or_update_stop_visits")
	defer child.Finish()

	defer monitoring.HandlePanic()

	m := guardian.referential.Model()

	svs := m.StopVisits().UnsafeFindAll()
	persistence := guardian.referential.ModelPersistenceDuration()
	vjs := make(map[model.VehicleJourneyId]struct{})
	stopVisitstoDelete := []*model.StopVisit{}

	child.SetTag("stop_visits_count", len(svs))
	for i := range svs {
		if svs[i].ReferenceTime().Before(guardian.Clock().Now().Add(persistence)) {
			vjs[svs[i].VehicleJourneyId] = struct{}{}
			stopVisitstoDelete = append(stopVisitstoDelete, svs[i])
			continue
		}

		if svs[i].IsCollected() {
			continue
		}

		simulator := NewActualAttributesSimulator(svs[i])
		simulator.SetClock(guardian.Clock())
		if simulator.Simulate() {
			svs[i].Save()
			if svs[i].IsArchivable() {
				sva := &model.StopVisitArchiver{
					Model:     guardian.referential.Model(),
					StopVisit: svs[i],
				}
				sva.Archive()
			}

		}
	}

	logger.Log.Debugf("Referential persistence deleting %d StopVisits", len(stopVisitstoDelete))
	m.StopVisits().DeleteMultiple(stopVisitstoDelete)

	for id := range vjs {
		if !m.StopVisits().VehicleJourneyHasStopVisits(id) {
			m.VehicleJourneys().DeleteById(id)
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
	return simulator.stopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected, schedules.Aimed})
}

func (simulator *ActualAttributesSimulator) AfterArrivalTime() bool {
	return simulator.Clock().Now().After(simulator.ArrivalTime())
}

func (simulator *ActualAttributesSimulator) DepartureTime() time.Time {
	return simulator.stopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected, schedules.Aimed})
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
		simulator.stopVisit.Schedules.SetArrivalTime(schedules.Actual, simulator.ArrivalTime())

		logger.Log.Debugf("Set StopVisit %s ArrivalStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_ARRIVAL_CANCELLED)

		if !simulator.AfterDepartureTime() {
			simulator.stopVisit.VehicleAtStop = true
			logger.Log.Debugf("Set StopVisit %s VehicleAtStop at true", simulator.stopVisit.Id())
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

		simulator.stopVisit.Schedules.SetDepartureTime(schedules.Actual, simulator.DepartureTime())
		simulator.stopVisit.VehicleAtStop = false

		logger.Log.Debugf("Set StopVisit %s DepartureStatus at %s", simulator.stopVisit.Id(), model.STOP_VISIT_DEPARTURE_CANCELLED)

		return true
	}
	return false
}
