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

	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/slite"
	"github.com/stretchr/testify/assert"
)

func createHTTPLiteServer(t *testing.T, returnedFile string, opts ...int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// force status code when provided
		if len(opts) != 0 {
			w.WriteHeader(opts[0])
		}
		file, err := os.Open(fmt.Sprintf("testdata/%s.json", returnedFile))
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
}

func createHTTPServer(t *testing.T, returnedFile string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}

		file, err := os.Open(fmt.Sprintf("testdata/%s.xml", returnedFile))
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
}

func Test_SIRILiteClient_StopMonitoringDelivery(t *testing.T) {
	assert := assert.New(t)

	// Create a test http server
	ts := createHTTPLiteServer(t, "stopmonitoring-lite-delivery")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})

	dest := &slite.SIRILiteStopMonitoring{}
	stopArea := "STIF:StopPoint:Q:41178:"
	query, err := httpClient.SIRILiteStopMonitoringRequest(dest, stopArea)

	assert.Nil(err)
	assert.Equal("MonitoringRef=STIF:StopPoint:Q:41178:", query)
	assert.Equal("STIF:StopPoint:Q:41178:", dest.Siri.ServiceDelivery.StopMonitoringDelivery[0].MonitoredStopVisit[0].MonitoringRef)
}

func Test_SIRILiteClient_StopMonitoringDelivery_With_Error_400(t *testing.T) {
	assert := assert.New(t)

	// Create a test http server with return code 400 and error payload
	ts := createHTTPLiteServer(t, "stopmonitoring-lite-delivery-error", 400)
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})

	dest := &slite.SIRILiteStopMonitoring{}
	stopArea := "STIF:StopPoint:Q:41178:"
	prettyQuery, err := httpClient.SIRILiteStopMonitoringRequest(dest, stopArea)

	assert.Equal("MonitoringRef=STIF:StopPoint:Q:41178:", prettyQuery)
	assert.Error(err, "request failed with status 400: La requête contient des identifiants qui sont inconnus")
}

func Test_SIRIClient_SOAP_CheckStatus(t *testing.T) {
	// Create a test http server
	ts := createHTTPServer(t, "checkstatus-response-soap")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRIClient()
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

func Test_SIRIClient_SOAP_CheckStatus_GzipResponse(t *testing.T) {
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
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRIClient()
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

func Test_SIRIClient_SOAP_StopMonitoring(t *testing.T) {
	// Create a test http server
	ts := createHTTPServer(t, "stopmonitoring-response-soap")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRIClient()
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

func Test_SIRIClient_Raw_CheckStatus(t *testing.T) {
	// Create a test http server
	ts := createHTTPServer(t, "checkstatus-response-raw")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{SiriEnvelopeType: RAW_SIRI_ENVELOPE, Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRIClient()
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

func Test_SIRIClient_Raw_CheckStatus_GzipResponse(t *testing.T) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/checkstatus-response-raw.xml.gz")
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
	httpClient := NewHTTPClient(HTTPClientOptions{SiriEnvelopeType: RAW_SIRI_ENVELOPE, Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRIClient()
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

// TODO: must be implemented later.
// func Test_SIRIClient_Raw_StopMonitoring(t *testing.T) {
// 	// Create a test http server
// 	ts := createHTTPServer(t, "stopmonitoring-response-raw")
// 	defer ts.Close()

// 	// Create and send request
// 	httpClient := NewHTTPClient(HTTPClientOptions{SiriEnvelopeType: RAW_SIRI_ENVELOPE, Urls: HTTPClientUrls{Url: ts.URL}})
// 	client := httpClient.SIRIClient()
// 	request := &siri.SIRIGetStopMonitoringRequest{
// 		RequestorRef: "Ara",
// 	}
// 	request.MessageIdentifier = "Ara:Message::6ba7b814-9dad-11d1-32-00c04fd430c8:LOC"
// 	request.MonitoringRef = "STIF:StopArea:SP:6ba7b814-9dad-11d1-32-00c04fd430c8"
// 	request.RequestTimestamp = time.Now()

// 	response, err := client.StopMonitoring(request)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Check the content of the response
// 	if expected := "NINOXE:default"; response.ProducerRef() != expected {
// 		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
// 	}

// 	if expected := "StopMonitoring:Test:0"; response.RequestMessageRef() != expected {
// 		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: %v", response.RequestMessageRef(), expected)
// 	}

// 	if expected := "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26"; response.ResponseMessageIdentifier() != expected {
// 		t.Errorf("Wrong ResponseMessageIdentifier in response:\n got: %v\n want: %v", response.ResponseMessageIdentifier(), expected)
// 	}

// 	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T08:01:20.227+02:00"); !response.ResponseTimestamp().Equal(expected) {
// 		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ResponseTimestamp(), expected)
// 	}
// }
