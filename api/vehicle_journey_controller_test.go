package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/model"
)

func checkVehicleJourneyResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareVehicleJourneyRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (vehicleJourney model.VehicleJourney, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Save()

	// Set the fake UUID generator
	model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())
	// Save a new vehicleJourney
	vehicleJourney = referential.Model().VehicleJourneys().New()
	referential.Model().VehicleJourneys().Save(&vehicleJourney)

	// Create a request
	address := []byte("/default/vehicle_journeys")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", vehicleJourney.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call APIHandler method and pass in our Request and ResponseRecorder.
	server.APIHandler(responseRecorder, request)

	return
}

func Test_VehicleJourneyController_Delete(t *testing.T) {
	// Send request
	vehicleJourney, responseRecorder, referential := prepareVehicleJourneyRequest("DELETE", true, nil, t)

	// Test response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().VehicleJourneys().Find(vehicleJourney.Id())
	if ok {
		t.Errorf("VehicleJourney shouldn't be found after DELETE request")
	}
	if expected, _ := vehicleJourney.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleJourneyController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	vehicleJourney, responseRecorder, referential := prepareVehicleJourneyRequest("PUT", true, body, t)

	// Check response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	// Test Results
	updatedVehicleJourney, ok := referential.Model().VehicleJourneys().Find(vehicleJourney.Id())
	if !ok {
		t.Errorf("VehicleJourney should be found after PUT request")
	}

	if expected, _ := updatedVehicleJourney.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleJourneyController_Show(t *testing.T) {
	// Send request
	vehicleJourney, responseRecorder, _ := prepareVehicleJourneyRequest("GET", true, nil, t)

	// Test response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := vehicleJourney.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleJourneyController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	_, responseRecorder, referential := prepareVehicleJourneyRequest("POST", false, body, t)

	// Check response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// vehicleJourney should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	vehicleJourney, ok := referential.Model().VehicleJourneys().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("VehicleJourney should be found after POST request")
	}
	if expected, _ := vehicleJourney.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_VehicleJourneyController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareVehicleJourneyRequest("GET", false, nil, t)

	// Test response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Attributes":{},"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","LineId":"","Name":"","References":{},"StopVisits":[]}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}
