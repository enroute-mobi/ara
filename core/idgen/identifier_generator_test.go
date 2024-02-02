package idgen

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_IdentifierGenerator_NewIdentifier(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		ReferenceSetting: "%{type}:%{uuid}:%{default}:%{id}:%{uuid}:%{code}",
	}
	generator := NewIdentifierGenerator(settings, uuid.NewFakeUUIDGenerator())
	attributes := IdentifierAttributes{
		Id:   "iD",
		Type: "Tp",
	}
	idf := generator.NewIdentifier(attributes)
	assert.Equal("Tp:6ba7b814-9dad-11d1-0-00c04fd430c8:iD:iD:6ba7b814-9dad-11d1-1-00c04fd430c8:iD", idf)
}
