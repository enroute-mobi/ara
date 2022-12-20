package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_Situation_Id(t *testing.T) {
	situation := Situation{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if situation.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("situation.Id() returns wrong value, got: %s, required: %s", situation.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_Situation_MarshalJSON(t *testing.T) {
	situation := Situation{
		id:     "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Origin: "test",
	}
	situation.Messages = append(situation.Messages, &Message{Content: "Joyeux Noel", Type: "Un Type"})
	expected := `{"Origin":"test","Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Messages":[{"MessageText":"Joyeux Noel","MessageType":"Un Type"}]}`
	jsonBytes, err := situation.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Situation.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Situation_UnmarshalJSON(t *testing.T) {
	text := `{
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" }
  }`

	situation := Situation{}
	err := json.Unmarshal([]byte(text), &situation)
	if err != nil {
		t.Fatal(err)
	}

	expectedObjectIds := []ObjectID{
		NewObjectID("reflex", "FR:77491:ZDE:34004:STIF"),
		NewObjectID("hastus", "sqypis"),
	}

	for _, expectedObjectId := range expectedObjectIds {
		objectId, found := situation.ObjectID(expectedObjectId.Kind())
		if !found {
			t.Errorf("Missing situation ObjectId '%s' after UnmarshalJSON()", expectedObjectId.Kind())
		}
		if !reflect.DeepEqual(expectedObjectId, objectId) {
			t.Errorf("Wrong situation ObjectId after UnmarshalJSON():\n got: %s\n want: %s", objectId, expectedObjectId)
		}
	}
}

func Test_Situation_Save(t *testing.T) {
	model := NewMemoryModel()
	situation := model.Situations().New()
	objectid := NewObjectID("kind", "value")
	situation.SetObjectID(objectid)

	if situation.model != model {
		t.Errorf("New situation model should be MemorySituation model")
	}

	ok := situation.Save()
	if !ok {
		t.Errorf("situation.Save() should succeed")
	}
	_, ok = model.Situations().Find(situation.Id())
	if !ok {
		t.Errorf("New situation should be found in MemorySituation")
	}
}

func Test_Situation_ObjectId(t *testing.T) {
	situation := Situation{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	situation.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	situation.SetObjectID(objectid)

	foundObjectId, ok := situation.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = situation.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(situation.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", situation.ObjectIDs())
	}
}

func Test_MemorySituations_New(t *testing.T) {
	situations := NewMemorySituations()

	situation := situations.New()
	if situation.Id() != "" {
		t.Errorf("New situation identifier should be an empty string, got: %s", situation.Id())
	}
}

func Test_MemorySituations_Save(t *testing.T) {
	situations := NewMemorySituations()

	situation := situations.New()

	if success := situations.Save(&situation); !success {
		t.Errorf("Save should return true")
	}

	if situation.Id() == "" {
		t.Errorf("New situation identifier shouldn't be an empty string")
	}
}

func Test_MemorySituations_Find_NotFound(t *testing.T) {
	situations := NewMemorySituations()
	_, ok := situations.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Situation isn't found")
	}
}

func Test_MemorySituations_Find(t *testing.T) {
	situations := NewMemorySituations()

	existingSituation := situations.New()
	situations.Save(&existingSituation)

	situationId := existingSituation.Id()

	situation, ok := situations.Find(situationId)
	if !ok {
		t.Errorf("Find should return true when situation is found")
	}
	if situation.Id() != situationId {
		t.Errorf("Find should return a situation with the given Id")
	}
}

func Test_MemorySituations_FindAll(t *testing.T) {
	situations := NewMemorySituations()

	for i := 0; i < 5; i++ {
		existingSituation := situations.New()
		situations.Save(&existingSituation)
	}

	foundSituations := situations.FindAll()

	if len(foundSituations) != 5 {
		t.Errorf("FindAll should return all situations")
	}
}

func Test_MemorySituations_Delete(t *testing.T) {
	situations := NewMemorySituations()
	existingSituation := situations.New()
	objectid := NewObjectID("kind", "value")
	existingSituation.SetObjectID(objectid)
	situations.Save(&existingSituation)

	situations.Delete(&existingSituation)

	_, ok := situations.Find(existingSituation.Id())
	if ok {
		t.Errorf("Deleted situation should not be findable")
	}
}
