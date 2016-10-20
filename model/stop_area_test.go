package model

import (
	"testing"
)

func Test_StopArea_Id(t *testing.T) {
	stopArea := StopArea{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if stopArea.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("StopArea.Id() returns wrong value, got: %s, required: %s", stopArea.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_MemoryStopAreas_New(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	stopArea := stopAreas.New()
	if stopArea.Id() != "" {
		t.Errorf("New StopArea identifier should be an empty string, got: %s", stopArea.Id())
	}
}

func Test_MemoryStopAreas_Save(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	stopArea := stopAreas.New()

	if success := stopAreas.Save(&stopArea); !success {
		t.Errorf("Save should return true")
	}

	if stopArea.Id() == "" {
		t.Errorf("New StopArea identifier shouldn't be an empty string")
	}
}

func Test_MemoryStopAreas_Find_NotFound(t *testing.T) {
	stopAreas := NewMemoryStopAreas()
	_, ok := stopAreas.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when StopArea isn't found")
	}
}

func Test_MemoryStopAreas_Find(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

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

func Test_MemoryStopAreas_FindAll(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	for i := 0; i < 5; i++ {
		existingStopArea := stopAreas.New()
		stopAreas.Save(&existingStopArea)
	}

	foundStopAreas := stopAreas.FindAll()

	if len(foundStopAreas) != 5 {
		t.Errorf("FindAll should return all stopAreas")
	}
}

func Test_MemoryStopAreas_Delete(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	existingStopArea := stopAreas.New()
	stopAreas.Save(&existingStopArea)

	stopAreas.Delete(&existingStopArea)

	_, ok := stopAreas.Find(existingStopArea.Id())
	if ok {
		t.Errorf("Deleted StopArea should not be findable")
	}
}
