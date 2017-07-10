package model

import (
	"encoding/json"
	"testing"
)

func Test_Operator_Id(t *testing.T) {
	operator := Operator{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if operator.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("operator.Id() returns wrong value, got: %s, required: %s", operator.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_Operator_MarshalJSON(t *testing.T) {
	operator := Operator{
		id:   "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Name: "OperatorName",
	}

	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Name":"OperatorName"}`
	jsonBytes, err := operator.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Operator.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Operator_UnmarshalJSON(t *testing.T) {
	test := `{
		"Name": "OperatorName"
		}`

	operator := Operator{}
	err := json.Unmarshal([]byte(test), &operator)

	if err != nil {
		t.Errorf("Error while Unmarshalling Operator %v", err)
	}

	if operator.Name != "OperatorName" {
		t.Errorf("Got wrong Operator Name want OperatorName got %v", operator.Name)
	}
}

func Test_Operator_Save(t *testing.T) {
	model := NewMemoryModel()
	operator := model.Operators().New()
	objectid := NewObjectID("kind", "value")
	operator.SetObjectID(objectid)

	if operator.model != model {
		t.Errorf("New operator model should be MemoryOperator model")
	}

	ok := operator.Save()
	if !ok {
		t.Errorf("operator.Save() should succeed")
	}
	_, ok = model.Operators().Find(operator.Id())
	if !ok {
		t.Errorf("New operator should be found in MemoryOperator")
	}
}

func Test_Operator_ObjectId(t *testing.T) {
	operator := Operator{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	operator.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	operator.SetObjectID(objectid)

	foundObjectId, ok := operator.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = operator.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(operator.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", operator.ObjectIDs())
	}
}

func Test_MemoryOperators_New(t *testing.T) {
	operators := NewMemoryOperators()

	operator := operators.New()
	if operator.Id() != "" {
		t.Errorf("New operator identifier should be an empty string, got: %s", operator.Id())
	}
}

func Test_MemoryOperators_Save(t *testing.T) {
	operators := NewMemoryOperators()

	operator := operators.New()

	if success := operators.Save(&operator); !success {
		t.Errorf("Save should return true")
	}

	if operator.Id() == "" {
		t.Errorf("New operator identifier shouldn't be an empty string")
	}
}

func Test_MemoryOperators_Find_NotFound(t *testing.T) {
	operators := NewMemoryOperators()
	_, ok := operators.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Operator isn't found")
	}
}

func Test_MemoryOperators_Find(t *testing.T) {
	operators := NewMemoryOperators()

	existingOperator := operators.New()
	operators.Save(&existingOperator)

	operatorId := existingOperator.Id()

	operator, ok := operators.Find(operatorId)
	if !ok {
		t.Errorf("Find should return true when operator is found")
	}
	if operator.Id() != operatorId {
		t.Errorf("Find should return a operator with the given Id")
	}
}

func Test_MemoryOperators_FindAll(t *testing.T) {
	operators := NewMemoryOperators()

	for i := 0; i < 5; i++ {
		existingOperator := operators.New()
		operators.Save(&existingOperator)
	}

	foundOperators := operators.FindAll()

	if len(foundOperators) != 5 {
		t.Errorf("FindAll should return all operators")
	}
}

func Test_MemoryOperators_Delete(t *testing.T) {
	operators := NewMemoryOperators()
	existingOperator := operators.New()
	objectid := NewObjectID("kind", "value")
	existingOperator.SetObjectID(objectid)
	operators.Save(&existingOperator)

	operators.Delete(&existingOperator)

	_, ok := operators.Find(existingOperator.Id())
	if ok {
		t.Errorf("Deleted operator should not be findable")
	}
}

func Test_MemoryOperators_Load(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Insert Data in the test db
	var databaseOperator = struct {
		Id            string `db:"id"`
		ReferentialId string `db:"referential_id"`
		Name          string `db:"name"`
		ObjectIDs     string `db:"object_ids"`
	}{
		Id:            "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialId: "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Name:          "operator",
		ObjectIDs:     `{"internal":"value"}`,
	}

	Database.AddTableWithName(databaseOperator, "operators")
	err := Database.Insert(&databaseOperator)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	operators := NewMemoryOperators()
	err = operators.Load("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if err != nil {
		t.Fatal(err)
	}

	operatorId := OperatorId(databaseOperator.Id)
	operator, ok := operators.Find(operatorId)
	if !ok {
		t.Fatal("Loaded Liness should be found")
	}

	if operator.id != operatorId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", operator.id, operatorId)
	}
	if operator.Name != "operator" {
		t.Errorf("Wrong Name:\n got: %v\n expected: operator", operator.Name)
	}
	if objectid, ok := operator.ObjectID("internal"); !ok || objectid.Value() != "value" {
		t.Errorf("Wrong ObjectID:\n got: %v:%v\n expected: \"internal\":\"value\"", objectid.Kind(), objectid.Value())
	}
}
