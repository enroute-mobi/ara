package model

import "testing"

func Test_Equal(t *testing.T) {
	type testStruct struct {
		A int
		B int
	}
	t1 := &testStruct{A: 2, B: 2}
	t2 := &testStruct{A: 2, B: 2}
	t3 := &testStruct{A: 2, B: 1}

	result, err := Equal(t1, t2)
	if err != nil {
		t.Fatalf("Error in Equal: %v", err)
	}
	if !result.Equal {
		t.Errorf("Equal should return true, result: %v", result)
	}

	result, err = Equal(t1, t3)
	if err != nil {
		t.Fatalf("Error in Equal: %v", err)
	}
	if result.Equal {
		t.Errorf("Equal should return false, result: %v", result)
	}
	b, ok := result.DiffMap["B"]
	if !ok {
		t.Errorf("DiffMap should contains B: %v", result.DiffMap)
	}
	if b != 1 {
		t.Errorf("Wrong value in diffMap, got: %v expected: 1", b)
	}
}

func Test_Equal_DiffignoreTag(t *testing.T) {
	type testStruct struct {
		A int
		B int `diffignore:"true"`
	}
	t1 := &testStruct{A: 2, B: 2}
	t2 := &testStruct{A: 2, B: 1}

	result, err := Equal(t1, t2)
	if err != nil {
		t.Fatalf("Error in Equal: %v", err)
	}
	if !result.Equal {
		t.Errorf("Equal should ignore unexported fields, result: %v", result)
	}
}

func Test_Equal_Unexported(t *testing.T) {
	type testStruct struct {
		A int
		a int
	}
	t1 := &testStruct{A: 2, a: 2}
	t2 := &testStruct{A: 2, a: 1}

	result, err := Equal(t1, t2)
	if err != nil {
		t.Fatalf("Error in Equal: %v", err)
	}
	if !result.Equal {
		t.Errorf("Equal should ignore unexported fields, result: %v", result)
	}
}

// Commented for now, SA have a field Origins which is a sync.RWMutex. It would need a tag to ignore but Equal isn't used for now
// func Test_Equal_StopAreas(t *testing.T) {
// 	model := NewMemoryModel()
// 	testTime := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)

// 	attributes := NewAttributes()
// 	attributes.Set("key", "value")

// 	references := NewReferences()
// 	obj := NewCode("codeSpace", "value")

// 	reference := Reference{Code: &obj}
// 	references.Set("key", reference)

// 	sa1 := &StopArea{
// 		model:           model,
// 		id:              "1234",
// 		nextCollectAt:   testTime,
// 		collectedAt:     testTime,
// 		CollectedUntil:  testTime,
// 		CollectedAlways: true,
// 		Name:            "Name",
// 		LineIds:         []LineId{"1234"},
// 		Attributes:      attributes,
// 		Origins:         NewStopAreaOrigins(),
// 		//		References:      references, Test fail if added
// 	}
// 	sa2 := &StopArea{
// 		model:           model,
// 		id:              "1234",
// 		nextCollectAt:   testTime,
// 		collectedAt:     testTime,
// 		CollectedUntil:  testTime,
// 		CollectedAlways: true,
// 		Name:            "Name",
// 		LineIds:         []LineId{"1234"},
// 		Attributes:      attributes,
// 		Origins:         NewStopAreaOrigins(),
// 		//		References:      references, Test fail if added
// 	}
// 	result, err := Equal(sa1, sa2)
// 	if err != nil {
// 		t.Fatalf("Error in Equal: %v", err)
// 	}
// 	if !result.Equal {
// 		t.Errorf("Equal should return true, result: %v", result)
// 	}

// 	sa2.Name = "Name2"
// 	result, err = Equal(sa1, sa2)
// 	if err != nil {
// 		t.Fatalf("Error in Equal: %v", err)
// 	}
// 	if result.Equal {
// 		t.Errorf("Equal should return false, result: %v", result)
// 	}
// }
