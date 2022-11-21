package model

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
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
		Schedules: NewStopVisitSchedules(),
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
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
    "StopAreaId": "6ba7b814-9dad-11d1-1-00c04fd430c8",
    "VehicleJourneyId": "6ba7b814-9dad-11d1-2-00c04fd430c8",
    "PassageOrder": 10
  }`

	stopVisit := StopVisit{}
	err := json.Unmarshal([]byte(text), &stopVisit)
	if err != nil {
		t.Fatal(err)
	}

	expectedObjectIds := []ObjectID{
		NewObjectID("reflex", "FR:77491:ZDE:34004:STIF"),
		NewObjectID("hastus", "sqypis"),
	}

	for _, expectedObjectId := range expectedObjectIds {
		objectId, found := stopVisit.ObjectID(expectedObjectId.Kind())
		if !found {
			t.Errorf("Missing StopVisit ObjectId '%s' after UnmarshalJSON()", expectedObjectId.Kind())
		}
		if !reflect.DeepEqual(expectedObjectId, objectId) {
			t.Errorf("Wrong StopVisit ObjectId after UnmarshalJSON():\n got: %s\n want: %s", objectId, expectedObjectId)
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
	model := NewMemoryModel()
	stopVisit := model.StopVisits().New()
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)
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
	_, ok = model.StopVisits().FindByObjectId(objectid)
	if !ok {
		t.Errorf("New StopVisit should be found by objectid in memoryStopVisits")
	}
	foundStopVisits := model.StopVisits().FindByVehicleJourneyId("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if len(foundStopVisits) == 0 || foundStopVisits[0].Id() != stopVisit.id {
		t.Errorf("New StopVisit should be found by vehicleJourneyId in memoryStopVisits")
	}
}

func Test_StopVisit_ObjectId(t *testing.T) {
	stopVisit := StopVisit{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	stopVisit.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)

	foundObjectId, ok := stopVisit.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = stopVisit.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(stopVisit.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", stopVisit.ObjectIDs())
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
		sv.Schedules.SetArrivalTime(STOP_VISIT_SCHEDULE_ACTUAL, time.Now().Add(-time.Duration(i)*time.Minute))
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
	objectid := NewObjectID("kind", "value")
	existingStopVisit.SetObjectID(objectid)
	stopVisits.Save(existingStopVisit)

	stopVisits.Delete(existingStopVisit)

	_, ok := stopVisits.Find(existingStopVisit.Id())
	if ok {
		t.Errorf("Deleted StopVisit should not be findable")
	}
	_, ok = stopVisits.FindByObjectId(objectid)
	if ok {
		t.Errorf("New StopVisit should not be findable by objectid")
	}
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
		ObjectIDs:        `{"internal":"value"}`,
		StopAreaId:       "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		VehicleJourneyId: "d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Schedules:        `[{"Kind":"expected","DepartureTime":"2017-08-17T10:45:55+02:00"}]`,
		PassageOrder:     1,
		Attributes:       "{}",
		References:       `{"Ref":{"Type":"Ref","ObjectId":{"kind":"value"}}}`,
	}

	Database.AddTableWithName(databaseStopVisit, "stop_visits")
	err := Database.Insert(&databaseStopVisit)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	model := NewMemoryModel()
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
	if objectid, ok := stopVisit.ObjectID("internal"); !ok || objectid.Value() != "value" {
		t.Errorf("Wrong ObjectID:\n got: %v:%v\n expected: \"internal\":\"value\"", objectid.Kind(), objectid.Value())
	}
	if stopVisit.PassageOrder != 1 {
		t.Errorf("StopVisit has wrong PassageOrder, got: %v want: 1", stopVisit.PassageOrder)
	}
	if ref, ok := stopVisit.Reference("Ref"); !ok || ref.Type != "Ref" || ref.ObjectId.Kind() != "kind" || ref.ObjectId.Value() != "value" {
		t.Errorf("Wrong References:\n got: %v\n expected Type: \"Ref\" and ObjectId: \"kind:value\"", ref)
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
