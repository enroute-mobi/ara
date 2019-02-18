package model

import "testing"

func Test_ObjectIDIndex_simple(t *testing.T) {
	index := NewObjectIdIndex()

	objectid := NewObjectID("kind", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.objectids = make(ObjectIDs)
	stopVisit.SetObjectID(objectid)

	index.Index(ModelId(stopVisit.Id()), stopVisit)

	foundStopVisit, ok := index.Find(objectid)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}
}

func Test_ObjectIDIndex_Multiple(t *testing.T) {
	index := NewObjectIdIndex()

	objectid := NewObjectID("kind", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.objectids = make(ObjectIDs)
	stopVisit.SetObjectID(objectid)

	index.Index(ModelId(stopVisit.Id()), stopVisit)

	objectid2 := NewObjectID("kind", "value2")
	stopVisit2 := &StopVisit{id: "stopVisitId2"}
	stopVisit2.objectids = make(ObjectIDs)
	stopVisit2.SetObjectID(objectid2)

	index.Index(ModelId(stopVisit2.Id()), stopVisit2)

	foundStopVisit, ok := index.Find(objectid)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}

	foundStopVisit, ok = index.Find(objectid2)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit2.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit2.id)
	}
}

func Test_ObjectIDIndex_Change(t *testing.T) {
	index := NewObjectIdIndex()

	objectid := NewObjectID("kind", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.objectids = make(ObjectIDs)
	stopVisit.SetObjectID(objectid)

	index.Index(ModelId(stopVisit.Id()), stopVisit)

	objectid2 := NewObjectID("kind", "value2")
	stopVisit.SetObjectID(objectid2)
	index.Index(ModelId(stopVisit.Id()), stopVisit)

	_, ok := index.Find(objectid)
	if ok {
		t.Error("Can find StopVisit after changing index: ", index)
	}
	foundStopVisit, ok := index.Find(objectid2)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}
}

func Test_ObjectIDIndex_Delete(t *testing.T) {
	index := NewObjectIdIndex()

	objectid := NewObjectID("kind", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.objectids = make(ObjectIDs)
	stopVisit.SetObjectID(objectid)

	index.Index(ModelId(stopVisit.Id()), stopVisit)
	index.Delete("stopVisitId")

	_, ok := index.Find(objectid)
	if ok {
		t.Error("Can find StopVisit after delete: ", index)
	}
}

var benchmarkObjectIDResult StopVisit
var benchmarkObjectIDResultId ModelId

func benchmarkObjectIDFindWithoutIndex(sv int, b *testing.B) {
	model := NewMemoryModel()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		objectid := NewObjectID("kind", DefaultUUIDGenerator().NewUUID())
		stopVisit.SetObjectID(objectid)
		stopVisit.Save()
	}
	stopVisit := model.StopVisits().New()
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
	stopVisit.Save()

	var foundStopVisit StopVisit
	for n := 0; n < b.N; n++ {
		foundStopVisit, _ = model.StopVisits().FindByObjectId(objectid)
	}
	benchmarkObjectIDResult = foundStopVisit
}

func benchmarkObjectIDFindWithIndex(sv int, b *testing.B) {
	model := NewMemoryModel()
	index := NewObjectIdIndex()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		objectid := NewObjectID("kind", DefaultUUIDGenerator().NewUUID())
		stopVisit.SetObjectID(objectid)
		stopVisit.Save()
		index.Index(ModelId(stopVisit.Id()), &stopVisit)
	}
	stopVisit := model.StopVisits().New()
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
	stopVisit.Save()
	index.Index(ModelId(stopVisit.Id()), &stopVisit)

	var foundStopVisit ModelId
	for n := 0; n < b.N; n++ {
		foundStopVisit, _ = index.Find(objectid)
	}
	benchmarkObjectIDResultId = foundStopVisit
}

func BenchmarkObjectIDFindWithIndex10(b *testing.B)   { benchmarkObjectIDFindWithIndex(9, b) }
func BenchmarkObjectIDFindWithIndex50(b *testing.B)   { benchmarkObjectIDFindWithIndex(49, b) }
func BenchmarkObjectIDFindWithIndex100(b *testing.B)  { benchmarkObjectIDFindWithIndex(99, b) }
func BenchmarkObjectIDFindWithIndex1000(b *testing.B) { benchmarkObjectIDFindWithIndex(999, b) }

func BenchmarkObjectIDFindWithoutIndex10(b *testing.B)   { benchmarkObjectIDFindWithoutIndex(9, b) }
func BenchmarkObjectIDFindWithoutIndex50(b *testing.B)   { benchmarkObjectIDFindWithoutIndex(49, b) }
func BenchmarkObjectIDFindWithoutIndex100(b *testing.B)  { benchmarkObjectIDFindWithoutIndex(99, b) }
func BenchmarkObjectIDFindWithoutIndex1000(b *testing.B) { benchmarkObjectIDFindWithoutIndex(999, b) }
