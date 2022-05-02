package model

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

func createTestIndex() *Index {
	extractor := func(instance ModelInstance) ModelId {
		stopVisit := (instance.(*StopVisit))
		return ModelId(stopVisit.VehicleJourneyId)
	}

	return NewIndex(extractor)
}

func Test_Index_simple(t *testing.T) {
	index := createTestIndex()

	stopVisit := &StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"}
	index.Index(stopVisit)

	foundStopVisits, ok := index.Find(ModelId("dummy"))
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if len(foundStopVisits) != 1 {
		t.Errorf("Wrong number of StopVisit found, got: %v wanted: 1. Find result: %v", len(foundStopVisits), foundStopVisits)
	}
}

func Test_Index_MultipleIndex(t *testing.T) {
	index := createTestIndex()

	stopVisit := &StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"}
	index.Index(stopVisit)
	index.Index(stopVisit)

	foundStopVisits, ok := index.Find(ModelId("dummy"))
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if len(foundStopVisits) != 1 {
		t.Errorf("Wrong number of StopVisit found, got: %v wanted: 1. Find result: %v", len(foundStopVisits), foundStopVisits)
	}
}

func Test_Index_Multiple(t *testing.T) {
	index := createTestIndex()

	stopVisit := &StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"}
	index.Index(stopVisit)

	stopVisit2 := &StopVisit{id: "stopVisitId2", VehicleJourneyId: "dummy"}
	index.Index(stopVisit2)

	foundStopVisits, ok := index.Find(ModelId("dummy"))
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if len(foundStopVisits) != 2 {
		t.Errorf("Wrong number of StopVisit found, got: %v wanted: 2. Find result: %v", len(foundStopVisits), foundStopVisits)
	}
}

func Test_Index_Change(t *testing.T) {
	index := createTestIndex()

	stopVisit := &StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"}
	index.Index(stopVisit)

	stopVisit.VehicleJourneyId = "dummy2"
	index.Index(stopVisit)

	_, ok := index.Find(ModelId("dummy"))
	if ok {
		t.Error("Can find StopVisit after changing index: ", index)
	}
	foundStopVisits, ok := index.Find(ModelId("dummy2"))
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if len(foundStopVisits) != 1 {
		t.Errorf("Wrong number of StopVisit found, got: %v wanted: 1. Find result: %v", len(foundStopVisits), foundStopVisits)
	}
}

func Test_Index_Delete(t *testing.T) {
	index := createTestIndex()

	stopVisit := &StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"}
	index.Index(stopVisit)
	index.Delete("stopVisitId")

	foundStopVisits, ok := index.Find(ModelId("dummy"))
	if ok {
		t.Error("Can find StopVisit after delete: ", index)
	}
	if len(foundStopVisits) != 0 {
		t.Errorf("Wrong number of StopVisit found, got: %v wanted: 0. Find result: %v", len(foundStopVisits), foundStopVisits)
	}
}

var benchmarkResult []*StopVisit
var benchmarkResultId []ModelId

func benchmarkFindWithoutIndex(sv int, b *testing.B) {
	model := NewMemoryModel()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		stopVisit.VehicleJourneyId = VehicleJourneyId(uuid.DefaultUUIDGenerator().NewUUID())
		stopVisit.Save()
	}
	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-0-00c04fd430c8"
	stopVisit.Save()

	var foundStopVisits []*StopVisit
	for n := 0; n < b.N; n++ {
		foundStopVisits = model.StopVisits().FindByVehicleJourneyId("6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
	benchmarkResult = foundStopVisits
}

func benchmarkFindWithIndex(sv int, b *testing.B) {
	model := NewMemoryModel()
	index := createTestIndex()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		stopVisit.VehicleJourneyId = VehicleJourneyId(uuid.DefaultUUIDGenerator().NewUUID())
		stopVisit.Save()
		index.Index(stopVisit)
	}
	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-0-00c04fd430c8"
	stopVisit.Save()
	index.Index(stopVisit)

	var foundStopVisits []ModelId
	for n := 0; n < b.N; n++ {
		foundStopVisits, _ = index.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
	benchmarkResultId = foundStopVisits
}

func benchmarkIndexing(sv int, b *testing.B) {
	model := NewMemoryModel()
	index := createTestIndex()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-0-00c04fd430c8"
		stopVisit.Save()
		index.Index(stopVisit)
	}
	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-0-00c04fd430c8"
	stopVisit.Save()
	index.Index(stopVisit)

	for n := 0; n < b.N; n++ {
		index.Index(stopVisit)
		index.Delete(ModelId(stopVisit.Id()))
	}
}

func BenchmarkFindWithIndex10(b *testing.B)   { benchmarkFindWithIndex(9, b) }
func BenchmarkFindWithIndex50(b *testing.B)   { benchmarkFindWithIndex(49, b) }
func BenchmarkFindWithIndex100(b *testing.B)  { benchmarkFindWithIndex(99, b) }
func BenchmarkFindWithIndex1000(b *testing.B) { benchmarkFindWithIndex(999, b) }

func BenchmarkFindWithoutIndex10(b *testing.B)   { benchmarkFindWithoutIndex(9, b) }
func BenchmarkFindWithoutIndex50(b *testing.B)   { benchmarkFindWithoutIndex(49, b) }
func BenchmarkFindWithoutIndex100(b *testing.B)  { benchmarkFindWithoutIndex(99, b) }
func BenchmarkFindWithoutIndex1000(b *testing.B) { benchmarkFindWithoutIndex(999, b) }

func BenchmarkIndexing10(b *testing.B)   { benchmarkIndexing(9, b) }
func BenchmarkIndexing50(b *testing.B)   { benchmarkIndexing(49, b) }
func BenchmarkIndexing100(b *testing.B)  { benchmarkIndexing(99, b) }
func BenchmarkIndexing1000(b *testing.B) { benchmarkIndexing(999, b) }
