package remote

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri"
)

func testSOAPFile(name string) (*os.File, error) {
	// Create a new SOAPEnvelope
	file, err := os.Open(fmt.Sprintf("testdata/%s-soap.xml", name))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func createHTTPServer(t *testing.T, returnedFile string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}

		file, err := testSOAPFile(returnedFile)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
}

func Test_SOAPClient_CheckStatus(t *testing.T) {
	// Create a test http server
	ts := createHTTPServer(t, "checkstatus-response")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientUrls{Url: ts.URL})
	client := httpClient.SOAPClient()
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      "Ara",
		RequestTimestamp:  time.Now(),
		MessageIdentifier: "Ara:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC",
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
		t.Errorf("Wrong ServiceStartedTime in response:\n got: %v\n want: %v", response.ServiceStartedTime(), expected)
	}
}

func Test_SOAPClient_CheckStatus_GzipResponse(t *testing.T) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/checkstatus-response-soap.xml.gz")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/xml")
		io.Copy(w, file)
	}))
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientUrls{Url: ts.URL})
	client := httpClient.SOAPClient()
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      "Ara",
		RequestTimestamp:  time.Now(),
		MessageIdentifier: "Ara:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC",
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

func Test_SOAPClient_StopMonitoring(t *testing.T) {
	// Create a test http server
	ts := createHTTPServer(t, "stopmonitoring-response")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientUrls{Url: ts.URL})
	client := httpClient.SOAPClient()
	request := &siri.SIRIGetStopMonitoringRequest{
		RequestorRef: "Ara",
	}
	request.MessageIdentifier = "Ara:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC"
	request.MonitoringRef = "STIF:StopArea:SP:6ba7b814-9dad-11d1-32-00c04fd430c8"
	request.RequestTimestamp = time.Now()

	response, err := client.StopMonitoring(request)
	if err != nil {
		log.Fatal(err)
	}

	// Check the content of the response
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}

	if expected := "StopMonitoring:Test:0"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: %v", response.RequestMessageRef(), expected)
	}

	if expected := "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier in response:\n got: %v\n want: %v", response.ResponseMessageIdentifier(), expected)
	}

	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T08:01:20.227+02:00"); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ResponseTimestamp(), expected)
	}
}
