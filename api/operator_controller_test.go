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

func checkOperatorResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareOperatorRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (operator model.Operator, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
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
	referential.Model().Operators().Save(&operator)

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

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	return

}

func Test_OperatorController_Delete(t *testing.T) {
	// Send request
	operator, responseRecorder, referential := prepareOperatorRequest("DELETE", true, nil, t)

	// Test response
	checkOperatorResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().Operators().Find(operator.Id())
	if ok {
		t.Errorf("Operator shouldn't be found after DELETE request")
	}
	if expected, _ := operator.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_OperatorController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Name":"OperatorName",
                  "Objectid": { "Kind": "value" }}`)
	operator, responseRecorder, referential := prepareOperatorRequest("PUT", true, body, t)

	// Check response
	checkOperatorResponseStatus(responseRecorder, t)

	// Test Results
	updatedOperator, ok := referential.Model().Operators().Find(operator.Id())
	if !ok {
		t.Errorf("Operator should be found after PUT request")
	}

	if expected, _ := updatedOperator.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_OperatorController_Show(t *testing.T) {
	// Send request
	operator, responseRecorder, _ := prepareOperatorRequest("GET", true, nil, t)

	// Test response
	checkOperatorResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := operator.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_OperatorController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{
		"Name":"OperatorName"}`)
	_, responseRecorder, referential := prepareOperatorRequest("POST", false, body, t)

	// Check response
	checkOperatorResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// operator should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	operator, ok := referential.Model().Operators().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("Operator should be found after POST request")
	}
	operatorMarshal, _ := operator.MarshalJSON()
	expected := `{"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8","Name":"OperatorName"}`
	if responseRecorder.Body.String() != string(expected) && string(operatorMarshal) != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_OperatorController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareOperatorRequest("GET", false, nil, t)

	// Test response
	checkOperatorResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}
