package model

import (
	"testing"
)

func Test_StopArea_Id(t *testing.T) {
	stopArea := StopArea{
		id: 42,
	}

	if stopArea.Id() != 42 {
		t.Errorf("StopArea.Id() returns wrong value, got: %v, required: %v", stopArea.Id(), 42)
	}
}

func Test_MemoryStopAreas_New(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	stopArea := stopAreas.New()
	if stopArea.Id() != 0 {
		t.Errorf("New StopArea identifier should be zero, got: %d", stopArea.Id())
	}
}

func Test_MemoryStopAreas_Save(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	stopArea := stopAreas.New()

	if success := stopAreas.Save(&stopArea); !success {
		t.Errorf("Save should return true")
	}

	if stopArea.Id() == 0 {
		t.Errorf("New StopArea identifier shouldn't be zero")
	}
}

func Test_MemoryStopAreas_Save_NextIdentifier(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	for expectedIdentifier := 1; expectedIdentifier < 10; expectedIdentifier++ {
		stopArea := stopAreas.New()
		stopAreas.Save(&stopArea)

		if stopArea.Id() != StopAreaId(expectedIdentifier) {
			t.Errorf("New StopArea identifier should be %v, got: %v", expectedIdentifier, stopArea.Id())
		}
	}
}

func Test_MemoryStopAreas_Find_NotFound(t *testing.T) {
	stopAreas := NewMemoryStopAreas()
	_, ok := stopAreas.Find(1)
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

func Test_MemoryStopAreas_Delete(t *testing.T) {
	stopAreas := NewMemoryStopAreas()

	existingStopArea := stopAreas.New()
	stopAreas.Save(&existingStopArea)

	stopAreas.Delete(&existingStopArea)

	_, ok := stopAreas.Find(existingStopArea.Id())
	if ok {
		t.Errorf("Deleted StopArea should be findable")
	}
}
