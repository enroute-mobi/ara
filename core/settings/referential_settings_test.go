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

	codes := referentialSettings.LoggerVerboseStopAreas()

	assert.NotEmpty(codes)
	assert.Equal(1, len(codes), "Should return a single Code for the moment")

	code := codes[0]
	assert.Equal("stif", code.CodeSpace(), "Code kind should be 'stif'")
	assert.Equal("STIF:StopPoint:Q:2342:", code.Value(), "Code kind should be 'STIF:StopPoint:Q:2342:'")
}

func Test_ReferentialSettings_LoggerVerboseStopAreas_WithWrongValue(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.s["logger.verbose.stop_areas"] = "wrong"

	codes := referentialSettings.LoggerVerboseStopAreas()
	assert.Empty(codes)
}

func Test_ReferentialSettings_LoggerDebugStopAreas_WithoutValue(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()

	codes := referentialSettings.LoggerVerboseStopAreas()
	assert.Empty(codes)
}

func Test_ModelRefreshTime_Without_setting(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()

	duration := referentialSettings.ModelRefreshTime()
	assert.Equalf(50*time.Second, duration, "should set at 50 seconds by default")
}

func Test_ModelRefreshTime_Below_30seconds(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.s["model.refresh_time"] = "10s"

	duration := referentialSettings.ModelRefreshTime()
	assert.Equalf(30*time.Second, duration, "should set minium duration at 30 seconds")
}

func Test_ModelRefreshTime_Abov_30seconds(t *testing.T) {
	assert := assert.New(t)

	referentialSettings := NewReferentialSettings()
	referentialSettings.s["model.refresh_time"] = "40s"

	duration := referentialSettings.ModelRefreshTime()
	assert.Equalf(40*time.Second, duration, "should set duration at 40 seconds")
}
