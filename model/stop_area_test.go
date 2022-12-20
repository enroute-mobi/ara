package model

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
)

func Test_StopArea_Id(t *testing.T) {
	stopArea := StopArea{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if stopArea.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("StopArea.Id() returns wrong value, got: %s, required: %s", stopArea.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_StopArea_Lines(t *testing.T) {
	model := NewMemoryModel()
	line := model.Lines().New()
	line.Save()
	stopArea := StopArea{
		model:   model,
		LineIds: []LineId{line.Id()},
	}

	if len(stopArea.Lines()) != 1 || stopArea.Lines()[0].Id() != line.Id() {
		t.Errorf("StopArea.Lines() returns wrong values, got: %v, required: [%s]", stopArea.Lines(), line.Id())
	}
}

func Test_StopArea_MarshalJSON(t *testing.T) {
	stopArea := StopArea{
		id:      "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Name:    "Test",
		Origins: NewStopAreaOrigins(),
	}
	stopArea.Origins.NewOrigin("partnerTest")
	expected := `{"Origins":{"partnerTest":true},"Name":"Test","CollectChildren":false,"CollectGeneralMessages":false,"CollectedAlways":false,"Monitored":false,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`
	jsonBytes, err := stopArea.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("StopArea.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_StopArea_UnmarshalJSON(t *testing.T) {
	text := `{
    "Name":"Test",
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
    "Lines": ["1234","5678"]
  }`

	stopArea := StopArea{}
	err := json.Unmarshal([]byte(text), &stopArea)
	if err != nil {
		t.Fatal(err)
	}

	if expected := "Test"; stopArea.Name != expected {
		t.Errorf("Wrong StopArea Name after UnmarshalJSON():\n got: %s\n want: %s", stopArea.Name, expected)
	}

	expectedObjectIds := []ObjectID{
		NewObjectID("reflex", "FR:77491:ZDE:34004:STIF"),
		NewObjectID("hastus", "sqypis"),
	}

	for _, expectedObjectId := range expectedObjectIds {
		objectId, found := stopArea.ObjectID(expectedObjectId.Kind())
		if !found {
			t.Errorf("Missing StopArea ObjectId '%s' after UnmarshalJSON()", expectedObjectId.Kind())
		}
		if !reflect.DeepEqual(expectedObjectId, objectId) {
			t.Errorf("Wrong StopArea ObjectId after UnmarshalJSON():\n got: %s\n want: %s", objectId, expectedObjectId)
		}
	}

	if len(stopArea.LineIds) != 2 || stopArea.LineIds[0] != LineId("1234") || stopArea.LineIds[1] != LineId("5678") {
		t.Errorf("Wrong StopArea Lines:\n got: %v\n want: [1234,5678]", stopArea.LineIds)
	}
}

func Test_StopArea_Save(t *testing.T) {
	model := NewMemoryModel()
	stopArea := model.StopAreas().New()
	objectid := NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)

	if stopArea.model != model {
		t.Errorf("New stopArea model should be memoryStopAreas model")
	}

	stopArea.Name = "Chatelet"
	ok := stopArea.Save()
	if !ok {
		t.Errorf("stopArea.Save() should succeed")
	}
	_, ok = model.StopAreas().Find(stopArea.Id())
	if !ok {
		t.Errorf("New StopArea should be found in memoryStopAreas")
	}
	_, ok = model.StopAreas().FindByObjectId(objectid)
	if !ok {
		t.Errorf("New StopArea should be found by objectid in memoryStopAreas")
	}
}

func Test_StopArea_ObjectId(t *testing.T) {
	stopArea := StopArea{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	stopArea.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)

	foundObjectId, ok := stopArea.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = stopArea.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(stopArea.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", stopArea.ObjectIDs())
	}
}

func Test_MemoryStopAreas_New(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	stopArea := stopAreas.New()
	if stopArea.Id() != "" {
		t.Errorf("New StopArea identifier should be an empty string, got: %s", stopArea.Id())
	}
}

func Test_MemoryStopAreas_Save(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	stopArea := stopAreas.New()

	if success := stopAreas.Save(stopArea); !success {
		t.Errorf("Save should return true")
	}

	if stopArea.Id() == "" {
		t.Errorf("New StopArea identifier shouldn't be an empty string")
	}
}

func Test_MemoryStopAreas_Find_NotFound(t *testing.T) {
	stopAreas := NewMemoryStopAreas()
	_, ok := stopAreas.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when StopArea isn't found")
	}
}

func Test_MemoryStopAreas_Find(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	existingStopArea := stopAreas.New()
	stopAreas.Save(existingStopArea)

	stopAreaId := existingStopArea.Id()

	stopArea, ok := stopAreas.Find(stopAreaId)
	if !ok {
		t.Errorf("Find should return true when StopArea is found")
	}
	if stopArea.Id() != stopAreaId {
		t.Errorf("Find should return a StopArea with the given Id")
	}
}

func Test_MemoryStopAreas_FindAll(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	for i := 0; i < 5; i++ {
		existingStopArea := stopAreas.New()
		stopAreas.Save(existingStopArea)
	}

	foundStopAreas := stopAreas.FindAll()

	if len(foundStopAreas) != 5 {
		t.Errorf("FindAll should return all stopAreas")
	}
}

func Test_MemoryStopAreas_Delete(t *testing.T) {
	stopAreas := NewMemoryStopAreas()
	existingStopArea := stopAreas.New()
	objectid := NewObjectID("kind", "value")
	existingStopArea.SetObjectID(objectid)
	stopAreas.Save(existingStopArea)

	stopAreas.Delete(existingStopArea)

	_, ok := stopAreas.Find(existingStopArea.Id())
	if ok {
		t.Errorf("Deleted StopArea should not be findable")
	}
	_, ok = stopAreas.FindByObjectId(objectid)
	if ok {
		t.Errorf("Deleted StopArea should not be findable by objectid")
	}
}

func Test_MemoryStopAreas_FindAscendants(t *testing.T) {
	stopAreas := NewMemoryStopAreas()
	stopArea := stopAreas.New()
	stopAreas.Save(stopArea)

	stopArea1 := stopAreas.New()
	stopArea1.ParentId = stopArea.id
	stopAreas.Save(stopArea1)

	stopArea2 := stopAreas.New()
	stopArea2.ParentId = stopArea1.id
	stopAreas.Save(stopArea2)

	foundStopAreas := stopAreas.FindAscendants(stopArea2.Id())
	if len(foundStopAreas) != 3 {
		t.Errorf("FindAscendants should return 3, got %v", len(foundStopAreas))
	}
}

func Test_MemoryStopAreas_FindFamily(t *testing.T) {
	stopAreas := NewMemoryStopAreas()
	stopArea := stopAreas.New()
	stopAreas.Save(stopArea)

	stopArea1 := stopAreas.New()
	stopArea1.ParentId = stopArea.id
	stopAreas.Save(stopArea1)

	stopArea2 := stopAreas.New()
	stopArea2.ParentId = stopArea1.id
	stopAreas.Save(stopArea2)

	stopArea3 := stopAreas.New()
	stopArea3.ParentId = stopArea2.id
	stopAreas.Save(stopArea3)

	stopArea4 := stopAreas.New()
	stopArea4.ParentId = stopArea1.id
	stopAreas.Save(stopArea4)

	stopArea5 := stopAreas.New()
	stopArea5.ParentId = stopArea.id
	stopAreas.Save(stopArea5)

	stopArea6 := stopAreas.New()
	stopArea6.ParentId = stopArea5.id
	stopAreas.Save(stopArea6)

	stopArea7 := stopAreas.New()
	stopAreas.Save(stopArea7)

	if len(stopAreas.FindFamily(stopArea.id)) != 7 {
		t.Errorf("FindFamily should find 6 StopAreas, got: %v", len(stopAreas.FindFamily(stopArea.id)))
	}
	if len(stopAreas.FindFamily(stopArea1.id)) != 4 {
		t.Errorf("FindFamily should find 3 StopAreas, got: %v", len(stopAreas.FindFamily(stopArea1.id)))
	}
	if len(stopAreas.FindFamily(stopArea7.id)) != 1 {
		t.Errorf("FindFamily should find 0 StopAreas, got: %v", len(stopAreas.FindFamily(stopArea7.id)))
	}
}

func Test_MemoryStopAreas_Load(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// Insert Data in the test db
	databaseStopArea := DatabaseStopArea{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ParentId: sql.NullString{
			String: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Valid:  true,
		},
		ReferentId: sql.NullString{
			String: "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Valid:  true,
		},
		ModelName:       "2017-01-01",
		Name:            "stopArea",
		ObjectIDs:       `{"internal":"value"}`,
		LineIds:         `["d0eebc99-9c0b","e0eebc99-9c0b"]`,
		Attributes:      "{}",
		References:      `{"Ref":{"Type":"Ref","ObjectId":{"kind":"value"}}}`,
		CollectedAlways: true,
		CollectChildren: true,
	}

	Database.AddTableWithName(databaseStopArea, "stop_areas")
	err := Database.Insert(&databaseStopArea)
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
	stopAreas := model.StopAreas().(*MemoryStopAreas)
	err = stopAreas.Load("referential")
	if err != nil {
		t.Fatal(err)
	}

	stopAreaId := StopAreaId(databaseStopArea.Id)
	stopArea, ok := stopAreas.Find(stopAreaId)
	if !ok {
		t.Fatalf("Loaded StopAreas should be found")
	}

	if stopArea.id != stopAreaId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", stopArea.id, stopAreaId)
	}
	if stopArea.ParentId != StopAreaId("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") {
		t.Errorf("Wrong ParentId:\n got: %v\n expected: c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", stopArea.ParentId)
	}
	if stopArea.ReferentId != StopAreaId("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") {
		t.Errorf("Wrong ReferentId:\n got: %v\n expected: c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", stopArea.ReferentId)
	}
	if stopArea.Name != "stopArea" {
		t.Errorf("Wrong Name:\n got: %v\n expected: stopArea", stopArea.Name)
	}
	if objectid, ok := stopArea.ObjectID("internal"); !ok || objectid.Value() != "value" {
		t.Errorf("Wrong ObjectID:\n got: %v:%v\n expected: \"internal\":\"value\"", objectid.Kind(), objectid.Value())
	}
	if !stopArea.CollectedAlways {
		t.Errorf("Wrong CollectedAlways:\n got: %v\n expected: true", stopArea.CollectedAlways)
	}
	now := clock.DefaultClock().Now()
	if stopArea.nextCollectAt.Before(now) || stopArea.nextCollectAt.After(now.Add(30*time.Second)) {
		t.Errorf("Wrong nextCollectAt:\n got: %v\n expected: between %v and %v", stopArea.nextCollectAt, now, now.Add(30*time.Second))
	}
	if !stopArea.CollectChildren {
		t.Errorf("Wrong CollectChildren:\n got: %v\n expected: true", stopArea.CollectChildren)
	}
	if len(stopArea.LineIds) != 2 {
		t.Fatalf("StopArea should have 2 LineIds, got: %v", len(stopArea.LineIds))
	}
	if stopArea.LineIds[0] != "d0eebc99-9c0b" || stopArea.LineIds[1] != "e0eebc99-9c0b" {
		t.Errorf("Wrong LineIds:\n got: %v\n expected: [d0eebc99-9c0b,e0eebc99-9c0b]", stopArea.LineIds)
	}
	if ref, ok := stopArea.Reference("Ref"); !ok || ref.Type != "Ref" || ref.ObjectId.Kind() != "kind" || ref.ObjectId.Value() != "value" {
		t.Errorf("Wrong References:\n got: %v\n expected Type: \"Ref\" and ObjectId: \"kind:value\"", ref)
	}
}
