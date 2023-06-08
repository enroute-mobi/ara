package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReferentialSettings_LoggerVerboseStopAreas(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.SetSetting("logger.verbose.stop_areas", "stif:STIF:StopPoint:Q:2342:")

	objectIds := referentialSettings.LoggerVerboseStopAreas()

	assert.NotEmpty(objectIds)
	assert.Equal(1, len(objectIds), "Should return a single ObjectID for the moment")

	objectId := objectIds[0]
	assert.Equal("stif", objectId.Kind(), "ObjectId kind should be 'stif'")
	assert.Equal("STIF:StopPoint:Q:2342:", objectId.Value(), "ObjectId kind should be 'STIF:StopPoint:Q:2342:'")
}

func Test_ReferentialSettings_LoggerVerboseStopAreas_WithWrongValue(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.s["logger.verbose.stop_areas"] = "wrong"

	objectIds := referentialSettings.LoggerVerboseStopAreas()
	assert.Empty(objectIds)
}

func Test_ReferentialSettings_LoggerDebugStopAreas_WithoutValue(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()

	objectIds := referentialSettings.LoggerVerboseStopAreas()
	assert.Empty(objectIds)
}
