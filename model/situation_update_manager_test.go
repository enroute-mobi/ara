package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func completeEvent(code Code, testTime time.Time) (event *SituationUpdateEvent) {
	period := &TimeRange{EndTime: testTime}

	event = &SituationUpdateEvent{
		RecordedAt:        testTime,
		SituationCode: code,
		Version:           1,
		ProducerRef:       "Ara",
		ValidityPeriods:   []*TimeRange{period},
		Keywords:          []string{"channel"},
	}

	summary := &SituationTranslatedString{
		DefaultValue: "Message Text",
	}

	event.Summary = summary
	event.Format = "format"

	return
}

func checkSituation(situation Situation, code Code, testTime time.Time) bool {
	summary := &SituationTranslatedString{
		DefaultValue: "Message Text",
	}

	period := &TimeRange{EndTime: testTime}

	testSituation := Situation{
		id:              situation.id,
		Summary:         summary,
		RecordedAt:      testTime,
		ValidityPeriods: []*TimeRange{period},
		Format:          "format",
		Keywords:        []string{"channel"},
		ProducerRef:     "Ara",
		Version:         1,
	}
	testSituation.model = situation.model
	testSituation.codes = make(Codes)
	testSituation.SetCode(code)
	testSituation.SetCode(NewCode("_default", code.HashValue()))

	return reflect.DeepEqual(situation, testSituation)
}

func Test_SituationUpdateManager_Update(t *testing.T) {
	assert := assert.New(t)
	code := NewCode("codeSpace", "value")
	testTime := time.Now()

	model := NewMemoryModel()
	situation := model.Situations().New()
	situation.SetCode(code)
	situation.SetCode(NewCode("_default", code.HashValue()))
	model.Situations().Save(&situation)

	manager := newSituationUpdateManager(model)
	event := completeEvent(code, testTime)
	manager.Update([]*SituationUpdateEvent{event})

	updatedSituation, _ := model.Situations().Find(situation.Id())
	assert.True(checkSituation(updatedSituation, code, testTime))
}

func Test_SituationUpdateManager_SameRecordedAt(t *testing.T) {
	assert := assert.New(t)
	code := NewCode("codeSpace", "value")
	testTime := time.Now()

	model := NewMemoryModel()
	situation := model.Situations().New()
	situation.SetCode(code)
	situation.RecordedAt = testTime
	situation.Keywords = []string{"situationChannel"}
	model.Situations().Save(&situation)

	manager := newSituationUpdateManager(model)
	event := completeEvent(code, testTime)
	manager.Update([]*SituationUpdateEvent{event})

	updatedSituation, _ := model.Situations().Find(situation.Id())
	assert.False(checkSituation(updatedSituation, code, testTime), "Situation should not be updated")
	assert.ElementsMatch(updatedSituation.Keywords, []string{"situationChannel"})
}

func Test_SituationUpdateManager_CreateSituation(t *testing.T) {
	assert := assert.New(t)
	code := NewCode("codeSpace", "value")
	testTime := time.Now()

	model := NewMemoryModel()

	manager := newSituationUpdateManager(model)
	event := completeEvent(code, testTime)
	manager.Update([]*SituationUpdateEvent{event})

	situations := model.Situations().FindAll()
	assert.Len(situations, 1, "Should find 1 situation")
	situation := situations[0]
	assert.True(checkSituation(situation, code, testTime))
}
