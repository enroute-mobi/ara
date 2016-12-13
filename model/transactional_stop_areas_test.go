package model

import (
	"testing"
)

func Test_TransactionalStopAreas_Find_NotFound(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	_, ok := stopAreas.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when StopArea isn't found")
	}
}

func Test_TransactionalStopAreas_Find_Model(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	existingStopArea := model.StopAreas().New()
	model.StopAreas().Save(&existingStopArea)

	stopAreaId := existingStopArea.Id()

	stopArea, ok := stopAreas.Find(stopAreaId)
	if !ok {
		t.Errorf("Find should return true when StopArea is found")
	}
	if stopArea.Id() != stopAreaId {
		t.Errorf("Find should return a StopArea with the given Id")
	}
}

func Test_TransactionalStopAreas_Find_Saved(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	existingStopArea := stopAreas.New()
	stopAreas.Save(&existingStopArea)

	stopAreaId := existingStopArea.Id()

	stopArea, ok := stopAreas.Find(stopAreaId)
	if !ok {
		t.Errorf("Find should return true when StopArea is found")
	}
	if stopArea.Id() != stopAreaId {
		t.Errorf("Find should return a StopArea with the given Id")
	}
}

func Test_TransactionalStopAreas_FindAll(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	for i := 0; i < 5; i++ {
		existingStopArea := stopAreas.New()
		stopAreas.Save(&existingStopArea)
	}

	foundStopAreas := stopAreas.FindAll()

	if len(foundStopAreas) != 5 {
		t.Errorf("FindAll should return all stopAreas")
	}
}

func Test_TransactionalStopAreas_Save(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	stopArea := stopAreas.New()
	objectid := NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)

	if success := stopAreas.Save(&stopArea); !success {
		t.Errorf("Save should return true")
	}
	if stopArea.Id() == "" {
		t.Errorf("New StopArea identifier shouldn't be an empty string")
	}
	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea shouldn't be saved before commit")
	}
	if _, ok := model.StopAreas().FindByObjectId(objectid); ok {
		t.Errorf("StopArea shouldn't be findable by objectid before commit")
	}
}

func Test_TransactionalStopAreas_Delete(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	existingStopArea := model.StopAreas().New()
	objectid := NewObjectID("kind", "value")
	existingStopArea.SetObjectID(objectid)
	model.StopAreas().Save(&existingStopArea)

	stopAreas.Delete(&existingStopArea)

	if _, ok := stopAreas.Find(existingStopArea.Id()); !ok {
		t.Errorf("StopArea should not be deleted before commit")
	}
	if _, ok := stopAreas.FindByObjectId(objectid); !ok {
		t.Errorf("StopArea should be findable by objectid before commit")
	}
}

func Test_TransactionalStopAreas_Commit(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	// Test Save
	stopArea := stopAreas.New()
	objectid := NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)
	stopAreas.Save(&stopArea)

	// Test Delete
	existingStopArea := model.StopAreas().New()
	secondObjectid := NewObjectID("kind2", "value")
	existingStopArea.SetObjectID(secondObjectid)
	model.StopAreas().Save(&existingStopArea)
	stopAreas.Delete(&existingStopArea)

	stopAreas.Commit()

	if _, ok := model.StopAreas().Find(stopArea.Id()); !ok {
		t.Errorf("StopArea should be saved after commit")
	}
	if _, ok := model.StopAreas().FindByObjectId(objectid); !ok {
		t.Errorf("StopArea should be findable by ObjectId after commit")
	}

	if _, ok := stopAreas.Find(existingStopArea.Id()); ok {
		t.Errorf("StopArea should be deleted after commit")
	}
	if _, ok := stopAreas.FindByObjectId(secondObjectid); ok {
		t.Errorf("StopArea shouldn't be findable by ObjectID after commit")
	}
}

func Test_TransactionalStopAreas_Rollback(t *testing.T) {
	model := NewMemoryModel()
	stopAreas := NewTransactionalStopAreas(model)

	stopArea := stopAreas.New()
	objectid := NewObjectID("kind", "value")
	stopArea.SetObjectID(objectid)
	stopAreas.Save(&stopArea)

	stopAreas.Rollback()
	stopAreas.Commit()

	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea should not be saved with a rollback")
	}
	if _, ok := model.StopAreas().FindByObjectId(objectid); ok {
		t.Errorf("StopArea should not be findable by ObjectID with a rollback")
	}
}
