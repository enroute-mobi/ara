package model

import (
	"encoding/json"
	"testing"

	"bitbucket.org/enroute-mobi/ara/audit"
	"github.com/stretchr/testify/assert"
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
	vehicle.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	vehicle.SetCode(code)

	expected := `{"Codes":{"codeSpace":"value"},"RecordedAtTime":"0001-01-01T00:00:00Z","ValidUntilTime":"0001-01-01T00:00:00Z","VehicleJourneyId":"Id","Longitude":1.2,"Latitude":3.4,"Bearing":5.6,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`
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
	code := NewCode("codeSpace", "value")
	vehicle.SetCode(code)

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

	_, ok = model.Vehicles().FindByCode(code)
	if !ok {
		t.Errorf("New vehicle should be found by code")
	}
}

func Test_Vehicle_Save_WithNextStopVisitId(t *testing.T) {
	assert := assert.New(t)

	model := NewMemoryModel()
	vehicle := model.Vehicles().New()
	stopVisit := model.StopVisits().New()
	stopVisit.Save()

	code := NewCode("codeSpace", "value")
	vehicle.SetCode(code)
	vehicle.NextStopVisitId = stopVisit.Id()
	ok := vehicle.Save()
	assert.True(ok)

	_, ok = model.Vehicles().FindByNextStopVisitId(stopVisit.Id())
	assert.Truef(ok, "Should find vehicle by Next stop visit Id")
}

func Test_Vehicle_NextStopVisitId_with_Updates(t *testing.T) {
	assert := assert.New(t)

	var ok bool

	model := NewMemoryModel()
	vehicleA := model.Vehicles().New()
	vehicleB := model.Vehicles().New()
	stopVisit1 := model.StopVisits().New()
	stopVisit1.Save()

	code := NewCode("codeSpace", "value")
	vehicleA.SetCode(code)
	vehicleB.SetCode(code)

	vehicleA.NextStopVisitId = stopVisit1.Id()
	ok = vehicleA.Save()
	assert.True(ok)

	ok = vehicleB.Save()
	assert.True(ok)

	vehicle, ok := model.Vehicles().FindByNextStopVisitId(stopVisit1.Id())
	assert.Equal(vehicleA, vehicle)
	assert.Truef(ok, "Should find vehicleA by nextStopVisit Id# 1")

	// Update the vehicleB with nextStopVisit => stopVisit1
	vehicleB.NextStopVisitId = stopVisit1.Id()
	ok = vehicleB.Save()
	assert.True(ok)

	vehicle, ok = model.Vehicles().FindByNextStopVisitId(stopVisit1.Id())
	assert.True(ok)
	assert.Equal(vehicleB, vehicle, "Should find vehicleB nexStopvisit with Id# 1")

	// VehicleA and VehicleB should have same nextStopVisitId
	assert.Equal(stopVisit1.Id(), vehicleA.NextStopVisitId)
	assert.Equal(stopVisit1.Id(), vehicleB.NextStopVisitId)
}

func Test_Vehicle_Code(t *testing.T) {
	vehicle := Vehicle{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	vehicle.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	vehicle.SetCode(code)

	foundCode, ok := vehicle.Code("codeSpace")
	if !ok {
		t.Errorf("Code should return true if Code exists")
	}
	if foundCode.Value() != code.Value() {
		t.Errorf("Code should return a correct Code:\n got: %v\n want: %v", foundCode, code)
	}

	_, ok = vehicle.Code("wrongkind")
	if ok {
		t.Errorf("Code should return false if Code doesn't exist")
	}

	if len(vehicle.Codes()) != 1 {
		t.Errorf("Codes should return an array with set Codes, got: %v", vehicle.Codes())
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
	code := NewCode("codeSpace", "value")
	existingVehicle.SetCode(code)
	vehicles.Save(existingVehicle)

	vehicles.Delete(existingVehicle)

	_, ok := vehicles.Find(existingVehicle.Id())
	if ok {
		t.Errorf("Deleted vehicle should not be findable")
	}
	_, ok = vehicles.FindByCode(code)
	if ok {
		t.Errorf("Deleted vehicle should not be findable by code")
	}
}

func Test_MemoryVehicles_Delete_WithNextStopVisitId(t *testing.T) {
	assert := assert.New(t)
	model := NewMemoryModel()

	vehicle := model.Vehicles().New()

	stopVisit := model.StopVisits().New()
	stopVisit.Save()

	code := NewCode("codeSpace", "value")
	vehicle.SetCode(code)
	vehicle.NextStopVisitId = stopVisit.Id()

	vehicle.Save()

	model.Vehicles().Delete(vehicle)

	_, ok := model.Vehicles().Find(vehicle.Id())
	assert.False(ok)

	_, ok = model.Vehicles().FindByNextStopVisitId(stopVisit.Id())
	assert.Falsef(ok, "Deleted vehicle should not be findable by next stopVisit id")
}

func Test_Save_BiqQuery(t *testing.T) {
	f := audit.NewFakeBigQuery()
	audit.SetCurrentBigQuery("ref", f)

	m := NewMemoryModel("ref")
	vehicles := NewMemoryVehicles()
	vehicles.model = m
	v := vehicles.New()
	code := NewCode("codeSpace", "value")
	v.SetCode(code)
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
