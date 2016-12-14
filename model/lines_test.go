package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_Line_Id(t *testing.T) {
	line := Line{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if line.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("Line.Id() returns wrong value, got: %s, required: %s", line.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

// WIP: Determine what to return in JSON
func Test_Line_MarshalJSON(t *testing.T) {
	line := Line{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`
	jsonBytes, err := line.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Line.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Line_UnmarshalJSON(t *testing.T) {
	text := `{
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" }
  }`

	line := Line{}
	err := json.Unmarshal([]byte(text), &line)
	if err != nil {
		t.Fatal(err)
	}

	expectedObjectIds := []ObjectID{
		NewObjectID("reflex", "FR:77491:ZDE:34004:STIF"),
		NewObjectID("hastus", "sqypis"),
	}

	for _, expectedObjectId := range expectedObjectIds {
		objectId, found := line.ObjectID(expectedObjectId.Kind())
		if !found {
			t.Errorf("Missing Line ObjectId '%s' after UnmarshalJSON()", expectedObjectId.Kind())
		}
		if !reflect.DeepEqual(expectedObjectId, objectId) {
			t.Errorf("Wrong Line ObjectId after UnmarshalJSON():\n got: %s\n want: %s", objectId, expectedObjectId)
		}
	}
}

func Test_Line_Save(t *testing.T) {
	model := NewMemoryModel()
	line := model.Lines().New()
	objectid := NewObjectID("kind", "value")
	line.SetObjectID(objectid)

	if line.model != model {
		t.Errorf("New line model should be memoryLines model")
	}

	ok := line.Save()
	if !ok {
		t.Errorf("line.Save() should succeed")
	}
	_, ok = model.Lines().Find(line.Id())
	if !ok {
		t.Errorf("New Line should be found in memoryLines")
	}
	_, ok = model.Lines().FindByObjectId(objectid)
	if !ok {
		t.Errorf("New Line should be found by objectid in memoryLines")
	}
}

func Test_Line_ObjectId(t *testing.T) {
	line := Line{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	line.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	line.SetObjectID(objectid)

	foundObjectId, ok := line.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = line.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(line.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", line.ObjectIDs())
	}
}

func Test_MemoryLines_New(t *testing.T) {
	lines := NewMemoryLines()

	line := lines.New()
	if line.Id() != "" {
		t.Errorf("New Line identifier should be an empty string, got: %s", line.Id())
	}
}

func Test_MemoryLines_Save(t *testing.T) {
	lines := NewMemoryLines()

	line := lines.New()

	if success := lines.Save(&line); !success {
		t.Errorf("Save should return true")
	}

	if line.Id() == "" {
		t.Errorf("New Line identifier shouldn't be an empty string")
	}
}

func Test_MemoryLines_Find_NotFound(t *testing.T) {
	lines := NewMemoryLines()
	_, ok := lines.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Line isn't found")
	}
}

func Test_MemoryLines_Find(t *testing.T) {
	lines := NewMemoryLines()

	existingLine := lines.New()
	lines.Save(&existingLine)

	lineId := existingLine.Id()

	line, ok := lines.Find(lineId)
	if !ok {
		t.Errorf("Find should return true when Line is found")
	}
	if line.Id() != lineId {
		t.Errorf("Find should return a Line with the given Id")
	}
}

func Test_MemoryLines_FindAll(t *testing.T) {
	lines := NewMemoryLines()

	for i := 0; i < 5; i++ {
		existingLine := lines.New()
		lines.Save(&existingLine)
	}

	foundLines := lines.FindAll()

	if len(foundLines) != 5 {
		t.Errorf("FindAll should return all lines")
	}
}

func Test_MemoryLines_Delete(t *testing.T) {
	lines := NewMemoryLines()
	existingLine := lines.New()
	objectid := NewObjectID("kind", "value")
	existingLine.SetObjectID(objectid)
	lines.Save(&existingLine)

	lines.Delete(&existingLine)

	_, ok := lines.Find(existingLine.Id())
	if ok {
		t.Errorf("Deleted Line should not be findable")
	}
	_, ok = lines.FindByObjectId(objectid)
	if ok {
		t.Errorf("Deleted Line should not be findable by objectid")
	}
}
