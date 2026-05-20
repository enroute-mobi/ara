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
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	// seed unused models
	if seed {
		for i := 0; i != 500; i++ {
			vehicleJourney := referential.model.VehicleJourneys().New()
			vjId := model.NewCode("codeSpace", fmt.Sprintf("vj%d", i))
			vehicleJourney.SetCode(vjId)
			vehicleJourney.LineId = line.Id()
			vehicleJourney.Save()

			// stopAreas & stopVisits
			for j := 0; j < 30; j++ {
				saId := model.NewCode("codeSpace", fmt.Sprintf("saId%d", j))
				stopArea := referential.Model().StopAreas().New()
				stopArea.SetCode(saId)
				stopArea.Save()

				stopVisit := referential.model.StopVisits().New()
				svId1 := model.NewCode("codeSpace", fmt.Sprintf("svId%d", j))
				stopVisit.SetCode(svId1)
				stopVisit.StopAreaId = stopArea.Id()
				stopVisit.VehicleJourneyId = vehicleJourney.Id()

				stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(time.Duration(float64(j*1e9))-10*time.Minute))
				stopVisit.PassageOrder = j
				stopVisit.Save()
			}
		}
	}

	// models for benchmark
	for i := 0; i != pc; i++ {
		vehicleJourney := referential.model.VehicleJourneys().New()
		vjId := model.NewCode("codeSpace", fmt.Sprintf("vj%d", i))
		vehicleJourney.SetCode(vjId)
		vehicleJourney.LineId = line.Id()
		vehicleJourney.Save()

		// stopAreas & stopVisits
		for j := 0; j < 30; j++ {
			saId := model.NewCode("codeSpace", fmt.Sprintf("saId%d", j))
			stopArea := referential.Model().StopAreas().New()
			stopArea.SetCode(saId)
			stopArea.Save()

			stopVisit := referential.model.StopVisits().New()
			svId1 := model.NewCode("codeSpace", fmt.Sprintf("svId%d", j))
			stopVisit.SetCode(svId1)
			stopVisit.StopAreaId = stopArea.Id()
			stopVisit.VehicleJourneyId = vehicleJourney.Id()

			stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(time.Duration(float64(j*1e9))-10*time.Minute))
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
