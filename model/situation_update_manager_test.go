package model

import "testing"

func Test_SituationUpdateManager_Update(t *testing.T) {
	model := NewMemoryModel()
	situation := model.Situations().New()
	objectid := NewObjectID("kind", "value")
	situation.SetObjectID(objectid)
	model.Situations().Save(&situation)
	manager := newSituationUpdateManager(model)

	event := &SituationUpdateEvent{}
	event.Version = 1
	event.SituationObjectID = objectid

	events := append([]*SituationUpdateEvent{}, event)

	manager.UpdateSituation(events)
	updatedSituation, _ := model.Situations().Find(situation.Id())
	if updatedSituation.Version != 1 {
		t.Errorf("Situation Version should be 1")
	}
}

func Test_SituationUpdateManager_CreateSituation(t *testing.T) {
	model := NewMemoryModel()
	tx := NewTransaction(model)

	defer tx.Close()

	event := &SituationUpdateEvent{}

	situationUpdater := NewSituationUpdater(tx, nil)
	situationUpdater.SetClock(NewFakeClock())
	situationUpdater.CreateSituationFromEvent(event)
	tx.Commit()

	nb := model.Situations().FindAll()
	if len(nb) != 1 {
		t.Errorf("Should find 1 situation %v", len(nb))
	}
}
