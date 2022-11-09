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

type RequestTestData struct {
	Server *Server
	Id     string
	Method string
	Action string
	Body   []byte
}

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

func createReferential() (*Server, *core.Referential) {
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	return server, referential
}

func createPartner(referential *core.Referential) *core.Partner {
	referential.Partners().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner := referential.Partners().New("first_partner")
	referential.Partners().Save(partner)

	return partner
}

func sendRequest(rdata *RequestTestData, t *testing.T) (responseRecorder *httptest.ResponseRecorder) {

	address := []byte("/default/partners")
	if rdata.Id != "" {
		address = append(address, fmt.Sprintf("/%s", rdata.Id)...)
		if rdata.Action != "" {
			address = append(address, fmt.Sprintf("/%s", rdata.Action)...)
		}
	}

	request, err := http.NewRequest(rdata.Method, string(address), bytes.NewReader(rdata.Body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.

	rdata.Server.HandleFlow(responseRecorder, request)
	return responseRecorder
}

func Test_PartnerController_Action_Subscriber(t *testing.T) {
	rdata := &RequestTestData{}
	server, referential := createReferential()
	partner := createPartner(referential)

	rdata.Id = string(partner.Id())
	rdata.Method = "GET"
	rdata.Action = "subscriptions"
	rdata.Server = server

	sub := partner.Subscriptions()
	sub.New("kind")

	responseRecorder := sendRequest(rdata, t)
	partnerCheckResponseStatus(responseRecorder, t)

	if len(partner.Subscriptions().FindAll()) != 1 {
		t.Errorf("Should find one subscription")
	}

}

func Test_PartnerController_Delete(t *testing.T) {
	server, referential := createReferential()
	partner := createPartner(referential)

	rdata := &RequestTestData{}
	rdata.Id = string(partner.Id())
	rdata.Method = "DELETE"
	rdata.Server = server

	responseRecorder := sendRequest(rdata, t)
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
	rdata := &RequestTestData{}
	server, referential := createReferential()
	partner := createPartner(referential)
	sub := partner.Subscriptions().New("Test")
	sub.Save()

	rdata.Body = []byte(`{ "Slug": "another_test", "Name": "test" }`)
	rdata.Method = "PUT"
	rdata.Server = server
	rdata.Id = string(partner.Id())

	responseRecorder := sendRequest(rdata, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	if updatedPartner == nil {
		t.Errorf("Partner should be found after PUT request")
	}

	if expected := core.PartnerSlug("another_test"); updatedPartner.Slug() != expected {
		t.Errorf("Partner slug should be updated after PUT request:\n got: %v\n want: %v", updatedPartner.Slug(), expected)
	}
	if expected := "test"; updatedPartner.Name != expected {
		t.Errorf("Partner name should be updated after PUT request:\n got: %v\n want: %v", updatedPartner.Slug(), expected)
	}
	if expected, _ := updatedPartner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
	if len(partner.Subscriptions().FindAll()) > 0 {
		t.Errorf("All subscription should be deleted after a partenr Edit")
	}
}

func Test_PartnerController_UpdateConnectorTypes(t *testing.T) {
	// Prepare and send request
	server, referential := createReferential()
	partner := createPartner(referential)

	rdata := &RequestTestData{}
	rdata.Id = string(partner.Id())
	rdata.Method = "PUT"
	rdata.Server = server
	rdata.Body = []byte(`{ "ConnectorTypes": ["test-check-status-client"] }`)

	responseRecorder := sendRequest(rdata, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	if updatedPartner == nil {
		t.Errorf("Partner should be found after PUT request")
	}

	if expected := core.PartnerSlug("first_partner"); updatedPartner.Slug() != expected {
		t.Errorf("Partner slug should be updated after PUT request:\n got: %v\n want: %v", updatedPartner.Slug(), expected)
	}

	if len(updatedPartner.ConnectorTypes) != 1 {
		t.Errorf("ConnectorTypes should have been updated by POST request:\n got: %v\n want: %v", updatedPartner.ConnectorTypes, []string{"test-check-status-client"})
	}
}

func Test_PartnerController_Show(t *testing.T) {
	// Send request

	server, referential := createReferential()
	partner := createPartner(referential)

	rdata := &RequestTestData{}
	rdata.Id = string(partner.Id())
	rdata.Method = "GET"
	rdata.Server = server

	responseRecorder := sendRequest(rdata, t)

	// Test response
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := partner.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Create(t *testing.T) {
	// Prepare and send request
	rdata := RequestTestData{}

	server, referential := createReferential()

	rdata.Method = "POST"
	rdata.Body = []byte(`{ "Slug": "test" }`)
	rdata.Server = server

	responseRecorder := sendRequest(&rdata, t)

	// Check response
	partnerCheckResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// partner should be 6ba7b814-9dad-11d1-1-00c04fd430c8

	if len(referential.Partners().FindAll()) != 1 {
		t.Errorf("Partner should be found after POST %v", referential.Partners().FindAll()[0].Id())
		return
	}

	partner := referential.Partners().FindAll()[0]
	if partner == nil {
		t.Errorf("Partner should be found after POST request, %v", referential.Partners().FindAll()[0].Id())
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
	server, _ := createReferential()

	rdata := &RequestTestData{}
	rdata.Id = ""
	rdata.Method = "POST"
	rdata.Server = server
	rdata.Body = []byte(`{ "Slug": "invalid_slug", "ConnectorTypes": ["test-validation-connector"] }`)

	responseRecorder := sendRequest(rdata, t)

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
	expected := `{"Slug":"invalid_slug","ConnectorTypes":["test-validation-connector"],"Errors":{"slug":["Invalid format"]}}`
	if responseRecorder.Body.String() != expected {
		t.Errorf("Wrong body for invalid POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_PartnerController_Index(t *testing.T) {
	server, referential := createReferential()
	createPartner(referential)

	rdata := &RequestTestData{}
	rdata.Id = ""
	rdata.Method = "GET"
	rdata.Server = server

	responseRecorder := sendRequest(rdata, t)
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"first_partner","PartnerStatus":{"OperationnalStatus":"unknown","RetryCount":0,"ServiceStartedAt":"0001-01-01T00:00:00Z"},"ConnectorTypes":[],"Settings":{}}]`
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

func Test_PartnerController_Save(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewRealUUIDGenerator())
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()
	status, refErr := referentials.SaveToDatabase()
	if status != 200 {
		t.Fatalf("Cannot save referentials to Database: %v", refErr)
	}

	// Initialize the partners manager
	referential.Partners().SetUUIDGenerator(uuid.NewRealUUIDGenerator())
	// Save a new partner
	partner := referential.Partners().New("First Partner")
	referential.Partners().Save(partner)

	// Create a request
	request, err := http.NewRequest("POST", "/default/partners/save", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	// Test response
	partnerCheckResponseStatus(responseRecorder, t)

	//Test Results
	selectPartners := []model.SelectPartner{}
	sqlQuery := "select * from partners"
	_, err = model.Database.Select(&selectPartners, sqlQuery)
	if err != nil {
		t.Fatalf("Error while fetching partners: %v", err)
	}
	if len(selectPartners) == 0 {
		t.Fatal("Partner should be found")
	}
	if selectPartners[0].Id != string(partner.Id()) {
		t.Errorf("Saved partner has wrong id, got: %v want: %v", selectPartners[0].Id, partner.Id())
	}
	if selectPartners[0].ReferentialId != string(referential.Id()) {
		t.Errorf("Saved partner has wrong referential id, got: %v want: %v", selectPartners[0].ReferentialId, referential.Id())
	}
}
