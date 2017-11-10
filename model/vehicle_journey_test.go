package model

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_VehicleJourney_Id(t *testing.T) {
	vehicleJourney := VehicleJourney{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if vehicleJourney.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("VehicleJourney.Id() returns wrong value, got: %s, required: %s", vehicleJourney.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

// WIP: Determine what to return in JSON
func Test_VehicleJourney_MarshalJSON(t *testing.T) {
	model := NewMemoryModel()
	generator := NewFakeUUIDGenerator()
	// Create a StopVisit
	model.StopVisits().SetUUIDGenerator(generator)
	stopVisit := model.StopVisits().New()
	stopVisit.VehicleJourneyId = "6ba7b814-9dad-11d1-1-00c04fd430c8"
	model.StopVisits().Save(&stopVisit)

	// Create the vehicleJourney
	model.VehicleJourneys().SetUUIDGenerator(generator)
	vehicleJourney := model.VehicleJourneys().New()
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.Save()

	expected := `{"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8","ObjectIDs":{"kind":"value"},"StopVisits":["6ba7b814-9dad-11d1-0-00c04fd430c8"],"Monitored":false}`
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
    "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
    "LineId": "6ba7b814-9dad-11d1-1-00c04fd430c8"
	}`

	vehicleJourney := VehicleJourney{}
	err := json.Unmarshal([]byte(text), &vehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	expectedObjectIds := []ObjectID{
		NewObjectID("reflex", "FR:77491:ZDE:34004:STIF"),
		NewObjectID("hastus", "sqypis"),
	}

	for _, expectedObjectId := range expectedObjectIds {
		objectId, found := vehicleJourney.ObjectID(expectedObjectId.Kind())
		if !found {
			t.Errorf("Missing VehicleJourney ObjectId '%s' after UnmarshalJSON()", expectedObjectId.Kind())
		}
		if !reflect.DeepEqual(expectedObjectId, objectId) {
			t.Errorf("Wrong VehicleJourney ObjectId after UnmarshalJSON():\n got: %s\n want: %s", objectId, expectedObjectId)
		}
	}

	if expected := LineId("6ba7b814-9dad-11d1-1-00c04fd430c8"); vehicleJourney.LineId != expected {
		t.Errorf("Wrong VehicleJourney LineId:\n got: %s\n want: %s", vehicleJourney.LineId, expected)
	}
}

func Test_VehicleJourney_Save(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourney := model.VehicleJourneys().New()
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)

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
	_, ok = model.VehicleJourneys().FindByObjectId(objectid)
	if !ok {
		t.Errorf("New VehicleJourney should be found by objectid in memoryVehicleJourneys")
	}
}

func Test_VehicleJourney_ObjectId(t *testing.T) {
	vehicleJourney := VehicleJourney{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	vehicleJourney.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)

	foundObjectId, ok := vehicleJourney.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = vehicleJourney.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(vehicleJourney.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", vehicleJourney.ObjectIDs())
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

	if success := vehicleJourneys.Save(&vehicleJourney); !success {
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
	vehicleJourneys.Save(&existingVehicleJourney)

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
		vehicleJourneys.Save(&existingVehicleJourney)
	}

	foundVehicleJourneys := vehicleJourneys.FindAll()

	if len(foundVehicleJourneys) != 5 {
		t.Errorf("FindAll should return all vehicleJourneys")
	}
}

func Test_MemoryVehicleJourneys_Delete(t *testing.T) {
	vehicleJourneys := NewMemoryVehicleJourneys()
	existingVehicleJourney := vehicleJourneys.New()
	objectid := NewObjectID("kind", "value")
	existingVehicleJourney.SetObjectID(objectid)
	vehicleJourneys.Save(&existingVehicleJourney)

	vehicleJourneys.Delete(&existingVehicleJourney)

	_, ok := vehicleJourneys.Find(existingVehicleJourney.Id())
	if ok {
		t.Errorf("Deleted VehicleJourney should not be findable")
	}
	_, ok = vehicleJourneys.FindByObjectId(objectid)
	if ok {
		t.Errorf("Deleted VehicleJourney should not be findable by objectid")
	}
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
		ObjectIDs:       `{"internal":"value"}`,
		LineId:          "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Attributes:      "{}",
		References:      "{}",
	}

	Database.AddTableWithName(databaseVehicleJourney, "vehicle_journeys")
	err := Database.Insert(&databaseVehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	model := NewMemoryModel()
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
	if objectid, ok := vehicleJourney.ObjectID("internal"); !ok || objectid.Value() != "value" {
		t.Errorf("Wrong ObjectID:\n got: %v:%v\n expected: \"internal\":\"value\"", objectid.Kind(), objectid.Value())
	}
	if vehicleJourney.LineId != "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" {
		t.Errorf("Wrong LineId:\n got: %v\n expected: c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", vehicleJourney.LineId)
	}
}
