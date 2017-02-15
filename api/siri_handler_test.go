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

func siriHandler_PrepareServer() (*Server, *core.Referential) {
	model.SetDefaultClock(model.NewFakeClock())
	defer model.SetDefaultClock(model.NewRealClock())

	// create a server with a fake clock and fake UUID generator
	server := NewTestServer()

	// Create the default referential with the appropriate connectors
	referential := server.CurrentReferentials().New("default")
	referential.Start()
	referential.Stop()

	partner := referential.Partners().New("partner")
	partner.Settings = map[string]string{
		"remote_url":           "",
		"remote_credential":    "",
		"remote_objectid_kind": "objectidKind",
		"local_credential":     "Edwig",
		"address":              "edwig.edwig",
	}
	partner.ConnectorTypes = []string{"siri-check-status-server", "siri-stop-monitoring-request-broadcaster"}
	partner.RefreshConnectors()
	siriPartner := core.NewSIRIPartner(partner)
	generator := core.NewFormatMessageIdentifierGenerator("Edwig:ResponseMessage::%v:LOC")
	generator.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	siriPartner.SetMessageIdentifierGenerator(generator)
	partner.Context().SetValue(core.SIRI_PARTNER, siriPartner)

	partner.Save()
	referential.Save()

	return server, referential
}

func siriHandler_Request(server *Server, soapEnvelope *siri.SOAPEnvelopeBuffer, t *testing.T) *httptest.ResponseRecorder {
	model.SetDefaultClock(model.NewFakeClock())
	defer model.SetDefaultClock(model.NewRealClock())

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

	return responseRecorder
}

func Test_SIRIHandler_CheckStatus(t *testing.T) {
	// Generate the request Body
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	request, err := siri.NewSIRICheckStatusRequest("Edwig",
		model.DefaultClock().Now(),
		"Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	soapEnvelope.WriteXML(request)

	server, _ := siriHandler_PrepareServer()
	responseRecorder := siriHandler_Request(server, soapEnvelope, t)

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

func Test_SIRIHandler_StopMonitoring(t *testing.T) {
	// Generate the request Body
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	request, err := siri.NewSIRIStopMonitoringRequest("Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC",
		"objectidValue",
		"Edwig",
		model.DefaultClock().Now()).BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	soapEnvelope.WriteXML(request)

	server, referential := siriHandler_PrepareServer()
	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "objectidValue")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.SetStopAreaId(stopArea.Id())
	stopVisit.SetObjectID(objectid)
	stopVisit.Save()

	responseRecorder := siriHandler_Request(server, soapEnvelope, t)

	// Check the response body is what we expect.
	response, err := siri.NewXMLStopMonitoringResponseFromContent(responseRecorder.Body.Bytes())
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

	expectedDate := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	if !response.ResponseTimestamp().Equal(expectedDate) {
		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ResponseTimestamp(), expectedDate)
	}
}
