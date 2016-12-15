package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_StopVisit_Id(t *testing.T) {
	stopVisit := StopVisit{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if stopVisit.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("StopVisit.Id() returns wrong value, got: %s, required: %s", stopVisit.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

// WIP: Determine what to return in JSON
// func Test_StopVisit_MarshalJSON(t *testing.T) {
// 	stopVisit := StopVisit{
// 		id:   "6ba7b814-9dad-11d1-0-00c04fd430c8",
// 		Name: "Test",
// 	}
// 	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Name":"Test"}`
// 	jsonBytes, err := stopVisit.MarshalJSON()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	jsonString := string(jsonBytes)
// 	if jsonString != expected {
// 		t.Errorf("StopVisit.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
// 	}
// }

func Test_StopVisit_UnmarshalJSON(t *testing.T) {
	text := `{
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" }
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
}

func Test_StopVisit_Save(t *testing.T) {
	model := NewMemoryModel()
	stopVisit := model.StopVisits().New()
	objectid := NewObjectID("kind", "value")
	stopVisit.SetObjectID(objectid)

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

	if success := stopVisits.Save(&stopVisit); !success {
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
	stopVisits.Save(&existingStopVisit)

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
		stopVisits.Save(&existingStopVisit)
	}

	foundStopVisits := stopVisits.FindAll()

	if len(foundStopVisits) != 5 {
		t.Errorf("FindAll should return all stopVisits")
	}
}

func Test_MemoryStopVisits_Delete(t *testing.T) {
	stopVisits := NewMemoryStopVisits()
	existingStopVisit := stopVisits.New()
	objectid := NewObjectID("kind", "value")
	existingStopVisit.SetObjectID(objectid)
	stopVisits.Save(&existingStopVisit)

	stopVisits.Delete(&existingStopVisit)

	_, ok := stopVisits.Find(existingStopVisit.Id())
	if ok {
		t.Errorf("Deleted StopVisit should not be findable")
	}
	_, ok = stopVisits.FindByObjectId(objectid)
	if ok {
		t.Errorf("New StopVisit should not be findable by objectid")
	}
}
