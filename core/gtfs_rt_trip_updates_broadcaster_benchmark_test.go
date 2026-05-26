package core

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func benchmarkHandleGtfs(pc int, seed bool, b *testing.B) {
	// Setup
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	// seed unused models
	if seed {
		for i := 0; i != 500; i++ {
			vehicleJourney := referential.model.VehicleJourneys().New()
			vjId := model.NewCode("internal", fmt.Sprintf("vehicle_jourey%d", i))
			vehicleJourney.SetCode(vjId)
			vehicleJourney.LineId = line.Id()
			vehicleJourney.Save()

			// stopAreas & stopVisits
			for j := 0; j < 30; j++ {
				saId := model.NewCode("internal", fmt.Sprintf("stop_area%d", j))
				stopArea := referential.Model().StopAreas().New()
				stopArea.SetCode(saId)
				stopArea.Save()

				stopVisit := referential.model.StopVisits().New()
				svId1 := model.NewCode("internal", fmt.Sprintf("stop_visit%d", j))
				stopVisit.SetCode(svId1)
				stopVisit.StopAreaId = stopArea.Id()
				stopVisit.VehicleJourneyId = vehicleJourney.Id()

				delta := time.Duration(float64(j) * float64(time.Minute))
				base := connector.Clock().Now().Add(-4 * time.Minute)
				stopVisit.Schedules.SetDepartureTime("actual", base.Add(delta))
				stopVisit.PassageOrder = j

				stopVisit.Save()
			}
		}
	}

	// models for benchmark
	for i := 0; i != pc; i++ {
		vehicleJourney := referential.model.VehicleJourneys().New()
		vjId := model.NewCode("internal", fmt.Sprintf("vj%d", i))
		vehicleJourney.SetCode(vjId)
		vehicleJourney.LineId = line.Id()
		vehicleJourney.Save()

		// stopAreas & stopVisits
		for j := 0; j < 30; j++ {
			saId := model.NewCode("internal", fmt.Sprintf("saId%d", j))
			stopArea := referential.Model().StopAreas().New()
			stopArea.SetCode(saId)
			stopArea.Save()

			stopVisit := referential.model.StopVisits().New()
			svId1 := model.NewCode("internal", fmt.Sprintf("svId%d", j))
			stopVisit.SetCode(svId1)
			stopVisit.StopAreaId = stopArea.Id()
			stopVisit.VehicleJourneyId = vehicleJourney.Id()

			delta := time.Duration(float64(j) * float64(time.Minute))
			base := connector.Clock().Now().Add(-10 * time.Minute)
			stopVisit.Schedules.SetDepartureTime("actual", base.Add(delta))

			stopVisit.PassageOrder = j

			stopVisit.Save()

		}
	}

	gtfsFeed := &gtfs.FeedMessage{}
	connector.HandleGtfs(gtfsFeed)
}

func BenchmarkHandleGtfs10(b *testing.B)    { benchmarkHandleGtfs(9, false, b) }
func BenchmarkHandleGtfs50(b *testing.B)    { benchmarkHandleGtfs(49, false, b) }
func BenchmarkHandleGtfs100(b *testing.B)   { benchmarkHandleGtfs(99, false, b) }
func BenchmarkHandleGtfs1000(b *testing.B)  { benchmarkHandleGtfs(999, false, b) }
func BenchmarkHandleGtfs10000(b *testing.B) { benchmarkHandleGtfs(9999, false, b) }

func BenchmarkHandleGtfsWithSeed10(b *testing.B)    { benchmarkHandleGtfs(9, true, b) }
func BenchmarkHandleGtfsWithSeed50(b *testing.B)    { benchmarkHandleGtfs(49, true, b) }
func BenchmarkHandleGtfsWithSeed100(b *testing.B)   { benchmarkHandleGtfs(99, true, b) }
func BenchmarkHandleGtfsWithSeed1000(b *testing.B)  { benchmarkHandleGtfs(999, true, b) }
func BenchmarkHandleGtfsWithSeed10000(b *testing.B) { benchmarkHandleGtfs(9999, true, b) }
