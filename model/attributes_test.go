package model

import "testing"

func Test_NewAttributes(t *testing.T) {
	attributes := NewAttributes()

	if len(attributes) != 0 {
		t.Errorf("New attributes should be empty")
	}
}

func Test_Attributes_Set(t *testing.T) {
	attributes := NewAttributes()

	attributes.Set("key", "value")

	if len(attributes) != 1 {
		t.Errorf("Attributes should have one entry")
	}

	if attributes["key"] != "value" {
		t.Errorf("'key' should be associated to 'value'")
	}
}

func Test_Attributes_Set_IgnoreEmptyValues(t *testing.T) {
	attributes := NewAttributes()

	attributes.Set("key", "")

	if _, ok := attributes["key"]; ok {
		t.Errorf("'key' should not be associated")
	}
}
