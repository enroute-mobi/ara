package model

import (
	"testing"
)

func Test_TransactionalLines_Find_NotFound(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	_, ok := lines.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Line isn't found")
	}
}

func Test_TransactionalLines_Find_Model(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	existingLine := model.Lines().New()
	model.Lines().Save(&existingLine)

	lineId := existingLine.Id()

	line, ok := lines.Find(lineId)
	if !ok {
		t.Errorf("Find should return true when Line is found")
	}
	if line.Id() != lineId {
		t.Errorf("Find should return a Line with the given Id")
	}
}

func Test_TransactionalLines_Find_Saved(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	existingLine := lines.New()
	lines.Save(&existingLine)

	lineId := existingLine.Id()

	line, ok := lines.Find(lineId)
	if !ok {
		t.Errorf("Find should return true when Line is found")
	}
	if line.Id() != lineId {
		t.Errorf("Find should return a Line with the given Id")
	}
}

func Test_TransactionalLines_FindAll(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	for i := 0; i < 5; i++ {
		existingLine := lines.New()
		lines.Save(&existingLine)
	}

	foundLines := lines.FindAll()

	if len(foundLines) != 5 {
		t.Errorf("FindAll should return all lines")
	}
}

func Test_TransactionalLines_Save(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	line := lines.New()
	objectid := NewObjectID("kind", "value")
	line.SetObjectID(objectid)

	if success := lines.Save(&line); !success {
		t.Errorf("Save should return true")
	}
	if line.Id() == "" {
		t.Errorf("New Line identifier shouldn't be an empty string")
	}
	if _, ok := model.Lines().Find(line.Id()); ok {
		t.Errorf("Line shouldn't be saved before commit")
	}
	if _, ok := model.Lines().FindByObjectId(objectid); ok {
		t.Errorf("Line shouldn't be found by objectid before commit")
	}
}

func Test_TransactionalLines_Delete(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	existingLine := model.Lines().New()
	objectid := NewObjectID("kind", "value")
	existingLine.SetObjectID(objectid)
	model.Lines().Save(&existingLine)

	lines.Delete(&existingLine)

	_, ok := lines.Find(existingLine.Id())
	if !ok {
		t.Errorf("Line should not be deleted before commit")
	}
	_, ok = lines.FindByObjectId(objectid)
	if !ok {
		t.Errorf("Line should be found by objectid before commit")
	}
}

func Test_TransactionalLines_Commit(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	// Test Save
	line := lines.New()
	objectid := NewObjectID("kind", "value")
	line.SetObjectID(objectid)
	lines.Save(&line)

	// Test Delete
	existingLine := model.Lines().New()
	secondObjectid := NewObjectID("kind", "value2")
	existingLine.SetObjectID(secondObjectid)
	model.Lines().Save(&existingLine)
	lines.Delete(&existingLine)

	lines.Commit()

	if _, ok := model.Lines().Find(line.Id()); !ok {
		t.Errorf("Line should be saved after commit")
	}
	if _, ok := model.Lines().FindByObjectId(objectid); !ok {
		t.Errorf("Line should be found by objectid after commit")
	}

	if _, ok := lines.Find(existingLine.Id()); ok {
		t.Errorf("Line should be deleted after commit")
	}
	if _, ok := lines.FindByObjectId(secondObjectid); ok {
		t.Errorf("Line should not be foundable by objectid after commit")
	}
}

func Test_TransactionalLines_Rollback(t *testing.T) {
	model := NewMemoryModel()
	lines := NewTransactionalLines(model)

	line := lines.New()
	objectid := NewObjectID("kind", "value")
	line.SetObjectID(objectid)
	lines.Save(&line)

	lines.Rollback()
	lines.Commit()

	if _, ok := model.Lines().Find(line.Id()); ok {
		t.Errorf("Line should not be saved with a rollback")
	}
	if _, ok := model.Lines().FindByObjectId(objectid); ok {
		t.Errorf("Line should not be foundable by objectid with a rollback")
	}
}
