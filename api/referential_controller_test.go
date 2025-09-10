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
	p "bitbucket.org/enroute-mobi/ara/core/partners"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokens(t *testing.T) {
	require := require.New(t)

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	server := &Server{}
	server.SetReferentials(referentials)

	// Save a new referential
	referential := referentials.New("first_referential")
	referentials.Save(referential)

	// Create a request
	request, err := http.NewRequest("GET", "/first_referential/partners", nil)
	require.NoError(err)

	responseRecorder := httptest.NewRecorder()
	request.SetPathValue("referential_slug", "first_referential")
	request.SetPathValue("model", "partners")
	server.handleReferentialModelIndex(responseRecorder, request)

	require.Equal(http.StatusUnauthorized, responseRecorder.Code)

	referential.Tokens = []string{"12345"}
	referential.ImportTokens = []string{"23456"}
	referential.Save()

	responseRecorder = httptest.NewRecorder()
	request.SetPathValue("referential_slug", "first_referential")
	request.SetPathValue("model", "partners")
	server.handleReferentialModelIndex(responseRecorder, request)

	require.Equal(http.StatusUnauthorized, responseRecorder.Code)

	request.Header.Set("Authorization", "Token token=12345")

	responseRecorder = httptest.NewRecorder()
	request.SetPathValue("referential_slug", "first_referential")
	request.SetPathValue("model", "partners")
	server.handleReferentialModelIndex(responseRecorder, request)

	require.Equal(http.StatusOK, responseRecorder.Code)

	request.Header.Set("Authorization", "Token token=23456")

	responseRecorder = httptest.NewRecorder()
	request.SetPathValue("referential_slug", "first_referential")
	request.SetPathValue("model", "partners")
	server.handleReferentialModelIndex(responseRecorder, request)

	require.Equal(http.StatusUnauthorized, responseRecorder.Code)
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

func referentialPrepareRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (referential *core.Referential, responseRecorder *httptest.ResponseRecorder, server *Server, request *http.Request) {
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
	return
}

func Test_ReferentialController_Delete(t *testing.T) {
	// Send request
	referential, responseRecorder, server, request := referentialPrepareRequest("DELETE", true, nil, t)
	request.SetPathValue("id", string(referential.Id()))
	server.handleReferentialDelete(responseRecorder, request)

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

	referential, responseRecorder, server, request := referentialPrepareRequest("PUT", true, body, t)
	request.SetPathValue("id", string(referential.Id()))
	server.handleReferentialUpdate(responseRecorder, request)

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
	expectedLoggerStopAreas := []model.Code{model.NewCode("stif", "STIF:StopPoint:Q:473947:")}
	assert.Equal(expectedLoggerStopAreas, updatedReferential.LoggerVerboseStopAreas())
}

func Test_ReferentialController_Show(t *testing.T) {
	// Send request
	referential, responseRecorder, server, request := referentialPrepareRequest("GET", true, nil, t)
	request.SetPathValue("model", string(referential.Id()))
	request.SetPathValue("referential_slug", "_referentials")
	server.handleReferentialGet(responseRecorder, request)

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
	_, responseRecorder, server, request := referentialPrepareRequest("POST", false, body, t)
	server.handleReferentialCreate(responseRecorder, request)

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
	_, responseRecorder, server, request := referentialPrepareRequest("POST", false, body, t)
	server.handleReferentialCreate(responseRecorder, request)

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
	_, responseRecorder, server, request := referentialPrepareRequest("POST", false, body, t)
	server.handleReferentialCreate(responseRecorder, request)

	// Check response
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusBadRequest)
	}
}

func Test_ReferentialController_Index(t *testing.T) {
	// Send request
	referential, responseRecorder, server, request := referentialPrepareRequest("GET", false, nil, t)
	server.handleReferentialIndex(responseRecorder, request)

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

	server.handleReferentialSave(responseRecorder, request)

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
		Codes:           `{"internal":"value"}`,
		LineIds:         `["d0eebc99-9c0b","e0eebc99-9c0b"]`,
		Attributes:      "{}",
		References:      `{"Ref":{"Type":"Ref","Code":{"kind":"value"}}}`,
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
	request.SetPathValue("id", string(referential.Id()))
	server.handleReferentialReload(responseRecorder, request)

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

func Test_ReferentialController_Reload_Partner(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// Insert Data in the test db
	databasePartner := model.DatabasePartner{
		Id:            "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialId: "6ba7b814-9dad-11d1-0000-00c04fd430c8",
		Slug:          "partner_slug",
		Name:          "Partner Name",
	}

	model.Database.AddTableWithName(databasePartner, "partners")

	err := model.Database.Insert(&databasePartner)
	require.NoError(err)

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGeneratorLegacy())
	// Save a new referential
	referential := referentials.New("referential")
	referentials.Save(referential)

	server := &Server{}
	server.SetReferentials(referentials)

	// Ensure No partners exist
	partners := referential.Partners().FindAll()
	assert.Empty(partners, "No partners should exist")

	// Create a reload request
	request, err := http.NewRequest("POST", fmt.Sprintf("/_referentials/%v/reload", referential.Id()), nil)
	require.NoError(err)

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	request.SetPathValue("id", string(referential.Id()))
	server.handleReferentialReload(responseRecorder, request)
	require.Equal(http.StatusOK, responseRecorder.Code)

	// Test Results
	partners = referential.Partners().FindAll()
	assert.Len(partners, 1, "Partners must be loaded after a referential is reloaded")
	assert.Equal(partners[0].Id(), p.Id("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"))
}
