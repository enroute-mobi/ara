package model

import (
	"reflect"
	"testing"
	"time"
)

func completeEvent(objectid ObjectID, testTime time.Time) (event *SituationUpdateEvent) {
	event = &SituationUpdateEvent{
		RecordedAt:        testTime,
		SituationObjectID: objectid,
		Version:           1,
		ProducerRef:       "Edwig",
	}

	message := &Message{
		Content:             "Message Text",
		Type:                "MessageType",
		NumberOfLines:       2,
		NumberOfCharPerLine: 20,
	}

	event.SituationAttributes = SituationAttributes{
		Format:     "format",
		Channel:    "channel",
		Messages:   []*Message{message},
		ValidUntil: testTime,
	}
	event.SituationAttributes.References = append(event.SituationAttributes.References, &Reference{ObjectId: &objectid, Id: "id", Type: "type"})

	return
}

func checkSituation(situation Situation, objectid ObjectID, testTime time.Time) bool {
	message := &Message{
		Content:             "Message Text",
		Type:                "MessageType",
		NumberOfLines:       2,
		NumberOfCharPerLine: 20,
	}

	testSituation := Situation{
		id:          situation.id,
		Messages:    []*Message{message},
		RecordedAt:  testTime,
		ValidUntil:  testTime,
		Format:      "format",
		Channel:     "channel",
		ProducerRef: "Edwig",
		Version:     1,
	}
	testSituation.model = situation.model
	testSituation.objectids = make(ObjectIDs)
	testSituation.SetObjectID(objectid)
	testSituation.SetObjectID(NewObjectID("_default", objectid.HashValue()))
	testSituation.References = append(testSituation.References, &Reference{ObjectId: &objectid, Id: "id", Type: "type"})

	return reflect.DeepEqual(situation, testSituation)
}

func Test_SituationUpdateManager_Update(t *testing.T) {
	objectid := NewObjectID("kind", "value")
	testTime := time.Now()

	model := NewMemoryModel()
	situation := model.Situations().New()
	situation.SetObjectID(objectid)
	situation.SetObjectID(NewObjectID("_default", objectid.HashValue()))
	model.Situations().Save(&situation)

	manager := newSituationUpdateManager(model)
	event := completeEvent(objectid, testTime)
	manager.UpdateSituation([]*SituationUpdateEvent{event})

	updatedSituation, _ := model.Situations().Find(situation.Id())
	if !checkSituation(updatedSituation, objectid, testTime) {
		t.Errorf("Situation is not properly updated:\n got: %v\n event: %v", updatedSituation, event)
	}
}

func Test_SituationUpdateManager_SameVersion(t *testing.T) {
	objectid := NewObjectID("kind", "value")
	testTime := time.Now()

	model := NewMemoryModel()
	situation := model.Situations().New()
	situation.SetObjectID(objectid)
	situation.Version = 1
	situation.Channel = "SituationChannel"
	model.Situations().Save(&situation)

	manager := newSituationUpdateManager(model)
	event := completeEvent(objectid, testTime)
	manager.UpdateSituation([]*SituationUpdateEvent{event})

	updatedSituation, _ := model.Situations().Find(situation.Id())
	if checkSituation(updatedSituation, objectid, testTime) {
		t.Errorf("Situation should not be updated:\n got: %v\n event: %v", updatedSituation, event)
	}
	if updatedSituation.Channel != "SituationChannel" {
		t.Errorf("Situation Channel should not have been updated:\n got: %v\n want: SituationChannel", updatedSituation.Channel)
	}
}

func Test_SituationUpdateManager_CreateSituation(t *testing.T) {
	objectid := NewObjectID("kind", "value")
	testTime := time.Now()

	model := NewMemoryModel()

	manager := newSituationUpdateManager(model)
	event := completeEvent(objectid, testTime)
	manager.UpdateSituation([]*SituationUpdateEvent{event})

	situations := model.Situations().FindAll()
	if len(situations) != 1 {
		t.Fatalf("Should find 1 situation, got %v", len(situations))
	}

	situation := situations[0]
	if !checkSituation(situation, objectid, testTime) {
		t.Errorf("Situation is not properly created:\n got: %v\n event: %v", situation, event)
	}
}
