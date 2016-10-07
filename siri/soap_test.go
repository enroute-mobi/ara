package siri

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/af83/edwig/api"
	"github.com/jonboulle/clockwork"
)

func Test_WrapSoap(t *testing.T) {
	expected := `<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
test
	</S:Body>
</S:Envelope>`
	if WrapSoap("test") != expected {
		t.Errorf("Error when wraping soap:\n got: %v\nwant: %v", WrapSoap("test"), expected)
	}
}

func Test_SOAPClient_CheckStatus(t *testing.T) {
	// Set the fake clock and UUID generator
	api.SetDefaultClock(clockwork.NewFakeClock())
	api.SetDefaultUUIDGenerator(api.NewFakeUUIDGenerator())

	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(CheckStatusHandler))
	defer ts.Close()

	// Create and send request
	client := NewSOAPClient(ts.URL)
	request := &SIRICheckStatusRequest{
		RequestorRef:      "Edwig",
		RequestTimestamp:  api.DefaultClock().Now(),
		MessageIdentifier: "Edwig:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC",
	}
	response, err := client.CheckStatus(request)
	if err != nil {
		log.Fatal(err)
	}

	// Check the content of the response
	if expected := "Edwig"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}

	if expected := "Edwig:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC"; response.RequestMessageRef() != expected {
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

func Test_SOAPClient_CheckStatus_GzipResponse(t *testing.T) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/checkstatus_response_compressed.xml.gz")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		w.Header().Set("Content-Encoding", "gzip")
		io.Copy(w, file)
	}))
	defer ts.Close()

	// Create and send request
	client := NewSOAPClient(ts.URL)
	request := &SIRICheckStatusRequest{
		RequestorRef:      "Edwig",
		RequestTimestamp:  api.DefaultClock().Now(),
		MessageIdentifier: "Edwig:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC",
	}
	response, err := client.CheckStatus(request)
	if err != nil {
		log.Fatal(err)
	}

	// Check a field in the response
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}
}

func Test_CheckStatusHandler(t *testing.T) {
	// Set the fake clock and UUID generator
	api.SetDefaultClock(clockwork.NewFakeClock())
	api.SetDefaultUUIDGenerator(api.NewFakeUUIDGenerator())

	// Generate the request Body
	requestBody := WrapSoap(NewSIRICheckStatusRequest("Edwig",
		api.DefaultClock().Now(),
		"Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML())

	// Create a request
	request, err := http.NewRequest("POST", "/siri", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckStatusHandler)

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
	response, err := NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
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
