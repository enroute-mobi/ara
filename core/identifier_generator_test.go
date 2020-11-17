package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_IdentifierGenerator_NewIdentifier(t *testing.T) {
	generator := NewIdentifierGenerator("%{type}:%{uuid}:%{default}:%{id}:%{uuid}")
	generator.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	attributes := IdentifierAttributes{
		Default: "Df",
		Id:      "iD",
		Type:    "Tp",
	}
	idf := generator.NewIdentifier(attributes)
	if expected := "Tp:6ba7b814-9dad-11d1-0-00c04fd430c8:Df:iD:6ba7b814-9dad-11d1-1-00c04fd430c8"; idf != expected {
		t.Errorf("Identifier should be %v, got: %v", expected, idf)
	}
}

func Test_IdentifierGenerator_NewIdentifier_WithoutSubstitution(t *testing.T) {
	generator := NewIdentifierGenerator("%{objectid}")
	attributes := IdentifierAttributes{
		ObjectId: "unchanged",
	}
	identifier := generator.NewIdentifier(attributes)
	if expected := "unchanged"; identifier != expected {
		t.Errorf("Identifier should be %v, got: %v", expected, identifier)
	}
}

func Test_IdentifierGenerator_NewIdentifier_WithSubstitution(t *testing.T) {
	generator := NewIdentifierGenerator("%{objectid//pattern/replacement}")
	attributes := IdentifierAttributes{
		ObjectId: "before-pattern-between-pattern-after",
	}
	identifier := generator.NewIdentifier(attributes)
	if expected := "before-replacement-between-replacement-after"; identifier != expected {
		t.Errorf("Identifier should be %v, got: %v", expected, identifier)
	}
}
