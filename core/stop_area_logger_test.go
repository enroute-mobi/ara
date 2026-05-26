package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_NewStopAreaLogger(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	referential.SetSetting("logger.verbose.stop_areas", "internal:value")

	stopArea := referential.Model().StopAreas().New()

	code := model.NewCode("internal", "value")
	stopArea.SetCode(code)

	logger := NewStopAreaLogger(referential, stopArea)
	assert.True(logger.IsVerbose(), "StopAreaLogger should be in verbose")
}

func Test_NewStopAreaLogger_WithMultipleCodes(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	referential.SetSetting("logger.verbose.stop_areas", "internal:value")

	stopArea := referential.Model().StopAreas().New()

	code := model.NewCode("internal", "value")
	stopArea.SetCode(code)

	stopArea.SetCode(model.NewCode("external", "value"))

	logger := NewStopAreaLogger(referential, stopArea)
	assert.True(logger.IsVerbose(), "StopAreaLogger should be in verbose")
}

func Test_NewStopAreaLogger_NoMatch(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	referential.SetSetting("logger.debug.stop_areas", "internal:value")

	stopArea := referential.Model().StopAreas().New()

	code := model.NewCode("external", "value")
	stopArea.SetCode(code)

	logger := NewStopAreaLogger(referential, stopArea)
	assert.False(logger.IsVerbose(), "StopAreaLogger should not be in verbose")
}
