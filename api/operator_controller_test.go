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

func checkOperatorResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareOperatorRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (operator *model.Operator, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new operator
	operator = referential.Model().Operators().New()
	referential.Model().Operators().Save(operator)

	// Create a request
	address := []byte("/default/operators")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", operator.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "operators")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(operator.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(operator.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(operator.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return

}

func Test_OperatorController_Delete(t *testing.T) {
	assert := assert.New(t)

	// Send request
	operator, responseRecorder, referential := prepareOperatorRequest("DELETE", true, nil, t)

	// Test response
	checkOperatorResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().Operators().Find(operator.Id())
	assert.False(ok, "Operator shouldn't be found after DELETE request")

	expected, err := operator.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_OperatorController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Name":"OperatorName",
                  "Code": { "CodeSpace": "value" }}`)
	operator, responseRecorder, referential := prepareOperatorRequest("PUT", true, body, t)

	// Check response
	checkOperatorResponseStatus(responseRecorder, t)

	// Test Results
	updatedOperator, ok := referential.Model().Operators().Find(operator.Id())
	assert.True(ok, "Operator should be found after PUT request")

	expected, err := updatedOperator.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_OperatorController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	operator, responseRecorder, _ := prepareOperatorRequest("GET", true, nil, t)

	// Test response
	checkOperatorResponseStatus(responseRecorder, t)

	//Test Results
	expectedOperator, _ := operator.MarshalJSON()
	assert.JSONEq(string(expectedOperator), responseRecorder.Body.String())
}

func Test_OperatorController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{
		"Name":"OperatorName"}`)
	_, responseRecorder, referential := prepareOperatorRequest("POST", false, body, t)

	// Check response
	checkOperatorResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// operator should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	_, ok := referential.Model().Operators().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "Operator should be found after POST request")

	expected := `{"Name":"OperatorName","Id":"6ba7b814-9dad-11d1-1-00c04fd430c8"}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_OperatorController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareOperatorRequest("GET", false, nil, t)

	// Test response
	checkOperatorResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	assert.JSONEq(expected, responseRecorder.Body.String())
}
