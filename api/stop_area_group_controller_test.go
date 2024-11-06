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

func checkStopAreaGroupResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)
	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareStopAreaGroupRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (stopAreaGroup *model.StopAreaGroup, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new stopAreaGroup
	stopAreaGroup = referential.Model().StopAreaGroups().New()
	referential.Model().StopAreaGroups().Save(stopAreaGroup)

	// Create a request
	address := []byte("/default/stop_area_groups")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", stopAreaGroup.Id())...)
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

func Test_StopAreaGroupController_Delete(t *testing.T) {
	assert := assert.New(t)

	// Send request
	stopAreaGroup, responseRecorder, referential := prepareStopAreaGroupRequest("DELETE", true, nil, t)

	// Test response
	checkStopAreaGroupResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().StopAreaGroups().Find(stopAreaGroup.Id())
	assert.False(ok, "StopAreaGroup shouldn't be found after DELETE request")

	expected, err := stopAreaGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(expected), responseRecorder.Body.String())
}

func Test_StopAreaGroupController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Name":"StopAreaGroupName", "ShortName":"short_name"}`)
	stopAreaGroup, responseRecorder, referential := prepareStopAreaGroupRequest("PUT", true, body, t)

	// Check response
	checkStopAreaGroupResponseStatus(responseRecorder, t)

	// Test Results
	updatedStopAreaGroup, ok := referential.Model().StopAreaGroups().Find(stopAreaGroup.Id())
	assert.True(ok, "StopAreaGroup should be found after PUT request")

	expected, err := updatedStopAreaGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(expected), responseRecorder.Body.String())
}

func Test_StopAreaGroupController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	stopAreaGroup, responseRecorder, _ := prepareStopAreaGroupRequest("GET", true, nil, t)

	// Test response
	checkStopAreaGroupResponseStatus(responseRecorder, t)

	//Test Results
	expected, err := stopAreaGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(expected), responseRecorder.Body.String()) //
}

func Test_StopAreaGroupController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{
		"Name":"StopAreaGroupName"}`)
	_, responseRecorder, referential := prepareStopAreaGroupRequest("POST", false, body, t)

	// Check response
	checkStopAreaGroupResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// stopAreaGroup should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	stopAreaGroup, ok := referential.Model().StopAreaGroups().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "StopAreaGroup should be found after POST request")

	expected := `{"Name":"StopAreaGroupName","Id":"6ba7b814-9dad-11d1-1-00c04fd430c8"}`
	stopAreaGroupMarshal, err := stopAreaGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(expected, string(stopAreaGroupMarshal))
}

func Test_StopAreaGroupController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareStopAreaGroupRequest("GET", false, nil, t)

	// Test response
	checkStopAreaGroupResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	assert.Equal(string(expected), responseRecorder.Body.String())
}
