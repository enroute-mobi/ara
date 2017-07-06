package model

import (
	"testing"
	"time"
)

// type StopArea struct {
// 	ObjectIDConsumer
// 	model Model

// 	id       StopAreaId
// 	ParentId StopAreaId `json:",omitempty"`

// 	NextCollectAt   time.Time
// collectedAt     time.Time
// CollectedUntil  time.Time
// CollectedAlways bool

// Name       string
// LineIds    []LineId `json:"Lines,omitempty"`
// Attributes Attributes
// References References
// 	// ...
// }

func Test_Equal_StopAreas(t *testing.T) {
	model := NewMemoryModel()
	testTime := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)

	attributes := NewAttributes()
	attributes.Set("key", "value")

	references := NewReferences()
	obj := NewObjectID("kind", "value")

	reference := Reference{ObjectId: &obj, Id: ""}
	references.Set("key", reference)

	sa1 := &StopArea{
		model:           model,
		id:              "1234",
		NextCollectAt:   testTime,
		collectedAt:     testTime,
		CollectedUntil:  testTime,
		CollectedAlways: true,
		Name:            "Name",
		LineIds:         []LineId{"1234"},
		Attributes:      attributes,
		References:      references,
	}
	sa2 := &StopArea{
		model:           model,
		id:              "1234",
		NextCollectAt:   testTime,
		collectedAt:     testTime,
		CollectedUntil:  testTime,
		CollectedAlways: true,
		Name:            "Name",
		LineIds:         []LineId{"1234"},
		Attributes:      attributes,
		References:      references,
	}
	result, ok := Equal(sa1, sa2)
	if !ok {
		t.Errorf("Equal should return true: %v\n result: %v", ok, result)
	}
	sa3 := &StopArea{
		model:           model,
		id:              "12345",
		NextCollectAt:   testTime,
		collectedAt:     testTime,
		CollectedUntil:  testTime,
		CollectedAlways: true,
		Name:            "Name",
		LineIds:         []LineId{"1234"},
		Attributes:      attributes,
		References:      references,
	}
	result, ok = Equal(sa1, sa3)
	if ok {
		t.Errorf("Equal should return false: %v\n result: %v", ok, result)
	}
	id, ok := result["id"]
	if !ok {
		t.Errorf("Equal should return a map with id, got: %v", result)
	}
	if id != StopAreaId("12345") {
		t.Errorf("id should be 12345, got: %v. Map: %v", id, result)
	}
}
