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

func Test_TransactionalModel_StopVisits(t *testing.T) {
	model := NewMemoryModel()
	existingStopVisit := model.StopVisits().New()
	model.StopVisits().Save(&existingStopVisit)
	stopVisitId := existingStopVisit.Id()

	transactionnalModel := NewTransactionalModel(model)

	if _, ok := transactionnalModel.StopVisits().Find(stopVisitId); !ok {
		t.Errorf("TransactionalModel should have same StopVisits as parent model")
	}
}

func Test_TransactionalModel_Commit(t *testing.T) {
	model := NewMemoryModel()
	transactionnalModel := NewTransactionalModel(model)

	stopArea := transactionnalModel.StopAreas().New()
	stopVisit := transactionnalModel.StopVisits().New()

	success := transactionnalModel.StopAreas().Save(&stopArea) && transactionnalModel.StopVisits().Save(&stopVisit)
	if !success {
		t.Errorf("Save should return true")
	}
	if stopArea.Id() == "" {
		t.Errorf("New StopArea identifier shouldn't be an empty string")
	}
	if stopVisit.Id() == "" {
		t.Errorf("New StopVisit identifier shouldn't be an empty string")
	}

	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea shouldn't be saved before commit")
	}
	if _, ok := model.StopVisits().Find(stopVisit.Id()); ok {
		t.Errorf("StopVisit shouldn't be saved before commit")
	}

	transactionnalModel.Commit()
	if _, ok := model.StopAreas().Find(stopArea.Id()); !ok {
		t.Errorf("StopArea should be saved after commit")
	}
	if _, ok := model.StopVisits().Find(stopVisit.Id()); !ok {
		t.Errorf("StopVisit should be saved after commit")
	}
}

func Test_TransactionalModel_Rollback(t *testing.T) {
	model := NewMemoryModel()
	transactionnalModel := NewTransactionalModel(model)

	stopArea := transactionnalModel.StopAreas().New()
	stopVisit := transactionnalModel.StopVisits().New()

	transactionnalModel.Rollback()
	transactionnalModel.Commit()

	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea shouldn't be saved with a rollback")
	}
	if _, ok := model.StopVisits().Find(stopVisit.Id()); ok {
		t.Errorf("StopVisit shouldn't be saved with a rollback")
	}
}
