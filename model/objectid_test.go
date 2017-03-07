package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_ObjectID_Kind(t *testing.T) {
	objectID := ObjectID{
		kind: "kind",
	}

	if expected := "kind"; objectID.Kind() != expected {
		t.Errorf("ObjectID.Kind() returns wrong value, got: %s, required: %s", objectID.Kind(), expected)
	}
}

func Test_ObjectID_Value(t *testing.T) {
	objectID := ObjectID{
		value: "value",
	}

	if expected := "value"; objectID.Value() != expected {
		t.Errorf("ObjectID.Value() returns wrong value, got: %s, required: %s", objectID.Value(), expected)
	}
}

func Test_NewObjectIDsFromMap(t *testing.T) {
	idmap := map[string]string{
		"reflex": "FR:77491:ZDE:34004:STIF",
		"hastus": "sqypis",
	}
	identifiers := NewObjectIDsFromMap(idmap)

	expectedIdentifiers := make(ObjectIDs)
	expectedIdentifiers["reflex"] = NewObjectID("reflex", "FR:77491:ZDE:34004:STIF")
	expectedIdentifiers["hastus"] = NewObjectID("hastus", "sqypis")

	if !reflect.DeepEqual(expectedIdentifiers, identifiers) {
		t.Errorf("Wrong unmarshalled identifers from %s\n want: %#v\n got: %#v", idmap, expectedIdentifiers, identifiers)
	}
}

func Test_ObjectIDs_UnmarshalJSON(t *testing.T) {
	text := `{ "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" }`
	identifiers := make(ObjectIDs)
	err := json.Unmarshal([]byte(text), &identifiers)
	if err != nil {
		t.Fatal(err)
	}

	expectedIdentifiers := make(ObjectIDs)
	expectedIdentifiers["reflex"] = NewObjectID("reflex", "FR:77491:ZDE:34004:STIF")
	expectedIdentifiers["hastus"] = NewObjectID("hastus", "sqypis")

	if !reflect.DeepEqual(expectedIdentifiers, identifiers) {
		t.Errorf("Wrong unmarshalled identifers from %s\n want: %#v\n got: %#v", text, expectedIdentifiers, identifiers)
	}
}
