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

func checkStopAreaResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareStopAreaRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (stopArea model.StopArea, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Save()

	// Set the fake UUID generator
	model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())
	// Save a new stopArea
	stopArea = referential.Model().StopAreas().New()
	stopArea.Name = "First StopArea"
	referential.Model().StopAreas().Save(&stopArea)

	// Create a request
	address := []byte("/default/stop_areas")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", stopArea.Id())...)
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

func Test_StopAreaController_Delete(t *testing.T) {
	// Send request
	stopArea, responseRecorder, referential := prepareStopAreaRequest("DELETE", true, nil, t)

	// Test response
	checkStopAreaResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().StopAreas().Find(stopArea.Id())
	if ok {
		t.Errorf("StopArea shouldn't be found after DELETE request")
	}
	if expected, _ := stopArea.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopAreaController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Name": "Yet another test" }`)
	stopArea, responseRecorder, referential := prepareStopAreaRequest("PUT", true, body, t)

	// Check response
	checkStopAreaResponseStatus(responseRecorder, t)

	// Test Results
	updatedStopArea, ok := referential.Model().StopAreas().Find(stopArea.Id())
	if !ok {
		t.Errorf("StopArea should be found after PUT request")
	}

	if expected := "Yet another test"; updatedStopArea.Name != expected {
		t.Errorf("StopArea name should be updated after PUT request:\n got: %v\n want: %v", updatedStopArea.Name, expected)
	}
	if expected, _ := updatedStopArea.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopAreaController_Show(t *testing.T) {
	// Send request
	stopArea, responseRecorder, _ := prepareStopAreaRequest("GET", true, nil, t)

	// Test response
	checkStopAreaResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := stopArea.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopAreaController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Name": "test" }`)
	_, responseRecorder, referential := prepareStopAreaRequest("POST", false, body, t)

	// Check response
	checkStopAreaResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// stopArea should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	stopArea, ok := referential.Model().StopAreas().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("StopArea should be found after POST request")
	}
	if expected := "test"; stopArea.Name != expected {
		t.Errorf("Invalid stopArea name after POST request:\n got: %v\n want: %v", stopArea.Name, expected)
	}
	if expected, _ := stopArea.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopAreaController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareStopAreaRequest("GET", false, nil, t)

	// Test response
	checkStopAreaResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Name":"First StopArea"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}
