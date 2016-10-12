package siri

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/checkstatus-soap-response.xml")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
	defer ts.Close()

	// Create and send request
	client := NewSOAPClient(ts.URL)
	request := &SIRICheckStatusRequest{
		RequestorRef:      "Edwig",
		RequestTimestamp:  time.Now(),
		MessageIdentifier: "Edwig:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC",
	}
	response, err := client.CheckStatus(request)
	if err != nil {
		log.Fatal(err)
	}

	// Check the content of the response
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}

	if expected := "CheckStatus:Test:0"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: %v", response.RequestMessageRef(), expected)
	}

	if expected := "c464f588-5128-46c8-ac3f-8b8a465692ab"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier in response:\n got: %v\n want: %v", response.ResponseMessageIdentifier(), expected)
	}

	if !response.Status() {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status())
	}

	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:58:34+02:00"); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ResponseTimestamp(), expected)
	}

	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T03:30:32+02:00"); !response.ServiceStartedTime().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ServiceStartedTime(), expected)
	}
}

func Test_SOAPClient_CheckStatus_GzipResponse(t *testing.T) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/checkstatus-soap-response.xml.gz")
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
		RequestTimestamp:  time.Now(),
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
