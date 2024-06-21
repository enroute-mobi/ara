package model

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_VehicleJourney_Id(t *testing.T) {
	vehicleJourney := VehicleJourney{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if vehicleJourney.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("VehicleJourney.Id() returns wrong value, got: %s, required: %s", vehicleJourney.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_VehicleJourney_MarshalJSON(t *testing.T) {
	model := NewTestMemoryModel()
	generator := uuid.NewFakeUUIDGenerator()
	// Create a StopVisit
	model.StopVisits().SetUUIDGenerator(generator)
	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-1-00c04fd430c8"
	model.StopVisits().Save(stopVisit)

	// Create the vehicleJourney
	model.VehicleJourneys().SetUUIDGenerator(generator)
	vehicleJourney := model.VehicleJourneys().New()
	code := NewCode("codeSpace", "value")
	vehicleJourney.SetCode(code)
	vehicleJourney.Save()

	expected := `{"Codes":{"codeSpace":"value"},"Monitored":false,"HasCompleteStopSequence":false,"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8","StopVisits":["6ba7b814-9dad-11d1-0-00c04fd430c8"]}`
	jsonBytes, err := vehicleJourney.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("VehicleJourney.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_VehicleJourney_UnmarshalJSON(t *testing.T) {
	text := `{
    "Codes": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
    "LineId": "6ba7b814-9dad-11d1-1-00c04fd430c8"
	}`

	vehicleJourney := VehicleJourney{}
	err := json.Unmarshal([]byte(text), &vehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	expectedCodes := []Code{
		NewCode("reflex", "FR:77491:ZDE:34004:STIF"),
		NewCode("hastus", "sqypis"),
	}

	for _, expectedCode := range expectedCodes {
		code, found := vehicleJourney.Code(expectedCode.CodeSpace())
		if !found {
			t.Errorf("Missing VehicleJourney Code '%s' after UnmarshalJSON()", expectedCode.CodeSpace())
		}
		if !reflect.DeepEqual(expectedCode, code) {
			t.Errorf("Wrong VehicleJourney Code after UnmarshalJSON():\n got: %s\n want: %s", code, expectedCode)
		}
	}

	if expected := LineId("6ba7b814-9dad-11d1-1-00c04fd430c8"); vehicleJourney.LineId != expected {
		t.Errorf("Wrong VehicleJourney LineId:\n got: %s\n want: %s", vehicleJourney.LineId, expected)
	}
}

func Test_VehicleJourney_Save(t *testing.T) {
	model := NewTestMemoryModel()
	vehicleJourney := model.VehicleJourneys().New()
	code := NewCode("codeSpace", "value")
	vehicleJourney.SetCode(code)

	if vehicleJourney.model != model {
		t.Errorf("New vehicleJourney model should be memoryVehicleJourneys model")
	}

	ok := vehicleJourney.Save()
	if !ok {
		t.Errorf("vehicleJourney.Save() should succeed")
	}
	_, ok = model.VehicleJourneys().Find(vehicleJourney.Id())
	if !ok {
		t.Errorf("New VehicleJourney should be found in memoryVehicleJourneys")
	}
	_, ok = model.VehicleJourneys().FindByCode(code)
	if !ok {
		t.Errorf("New VehicleJourney should be found by code in memoryVehicleJourneys")
	}
}

func Test_VehicleJourney_Code(t *testing.T) {
	vehicleJourney := VehicleJourney{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	vehicleJourney.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	vehicleJourney.SetCode(code)

	foundCode, ok := vehicleJourney.Code("codeSpace")
	if !ok {
		t.Errorf("Code should return true if Code exists")
	}
	if foundCode.Value() != code.Value() {
		t.Errorf("Code should return a correct Code:\n got: %v\n want: %v", foundCode, code)
	}

	_, ok = vehicleJourney.Code("wrongkind")
	if ok {
		t.Errorf("Code should return false if Code doesn't exist")
	}

	if len(vehicleJourney.Codes()) != 1 {
		t.Errorf("Codes should return an array with set Codes, got: %v", vehicleJourney.Codes())
	}
}

func Test_MemoryVehicleJourneys_New(t *testing.T) {
	vehicleJourneys := NewMemoryVehicleJourneys()

	vehicleJourney := vehicleJourneys.New()
	if vehicleJourney.Id() != "" {
		t.Errorf("New VehicleJourney identifier should be an empty string, got: %s", vehicleJourney.Id())
	}
}

func Test_MemoryVehicleJourneys_Save(t *testing.T) {
	vehicleJourneys := NewMemoryVehicleJourneys()

	vehicleJourney := vehicleJourneys.New()

	if success := vehicleJourneys.Save(vehicleJourney); !success {
		t.Errorf("Save should return true")
	}

	if vehicleJourney.Id() == "" {
		t.Errorf("New VehicleJourney identifier shouldn't be an empty string")
	}
}

func Test_MemoryVehicleJourneys_Find_NotFound(t *testing.T) {
	vehicleJourneys := NewMemoryVehicleJourneys()
	_, ok := vehicleJourneys.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when VehicleJourney isn't found")
	}
}

func Test_MemoryVehicleJourneys_Find(t *testing.T) {
	vehicleJourneys := NewMemoryVehicleJourneys()

	existingVehicleJourney := vehicleJourneys.New()
	vehicleJourneys.Save(existingVehicleJourney)

	vehicleJourneyId := existingVehicleJourney.Id()

	vehicleJourney, ok := vehicleJourneys.Find(vehicleJourneyId)
	if !ok {
		t.Errorf("Find should return true when VehicleJourney is found")
	}
	if vehicleJourney.Id() != vehicleJourneyId {
		t.Errorf("Find should return a VehicleJourney with the given Id")
	}
}

func Test_MemoryVehicleJourneys_FindAll(t *testing.T) {
	vehicleJourneys := NewMemoryVehicleJourneys()

	for i := 0; i < 5; i++ {
		existingVehicleJourney := vehicleJourneys.New()
		vehicleJourneys.Save(existingVehicleJourney)
	}

	foundVehicleJourneys := vehicleJourneys.FindAll()

	if len(foundVehicleJourneys) != 5 {
		t.Errorf("FindAll should return all vehicleJourneys")
	}
}

func Test_MemoryVehicleJourneys_Delete(t *testing.T) {
	assert := assert.New(t)

	vehicleJourneys := NewMemoryVehicleJourneys()
	existingVehicleJourney := vehicleJourneys.New()
	code := NewCode("codeSpace", "value")
	existingVehicleJourney.SetCode(code)
	vehicleJourneys.Save(existingVehicleJourney)

	vehicleJourneys.SetFullVehicleJourneyBySubscriptionId("subscription1", existingVehicleJourney.id)
	vehicleJourneys.SetFullVehicleJourneyBySubscriptionId("subscription2", existingVehicleJourney.id)
	vehicleJourneys.Delete(existingVehicleJourney)

	_, ok := vehicleJourneys.Find(existingVehicleJourney.Id())
	assert.False(ok, "Deleted VehicleJourney should not be findable")

	_, ok = vehicleJourneys.FindByCode(code)
	assert.False(ok, "Deleted VehicleJourney should not be findable by code")

	ok = vehicleJourneys.FullVehicleJourneyExistBySubscriptionId("subscription1", existingVehicleJourney.id)
	assert.False(ok, "Deleted VehicleJourney should not exist in full broadcasted list for subscription 1")

	ok = vehicleJourneys.FullVehicleJourneyExistBySubscriptionId("subscription2", existingVehicleJourney.id)
	assert.False(ok, "Deleted VehicleJourney should not exist in full broadcasted list for subscription 2")

	assert.Equal(vehicleJourneys.TestLenFullVehicleJourneyBySubscriptionId(), 0, "List of full broadcasted Vehicle journey must be empty")
}

func Test_MemoryVehicleJourneys_Load(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Insert Data in the test db
	databaseVehicleJourney := DatabaseVehicleJourney{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "vehicleJourney",
		Codes:           `{"internal":"value"}`,
		LineId:          "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Attributes:      "{}",
		References:      `{"Ref":{"Type":"Ref","Code":{"kind":"value"}}}`,
	}

	Database.AddTableWithName(databaseVehicleJourney, "vehicle_journeys")
	err := Database.Insert(&databaseVehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	model := NewTestMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	vehicleJourneys := model.VehicleJourneys().(*MemoryVehicleJourneys)
	err = vehicleJourneys.Load("referential")
	if err != nil {
		t.Fatal(err)
	}

	vehicleJourneyId := VehicleJourneyId(databaseVehicleJourney.Id)
	vehicleJourney, ok := vehicleJourneys.Find(vehicleJourneyId)
	if !ok {
		t.Fatalf("Loaded VehicleJourneys should be found")
	}

	if vehicleJourney.id != vehicleJourneyId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", vehicleJourney.id, vehicleJourneyId)
	}
	if vehicleJourney.Name != "vehicleJourney" {
		t.Errorf("Wrong Name:\n got: %v\n expected: vehicleJourney", vehicleJourney.Name)
	}
	if code, ok := vehicleJourney.Code("internal"); !ok || code.Value() != "value" {
		t.Errorf("Wrong Code:\n got: %v:%v\n expected: \"internal\":\"value\"", code.CodeSpace(), code.Value())
	}
	if vehicleJourney.LineId != "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" {
		t.Errorf("Wrong LineId:\n got: %v\n expected: c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", vehicleJourney.LineId)
	}
	if ref, ok := vehicleJourney.Reference("Ref"); !ok || ref.Type != "Ref" || ref.Code.CodeSpace() != "kind" || ref.Code.Value() != "value" {
		t.Errorf("Wrong References:\n got: %v\n expected Type: \"Ref\" and Code: \"codeSpace:value\"", ref)
	}
}
