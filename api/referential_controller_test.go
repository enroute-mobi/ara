package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	server := &Server{}
	server.SetReferentials(referentials)

	// Save a new referential
	referential := referentials.New("first_referential")
	referentials.Save(referential)

	// Create a request
	request, err := http.NewRequest("Get", "/first_referential/partners", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	server.HandleFlow(responseRecorder, request)

	if status := responseRecorder.Code; status == http.StatusOK {
		t.Errorf("Handler returned wrong status code: %v", status)
		panic(responseRecorder.Body)
	}

	referential.Tokens = []string{"12345"}
	referential.ImportTokens = []string{"23456"}
	referential.Save()

	responseRecorder = httptest.NewRecorder()
	server.HandleFlow(responseRecorder, request)

	if status := responseRecorder.Code; status == http.StatusOK {
		t.Errorf("Handler returned wrong status code: %v", status)
	}

	request.Header.Set("Authorization", "Token token=12345")

	responseRecorder = httptest.NewRecorder()
	server.HandleFlow(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: %v", status)
	}

	request.Header.Set("Authorization", "Token token=23456")

	responseRecorder = httptest.NewRecorder()
	server.HandleFlow(responseRecorder, request)

	if status := responseRecorder.Code; status == http.StatusOK {
		t.Errorf("Handler returned wrong status code: %v", status)
	}
}

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
	referential = referentials.New("first_referential")
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
	assert := assert.New(t)
	// Prepare and send request
	body := []byte(`{
"Slug": "another_test",
"OrganisationId": "test",
"Name": "test name",
"Settings": {"model.refresh_time": "4h", "logger.verbose.stop_areas": "stif:STIF:StopPoint:Q:473947:"}
}`)

	referential, responseRecorder, server := referentialPrepareRequest("PUT", true, body, t)

	// Check response
	referentialCheckResponseStatus(responseRecorder, t)

	// Test Results
	updatedReferential := server.CurrentReferentials().Find(referential.Id())
	assert.NotNil(updatedReferential, "Referential should be found after PUT request")

	assert.Equal(core.ReferentialSlug("another_test"), updatedReferential.Slug())
	assert.Equal("test", updatedReferential.OrganisationId)
	assert.Equal("test name", updatedReferential.Name)

	jsonReferential, _ := updatedReferential.MarshalJSON()
	assert.JSONEq(responseRecorder.Body.String(), string(jsonReferential), "Wrong body for PUT response request")

	// Settings must be set
	assert.Equal(4*time.Hour, updatedReferential.ModelRefreshTime())
	expectedLoggerStopAreas := []model.ObjectID{model.NewObjectID("stif", "STIF:StopPoint:Q:473947:")}
	assert.Equal(expectedLoggerStopAreas, updatedReferential.LoggerVerboseStopAreas())
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
	body := []byte(`{"Slug":"first_referential"}`)
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
	referential := referentials.New("first_referential")
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

func Test_ReferentialController_Reload(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// Insert Data in the test db
	databaseStopArea := model.DatabaseStopArea{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "1984-04-04",
		Name:            "stopArea",
		ObjectIDs:       `{"internal":"value"}`,
		LineIds:         `["d0eebc99-9c0b","e0eebc99-9c0b"]`,
		Attributes:      "{}",
		References:      `{"Ref":{"Type":"Ref","ObjectId":{"kind":"value"}}}`,
		CollectedAlways: true,
		CollectChildren: true,
	}

	model.Database.AddTableWithName(databaseStopArea, "stop_areas")
	err := model.Database.Insert(&databaseStopArea)
	if err != nil {
		t.Fatal(err)
	}

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewRealUUIDGenerator())
	// Save a new referential
	referential := referentials.New("referential")
	referentials.Save(referential)

	server := &Server{}
	server.SetReferentials(referentials)
	// Create a request
	request, err := http.NewRequest("POST", fmt.Sprintf("/_referentials/%v/reload", referential.Id()), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	// Test response
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Fatalf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	//Test Results
	sas := referential.Model().StopAreas().FindAll()

	if len(sas) != 1 {
		t.Fatalf("After reload, referential should have one stopArea, got %v", len(sas))
	}
	if sas[0].Id() != "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" {
		t.Errorf("Loaded stopArea have wrong id:\n got: %v\n want: a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", sas[0].Id())
	}
}
