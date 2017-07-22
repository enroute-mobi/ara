package core

import "testing"

func Test_IdentifierGenerator_NewIdentifier(t *testing.T) {
	generator := NewIdentifierGenerator("%{type}:%{default}:%{id}:%{uuid}")
	attributes := IdentifierAttributes{
		Default: "Df",
		Id:      "iD",
		Type:    "Tp",
		UUID:    "uI",
	}
	if idf := generator.NewIdentifier(attributes); idf != "Tp:Df:iD:uI" {
		t.Errorf("Identifier should be Tp:Df:iD:uI, got: %v", idf)
	}
}
