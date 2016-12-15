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

func checkStopVisitResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareStopVisitRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (stopVisit model.StopVisit, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Save()

	// Set the fake UUID generator
	model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())
	// Save a new stopVisit
	stopVisit = referential.Model().StopVisits().New()
	referential.Model().StopVisits().Save(&stopVisit)

	// Create a request
	address := []byte("/default/stop_visits")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", stopVisit.Id())...)
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

func Test_StopVisitController_Delete(t *testing.T) {
	// Send request
	stopVisit, responseRecorder, referential := prepareStopVisitRequest("DELETE", true, nil, t)

	// Test response
	checkStopVisitResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().StopVisits().Find(stopVisit.Id())
	if ok {
		t.Errorf("StopVisit shouldn't be found after DELETE request")
	}
	if expected, _ := stopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" } }`)
	stopVisit, responseRecorder, referential := prepareStopVisitRequest("PUT", true, body, t)

	// Check response
	checkStopVisitResponseStatus(responseRecorder, t)

	// Test Results
	updatedStopVisit, ok := referential.Model().StopVisits().Find(stopVisit.Id())
	if !ok {
		t.Errorf("StopVisit should be found after PUT request")
	}

	if expected, _ := updatedStopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Show(t *testing.T) {
	// Send request
	stopVisit, responseRecorder, _ := prepareStopVisitRequest("GET", true, nil, t)

	// Test response
	checkStopVisitResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := stopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" } }`)
	_, responseRecorder, referential := prepareStopVisitRequest("POST", false, body, t)

	// Check response
	checkStopVisitResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// stopVisit should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	stopVisit, ok := referential.Model().StopVisits().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("StopVisit should be found after POST request")
	}
	if expected, _ := stopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareStopVisitRequest("GET", false, nil, t)

	// Test response
	checkStopVisitResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"ArrivalStatus":"","DepartureStatus":"","Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","PassageOrder":0,"Schedules":[],"StopArea":"","VehicleJourney":""}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}
