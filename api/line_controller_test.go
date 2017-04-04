package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/model"
)

func checkLineResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareLineRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (line model.Line, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Save()

	// Set the fake UUID generator
	model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())
	// Save a new line
	line = referential.Model().Lines().New()
	referential.Model().Lines().Save(&line)

	// Create a request
	address := []byte("/default/lines")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", line.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call APIHandler method and pass in our Request and ResponseRecorder.
	server.APIHandler(responseRecorder, request)

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
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	line, responseRecorder, referential := prepareLineRequest("PUT", true, body, t)

	// Check response
	checkLineResponseStatus(responseRecorder, t)

	// Test Results
	updatedLine, ok := referential.Model().Lines().Find(line.Id())
	if !ok {
		t.Errorf("Line should be found after PUT request")
	}

	if expected, _ := updatedLine.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_LineController_Show(t *testing.T) {
	// Send request
	line, responseRecorder, _ := prepareLineRequest("GET", true, nil, t)

	// Test response
	checkLineResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := line.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_LineController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ 	"References" : {
		"JourneyPattern":{"ObjectId":{"lol":"lel"}, "Id":"42"}
	},
	"ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	_, responseRecorder, referential := prepareLineRequest("POST", false, body, t)

	// Check response
	checkLineResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// line should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	line, ok := referential.Model().Lines().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	if !ok {
		t.Errorf("Line should be found after POST request")
	}
	lineMarshal, _ := line.MarshalJSON()
	expected := `{"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8","ObjectIDs":{"reflex":"FR:77491:ZDE:34004:STIF"},"References":{"JourneyPattern":{"ObjectId":{"lol":"lel"},"Id":"42"}}}`
	if responseRecorder.Body.String() != string(expected) && string(lineMarshal) != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_LineController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareLineRequest("GET", false, nil, t)

	// Test response
	checkLineResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_LineController_FindLine(t *testing.T) {
	memoryModel := model.NewMemoryModel()
	tx := model.NewTransaction(memoryModel)
	defer tx.Close()

	line := memoryModel.Lines().New()
	objectid := model.NewObjectID("kind", "value")
	line.SetObjectID(objectid)
	memoryModel.Lines().Save(&line)

	controller := &LineController{}

	_, ok := controller.findLine(tx, "kind:value")
	if !ok {
		t.Error("Can't find Line by ObjectId")
	}

	_, ok = controller.findLine(tx, string(line.Id()))
	if !ok {
		t.Error("Can't find Line by Id")
	}
}
