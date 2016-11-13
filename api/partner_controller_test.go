package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/af83/edwig/model"
)

func partnerCheckResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "text/xml")
	}
}

func partnerPrepareRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (partner *model.Partner, responseRecorder *httptest.ResponseRecorder, referential *model.Referential) {
	// Create a referential
	referentials := model.NewMemoryReferentials()
	referential = referentials.New("default")
	referential.Save()
	// Create a partnerController
	controller := NewPartnerController()
	controller.SetReferential(referential)

	// Initialize the partners manager
	referential.Partners().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	// Save a new partner
	partner = referential.Partners().New()
	partner.Name = "First Partner"
	referential.Partners().Save(partner)

	// Create a request
	address := []byte("/partners")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", partner.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call ServeHTTP method and pass in our Request and ResponseRecorder.
	controller.ServeHTTP(responseRecorder, request)

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
	body := []byte(`{ "Name": "Yet another test" }`)
	partner, responseRecorder, referential := partnerPrepareRequest("PUT", true, body, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	if updatedPartner == nil {
		t.Errorf("Partner should be found after PUT request")
	}

	if expected := "Yet another test"; updatedPartner.Name != expected {
		t.Errorf("Partner name should be updated after PUT request:\n got: %v\n want: %v", updatedPartner.Name, expected)
	}
	if expected, _ := updatedPartner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
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
	body := []byte(`{ "Name": "test" }`)
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
	if expected := "test"; partner.Name != expected {
		t.Errorf("Invalid partner name after POST request:\n got: %v\n want: %v", partner.Name, expected)
	}
	if expected, _ := partner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := partnerPrepareRequest("GET", false, nil, t)

	// Test response
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Name":"First Partner"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}
