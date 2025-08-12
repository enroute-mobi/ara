package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
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
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Name": "test" }`)
	_, responseRecorder, referential := prepareStopAreaRequest("POST", false, body, t)

	// Check response
	checkStopAreaResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// stopArea should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	stopArea, ok := referential.Model().StopAreas().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "StopArea should be found after POST request")
	assert.Equal("test", stopArea.Name)

	expectedStopArea, err := stopArea.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expectedStopArea), responseRecorder.Body.String())
}

func Test_StopAreaController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareStopAreaRequest("GET", false, nil, t)

	// Test response
	checkStopAreaResponseStatus(responseRecorder, t)

	//Test Results
	expected := `{"Models":
  [{"Origins": {},
    "Name": "First StopArea",
    "CollectChildren": false,
    "CollectSituations": false,
    "CollectedAlways": true,
    "Monitored": false,
    "Id": "6ba7b814-9dad-11d1-0-00c04fd430c8"}],
 	"Pagination": {"CurrentPage": 1, "PerPage": 1, "TotalPages": 1}}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_StopAreaController_FindStopArea(t *testing.T) {
	assert := assert.New(t)

	ref := core.NewMemoryReferentials().New("test")

	stopArea := ref.Model().StopAreas().New()
	code := model.NewCode("codeSpace", "value")
	stopArea.SetCode(code)
	ref.Model().StopAreas().Save(stopArea)

	controller := &StopAreaController{
		referential: ref,
	}

	_, ok := controller.findStopArea("codeSpace:value")
	assert.True(ok, "Can't find StopArea by Code")

	_, ok = controller.findStopArea(string(stopArea.Id()))
	assert.True(ok, "Can't find StopArea by Id")
}

func Test_StopAreaController_Index_Paginated_With_Name_Order(t *testing.T) {
	assert := assert.New(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())

	// Create and save 2 new stopAreas
	stopArea := referential.Model().StopAreas().New()
	stopArea.Name = "A"
	referential.Model().StopAreas().Save(stopArea)

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.Name = "B"
	referential.Model().StopAreas().Save(stopArea2)

	all := referential.Model().StopAreas().FindAll()
	assert.Len(all, 2)

	// Create a request
	path := path.Join("default", "stop_areas")

	params := url.Values{}
	params.Add("page", "1")
	params.Add("per_page", "2")
	params.Add("order", "name")

	u, _ := URI("", path, params)

	request, _ := http.NewRequest("GET", u.String(), nil)
	request.Header.Set("Authorization", "Token token=testToken")
	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "stop_areas")

	// Create a ResponseRecorder and send request
	responseRecorder := httptest.NewRecorder()
	server.handleReferentialModelIndex(responseRecorder, request)

	res := responseRecorder.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(err)

	var paginatedResource PaginatedResource[model.StopArea]
	err = json.Unmarshal(data, &paginatedResource)
	assert.NoError(err)

	stopAreas := paginatedResource.Models
	assert.Len(stopAreas, 2, "All stopAreas should be provided with page: 1 and per_page: 2")
	// StopAreas should be ordered by Name ascending
	assert.Equal("A", stopAreas[0].Name)
	assert.Equal("B", stopAreas[1].Name)

	// Order by direction desc
	params = url.Values{}
	params.Add("page", "1")
	params.Add("per_page", "2")
	params.Add("order", "name")
	params.Add("direction", "desc")

	u, _ = URI("", path, params)

	request, _ = http.NewRequest("GET", u.String(), nil)
	request.Header.Set("Authorization", "Token token=testToken")
	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "stop_areas")

	// Create a ResponseRecorder and send request
	responseRecorder = httptest.NewRecorder()
	server.handleReferentialModelIndex(responseRecorder, request)

	res = responseRecorder.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	assert.NoError(err)

	err = json.Unmarshal(data, &paginatedResource)
	assert.NoError(err)

	stopAreas = paginatedResource.Models
	assert.Len(stopAreas, 2, "All stopAreas should be provided with page: 1 and per_page: 2")
	// StopAreas should be ordered by Name descending
	assert.Equal("B", stopAreas[0].Name)
	assert.Equal("A", stopAreas[1].Name)
}
