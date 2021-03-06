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

func referentialCheckResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func referentialPrepareRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (referential *core.Referential, responseRecorder *httptest.ResponseRecorder, server *Server) {
	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new referential
	referential = referentials.New("First Referential")
	referentials.Save(referential)

	server = &Server{}
	server.SetReferentials(referentials)
	// Create a request
	address := []byte("/_referentials")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", referential.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	return
}

func Test_ReferentialController_Delete(t *testing.T) {
	// Send request
	referential, responseRecorder, server := referentialPrepareRequest("DELETE", true, nil, t)

	// Test response
	referentialCheckResponseStatus(responseRecorder, t)

	//Test Results
	testedReferential := server.CurrentReferentials().Find(referential.Id())
	if testedReferential != nil {
		t.Errorf("Referential shouldn't be found after DELETE request")
	}
	if expected, _ := referential.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_ReferentialController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Slug": "Yet another test" }`)
	referential, responseRecorder, server := referentialPrepareRequest("PUT", true, body, t)

	// Check response
	referentialCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedReferential := server.CurrentReferentials().Find(referential.Id())
	if updatedReferential == nil {
		t.Errorf("Referential should be found after PUT request")
	}

	if expected := core.ReferentialSlug("Yet another test"); updatedReferential.Slug() != expected {
		t.Errorf("Referential slug should be updated after PUT request:\n got: %v\n want: %v", updatedReferential.Slug(), expected)
	}
	if expected, _ := updatedReferential.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_ReferentialController_Update_ExistingSlug(t *testing.T) {
	// Prepare and send request
	body := []byte(`{"Slug":"referential"}`)
	referentialPrepareRequest("POST", false, body, t)

	body = []byte(`{"Slug":"referential"}`)
	_, responseRecorder, _ := referentialPrepareRequest("POST", true, body, t)
	// Check response
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusBadRequest)
	}

	// Check response
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusBadRequest)
	}
}

func Test_ReferentialController_Show(t *testing.T) {
	// Send request
	referential, responseRecorder, _ := referentialPrepareRequest("GET", true, nil, t)

	// Test response
	referentialCheckResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := referential.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_ReferentialController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Slug": "test" }`)
	_, responseRecorder, server := referentialPrepareRequest("POST", false, body, t)

	// Check response
	referentialCheckResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// referential should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	referential := server.CurrentReferentials().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if referential == nil {
		t.Fatal("Referential should be found after POST request")
	}
	if expected := core.ReferentialSlug("test"); referential.Slug() != expected {
		t.Errorf("Invalid referential slug after POST request:\n got: %v\n want: %v", referential.Slug(), expected)
	}
	if expected, _ := referential.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_ReferentialController_Create_Invalid(t *testing.T) {
	// Prepare and send request
	body := []byte(`{}`)
	_, responseRecorder, _ := referentialPrepareRequest("POST", false, body, t)

	// Check response
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusBadRequest)
	}
	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}

	// Test Results
	expected := `{"Errors":{"Slug":["Can't be empty"]}}`
	if responseRecorder.Body.String() != expected {
		t.Errorf("Wrong body for invalid POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_ReferentialController_Create_ExistingSlug(t *testing.T) {
	// Prepare and send request
	body := []byte(`{"Slug":"First Referential"}`)
	_, responseRecorder, _ := referentialPrepareRequest("POST", false, body, t)

	// Check response
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusBadRequest)
	}
}

func Test_ReferentialController_Index(t *testing.T) {
	// Send request
	referential, responseRecorder, _ := referentialPrepareRequest("GET", false, nil, t)

	// Test response
	referentialCheckResponseStatus(responseRecorder, t)

	//Test Results
	referentialJSON, _ := referential.MarshalJSON()
	expected := fmt.Sprintf("[%s]", referentialJSON)

	if responseRecorder.Body.String() != expected {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), expected)
	}
}

func Test_ReferentialController_Save(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewRealUUIDGenerator())
	// Save a new referential
	referential := referentials.New("First Referential")
	referentials.Save(referential)

	server := &Server{}
	server.SetReferentials(referentials)
	// Create a request
	request, err := http.NewRequest("POST", "/_referentials/save", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	// Test response
	referentialCheckResponseStatus(responseRecorder, t)

	//Test Results
	referentials2 := core.NewMemoryReferentials()
	err = referentials2.Load()
	if err != nil {
		t.Fatal(err)
	}

	if r := referentials2.Find(core.ReferentialId(referential.Id())); r == nil {
		t.Errorf("Loaded Referentials should be found")
	}
}
