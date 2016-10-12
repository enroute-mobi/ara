package api

import (
	"github.com/af83/edwig/siri"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func NewTestServer(clock clockwork.Clock) *Server {
	server := Server{}
	server.SetClock(clock)
	server.startedTime = server.Clock().Now()
	return &server
}

func Test_CheckStatusHandler(t *testing.T) {
	server := NewTestServer(clockwork.NewFakeClock())

	// Set the fake clock and UUID generator
	server.SetUUIDGenerator(NewFakeUUIDGenerator())

	// Generate the request Body
	requestBody := siri.WrapSoap(siri.NewSIRICheckStatusRequest("Edwig",
		DefaultClock().Now(),
		"Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML())

	// Create a request
	request, err := http.NewRequest("POST", "/siri", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.checkStatusHandler)

	// Call ServeHTTP method and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(responseRecorder, request)

	// Check the status code is what we expect.
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "text/xml" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "text/xml")
	}

	// Check the response body is what we expect.
	response, err := siri.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if expected := "/siri"; response.Address() != expected {
		t.Errorf("Wrong Address in response:\n got: %v\n want: %v", response.Address(), expected)
	}

	if expected := "Edwig"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}

	if expected := "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: %v", response.RequestMessageRef(), expected)
	}

	if expected := "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier in response:\n got: %v\n want: %v", response.ResponseMessageIdentifier(), expected)
	}

	if !response.Status() {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status())
	}

	expectedDate := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	if !response.ResponseTimestamp().Equal(expectedDate) {
		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ResponseTimestamp(), expectedDate)
	}

	if !response.ServiceStartedTime().Equal(expectedDate) {
		t.Errorf("Wrong ServiceStartedTime in response:\n got: %v\n want: %v", response.ServiceStartedTime(), expectedDate)
	}
}
