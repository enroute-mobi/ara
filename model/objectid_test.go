package model

import "testing"

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
