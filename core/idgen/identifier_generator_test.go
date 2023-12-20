package idgen

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_IdentifierGenerator_NewIdentifier(t *testing.T) {
	generator := NewIdentifierGenerator("%{type}:%{uuid}:%{default}:%{id}:%{uuid}:%{code}", uuid.NewFakeUUIDGenerator())
	attributes := IdentifierAttributes{
		Id:   "iD",
		Type: "Tp",
	}
	idf := generator.NewIdentifier(attributes)
	if expected := "Tp:6ba7b814-9dad-11d1-0-00c04fd430c8:iD:iD:6ba7b814-9dad-11d1-1-00c04fd430c8:iD"; idf != expected {
		t.Errorf("Identifier should be %v, got: %v", expected, idf)
	}
}
