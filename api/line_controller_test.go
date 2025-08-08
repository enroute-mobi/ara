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

func checkLineResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareLineRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (line *model.Line, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new line
	line = referential.Model().Lines().New()
	referential.Model().Lines().Save(line)

	// Create a request
	address := []byte("/default/lines")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", line.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "lines")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(line.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(line.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(line.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return
}

func Test_LineController_Delete(t *testing.T) {
	// Send request
	line, responseRecorder, referential := prepareLineRequest("DELETE", true, nil, t)

	// Test response
	checkLineResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().Lines().Find(line.Id())
	if ok {
		t.Errorf("Line shouldn't be found after DELETE request")
	}
	if expected, _ := line.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_LineController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Codes": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	line, responseRecorder, referential := prepareLineRequest("PUT", true, body, t)

	// Check response
	checkLineResponseStatus(responseRecorder, t)

	updatedLine, ok := referential.Model().Lines().Find(line.Id())
	assert.True(ok, "Line should be found after PUT request")

	expectedLine, err := updatedLine.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expectedLine), responseRecorder.Body.String())
}

func Test_LineController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	line, responseRecorder, _ := prepareLineRequest("GET", true, nil, t)

	// Test response
	checkLineResponseStatus(responseRecorder, t)

	expectedLine, err := line.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expectedLine), responseRecorder.Body.String())
}

func Test_LineController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ 	"References" : {
		"JourneyPattern":{"Code":{"lol":"lel"}, "Id":"42"}
	},
	"Codes": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	_, responseRecorder, referential := prepareLineRequest("POST", false, body, t)

	// Check response
	checkLineResponseStatus(responseRecorder, t)

	// Using the fake uuid generator, the uuid of the created
	// line should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	_, ok := referential.Model().Lines().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "Line should be found after POST request")

	expected := `{"CollectSituations":false,"Codes":{"reflex":"FR:77491:ZDE:34004:STIF"},"References":{"JourneyPattern":{"Code":{"lol":"lel"}}},"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8"}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_LineController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareLineRequest("GET", false, nil, t)

	// Test response
	checkLineResponseStatus(responseRecorder, t)

	//Test Result
	expected := `{"Models":
  [{"CollectSituations": false,
    "Id": "6ba7b814-9dad-11d1-0-00c04fd430c8"}],
    "Pagination": {"CurrentPage": 1, "PerPage": 1, "TotalCount": 1, "TotalPages": 1}}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_LineController_FindLine(t *testing.T) {
	assert := assert.New(t)

	ref := core.NewMemoryReferentials().New("test")

	line := ref.Model().Lines().New()
	code := model.NewCode("codeSpace", "stif:value")
	line.SetCode(code)
	ref.Model().Lines().Save(line)

	controller := &LineController{
		referential: ref,
	}

	_, ok := controller.findLine("codeSpace:stif:value")
	assert.True(ok, "Can't find Line by Code")

	_, ok = controller.findLine(string(line.Id()))
	assert.True(ok, "Can't find Line by Id")
}
