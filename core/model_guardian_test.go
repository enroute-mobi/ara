package core

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Test_ModelGuardian_RefreshStopAreas_RequestedAt(t *testing.T) {
	assert := assert.New(t)
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
	referential.Model().StopAreas().Save(stopArea)
	stopAreaId := stopArea.Id()

	// Advance time
	fakeClock.Advance(11 * time.Second)
	referential.modelGuardian.refreshStopAreas(spanCtx)

	updatedStopArea, ok := referential.Model().StopAreas().Find(stopAreaId)
	assert.Truef(ok, "StopArea should still be found after guardian work")
	assert.Truef(updatedStopArea.NextCollectAt().After(fakeClock.Now()),
		"StopArea should have NextCollectAt before fakeClock %v, got: %v",
		fakeClock.Now(), updatedStopArea.NextCollectAt())

	// Advance time
	fakeClock.Advance(61 * time.Second)
	updatedStopArea, _ = referential.Model().StopAreas().Find(stopAreaId)
	assert.True(updatedStopArea.NextCollectAt().Before(fakeClock.Now()))
}

func Test_ModelGuardian_RandDuration_Without_Refresh_setting(t *testing.T) {
	assert := assert.New(t)
	referential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(referential)

	randDuration := referential.ModelGuardian().randDuration()
	assert.InDeltaf(time.Duration(s.DEFAULT_MODEL_REFRESH_TIME).Seconds(), randDuration.Seconds(),
		10.0,
		"should be between -10s/+10s range from the Default model.refresh_time of 50s")
}

func Test_ModelGuardian_RandDuration_With_Refresh_setting(t *testing.T) {
	assert := assert.New(t)
	referential := referentials.New(ReferentialSlug("referential"))
	referential.SetSetting("model.refresh_time", "45s")
	referentials.Save(referential)

	randDuration := referential.ModelGuardian().randDuration()
	assert.InDeltaf(45.0, randDuration.Seconds(),
		10.0,
		"should be between -10s/+10s range from the model.refresh_time")
}

func Test_ModelGuardian_RefreshStopAreas_CollectedUntil(t *testing.T) {
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()
	testSpan, spanCtx := tracer.StartSpanFromContext(ctx, "test.span")
	defer testSpan.Finish()

	referential := &Referential{
		model:          model.NewTestMemoryModel(),
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
	stopVisit.Schedules = schedules.NewStopVisitSchedules()

	stopVisit.DepartureStatus = model.STOP_VISIT_DEPARTURE_ONTIME
	stopVisit.ArrivalStatus = model.STOP_VISIT_ARRIVAL_ONTIME

	stopVisit.Schedules.SetArrivalTime(schedules.Aimed, fakeClock.Now().Add(1*time.Minute))
	stopVisit.Schedules.SetDepartureTime(schedules.Aimed, fakeClock.Now().Add(10*time.Minute))

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
	ctx := context.Background()

	mt := mocktracer.Start()
	defer mt.Stop()
	testSpan, spanCtx := tracer.StartSpanFromContext(ctx, "test.span")
	defer testSpan.Finish()

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
	sv1.Schedules = schedules.NewStopVisitSchedules()
	sv1.Schedules.SetArrivalTime(schedules.Actual, fakeClock.Now().Add(-1*time.Hour))
	sv1.Save()

	sv2 := referential.Model().StopVisits().New()
	sv2.VehicleJourneyId = vj1.Id()
	sv2.Schedules = schedules.NewStopVisitSchedules()
	sv2.Schedules.SetArrivalTime(schedules.Actual, fakeClock.Now().Add(-1*time.Hour))
	sv2.Save()

	sv3 := referential.Model().StopVisits().New()
	sv3.VehicleJourneyId = vj2.Id()
	sv3.Schedules = schedules.NewStopVisitSchedules()
	sv3.Schedules.SetArrivalTime(schedules.Actual, fakeClock.Now().Add(-1*time.Hour))
	sv3.Save()

	sv4 := referential.Model().StopVisits().New()
	sv4.VehicleJourneyId = vj2.Id()
	sv4.Schedules = schedules.NewStopVisitSchedules()
	sv4.Schedules.SetArrivalTime(schedules.Actual, fakeClock.Now())
	sv4.Save()

	referential.modelGuardian.cleanOrUpdateStopVisits(spanCtx)

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
