package model

import "testing"

func Test_TransactionalModel_StopAreas(t *testing.T) {
	model := NewMemoryModel()
	existingStopArea := model.StopAreas().New()
	model.StopAreas().Save(&existingStopArea)
	stopAreaId := existingStopArea.Id()

	transactionnalModel := NewTransactionalModel(model)

	if _, ok := transactionnalModel.StopAreas().Find(stopAreaId); !ok {
		t.Errorf("TransactionalModel should have same StopAreas as parent model")
	}
}

func Test_TransactionalModel_Commit(t *testing.T) {
	model := NewMemoryModel()
	transactionnalModel := NewTransactionalModel(model)

	stopArea := transactionnalModel.StopAreas().New()

	if success := transactionnalModel.StopAreas().Save(&stopArea); !success {
		t.Errorf("Save should return true")
	}
	if stopArea.Id() == "" {
		t.Errorf("New StopArea identifier shouldn't be an empty string")
	}

	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea shouldn't be saved before commit")
	}

	transactionnalModel.Commit()
	if _, ok := model.StopAreas().Find(stopArea.Id()); !ok {
		t.Errorf("StopArea should be saved before commit")
	}
}

func Test_TransactionalModel_Rollback(t *testing.T) {
	model := NewMemoryModel()
	transactionnalModel := NewTransactionalModel(model)

	stopArea := transactionnalModel.StopAreas().New()

	transactionnalModel.Rollback()
	transactionnalModel.Commit()

	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea shouldn't be saved with a rollback")
	}
}
