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

func partnerCheckResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func partnerPrepareRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (partner *core.Partner, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Initialize the partners manager
	referential.Partners().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	// Save a new partner
	partner = referential.Partners().New("First Partner")
	referential.Partners().Save(partner)

	// Create a request
	address := []byte("/default/partners")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", partner.Id())...)
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

func Test_PartnerController_Delete(t *testing.T) {
	// Send request
	partner, responseRecorder, referential := partnerPrepareRequest("DELETE", true, nil, t)

	// Test response
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	testedPartner := referential.Partners().Find(partner.Id())
	if testedPartner != nil {
		t.Errorf("Partner shouldn't be found after DELETE request")
	}
	if expected, _ := partner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Slug": "Yet another test" }`)
	partner, responseRecorder, referential := partnerPrepareRequest("PUT", true, body, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	if updatedPartner == nil {
		t.Errorf("Partner should be found after PUT request")
	}

	if expected := core.PartnerSlug("Yet another test"); updatedPartner.Slug() != expected {
		t.Errorf("Partner slug should be updated after PUT request:\n got: %v\n want: %v", updatedPartner.Slug(), expected)
	}
	if expected, _ := updatedPartner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_UpdateConnectorTypes(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ConnectorTypes": ["test"] }`)
	partner, responseRecorder, referential := partnerPrepareRequest("PUT", true, body, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	if updatedPartner == nil {
		t.Errorf("Partner should be found after PUT request")
	}

	if expected := core.PartnerSlug("First Partner"); updatedPartner.Slug() != expected {
		t.Errorf("Partner slug should be updated after PUT request:\n got: %v\n want: %v", updatedPartner.Slug(), expected)
	}

	if len(updatedPartner.ConnectorTypes) != 1 {
		t.Errorf("ConnectorTypes should have been updated by POST request:\n got: %v\n want: %v", updatedPartner.ConnectorTypes, []string{"test"})
	}
}

func Test_PartnerController_Show(t *testing.T) {
	// Send request
	partner, responseRecorder, _ := partnerPrepareRequest("GET", true, nil, t)

	// Test response
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := partner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Slug": "test" }`)
	_, responseRecorder, referential := partnerPrepareRequest("POST", false, body, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// partner should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	partner := referential.Partners().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if partner == nil {
		t.Errorf("Partner should be found after POST request")
	}
	if expected := core.PartnerSlug("test"); partner.Slug() != expected {
		t.Errorf("Invalid partner slug after POST request:\n got: %v\n want: %v", partner.Slug(), expected)
	}
	if expected, _ := partner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Create_Invalid(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Slug": "InvalidSlug", "ConnectorTypes": ["test-validation-connector"] }`)
	_, responseRecorder, _ := partnerPrepareRequest("POST", false, body, t)

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
	expected := `{"Slug":"InvalidSlug","ConnectorTypes":["test-validation-connector"],"Errors":{"slug":["Invalid format"]}}`
	if responseRecorder.Body.String() != expected {
		t.Errorf("Wrong body for invalid POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := partnerPrepareRequest("GET", false, nil, t)

	// Test response
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"First Partner","OperationnalStatus":"unknown","ConnectorTypes":[],"Settings":{}}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_FindPartner(t *testing.T) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	referential := referentials.New("default")
	referential.Save()

	// Save a new partner
	partner := referential.Partners().New("First Partner")
	referential.Partners().Save(partner)

	controller := &PartnerController{referential: referential}

	foundPartner := controller.findPartner("First Partner")
	if foundPartner == nil {
		t.Error("Can't find Partner by Slug")
	}

	foundPartner = controller.findPartner(string(partner.Id()))
	if foundPartner == nil {
		t.Error("Can't find Partner by Id")
	}
}
