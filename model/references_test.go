package model

import "testing"

func Test_NewReferences(t *testing.T) {

	references := NewReferences()

	if len(references) != 0 {
		t.Errorf("New references should be empty")
	}

}

func Test_References_Set(t *testing.T) {
	references := NewReferences()
	obj := NewObjectID("kind", "value")

	reference := Reference{ObjectId: &obj, Id: ""}
	references.Set("key", reference)

	if len(references) != 1 {
		t.Errorf("references should have one entry")
	}

	if references["key"] != reference {
		t.Errorf("'key' should be associated to 'reference'")
	}
}

func Test_References_Set_IgnoreEmptyValues(t *testing.T) {

	references := NewReferences()

	references.Set("key", Reference{})

	if _, ok := references["key"]; ok {
		t.Errorf("'key' should not be associated")
	}

}
