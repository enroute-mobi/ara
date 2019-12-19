package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/edwig/model"
)

func timeCheckResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func timePrepareRequest(method string, advance bool, body []byte, t *testing.T) (responseRecorder *httptest.ResponseRecorder, server *Server) {
	server = &Server{}
	server.SetClock(model.NewFakeClock())

	// Create a request
	var address string
	if advance {
		address = "/_time/advance"
	} else {
		address = "/_time"
	}

	request, err := http.NewRequest(method, address, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	return
}

func Test_TimeController_Get(t *testing.T) {
	// Send request
	responseRecorder, _ := timePrepareRequest("GET", false, nil, t)

	// Test response
	timeCheckResponseStatus(responseRecorder, t)

	//Test Results
	if expected := `{ "time": "1984-04-04T00:00:00.000Z" }`; responseRecorder.Body.String() != expected {
		t.Errorf("Wrong body for GET response request:\n got: %v\n want: %v", responseRecorder.Body.String(), expected)
	}
}

func Test_TimeController_Advance(t *testing.T) {
	// Send request
	body := []byte(`{ "duration": "1s" }`)
	responseRecorder, server := timePrepareRequest("POST", true, body, t)

	// Test response
	timeCheckResponseStatus(responseRecorder, t)

	//Test Results
	if expected := `{ "time": "1984-04-04T00:00:01.000Z" }`; responseRecorder.Body.String() != expected {
		t.Errorf("Wrong body for GET response request:\n got: %v\n want: %v", responseRecorder.Body.String(), expected)
	}
	if expected := time.Date(1984, time.April, 4, 0, 0, 1, 0, time.UTC); !server.Clock().Now().Equal(expected) {
		t.Errorf("Server Clock should have advanced:\n got: %v\n want: %v", server.Clock().Now(), expected)
	}
}
