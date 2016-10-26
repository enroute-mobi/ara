package model

import "testing"

func Test_NewTransaction(t *testing.T) {
	model := NewMemoryModel()
	transaction := NewTransaction(model)

	if transaction.Status() != PENDING {
		t.Errorf("Transaction status should be PENDING when created, got: %v", transaction.Status())
	}
}

func Test_Transaction_Commit(t *testing.T) {
	model := NewMemoryModel()
	transaction := NewTransaction(model)

	stopArea := transaction.Model().StopAreas().New()
	if success := transaction.Model().StopAreas().Save(&stopArea); !success {
		t.Errorf("Save should return true")
	}

	err := transaction.Commit()
	if err != nil {
		t.Errorf("Commit should not return errors")
	}
	if transaction.Status() != COMMIT {
		t.Errorf("Transaction status should be COMMIT after Commit, got: %v", transaction.Status())
	}
	if _, ok := model.StopAreas().Find(stopArea.Id()); !ok {
		t.Errorf("StopArea should be saved in the model after Commit")
	}
}

func Test_Transaction_Rollback(t *testing.T) {
	model := NewMemoryModel()
	transaction := NewTransaction(model)

	stopArea := transaction.Model().StopAreas().New()
	if success := transaction.Model().StopAreas().Save(&stopArea); !success {
		t.Errorf("Save should return true")
	}

	err := transaction.Rollback()
	if err != nil {
		t.Errorf("Rollback should not return errors")
	}
	if transaction.Status() != ROLLBACK {
		t.Errorf("Transaction status should be ROLLBACK after Rollback, got: %v", transaction.Status())
	}

	err = transaction.Commit()
	if err != nil {
		t.Errorf("Commit should not return errors")
	}
	if _, ok := model.StopAreas().Find(stopArea.Id()); ok {
		t.Errorf("StopArea should not be saved in the model after Rollback")
	}
}

func Test_Transaction_Close(t *testing.T) {
	model := NewMemoryModel()
	transaction := NewTransaction(model)

	err := transaction.Close()
	if err != nil {
		t.Errorf("Close should not return errors")
	}
	if transaction.Status() != ROLLBACK {
		t.Errorf("When closing transaction, status should be ROLLBACK when previous was PENDING, got: %v", transaction.Status())
	}

	err = transaction.Close()
	if err != nil {
		t.Errorf("Close should not return errors")
	}
	if transaction.Status() != ROLLBACK {
		t.Errorf("When closing transaction, status should be ROLLBACK when previous was ROLLBACK, got: %v", transaction.Status())
	}

	err = transaction.Commit()
	if err != nil {
		t.Errorf("Commit should not return errors")
	}
	err = transaction.Close()
	if err != nil {
		t.Errorf("Close should not return errors")
	}
	if transaction.Status() != COMMIT {
		t.Errorf("When closing transaction, status should be COMMIT when previous was COMMIT, got: %v", transaction.Status())
	}
}
