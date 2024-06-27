package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_NewStopAreaLogger(t *testing.T) {
	assert := assert.New(t)

	referential := referentials.New("default")
	referential.SetSetting("logger.verbose.stop_areas", "codeSpace:value")

	memoryModel := model.NewTestMemoryModel()
	stopArea := memoryModel.StopAreas().New()

	code := model.NewCode("codeSpace", "value")
	stopArea.SetCode(code)

	logger := NewStopAreaLogger(referential, stopArea)
	assert.True(logger.IsVerbose(), "StopAreaLogger should be in verbose")
}

func Test_NewStopAreaLogger_WithMultipleCodes(t *testing.T) {
	assert := assert.New(t)

	referential := referentials.New("default")
	referential.SetSetting("logger.verbose.stop_areas", "codeSpace:value")

	memoryModel := model.NewTestMemoryModel()
	stopArea := memoryModel.StopAreas().New()

	code := model.NewCode("codeSpace", "value")
	stopArea.SetCode(code)

	stopArea.SetCode(model.NewCode("second", "value"))

	logger := NewStopAreaLogger(referential, stopArea)
	assert.True(logger.IsVerbose(), "StopAreaLogger should be in verbose")
}

func Test_NewStopAreaLogger_NoMatch(t *testing.T) {
	assert := assert.New(t)

	referential := referentials.New("default")
	referential.SetSetting("logger.debug.stop_areas", "codeSpace:value")

	memoryModel := model.NewTestMemoryModel()
	stopArea := memoryModel.StopAreas().New()

	code := model.NewCode("no", "match")
	stopArea.SetCode(code)

	logger := NewStopAreaLogger(referential, stopArea)
	assert.False(logger.IsVerbose(), "StopAreaLogger should not be in verbose")
}
