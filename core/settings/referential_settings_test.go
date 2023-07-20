package settings

import (
	"testing"
	"time"

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

func Test_ModelRefreshTime_Without_setting(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()

	duration := referentialSettings.ModelRefreshTime()
	assert.Equalf(time.Duration(50_000_000_000), duration, "should set at 50 seconds by default")
}

func Test_ModelRefreshTime_Below_30seconds(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.s["model.refresh_time"] = "10s"

	duration := referentialSettings.ModelRefreshTime()
	assert.Equalf(time.Duration(30_000_000_000), duration, "should set minium duration at 30 seconds")
}

func Test_ModelRefreshTime_Abov_30seconds(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.s["model.refresh_time"] = "40s"

	duration := referentialSettings.ModelRefreshTime()
	assert.Equalf(time.Duration(40_000_000_000), duration, "should set duration at 40 seconds")
}
