package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_StopArea_Id(t *testing.T) {
	stopArea := StopArea{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if stopArea.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("StopArea.Id() returns wrong value, got: %s, required: %s", stopArea.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_StopArea_MarshalJSON(t *testing.T) {
	stopArea := StopArea{
		id:   "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Name: "Test",
	}
	expected := `{"CollectedAlways":false,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Name":"Test"}`
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
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" }
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

	if success := stopAreas.Save(&stopArea); !success {
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
	stopAreas.Save(&existingStopArea)

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
		stopAreas.Save(&existingStopArea)
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
	stopAreas.Save(&existingStopArea)

	stopAreas.Delete(&existingStopArea)

	_, ok := stopAreas.Find(existingStopArea.Id())
	if ok {
		t.Errorf("Deleted StopArea should not be findable")
	}
	_, ok = stopAreas.FindByObjectId(objectid)
	if ok {
		t.Errorf("Deleted StopArea should not be findable by objectid")
	}
}
