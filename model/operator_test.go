package model

import (
	"encoding/json"
	"testing"
	"time"
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

	expected := `{"Name":"OperatorName","Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`
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
	code := NewCode("codeSpace", "value")
	operator.SetCode(code)

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

func Test_Operator_Code(t *testing.T) {
	operator := Operator{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	operator.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	operator.SetCode(code)

	foundCode, ok := operator.Code("codeSpace")
	if !ok {
		t.Errorf("Code should return true if Code exists")
	}
	if foundCode.Value() != code.Value() {
		t.Errorf("Code should return a correct Code:\n got: %v\n want: %v", foundCode, code)
	}

	_, ok = operator.Code("wrongkind")
	if ok {
		t.Errorf("Code should return false if Code doesn't exist")
	}

	if len(operator.Codes()) != 1 {
		t.Errorf("Codes should return an array with set Codes, got: %v", operator.Codes())
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

	if success := operators.Save(operator); !success {
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
	operators.Save(existingOperator)

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
		operators.Save(existingOperator)
	}

	foundOperators := operators.FindAll()

	if len(foundOperators) != 5 {
		t.Errorf("FindAll should return all operators")
	}
}

func Test_MemoryOperators_Delete(t *testing.T) {
	operators := NewMemoryOperators()
	existingOperator := operators.New()
	code := NewCode("codeSpace", "value")
	existingOperator.SetCode(code)
	operators.Save(existingOperator)

	operators.Delete(existingOperator)

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
		Id              string `db:"id"`
		ReferentialSlug string `db:"referential_slug"`
		Name            string `db:"name"`
		Codes           string `db:"codes"`
		ModelName       string `db:"model_name"`
	}{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		Name:            "operator",
		Codes:           `{"internal":"value"}`,
		ModelName:       "2017-01-01",
	}

	Database.AddTableWithName(databaseOperator, "operators")
	err := Database.Insert(&databaseOperator)
	if err != nil {
		t.Fatal(err)
	}
	model := NewMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}

	operators := model.Operators().(*MemoryOperators)

	// Fetch data from the db
	err = operators.Load("referential")
	if err != nil {
		t.Fatal(err)
	}

	operatorId := OperatorId(databaseOperator.Id)
	operator, ok := operators.Find(operatorId)
	if !ok {
		t.Fatal("Loaded Operator should be found")
	}

	if operator.id != operatorId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", operator.id, operatorId)
	}
	if operator.Name != "operator" {
		t.Errorf("Wrong Name:\n got: %v\n expected: operator", operator.Name)
	}
	if code, ok := operator.Code("internal"); !ok || code.Value() != "value" {
		t.Errorf("Wrong Code:\n got: %v:%v\n expected: \"internal\":\"value\"", code.CodeSpace(), code.Value())
	}
}
