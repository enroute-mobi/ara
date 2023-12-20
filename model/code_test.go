package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_Code_CodeSpace(t *testing.T) {
	code := Code{
		codeSpace: "codeSpace",
	}

	if expected := "codeSpace"; code.CodeSpace() != expected {
		t.Errorf("Code.CodeSpace() returns wrong value, got: %s, required: %s", code.CodeSpace(), expected)
	}
}

func Test_Code_Value(t *testing.T) {
	code := Code{
		value: "value",
	}

	if expected := "value"; code.Value() != expected {
		t.Errorf("Code.Value() returns wrong value, got: %s, required: %s", code.Value(), expected)
	}
}

func Test_Code_String(t *testing.T) {
	code := Code{
		codeSpace: "codeSpace",
		value:     "value",
	}
	if expected := "codeSpace:value"; code.String() != expected {
		t.Errorf("Code.String() returns wrong value, got: %s, required: %s", code.String(), expected)
	}
}

func Test_NewCodesFromMap(t *testing.T) {
	idmap := map[string]string{
		"reflex": "FR:77491:ZDE:34004:STIF",
		"hastus": "sqypis",
	}
	identifiers := NewCodesFromMap(idmap)

	expectedIdentifiers := make(Codes)
	expectedIdentifiers["reflex"] = NewCode("reflex", "FR:77491:ZDE:34004:STIF")
	expectedIdentifiers["hastus"] = NewCode("hastus", "sqypis")

	if !reflect.DeepEqual(expectedIdentifiers, identifiers) {
		t.Errorf("Wrong unmarshalled identifers from %s\n want: %#v\n got: %#v", idmap, expectedIdentifiers, identifiers)
	}
}

func Test_Codes_UnmarshalJSON(t *testing.T) {
	text := `{ "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" }`
	identifiers := make(Codes)
	err := json.Unmarshal([]byte(text), &identifiers)
	if err != nil {
		t.Fatal(err)
	}

	expectedIdentifiers := make(Codes)
	expectedIdentifiers["reflex"] = NewCode("reflex", "FR:77491:ZDE:34004:STIF")
	expectedIdentifiers["hastus"] = NewCode("hastus", "sqypis")

	if !reflect.DeepEqual(expectedIdentifiers, identifiers) {
		t.Errorf("Wrong unmarshalled identifers from %s\n want: %#v\n got: %#v", text, expectedIdentifiers, identifiers)
	}
}

func Test_Code_ToSlice(t *testing.T) {
	m := map[string]string{
		"codeSpace":  "value",
		"codeSpace2": "value2",
	}
	objs := NewCodesFromMap(m)
	s := objs.ToSlice()
	if len(s) != 2 {
		t.Errorf("Wrong number of entries in code slice, want: 2 got: %v", len(s))
	}
	if s[0] != "codeSpace:value" && s[1] != "codeSpace:value" {
		t.Errorf("We should find 'kind:value' in result slice, got %v", s)
	}
}
