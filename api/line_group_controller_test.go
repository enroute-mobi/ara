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

func checkLineGroupResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)
	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareLineGroupRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (lineGroup *model.LineGroup, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new lineGroup
	lineGroup = referential.Model().LineGroups().New()
	referential.Model().LineGroups().Save(lineGroup)

	// Create a request
	address := []byte("/default/line_groups")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", lineGroup.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "line_groups")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(lineGroup.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(lineGroup.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(lineGroup.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return

}

func Test_LineGroupController_Delete(t *testing.T) {
	assert := assert.New(t)

	// Send request
	lineGroup, responseRecorder, referential := prepareLineGroupRequest("DELETE", true, nil, t)

	// Test response
	checkLineGroupResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().LineGroups().Find(lineGroup.Id())
	assert.False(ok, "LineGroup shouldn't be found after DELETE request")

	expected, err := lineGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(expected), responseRecorder.Body.String())
}

func Test_LineGroupController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Name":"LineGroupName", "ShortName":"short_name"}`)
	lineGroup, responseRecorder, referential := prepareLineGroupRequest("PUT", true, body, t)

	// Check response
	checkLineGroupResponseStatus(responseRecorder, t)

	// Test Results
	updatedLineGroup, ok := referential.Model().LineGroups().Find(lineGroup.Id())
	assert.True(ok, "LineGroup should be found after PUT request")

	expected, err := updatedLineGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(expected), responseRecorder.Body.String())
}

func Test_LineGroupController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	lineGroup, responseRecorder, _ := prepareLineGroupRequest("GET", true, nil, t)

	// Test response
	checkLineGroupResponseStatus(responseRecorder, t)

	//Test Results
	expected, err := lineGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(string(expected), responseRecorder.Body.String()) //
}

func Test_LineGroupController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{
		"Name":"LineGroupName"}`)
	_, responseRecorder, referential := prepareLineGroupRequest("POST", false, body, t)

	// Check response
	checkLineGroupResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// lineGroup should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	lineGroup, ok := referential.Model().LineGroups().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "LineGroup should be found after POST request")

	expected := `{"Name":"LineGroupName","Id":"6ba7b814-9dad-11d1-1-00c04fd430c8"}`
	lineGroupMarshal, err := lineGroup.MarshalJSON()
	assert.NoError(err)
	assert.Equal(expected, string(lineGroupMarshal))
}

func Test_LineGroupController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareLineGroupRequest("GET", false, nil, t)

	// Test response
	checkLineGroupResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	assert.Equal(string(expected), responseRecorder.Body.String())
}
