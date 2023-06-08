package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_NewStopAreaLogger(t *testing.T) {
	assert := assert.New(t)

	referential := referentials.New("default")
	referential.SetSetting("logger.verbose.stop_areas", "kind:value")

	memoryModel := model.NewMemoryModel()
	stopArea := memoryModel.StopAreas().New()

	objectid := model.NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)

	logger := NewStopAreaLogger(referential, stopArea)
	assert.True(logger.IsVerbose(), "StopAreaLogger should be in verbose")
}

func Test_NewStopAreaLogger_WithMultipleObjectIds(t *testing.T) {
	assert := assert.New(t)

	referential := referentials.New("default")
	referential.SetSetting("logger.verbose.stop_areas", "kind:value")

	memoryModel := model.NewMemoryModel()
	stopArea := memoryModel.StopAreas().New()

	objectid := model.NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)

	stopArea.SetObjectID(model.NewObjectID("second", "value"))

	logger := NewStopAreaLogger(referential, stopArea)
	assert.True(logger.IsVerbose(), "StopAreaLogger should be in verbose")
}

func Test_NewStopAreaLogger_NoMatch(t *testing.T) {
	assert := assert.New(t)

	referential := referentials.New("default")
	referential.SetSetting("logger.debug.stop_areas", "kind:value")

	memoryModel := model.NewMemoryModel()
	stopArea := memoryModel.StopAreas().New()

	objectid := model.NewObjectID("no", "match")
	stopArea.SetObjectID(objectid)

	logger := NewStopAreaLogger(referential, stopArea)
	assert.False(logger.IsVerbose(), "StopAreaLogger should not be in verbose")
}
