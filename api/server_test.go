package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func NewTestServer(clock model.Clock) *Server {
	server := Server{}
	referentials := core.NewMemoryReferentials()
	server.SetReferentials(referentials)
	server.SetClock(clock)
	server.startedTime = server.Clock().Now()
	return &server
}

func Test_CheckStatusHandler(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())
	defer model.SetDefaultClock(model.NewRealClock())

	// create a server with a fake clock and fake UUID generator
	server := NewTestServer(model.NewFakeClock())

	// Create the default referential with the appropriate connectors
	referential := server.CurrentReferentials().New("default")
	referential.Start()
	referential.Stop()

	partner := referential.Partners().New("partner")
	partner.Settings = map[string]string{
		"remote_url":        "",
		"remote_credential": "",
		"local_credential":  "Edwig",
		"Address":           "edwig.edwig",
	}
	partner.ConnectorTypes = []string{"siri-check-status-client"}
	partner.RefreshConnectors()
	siriPartner := core.NewSIRIPartner(partner)
	generator := core.NewFormatMessageIdentifierGenerator("Edwig:ResponseMessage::%v:LOC")
	generator.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	siriPartner.SetMessageIdentifierGenerator(generator)
	partner.Context().SetValue(core.SIRI_PARTNER, siriPartner)

	partner.Save()

	referential.Save()

	// Generate the request Body
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(siri.NewSIRICheckStatusRequest("Edwig",
		model.DefaultClock().Now(),
		"Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML())

	// Create a request
	request, err := http.NewRequest("POST", "/default/siri", soapEnvelope)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.APIHandler)

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

	if expected := "edwig.edwig"; response.Address() != expected {
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
