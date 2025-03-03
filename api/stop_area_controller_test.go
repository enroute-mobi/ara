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

func checkStopAreaResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareStopAreaRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (stopArea *model.StopArea, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new stopArea
	stopArea = referential.Model().StopAreas().New()
	stopArea.Name = "First StopArea"
	referential.Model().StopAreas().Save(stopArea)

	// Create a request
	address := []byte("/default/stop_areas")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", stopArea.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "stop_areas")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(stopArea.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(stopArea.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(stopArea.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return
}

func Test_StopAreaController_Delete(t *testing.T) {
	assert := assert.New(t)
	// Send request
	stopArea, responseRecorder, referential := prepareStopAreaRequest("DELETE", true, nil, t)

	// Test response
	checkStopAreaResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().StopAreas().Find(stopArea.Id())
	assert.False(ok, "StopArea shouldn't be found after DELETE request")

	stopAreaBytes, err := stopArea.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(stopAreaBytes), responseRecorder.Body.String())
}

func Test_StopAreaController_Update(t *testing.T) {
	assert := assert.New(t)
	// Prepare and send request
	body := []byte(`{ "Name": "Yet another test" }`)
	stopArea, responseRecorder, referential := prepareStopAreaRequest("PUT", true, body, t)

	// Check response
	checkStopAreaResponseStatus(responseRecorder, t)

	// Test Results
	updatedStopArea, ok := referential.Model().StopAreas().Find(stopArea.Id())
	assert.True(ok)
	assert.Equal("Yet another test", updatedStopArea.Name, "StopArea name should be updated after PUT request")

	updatedStopAreaBytes, err := updatedStopArea.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(updatedStopAreaBytes), responseRecorder.Body.String())
}

func Test_StopAreaController_Show(t *testing.T) {
	assert := assert.New(t)
	// Send request
	stopArea, responseRecorder, _ := prepareStopAreaRequest("GET", true, nil, t)

	// Test response
	checkStopAreaResponseStatus(responseRecorder, t)

	//Test Results
	stopAreaBytes, err := stopArea.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(stopAreaBytes), responseRecorder.Body.String())
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
	expected := `[{"Origins":{},"Name":"First StopArea","CollectChildren":false,"CollectSituations":false,"CollectedAlways":true,"Monitored":false,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopAreaController_FindStopArea(t *testing.T) {
	ref := core.NewMemoryReferentials().New("test")

	stopArea := ref.Model().StopAreas().New()
	code := model.NewCode("codeSpace", "value")
	stopArea.SetCode(code)
	ref.Model().StopAreas().Save(stopArea)

	controller := &StopAreaController{
		referential: ref,
	}

	_, ok := controller.findStopArea("codeSpace:value")
	if !ok {
		t.Error("Can't find StopArea by Code")
	}

	_, ok = controller.findStopArea(string(stopArea.Id()))
	if !ok {
		t.Error("Can't find StopArea by Id")
	}
}
