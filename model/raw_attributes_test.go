package model

import "testing"

func Test_NewRawAttributes(t *testing.T) {
	attributes := NewRawAttributes()

	if len(attributes) != 0 {
		t.Errorf("New attributes should be empty")
	}
}

func Test_RawAttributes_Set(t *testing.T) {
	attributes := NewRawAttributes()

	attributes.Set("key", "value")

	if len(attributes) != 1 {
		t.Errorf("RawAttributes should have one entry")
	}

	if attributes["key"] != "value" {
		t.Errorf("'key' should be associated to 'value'")
	}
}

func Test_RawAttributes_Set_IgnoreEmptyValues(t *testing.T) {
	attributes := NewRawAttributes()

	attributes.Set("key", "")

	if _, ok := attributes["key"]; ok {
		t.Errorf("'key' should not be associated")
	}
}
