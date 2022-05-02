package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func checkVehicleResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareVehicleRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (vehicle *model.Vehicle, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new vehicle
	vehicle = referential.Model().Vehicles().New()
	vehicle.Longitude = 1.2
	vehicle.Latitude = 3.4
	vehicle.Bearing = 5.6
	referential.Model().Vehicles().Save(vehicle)

	// Create a request
	address := []byte("/default/vehicles")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", vehicle.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	return
}

func Test_VehicleController_Delete(t *testing.T) {
	// Send request
	vehicle, responseRecorder, referential := prepareVehicleRequest("DELETE", true, nil, t)

	// Test response
	checkVehicleResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().Vehicles().Find(vehicle.Id())
	if ok {
		t.Errorf("Vehicle shouldn't be found after DELETE request")
	}
	if expected, _ := vehicle.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Bearing": 13.3 }`)
	vehicle, responseRecorder, referential := prepareVehicleRequest("PUT", true, body, t)

	// Check response
	checkVehicleResponseStatus(responseRecorder, t)

	// Test Results
	updatedVehicle, ok := referential.Model().Vehicles().Find(vehicle.Id())
	if !ok {
		t.Errorf("Vehicle should be found after PUT request")
	}

	if expected := 13.3; updatedVehicle.Bearing != expected {
		t.Errorf("Vehicle bearing should be updated after PUT request:\n got: %v\n want: %v", updatedVehicle.Bearing, expected)
	}
	if expected, _ := updatedVehicle.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleController_Show(t *testing.T) {
	// Send request
	vehicle, responseRecorder, _ := prepareVehicleRequest("GET", true, nil, t)

	// Test response
	checkVehicleResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := vehicle.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Bearing": 1.234 }`)
	_, responseRecorder, referential := prepareVehicleRequest("POST", false, body, t)

	// Check response
	checkVehicleResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// vehicle should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	vehicle, ok := referential.Model().Vehicles().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("Vehicle should be found after POST request")
	}
	if expected := 1.234; vehicle.Bearing != expected {
		t.Errorf("Invalid vehicle bearing after POST request:\n got: %v\n want: %v", vehicle.Bearing, expected)
	}
	if expected, _ := vehicle.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareVehicleRequest("GET", false, nil, t)

	// Test response
	checkVehicleResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Longitude":1.2,"Latitude":3.4,"Bearing":5.6,"ValidUntilTime":"0001-01-01T00:00:00Z","RecordedAtTime":"0001-01-01T00:00:00Z"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleController_FindVehicle(t *testing.T) {
	ref := core.NewMemoryReferentials().New("test")
	vehicle := ref.Model().Vehicles().New()
	objectid := model.NewObjectID("kind", "value")
	vehicle.SetObjectID(objectid)
	ref.Model().Vehicles().Save(vehicle)

	controller := &VehicleController{
		referential: ref,
	}

	_, ok := controller.findVehicle("kind:value")
	if !ok {
		t.Error("Can't find Vehicle by ObjectId")
	}

	_, ok = controller.findVehicle(string(vehicle.Id()))
	if !ok {
		t.Error("Can't find Vehicle by Id")
	}
}
