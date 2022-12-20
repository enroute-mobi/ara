package model

import (
	"encoding/json"
	"testing"

	"bitbucket.org/enroute-mobi/ara/audit"
)

func Test_Vehicle_Id(t *testing.T) {
	vehicle := Vehicle{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if vehicle.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("vehicle.Id() returns wrong value, got: %s, required: %s", vehicle.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_Vehicle_MarshalJSON(t *testing.T) {
	vehicle := Vehicle{
		id:               "6ba7b814-9dad-11d1-0-00c04fd430c8",
		VehicleJourneyId: "Id",
		Longitude:        1.2,
		Latitude:         3.4,
		Bearing:          5.6,
	}
	vehicle.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	vehicle.SetObjectID(objectid)

	expected := `{"ObjectIDs":{"kind":"value"},"RecordedAtTime":"0001-01-01T00:00:00Z","ValidUntilTime":"0001-01-01T00:00:00Z","VehicleJourneyId":"Id","Longitude":1.2,"Latitude":3.4,"Bearing":5.6,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`
	jsonBytes, err := vehicle.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Vehicle.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Vehicle_UnmarshalJSON(t *testing.T) {
	test := `{
		"VehicleJourneyId": "Id",
		"Longitude": 1.3,
		"Latitude": 2.4,
		"Bearing": 5.6
}`

	vehicle := Vehicle{}
	err := json.Unmarshal([]byte(test), &vehicle)

	if err != nil {
		t.Errorf("Error while Unmarshalling Vehicle %v", err)
	}

	if vehicle.VehicleJourneyId != "Id" {
		t.Errorf("Got wrong Vehicle VehicleJourneyId want Id got %v", vehicle.VehicleJourneyId)
	}
	if vehicle.Longitude != 1.3 {
		t.Errorf("Got wrong Vehicle Longitude want 1.3 got %v", vehicle.Longitude)
	}
	if vehicle.Latitude != 2.4 {
		t.Errorf("Got wrong Vehicle Latitude want 2.4 got %v", vehicle.Latitude)
	}
	if vehicle.Bearing != 5.6 {
		t.Errorf("Got wrong Vehicle Bearing want 5.6 got %v", vehicle.Bearing)
	}
}

func Test_Vehicle_Save(t *testing.T) {
	model := NewMemoryModel()
	vehicle := model.Vehicles().New()
	objectid := NewObjectID("kind", "value")
	vehicle.SetObjectID(objectid)

	if vehicle.model != model {
		t.Errorf("New vehicle model should be MemoryVehicle model")
	}

	ok := vehicle.Save()
	if !ok {
		t.Errorf("vehicle.Save() should succeed")
	}
	_, ok = model.Vehicles().Find(vehicle.Id())
	if !ok {
		t.Errorf("New vehicle should be found in MemoryVehicle")
	}

	_, ok = model.Vehicles().FindByObjectId(objectid)
	if !ok {
		t.Errorf("New vehicle should be found by objectId")
	}
}

func Test_Vehicle_ObjectId(t *testing.T) {
	vehicle := Vehicle{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	vehicle.objectids = make(ObjectIDs)
	objectid := NewObjectID("kind", "value")
	vehicle.SetObjectID(objectid)

	foundObjectId, ok := vehicle.ObjectID("kind")
	if !ok {
		t.Errorf("ObjectID should return true if ObjectID exists")
	}
	if foundObjectId.Value() != objectid.Value() {
		t.Errorf("ObjectID should return a correct ObjectID:\n got: %v\n want: %v", foundObjectId, objectid)
	}

	_, ok = vehicle.ObjectID("wrongkind")
	if ok {
		t.Errorf("ObjectID should return false if ObjectID doesn't exist")
	}

	if len(vehicle.ObjectIDs()) != 1 {
		t.Errorf("ObjectIDs should return an array with set ObjectIDs, got: %v", vehicle.ObjectIDs())
	}
}

func Test_MemoryVehicles_New(t *testing.T) {
	vehicles := NewMemoryVehicles()

	vehicle := vehicles.New()
	if vehicle.Id() != "" {
		t.Errorf("New vehicle identifier should be an empty string, got: %s", vehicle.Id())
	}
}

func Test_MemoryVehicles_Save(t *testing.T) {
	vehicles := NewMemoryVehicles()

	vehicle := vehicles.New()

	if success := vehicles.Save(vehicle); !success {
		t.Errorf("Save should return true")
	}

	if vehicle.Id() == "" {
		t.Errorf("New vehicle identifier shouldn't be an empty string")
	}
}

func Test_MemoryVehicles_Find_NotFound(t *testing.T) {
	vehicles := NewMemoryVehicles()
	_, ok := vehicles.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Vehicle isn't found")
	}
}

func Test_MemoryVehicles_Find(t *testing.T) {
	vehicles := NewMemoryVehicles()

	existingVehicle := vehicles.New()
	vehicles.Save(existingVehicle)

	vehicleId := existingVehicle.Id()

	vehicle, ok := vehicles.Find(vehicleId)
	if !ok {
		t.Errorf("Find should return true when vehicle is found")
	}
	if vehicle.Id() != vehicleId {
		t.Errorf("Find should return a vehicle with the given Id")
	}
}

func Test_MemoryVehicles_FindAll(t *testing.T) {
	vehicles := NewMemoryVehicles()

	for i := 0; i < 5; i++ {
		existingVehicle := vehicles.New()
		vehicles.Save(existingVehicle)
	}

	foundVehicles := vehicles.FindAll()

	if len(foundVehicles) != 5 {
		t.Errorf("FindAll should return all vehicles")
	}
}

func Test_MemoryVehicles_Delete(t *testing.T) {
	vehicles := NewMemoryVehicles()
	existingVehicle := vehicles.New()
	objectid := NewObjectID("kind", "value")
	existingVehicle.SetObjectID(objectid)
	vehicles.Save(existingVehicle)

	vehicles.Delete(existingVehicle)

	_, ok := vehicles.Find(existingVehicle.Id())
	if ok {
		t.Errorf("Deleted vehicle should not be findable")
	}
	_, ok = vehicles.FindByObjectId(objectid)
	if ok {
		t.Errorf("Deleted vehicle should not be findable by objectid")
	}
}

func Test_Save_BiqQuery(t *testing.T) {
	f := audit.NewFakeBigQuery()
	audit.SetCurrentBigQuery("ref", f)

	m := NewMemoryModel("ref")
	vehicles := NewMemoryVehicles()
	vehicles.model = m
	v := vehicles.New()
	objectid := NewObjectID("kind", "value")
	v.SetObjectID(objectid)
	v.Latitude = 1.0
	vehicles.Save(v)

	if len(f.VehicleEvents()) != 1 {
		t.Error("New VehicleJourney save should have send a BQ message")
	}

	v2, _ := vehicles.Find(v.id)
	v2.Latitude = 2.0
	vehicles.Save(v2)

	if len(f.VehicleEvents()) != 2 {
		t.Error("VehicleJourney modification save should have send a BQ message")
	}

	v3, _ := vehicles.Find(v.id)
	vehicles.Save(v3)

	if len(f.VehicleEvents()) != 2 {
		t.Error("VehicleJourney save without modifications should not have send a BQ message")
	}

}
