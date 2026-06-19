package model

import (
	"strconv"
	"testing"
)

// Vehicle index benchmarks. They use only the stable public API, so the same
// file can run against an older revision (the hand-rolled byNextStopVisitId map)
// and HEAD (the registered ByNextStopVisit OneToOne index) to compare the two.
//
// Run (the model package's TestMain needs the test config):
//
//	ARA_ENV=test ARA_CONFIG=$PWD/config ARA_ROOT=$PWD \
//	  go test -run='^$' -bench=BenchmarkVehicle -benchmem -count=6 ./model/
//
// Compare two revisions with benchstat:
//
//	# baseline revision (e.g. via a git worktree), capturing only the result lines:
//	ARA_ENV=test ... go test -run='^$' -bench=BenchmarkVehicle -benchmem -count=6 ./model/ \
//	  | grep '^Benchmark' > old.txt
//	# this revision:
//	ARA_ENV=test ... go test -run='^$' -bench=BenchmarkVehicle -benchmem -count=6 ./model/ \
//	  | grep '^Benchmark' > new.txt
//	benchstat old.txt new.txt
func benchPopulateVehicles(n int) *memoryVehicles {
	vehicles := NewMemoryVehicles().(*memoryVehicles)
	for i := 0; i < n; i++ {
		v := vehicles.New()
		v.id = VehicleId("veh-" + strconv.Itoa(i))
		v.LineId = LineId("line-" + strconv.Itoa(i%50))
		v.VehicleJourneyId = VehicleJourneyId("vj-" + strconv.Itoa(i))
		v.NextStopVisitId = StopVisitId("sv-" + strconv.Itoa(i))
		vehicles.Save(v)
	}
	return vehicles
}

// Save advancing the next stop each cycle — exercises the eviction/reindex path
// that differs between the hand-rolled byNextStopVisitId map and the registered
// ByNextStopVisit OneToOne index.
func BenchmarkVehicleSaveAdvanceNextStop(b *testing.B) {
	vehicles := NewMemoryVehicles().(*memoryVehicles)
	v := vehicles.New()
	v.id = VehicleId("veh-0")
	v.LineId = LineId("line-0")
	v.VehicleJourneyId = VehicleJourneyId("vj-0")
	v.NextStopVisitId = StopVisitId("sv-a")
	vehicles.Save(v)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i&1 == 0 {
			v.NextStopVisitId = StopVisitId("sv-b")
		} else {
			v.NextStopVisitId = StopVisitId("sv-a")
		}
		vehicles.Save(v)
	}
}

// Save of a steady vehicle (next stop unchanged) — the common per-collect case.
func BenchmarkVehicleSaveSteady(b *testing.B) {
	vehicles := NewMemoryVehicles().(*memoryVehicles)
	v := vehicles.New()
	v.id = VehicleId("veh-0")
	v.LineId = LineId("line-0")
	v.VehicleJourneyId = VehicleJourneyId("vj-0")
	v.NextStopVisitId = StopVisitId("sv-0")
	vehicles.Save(v)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vehicles.Save(v)
	}
}

// Insert path: a fresh vehicle each iteration (index grows).
func BenchmarkVehicleSaveInsert(b *testing.B) {
	vehicles := NewMemoryVehicles().(*memoryVehicles)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := vehicles.New()
		v.id = VehicleId("veh-" + strconv.Itoa(i))
		v.LineId = LineId("line-0")
		v.VehicleJourneyId = VehicleJourneyId("vj-" + strconv.Itoa(i))
		v.NextStopVisitId = StopVisitId("sv-" + strconv.Itoa(i))
		vehicles.Save(v)
	}
}

func BenchmarkVehicleFindByNextStopVisitId(b *testing.B) {
	vehicles := benchPopulateVehicles(1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vehicles.FindByNextStopVisitId(StopVisitId("sv-500"))
	}
}
