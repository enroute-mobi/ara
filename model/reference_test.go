package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Reference_UnmarshalJSON_With_Code(t *testing.T) {
	assert := assert.New(t)
	test := []byte(`{"Type": "OperatorRef", "Code": { "internal": "test" } }`)

	reference := Reference{}
	err := json.Unmarshal(test, &reference)
	assert.Nil(err)
	assert.Equal("OperatorRef", reference.Type)
	assert.Equal("internal", reference.Code.CodeSpace())
	assert.Equal("test", reference.Code.Value())
}
