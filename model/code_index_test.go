package model

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_CodeIndex_simple(t *testing.T) {
	index := NewCodeIndex()

	code := NewCode("codeSpace", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.codes = make(Codes)
	stopVisit.SetCode(code)

	index.Index(stopVisit)

	foundStopVisit, ok := index.Find(code)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}
}

func Test_CodeIndex_MultipleIndex(t *testing.T) {
	index := NewCodeIndex()

	code := NewCode("codeSpace", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.codes = make(Codes)
	stopVisit.SetCode(code)

	index.Index(stopVisit)
	index.Index(stopVisit)

	foundStopVisit, ok := index.Find(code)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}
}

func Test_CodeIndex_Multiple(t *testing.T) {
	index := NewCodeIndex()

	code := NewCode("codeSpace", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.codes = make(Codes)
	stopVisit.SetCode(code)

	index.Index(stopVisit)

	code2 := NewCode("codeSpace", "value2")
	stopVisit2 := &StopVisit{id: "stopVisitId2"}
	stopVisit2.codes = make(Codes)
	stopVisit2.SetCode(code2)

	index.Index(stopVisit2)

	foundStopVisit, ok := index.Find(code)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}

	foundStopVisit, ok = index.Find(code2)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit2.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit2.id)
	}
}

func Test_CodeIndex_Change(t *testing.T) {
	index := NewCodeIndex()

	code := NewCode("codeSpace", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.codes = make(Codes)
	stopVisit.SetCode(code)

	index.Index(stopVisit)

	code2 := NewCode("codeSpace", "value2")
	stopVisit.SetCode(code2)
	index.Index(stopVisit)

	_, ok := index.Find(code)
	if ok {
		t.Error("Can find StopVisit after changing index: ", index)
	}
	foundStopVisit, ok := index.Find(code2)
	if !ok {
		t.Error("Can't find StopVisit after index: ", index)
	}
	if StopVisitId(foundStopVisit) != stopVisit.id {
		t.Errorf("Wrong Id returned, got: %v want: %v", foundStopVisit, stopVisit.id)
	}
}

func Test_CodeIndex_Delete(t *testing.T) {
	index := NewCodeIndex()

	code := NewCode("codeSpace", "value")
	stopVisit := &StopVisit{id: "stopVisitId"}
	stopVisit.codes = make(Codes)
	stopVisit.SetCode(code)

	index.Index(stopVisit)
	index.Delete("stopVisitId")

	_, ok := index.Find(code)
	if ok {
		t.Error("Can find StopVisit after delete: ", index)
	}
}

var benchmarkCodeResult *StopVisit
var benchmarkCodeResultId ModelId

func benchmarkCodeFindWithoutIndex(sv int, b *testing.B) {
	model := NewMemoryModel()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		code := NewCode("codeSpace", uuid.DefaultUUIDGenerator().NewUUID())
		stopVisit.SetCode(code)
		stopVisit.Save()
	}
	stopVisit := model.StopVisits().New()
	code := NewCode("codeSpace", "value")
	stopVisit.SetCode(code)
	stopVisit.Save()

	var foundStopVisit *StopVisit
	for n := 0; n < b.N; n++ {
		foundStopVisit, _ = model.StopVisits().FindByCode(code)
	}
	benchmarkCodeResult = foundStopVisit
}

func benchmarkCodeFindWithIndex(sv int, b *testing.B) {
	model := NewMemoryModel()
	index := NewCodeIndex()

	for i := 0; i < sv; i++ {
		stopVisit := model.StopVisits().New()
		code := NewCode("codeSpace", uuid.DefaultUUIDGenerator().NewUUID())
		stopVisit.SetCode(code)
		stopVisit.Save()
		index.Index(stopVisit)
	}
	stopVisit := model.StopVisits().New()
	code := NewCode("codeSpace", "value")
	stopVisit.SetCode(code)
	stopVisit.Save()
	index.Index(stopVisit)

	var foundStopVisit ModelId
	for n := 0; n < b.N; n++ {
		foundStopVisit, _ = index.Find(code)
	}
	benchmarkCodeResultId = foundStopVisit
}

func BenchmarkCodeFindWithIndex10(b *testing.B)   { benchmarkCodeFindWithIndex(9, b) }
func BenchmarkCodeFindWithIndex50(b *testing.B)   { benchmarkCodeFindWithIndex(49, b) }
func BenchmarkCodeFindWithIndex100(b *testing.B)  { benchmarkCodeFindWithIndex(99, b) }
func BenchmarkCodeFindWithIndex1000(b *testing.B) { benchmarkCodeFindWithIndex(999, b) }

func BenchmarkCodeFindWithoutIndex10(b *testing.B)   { benchmarkCodeFindWithoutIndex(9, b) }
func BenchmarkCodeFindWithoutIndex50(b *testing.B)   { benchmarkCodeFindWithoutIndex(49, b) }
func BenchmarkCodeFindWithoutIndex100(b *testing.B)  { benchmarkCodeFindWithoutIndex(99, b) }
func BenchmarkCodeFindWithoutIndex1000(b *testing.B) { benchmarkCodeFindWithoutIndex(999, b) }
