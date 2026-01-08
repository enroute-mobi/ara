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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkFacilityResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareFacilityRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (facility *model.Facility, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new facility
	facility = referential.Model().Facilities().New()
	referential.Model().Facilities().Save(facility)

	// Create a request
	address := []byte("/default/facilities")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", facility.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "facilities")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(facility.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(facility.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(facility.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return

}

func Test_FacilityController_Delete(t *testing.T) {
	assert := assert.New(t)

	// Send request
	facility, responseRecorder, referential := prepareFacilityRequest("DELETE", true, nil, t)

	// Test response
	checkFacilityResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().Facilities().Find(facility.Id())
	assert.False(ok, "Facility shouldn't be found after DELETE request")

	expected, err := facility.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_FacilityController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Name":"FacilityName",
                  "Code": { "CodeSpace": "value" }}`)
	facility, responseRecorder, referential := prepareFacilityRequest("PUT", true, body, t)

	// Check response
	checkFacilityResponseStatus(responseRecorder, t)

	// Test Results
	updatedFacility, ok := referential.Model().Facilities().Find(facility.Id())
	assert.True(ok, "Facility should be found after PUT request")

	expected, err := updatedFacility.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_FacilityController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	facility, responseRecorder, _ := prepareFacilityRequest("GET", true, nil, t)

	// Test response
	checkFacilityResponseStatus(responseRecorder, t)

	//Test Results
	expectedFacility, _ := facility.MarshalJSON()
	assert.JSONEq(string(expectedFacility), responseRecorder.Body.String())
}

func Test_FacilityController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{
		"Status":"Unknown"}`)
	_, responseRecorder, referential := prepareFacilityRequest("POST", false, body, t)

	// Check response
	checkFacilityResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// facility should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	_, ok := referential.Model().Facilities().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "Facility should be found after POST request")

	expected := `{"Status":"unknown","Origin":"","Id":"6ba7b814-9dad-11d1-1-00c04fd430c8"}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_FacilityController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareFacilityRequest("GET", false, nil, t)

	// Test response
	checkFacilityResponseStatus(responseRecorder, t)

	//Test Results
	expected := `{"Models":
  [{
    "Id": "6ba7b814-9dad-11d1-0-00c04fd430c8",
    "Origin":"",
    "Status":"unknown"}],
    "Pagination": {"CurrentPage": 1, "PerPage": 1, "TotalCount": 1, "TotalPages": 1}}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}
