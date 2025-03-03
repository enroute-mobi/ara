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

func checkSituationResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareSituationRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (situation model.Situation, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new situation
	situation = referential.Model().Situations().New()
	referential.Model().Situations().Save(&situation)

	// Create a request
	address := []byte("/default/situations")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", situation.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "situations")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(situation.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(situation.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(situation.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return
}

func Test_SituationController_Delete(t *testing.T) {
	// Send request
	situation, responseRecorder, referential := prepareSituationRequest("DELETE", true, nil, t)

	// Test response
	checkSituationResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().Situations().Find(situation.Id())
	if ok {
		t.Errorf("Situation shouldn't be found after DELETE request")
	}
	if expected, _ := situation.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_SituationController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Summary":{"DefaultValue":"Noel"},
"Codes": { "reflex": "FR:77491:ZDE:34004:STIF" },
"IgnoreValidation":true }`)
	situation, responseRecorder, referential := prepareSituationRequest("PUT", true, body, t)

	// Check response
	checkSituationResponseStatus(responseRecorder, t)

	// Test Results
	updatedSituation, ok := referential.Model().Situations().Find(situation.Id())
	if !ok {
		t.Errorf("Situation should be found after PUT request")
	}

	if expected, _ := updatedSituation.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_SituationController_Show(t *testing.T) {
	// Send request
	situation, responseRecorder, _ := prepareSituationRequest("GET", true, nil, t)

	// Test response
	checkSituationResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := situation.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_SituationController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "Summary":{"DefaultValue":"Noel"},
                "Affects" : [{"LineId":"lol","Type": "Line"}],
		"Codes": { "reflex": "FR:77491:ZDE:34004:STIF" },
                "IgnoreValidation":true }`)
	_, responseRecorder, referential := prepareSituationRequest("POST", false, body, t)

	// Check response
	checkSituationResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// situation should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	situation, ok := referential.Model().Situations().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("Situation should be found after POST request")
	}
	situationMarshal, _ := situation.MarshalJSON()
	expected := `{"Codes":{"reflex":"FR:77491:ZDE:34004:STIF"},"Origin":"","ValidityPeriods":null,"PublicationWindows":null,"Summary":{"DefaultValue":"Noel"},"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8","Affects":[{"Type":"Line","LineId":"lol"}]}`
	if responseRecorder.Body.String() != string(expected) && string(situationMarshal) != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_SituationController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareSituationRequest("GET", false, nil, t)

	// Test response
	checkSituationResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Origin":"","ValidityPeriods":null,"PublicationWindows":null,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}
