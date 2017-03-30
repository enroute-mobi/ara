package model

import "testing"

func Test_TransactionalSituation_Find_NotFound(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	_, ok := situations.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when situations isn't found")
	}
}

func Test_TransactionalSituations_Find_Model(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	existingSituation := model.Situations().New()
	model.Situations().Save(&existingSituation)

	situationId := existingSituation.Id()

	situation, ok := situations.Find(situationId)
	if !ok {
		t.Errorf("Find should return true when Situation is found")
	}
	if situation.Id() != situationId {
		t.Errorf("Find should return a Situation with the given Id")
	}
}

func Test_TransactionalSituations_Find_Saved(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	existingSituation := situations.New()
	situations.Save(&existingSituation)

	situationId := existingSituation.Id()

	situation, ok := situations.Find(situationId)
	if !ok {
		t.Errorf("Find should return true when Situation is found")
	}
	if situation.Id() != situationId {
		t.Errorf("Find should return a Situation with the given Id")
	}
}

func Test_TransactionSituations_FindAll(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	for i := 0; i < 5; i++ {
		existingSituation := situations.New()
		situations.Save(&existingSituation)
	}

	foundSituations := situations.FindAll()

	if len(foundSituations) != 5 {
		t.Errorf("FindAll should return all Situations")
	}

}

func Test_TransactionalSituations_Save(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	situation := situations.New()
	objectid := NewObjectID("kind", "value")
	situation.SetObjectID(objectid)

	if success := situations.Save(&situation); !success {
		t.Errorf("Save should return true")
	}
	if situation.Id() == "" {
		t.Errorf("New situation identifier shouldn't be an empty string")
	}
	if _, ok := model.Situations().Find(situation.Id()); ok {
		t.Errorf("situation shouldn't be saved before commit")
	}
}

func Test_TransactionalSituations_Delete(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	existingSituation := model.Situations().New()
	model.Situations().Save(&existingSituation)

	situations.Delete(&existingSituation)

	_, ok := situations.Find(existingSituation.Id())
	if !ok {
		t.Errorf("Situation should not be deleted before commit")
	}
}

func Test_TransactionalSituations_Commit(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	// Test Save
	situation := situations.New()
	situations.Save(&situation)

	// Test Delete
	existingSituation := model.Situations().New()
	model.Situations().Save(&existingSituation)
	situations.Delete(&existingSituation)

	situations.Commit()

	if _, ok := model.Situations().Find(situation.Id()); !ok {
		t.Errorf("Situation should be saved after commit")
	}
	if _, ok := situations.Find(existingSituation.Id()); ok {
		t.Errorf("Situation should be deleted after commit")
	}
}

func Test_TransactionalSituations_Rollback(t *testing.T) {
	model := NewMemoryModel()
	situations := NewTransactionalSituations(model)

	situation := situations.New()
	situations.Save(&situation)

	situations.Rollback()
	situations.Commit()

	if _, ok := model.Situations().Find(situation.Id()); ok {
		t.Errorf("Situation should not be saved with a rollback")
	}
}
