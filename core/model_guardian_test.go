package core

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Test_ModelGuardian_RefreshStopAreas_RequestedAt(t *testing.T) {
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()
	testSpan, spanCtx := tracer.StartSpanFromContext(ctx, "test.span")
	defer testSpan.Finish()

	referential := &Referential{
		model:          model.NewMemoryModel(),
		collectManager: NewTestCollectManager(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	fakeClock := clock.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	referential.Model().StopAreas().Save(stopArea)
	stopAreaId := stopArea.Id()

	// Advance time
	fakeClock.Advance(11 * time.Second)
	referential.modelGuardian.refreshStopAreas(spanCtx)

	updatedStopArea, ok := referential.Model().StopAreas().Find(stopAreaId)
	if !ok {
		t.Fatal("StopArea should still be found after guardian work")
	}

	if updatedStopArea.NextCollectAt().Before(fakeClock.Now()) {
		t.Errorf("StopArea should have NextCollectAt before fakeClock %v, got: %v", fakeClock.Now(), updatedStopArea.NextCollectAt())
	}
}

func Test_ModelGuardian_RefreshStopAreas_CollectedUntil(t *testing.T) {
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()
	testSpan, spanCtx := tracer.StartSpanFromContext(ctx, "test.span")
	defer testSpan.Finish()

	referential := &Referential{
		model:          model.NewMemoryModel(),
		collectManager: NewTestCollectManager(),
	}
	referential.modelGuardian = NewModelGuardian(referential)

	fakeClock := clock.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	stopArea.CollectedAlways = false
	stopArea.CollectedUntil = fakeClock.Now().Add(15 * time.Minute)
	stopArea.NextCollect(fakeClock.Now())
	referential.Model().StopAreas().Save(stopArea)

	referential.modelGuardian.refreshStopAreas(spanCtx)

	updatedStopArea, ok := referential.Model().StopAreas().Find(stopArea.Id())
	if !ok {
		t.Fatal("StopArea not found after guardian work")
	}

	if updatedStopArea.NextCollectAt().Before(fakeClock.Now()) {
		t.Errorf("StopArea should have NextCollectAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.NextCollectAt())
	}

	nextCollectAt := updatedStopArea.NextCollectAt()

	fakeClock.Advance(15*time.Minute + time.Second)

	referential.modelGuardian.refreshStopAreas(spanCtx)

	updatedStopArea, ok = referential.Model().StopAreas().Find(stopArea.Id())
	if !ok {
		t.Error("StopArea should still be found after guardian work")
	}
	if updatedStopArea.NextCollectAt().After(nextCollectAt) {
		t.Errorf("StopArea should have NextCollectAt set at %v, got: %v", fakeClock.Now(), updatedStopArea.NextCollectAt())
	}
}

func Test_ModelGuardian_Run_cleanOrUpdateStopVisits(t *testing.T) {
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()
	testSpan, spanCtx := tracer.StartSpanFromContext(ctx, "test.span")
	defer testSpan.Finish()

	referential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(referential)

	fakeClock := clock.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules = model.NewStopVisitSchedules()

	stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_ONTIME
	stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_ONTIME

	stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_AIMED, fakeClock.Now().Add(1*time.Minute))
	stopVisit.Schedules.SetDepartureTime(model.STOP_VISIT_SCHEDULE_AIMED, fakeClock.Now().Add(10*time.Minute))

	stopVisit.Save()

	fakeClock.Advance(1*time.Minute + 1*time.Second)
	referential.modelGuardian.cleanOrUpdateStopVisits(spanCtx)

	stopVisit, _ = referential.Model().StopVisits().Find(stopVisit.Id())
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; stopVisit.ArrivalStatus != expected {
		t.Errorf("Wrong StopVisit ArrivalStatus at %s\n want: %#v\n got: %#v", fakeClock.Now(), expected, stopVisit.ArrivalStatus)
	}
	if !stopVisit.VehicleAtStop {
		t.Errorf("Wrong StopVisit VehicleAtStop at %s\n want: %#v\n got: %#v", fakeClock.Now(), true, stopVisit.VehicleAtStop)
	}

	fakeClock.Advance(10 * time.Minute)
	referential.modelGuardian.cleanOrUpdateStopVisits(spanCtx)

	stopVisit, _ = referential.Model().StopVisits().Find(stopVisit.Id())
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; stopVisit.ArrivalStatus != expected {
		t.Errorf("Wrong StopVisit ArrivalStatus at %s\n want: %#v\n got: %#v", fakeClock.Now(), expected, stopVisit.ArrivalStatus)
	}
	if expected := model.STOP_VISIT_DEPARTURE_DEPARTED; stopVisit.DepartureStatus != expected {
		t.Errorf("Wrong StopVisit DepartureStatus at %s\n want: %#v\n got: %#v", fakeClock.Now(), expected, stopVisit.DepartureStatus)
	}
	if stopVisit.VehicleAtStop {
		t.Errorf("Wrong StopVisit VehicleAtStop at %s\n want: %#v\n got: %#v", fakeClock.Now(), false, stopVisit.VehicleAtStop)
	}
}

func Test_ModelGuardian_Run_cleanOrUpdateStopVisits_Clean(t *testing.T) {
	referentials := NewMemoryReferentials()

	referential := referentials.New(ReferentialSlug("referential"))
	referential.SetSetting(s.MODEL_PERSISTENCE, "30M")
	referentials.Save(referential)

	fakeClock := clock.NewFakeClock()
	referential.ModelGuardian().SetClock(fakeClock)

	vj1 := referential.Model().VehicleJourneys().New()
	vj1.Save()

	vj2 := referential.Model().VehicleJourneys().New()
	vj2.Save()

	sv1 := referential.Model().StopVisits().New()
	sv1.VehicleJourneyId = vj1.Id()
	sv1.Schedules = model.NewStopVisitSchedules()
	sv1.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, fakeClock.Now().Add(-1*time.Hour))
	sv1.Save()

	sv2 := referential.Model().StopVisits().New()
	sv2.VehicleJourneyId = vj1.Id()
	sv2.Schedules = model.NewStopVisitSchedules()
	sv2.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, fakeClock.Now().Add(-1*time.Hour))
	sv2.Save()

	sv3 := referential.Model().StopVisits().New()
	sv3.VehicleJourneyId = vj2.Id()
	sv3.Schedules = model.NewStopVisitSchedules()
	sv3.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, fakeClock.Now().Add(-1*time.Hour))
	sv3.Save()

	sv4 := referential.Model().StopVisits().New()
	sv4.VehicleJourneyId = vj2.Id()
	sv4.Schedules = model.NewStopVisitSchedules()
	sv4.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, fakeClock.Now())
	sv4.Save()

	referential.modelGuardian.cleanOrUpdateStopVisits()

	_, ok := referential.Model().StopVisits().Find(sv1.Id())
	if ok {
		t.Errorf("Shoudldn't find StopVisit after clean")
	}
	_, ok = referential.Model().StopVisits().Find(sv2.Id())
	if ok {
		t.Errorf("Shoudldn't find StopVisit after clean")
	}
	_, ok = referential.Model().StopVisits().Find(sv3.Id())
	if ok {
		t.Errorf("Shoudldn't find StopVisit after clean")
	}
	_, ok = referential.Model().StopVisits().Find(sv4.Id())
	if !ok {
		t.Errorf("Shoudld find StopVisit after clean")
	}
	_, ok = referential.Model().VehicleJourneys().Find(vj1.Id())
	if ok {
		t.Errorf("Shoudldn't find VehicleJourney after clean")
	}
	_, ok = referential.Model().VehicleJourneys().Find(vj2.Id())
	if !ok {
		t.Errorf("Shoudld find VehicleJourney after clean")
	}
}
