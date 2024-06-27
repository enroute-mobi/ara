package model

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"github.com/stretchr/testify/assert"
)

func Test_StopVisit_Id(t *testing.T) {
	stopVisit := StopVisit{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if stopVisit.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("StopVisit.Id() returns wrong value, got: %s, required: %s", stopVisit.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_StopVisit_MarshalJSON(t *testing.T) {
	stopVisit := StopVisit{
		id:        "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Schedules: schedules.NewStopVisitSchedules(),
	}
	expected := `{"Origin":"","VehicleAtStop":false,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Collected":false}`
	jsonBytes, err := stopVisit.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("StopVisit.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_StopVisit_UnmarshalJSON(t *testing.T) {
	text := `{
    "Codes": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
    "StopAreaId": "6ba7b814-9dad-11d1-1-00c04fd430c8",
    "VehicleJourneyId": "6ba7b814-9dad-11d1-2-00c04fd430c8",
    "PassageOrder": 10
  }`

	stopVisit := StopVisit{}
	err := json.Unmarshal([]byte(text), &stopVisit)
	if err != nil {
		t.Fatal(err)
	}

	expectedCodes := []Code{
		NewCode("reflex", "FR:77491:ZDE:34004:STIF"),
		NewCode("hastus", "sqypis"),
	}

	for _, expectedCode := range expectedCodes {
		code, found := stopVisit.Code(expectedCode.CodeSpace())
		if !found {
			t.Errorf("Missing StopVisit Code '%s' after UnmarshalJSON()", expectedCode.CodeSpace())
		}
		if !reflect.DeepEqual(expectedCode, code) {
			t.Errorf("Wrong StopVisit Code after UnmarshalJSON():\n got: %s\n want: %s", code, expectedCode)
		}
	}

	if expected := StopAreaId("6ba7b814-9dad-11d1-1-00c04fd430c8"); stopVisit.StopAreaId != expected {
		t.Errorf("Wrong StopVisit StopAreaId:\n got: %s\n want: %s", stopVisit.StopAreaId, expected)
	}

	if expected := VehicleJourneyId("6ba7b814-9dad-11d1-2-00c04fd430c8"); stopVisit.VehicleJourneyId != expected {
		t.Errorf("Wrong StopVisit VehicleJourneyId:\n got: %s\n want: %s", stopVisit.VehicleJourneyId, expected)
	}

	if expected := 10; stopVisit.PassageOrder != expected {
		t.Errorf("Wrong StopVisit PassageOrder:\n got: %v\n want: %v", stopVisit.PassageOrder, expected)
	}
}

func Test_StopVisit_Save(t *testing.T) {
	model := NewTestMemoryModel()
	stopVisit := model.StopVisits().New()
	code := NewCode("codeSpace", "value")
	stopVisit.SetCode(code)
	stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-0-00c04fd430c8"

	if stopVisit.model != model {
		t.Errorf("New stopVisit model should be memoryStopVisits model")
	}

	ok := stopVisit.Save()
	if !ok {
		t.Errorf("stopVisit.Save() should succeed")
	}
	_, ok = model.StopVisits().Find(stopVisit.Id())
	if !ok {
		t.Errorf("New StopVisit should be found in memoryStopVisits")
	}
	_, ok = model.StopVisits().FindByCode(code)
	if !ok {
		t.Errorf("New StopVisit should be found by code in memoryStopVisits")
	}
	foundStopVisits := model.StopVisits().FindByVehicleJourneyId("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if len(foundStopVisits) == 0 || foundStopVisits[0].Id() != stopVisit.id {
		t.Errorf("New StopVisit should be found by vehicleJourneyId in memoryStopVisits")
	}
}

func Test_StopVisit_Code(t *testing.T) {
	stopVisit := StopVisit{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	stopVisit.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	stopVisit.SetCode(code)

	foundCode, ok := stopVisit.Code("codeSpace")
	if !ok {
		t.Errorf("Code should return true if Code exists")
	}
	if foundCode.Value() != code.Value() {
		t.Errorf("Code should return a correct Code:\n got: %v\n want: %v", foundCode, code)
	}

	_, ok = stopVisit.Code("wrongkind")
	if ok {
		t.Errorf("Code should return false if Code doesn't exist")
	}

	if len(stopVisit.Codes()) != 1 {
		t.Errorf("Codes should return an array with set Codes, got: %v", stopVisit.Codes())
	}
}

func Test_MemoryStopVisits_New(t *testing.T) {
	stopVisits := NewMemoryStopVisits()

	stopVisit := stopVisits.New()
	if stopVisit.Id() != "" {
		t.Errorf("New StopVisit identifier should be an empty string, got: %s", stopVisit.Id())
	}
}

func Test_MemoryStopVisits_Save(t *testing.T) {
	stopVisits := NewMemoryStopVisits()

	stopVisit := stopVisits.New()

	if success := stopVisits.Save(stopVisit); !success {
		t.Errorf("Save should return true")
	}

	if stopVisit.Id() == "" {
		t.Errorf("New StopVisit identifier shouldn't be an empty string")
	}
}

func Test_MemoryStopVisits_Find_NotFound(t *testing.T) {
	stopVisits := NewMemoryStopVisits()
	_, ok := stopVisits.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when StopVisit isn't found")
	}
}

func Test_MemoryStopVisits_Find(t *testing.T) {
	stopVisits := NewMemoryStopVisits()

	existingStopVisit := stopVisits.New()
	stopVisits.Save(existingStopVisit)

	stopVisitId := existingStopVisit.Id()

	stopVisit, ok := stopVisits.Find(stopVisitId)
	if !ok {
		t.Errorf("Find should return true when StopVisit is found")
	}
	if stopVisit.Id() != stopVisitId {
		t.Errorf("Find should return a StopVisit with the given Id")
	}
}

func Test_MemoryStopVisits_FindAll(t *testing.T) {
	stopVisits := NewMemoryStopVisits()

	for i := 0; i < 5; i++ {
		existingStopVisit := stopVisits.New()
		stopVisits.Save(existingStopVisit)
	}

	foundStopVisits := stopVisits.FindAll()

	if len(foundStopVisits) != 5 {
		t.Errorf("FindAll should return all stopVisits")
	}
}

func Test_MemoryStopVisits_FindAllAfter(t *testing.T) {
	stopVisits := NewMemoryStopVisits()

	for i := 0; i < 5; i++ {
		sv := stopVisits.New()
		sv.Schedules.SetArrivalTime(schedules.Actual, time.Now().Add(-time.Duration(i)*time.Minute))
		stopVisits.Save(sv)
	}

	foundStopVisits := stopVisits.FindAllAfter(time.Now().Add(-150 * time.Second))

	if len(foundStopVisits) != 3 {
		t.Errorf("FindAll should return 3 stopVisits: %v", stopVisits.byIdentifier)
	}
}

func Test_MemoryStopVisits_Delete(t *testing.T) {
	stopVisits := NewMemoryStopVisits()
	existingStopVisit := stopVisits.New()
	code := NewCode("codeSpace", "value")
	existingStopVisit.SetCode(code)
	stopVisits.Save(existingStopVisit)

	stopVisits.Delete(existingStopVisit)

	_, ok := stopVisits.Find(existingStopVisit.Id())
	if ok {
		t.Errorf("Deleted StopVisit should not be findable")
	}
	_, ok = stopVisits.FindByCode(code)
	if ok {
		t.Errorf("New StopVisit should not be findable by code")
	}
}

func Test_MemoryStopVisits_DeleteMultiple(t *testing.T) {
	assert := assert.New(t)

	stopVisits := NewMemoryStopVisits()
	stopVisit1 := stopVisits.New()
	stopVisit2 := stopVisits.New()
	code1 := NewCode("codeSpace", "value1")
	code2 := NewCode("codeSpace", "value2")
	stopVisit1.SetCode(code1)
	stopVisit2.SetCode(code2)
	stopVisits.Save(stopVisit1)
	stopVisits.Save(stopVisit2)

	toDelete := []*StopVisit{stopVisit1, stopVisit2}
	stopVisits.DeleteMultiple(toDelete)

	_, ok := stopVisits.Find(stopVisit1.Id())
	assert.False(ok, "Deleted StopVisit1 should not be findable")
	_, ok = stopVisits.Find(stopVisit2.Id())
	assert.False(ok, "Deleted StopVisit1 should not be findable")

	_, ok = stopVisits.FindByCode(code1)
	assert.False(ok, "New StopVisit should not be findable by code1")
	_, ok = stopVisits.FindByCode(code2)
	assert.False(ok, "New StopVisit should not be findable by code2")
}

func Test_MemoryStopVisits_Load(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	testTime := time.Now()
	// Insert Data in the test db
	databaseStopVisit := DatabaseStopVisit{
		Id:               "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug:  "referential",
		ModelName:        "2017-01-01",
		Codes:            `{"internal":"value"}`,
		StopAreaId:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		VehicleJourneyId: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Schedules:        `[{"CodeSpace":"expected","DepartureTime":"2017-08-17T10:45:55+02:00"}]`,
		PassageOrder:     1,
		Attributes:       "{}",
		References:       `{"Ref":{"Type":"Ref","Code":{"kind":"value"}}}`,
	}

	Database.AddTableWithName(databaseStopVisit, "stop_visits")
	err := Database.Insert(&databaseStopVisit)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	model := NewTestMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	stopVisits := model.StopVisits().(*MemoryStopVisits)
	err = stopVisits.Load("referential")
	if err != nil {
		t.Fatal(err)
	}

	stopVisitId := StopVisitId(databaseStopVisit.Id)
	stopVisit, ok := stopVisits.Find(stopVisitId)
	if !ok {
		t.Fatalf("Loaded StopVisits should be found")
	}

	if stopVisit.id != stopVisitId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", stopVisit.id, stopVisitId)
	}
	if code, ok := stopVisit.Code("internal"); !ok || code.Value() != "value" {
		t.Errorf("Wrong Code:\n got: %v:%v\n expected: \"internal\":\"value\"", code.CodeSpace(), code.Value())
	}
	if stopVisit.PassageOrder != 1 {
		t.Errorf("StopVisit has wrong PassageOrder, got: %v want: 1", stopVisit.PassageOrder)
	}
	if ref, ok := stopVisit.Reference("Ref"); !ok || ref.Type != "Ref" || ref.Code.CodeSpace() != "kind" || ref.Code.Value() != "value" {
		t.Errorf("Wrong References:\n got: %v\n expected Type: \"Ref\" and Code: \"codeSpace:value\"", ref)
	}
	svs := stopVisit.Schedules.Schedule("expected")
	if svs == nil {
		t.Fatal("StopVisit should have an 'expected' Schedule")
	}
	if svs.DepartureTime().Equal(testTime) {
		t.Errorf("StopVisitSchedule should have Departure time %v, got: %v", testTime, svs.DepartureTime())
	}
}

var r []byte //against Compilator optimisation

func benchmarkStopVisitsMarshal10000(sv int, b *testing.B) {
	stopVisits := NewMemoryStopVisits()

	for i := 0; i != sv; i++ {
		stopVisit := stopVisits.New()
		stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-0-00c04fd430c8"
		stopVisits.Save(stopVisit)
	}

	for n := 0; n < b.N; n++ {
		jsonBytes, _ := json.Marshal(stopVisits.FindAll())
		r = jsonBytes
	}
}

func BenchmarkStopVisitsMarshal10000(b *testing.B) { benchmarkStopVisitsMarshal10000(10000, b) }
