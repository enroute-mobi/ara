package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func completeEvent(objectid ObjectID, testTime time.Time) (event *SituationUpdateEvent) {
	period := &TimeRange{EndTime: testTime}

	event = &SituationUpdateEvent{
		RecordedAt:        testTime,
		SituationObjectID: objectid,
		Version:           1,
		ProducerRef:       "Ara",
		ValidityPeriods:   []*TimeRange{period},
		Keywords:          []string{"channel"},
	}

	summary := &SituationTranslatedString{
		DefaultValue: "Message Text",
	}

	event.Summary = summary
	event.SituationAttributes = SituationAttributes{
		Format: "format",
	}

	return
}

func checkSituation(situation Situation, objectid ObjectID, testTime time.Time) bool {
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
	testSituation.objectids = make(ObjectIDs)
	testSituation.SetObjectID(objectid)
	testSituation.SetObjectID(NewObjectID("_default", objectid.HashValue()))

	return reflect.DeepEqual(situation, testSituation)
}

func Test_SituationUpdateManager_Update(t *testing.T) {
	assert := assert.New(t)
	objectid := NewObjectID("kind", "value")
	testTime := time.Now()

	model := NewMemoryModel()
	situation := model.Situations().New()
	situation.SetObjectID(objectid)
	situation.SetObjectID(NewObjectID("_default", objectid.HashValue()))
	model.Situations().Save(&situation)

	manager := newSituationUpdateManager(model)
	event := completeEvent(objectid, testTime)
	manager.Update([]*SituationUpdateEvent{event})

	updatedSituation, _ := model.Situations().Find(situation.Id())
	assert.True(checkSituation(updatedSituation, objectid, testTime))
}

func Test_SituationUpdateManager_SameRecordedAt(t *testing.T) {
	assert := assert.New(t)
	objectid := NewObjectID("kind", "value")
	testTime := time.Now()

	model := NewMemoryModel()
	situation := model.Situations().New()
	situation.SetObjectID(objectid)
	situation.RecordedAt = testTime
	situation.Keywords = []string{"situationChannel"}
	model.Situations().Save(&situation)

	manager := newSituationUpdateManager(model)
	event := completeEvent(objectid, testTime)
	manager.Update([]*SituationUpdateEvent{event})

	updatedSituation, _ := model.Situations().Find(situation.Id())
	assert.False(checkSituation(updatedSituation, objectid, testTime), "Situation should not be updated")
	assert.ElementsMatch(updatedSituation.Keywords, []string{"situationChannel"})
}

func Test_SituationUpdateManager_CreateSituation(t *testing.T) {
	assert := assert.New(t)
	objectid := NewObjectID("kind", "value")
	testTime := time.Now()

	model := NewMemoryModel()

	manager := newSituationUpdateManager(model)
	event := completeEvent(objectid, testTime)
	manager.Update([]*SituationUpdateEvent{event})

	situations := model.Situations().FindAll()
	assert.Len(situations, 1, "Should find 1 situation")
	situation := situations[0]
	assert.True(checkSituation(situation, objectid, testTime))
}
