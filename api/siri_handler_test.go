package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func siriHandler_PrepareServer(envelopeType string) (*Server, *core.Referential) {
	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// create a server with a fake clock and fake UUID generator
	server := NewTestServer()

	// Create the default referential with the appropriate connectors
	referential := server.CurrentReferentials().New("default")

	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_url":                             "",
		"remote_credential":                      "",
		"remote_objectid_kind":                   "objectidKind",
		"local_credential":                       "Ara",
		"local_url":                              "http://ara",
		"generators.message_identifier":          "Ara:Message::%{uuid}:LOC",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
		"generators.data_frame_identifier":       "RATPDev:DataFrame::%{id}:LOC",
		"siri.envelope":                          envelopeType,
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{
		"siri-check-status-server",
		"siri-stop-monitoring-request-broadcaster",
		"siri-service-request-broadcaster",
		"siri-stop-monitoring-subscription-collector",
		"siri-general-message-subscription-collector",
		"siri-estimated-timetable-request-broadcaster",
		"siri-lines-discovery-request-broadcaster",
	}
	partner.RefreshConnectors()

	partner.Save()
	referential.Save()
	referential.Start()
	referential.Stop()

	return server, referential
}

func siriHandler_Request(server *Server, buffer remote.Buffer, t *testing.T) *httptest.ResponseRecorder {
	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// Create a request
	request, err := http.NewRequest("POST", "/default/siri", buffer)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleFlow)

	// Call ServeHTTP method and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(responseRecorder, request)

	// Check the status code is what we expect.
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "text/xml; charset=utf-8" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "text/xml; charset=utf-8")
	}

	return responseRecorder
}

func Test_SIRIHandler_SOAP(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer("")
	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if !response.Status() {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status())
	}
}

func Test_SIRIHandler_SOAPResponse(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer(remote.SOAP_SIRI_ENVELOPE)
	responseRecorder := siriHandler_Request(server, buffer, t)

	_, err = remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	if err != nil {
		t.Errorf("We should receive a SOAP response, got error: %v", err)
	}
}

func Test_SIRIHandler_RawResponse(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer(remote.RAW_SIRI_ENVELOPE)
	responseRecorder := siriHandler_Request(server, buffer, t)

	_, err = remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	if err == nil {
		t.Errorf("NewSIRIEnvelope with SOAP option should return an error")
	}

	responseRecorder = siriHandler_Request(server, buffer, t) // Making the request again as the reader should be empty
	_, err = remote.NewSIRIEnvelope(responseRecorder.Body, remote.RAW_SIRI_ENVELOPE)
	if err != nil {
		t.Errorf("We shouldn't get an error while trying to create a raw envelope, got: %v", err)
	}
}

func Test_SIRIHandler_Raw(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.RAW_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer("")
	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if !response.Status() {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status())
	}
}

func Test_SIRIHandler_CheckStatus(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer("")
	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if expected := "http://ara"; response.Address() != expected {
		t.Errorf("Wrong Address in response:\n got: %v\n want: %v", response.Address(), expected)
	}

	if expected := "Ara"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}

	if expected := "Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: %v", response.RequestMessageRef(), expected)
	}

	if expected := "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; response.ResponseMessageIdentifier() != expected {
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

func Test_SIRIHandler_CheckStatus_Gzip(t *testing.T) {
	server, _ := siriHandler_PrepareServer("")

	// Create a request
	file, err := os.Open("testdata/checkstatus-soap-request.xml.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	request, err := http.NewRequest("POST", "/default/siri", file)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleFlow)

	// Call ServeHTTP method and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(responseRecorder, request)

	// Check the status code is what we expect.
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if !response.Status() {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status())
	}
}

func Test_SIRIHandler_StopMonitoring(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRIGetStopMonitoringRequest("Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC",
		"objectidValue",
		"Ara",
		clock.DefaultClock().Now()).BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(request)

	server, referential := siriHandler_PrepareServer("")
	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "objectidValue")
	stopArea.SetObjectID(objectid)
	stopArea.Monitored = true
	stopArea.Save()

	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Monitored = true
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(2*time.Hour))
	stopVisit.SetObjectID(model.NewObjectID("objectidKind", "second"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Save()

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(1*time.Hour))
	stopVisit2.SetObjectID(model.NewObjectID("objectidKind", "first"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.Save()

	pastStopVisit := referential.Model().StopVisits().New()
	pastStopVisit.StopAreaId = stopArea.Id()
	pastStopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(-1*time.Hour))
	pastStopVisit.SetObjectID(model.NewObjectID("objectidKind", "past"))
	pastStopVisit.VehicleJourneyId = vehicleJourney.Id()
	pastStopVisit.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLStopMonitoringResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	delivery := response.StopMonitoringDeliveries()[0]
	if len(delivery.XMLMonitoredStopVisits()) != 2 {
		t.Fatalf("Past StopVisit should be ignored, got %v stopVisits", len(delivery.XMLMonitoredStopVisits()))
	}

	if !delivery.XMLMonitoredStopVisits()[1].ActualArrivalTime().After(delivery.XMLMonitoredStopVisits()[0].ActualArrivalTime()) {
		t.Errorf("Stop visits are not chronollogicaly ordered ")
	}

	if expected := "http://ara"; response.Address() != expected {
		t.Errorf("Wrong Address in response:\n got: %v\n want: %v", response.Address(), expected)
	}

	if expected := "Ara"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: %v", response.ProducerRef(), expected)
	}

	if expected := "Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: %v", response.RequestMessageRef(), expected)
	}

	if expected := "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier in response:\n got: %v\n want: %v", response.ResponseMessageIdentifier(), expected)
	}

	expectedDate := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	if !response.ResponseTimestamp().Equal(expectedDate) {
		t.Errorf("Wrong ResponseTimestamp in response:\n got: %v\n want: %v", response.ResponseTimestamp(), expectedDate)
	}
}

func Test_SIRIHandler_SiriService(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/siri-service-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")
	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "stopArea1")
	stopArea.SetObjectID(objectid)
	stopArea.Monitored = true
	stopArea.Save()

	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Monitored = true
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(2*time.Hour))
	stopVisit.SetObjectID(model.NewObjectID("objectidKind", "second"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Save()

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(1*time.Hour))
	stopVisit2.SetObjectID(model.NewObjectID("objectidKind", "first"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.Save()

	pastStopVisit := referential.Model().StopVisits().New()
	pastStopVisit.StopAreaId = stopArea.Id()
	pastStopVisit.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(-1*time.Hour))
	pastStopVisit.SetObjectID(model.NewObjectID("objectidKind", "past"))
	pastStopVisit.VehicleJourneyId = vehicleJourney.Id()
	pastStopVisit.Save()

	stopArea2 := referential.Model().StopAreas().New()
	objectid2 := model.NewObjectID("objectidKind", "stopArea2")
	stopArea2.SetObjectID(objectid2)
	stopArea2.Monitored = true
	stopArea2.Save()

	line2 := referential.Model().Lines().New()
	line2.SetObjectID(objectid2)
	line2.Save()

	vehicleJourney2 := referential.Model().VehicleJourneys().New()
	vehicleJourney2.SetObjectID(objectid2)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Monitored = true
	vehicleJourney2.Save()

	stopVisit3 := referential.Model().StopVisits().New()
	stopVisit3.StopAreaId = stopArea2.Id()
	stopVisit3.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(2*time.Hour))
	stopVisit3.SetObjectID(model.NewObjectID("objectidKind", "third"))
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.Save()

	stopVisit4 := referential.Model().StopVisits().New()
	stopVisit4.StopAreaId = stopArea2.Id()
	stopVisit4.Schedules.SetArrivalTime(model.STOP_VISIT_SCHEDULE_ACTUAL, referential.Clock().Now().Add(1*time.Hour))
	stopVisit4.SetObjectID(model.NewObjectID("objectidKind", "fourth"))
	stopVisit4.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit4.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	// responseRecorder.Body.String()
	envelope, err := remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}
	responseBody := envelope.Body().String()

	// TEMP: Find a better way to test?
	expectedResponseBody := `<sw:GetSiriServiceResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer>
		<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>Ara</siri:ProducerRef>
		<siri:ResponseMessageIdentifier>Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
		<siri:Status>true</siri:Status>
		<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
			<siri:MonitoringRef>stopArea1</siri:MonitoringRef>
			<siri:Status>true</siri:Status>
			<siri:MonitoredStopVisit>
				<siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
				<siri:ItemIdentifier>first</siri:ItemIdentifier>
				<siri:MonitoringRef>stopArea1</siri:MonitoringRef>
				<siri:MonitoredVehicleJourney>
					<siri:LineRef>stopArea1</siri:LineRef>
					<siri:FramedVehicleJourneyRef>
						<siri:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</siri:DataFrameRef>
						<siri:DatedVehicleJourneyRef>stopArea1</siri:DatedVehicleJourneyRef>
					</siri:FramedVehicleJourneyRef>
					<siri:Monitored>true</siri:Monitored>
					<siri:MonitoredCall>
						<siri:StopPointRef>stopArea1</siri:StopPointRef>
						<siri:VehicleAtStop>false</siri:VehicleAtStop>
						<siri:ActualArrivalTime>1984-04-04T01:00:00.000Z</siri:ActualArrivalTime>
					</siri:MonitoredCall>
				</siri:MonitoredVehicleJourney>
			</siri:MonitoredStopVisit>
			<siri:MonitoredStopVisit>
				<siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
				<siri:ItemIdentifier>second</siri:ItemIdentifier>
				<siri:MonitoringRef>stopArea1</siri:MonitoringRef>
				<siri:MonitoredVehicleJourney>
					<siri:LineRef>stopArea1</siri:LineRef>
					<siri:FramedVehicleJourneyRef>
						<siri:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</siri:DataFrameRef>
						<siri:DatedVehicleJourneyRef>stopArea1</siri:DatedVehicleJourneyRef>
					</siri:FramedVehicleJourneyRef>
					<siri:Monitored>true</siri:Monitored>
					<siri:MonitoredCall>
						<siri:StopPointRef>stopArea1</siri:StopPointRef>
						<siri:VehicleAtStop>false</siri:VehicleAtStop>
						<siri:ActualArrivalTime>1984-04-04T02:00:00.000Z</siri:ActualArrivalTime>
					</siri:MonitoredCall>
				</siri:MonitoredVehicleJourney>
			</siri:MonitoredStopVisit>
		</siri:StopMonitoringDelivery>
		<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>GetSIRIStopMonitoring:Test:0</siri:RequestMessageRef>
			<siri:MonitoringRef>stopArea2</siri:MonitoringRef>
			<siri:Status>true</siri:Status>
			<siri:MonitoredStopVisit>
				<siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
				<siri:ItemIdentifier>fourth</siri:ItemIdentifier>
				<siri:MonitoringRef>stopArea2</siri:MonitoringRef>
				<siri:MonitoredVehicleJourney>
					<siri:LineRef>stopArea2</siri:LineRef>
					<siri:FramedVehicleJourneyRef>
						<siri:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</siri:DataFrameRef>
						<siri:DatedVehicleJourneyRef>stopArea2</siri:DatedVehicleJourneyRef>
					</siri:FramedVehicleJourneyRef>
					<siri:Monitored>true</siri:Monitored>
					<siri:MonitoredCall>
						<siri:StopPointRef>stopArea2</siri:StopPointRef>
						<siri:VehicleAtStop>false</siri:VehicleAtStop>
						<siri:ActualArrivalTime>1984-04-04T01:00:00.000Z</siri:ActualArrivalTime>
					</siri:MonitoredCall>
				</siri:MonitoredVehicleJourney>
			</siri:MonitoredStopVisit>
			<siri:MonitoredStopVisit>
				<siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
				<siri:ItemIdentifier>third</siri:ItemIdentifier>
				<siri:MonitoringRef>stopArea2</siri:MonitoringRef>
				<siri:MonitoredVehicleJourney>
					<siri:LineRef>stopArea2</siri:LineRef>
					<siri:FramedVehicleJourneyRef>
						<siri:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</siri:DataFrameRef>
						<siri:DatedVehicleJourneyRef>stopArea2</siri:DatedVehicleJourneyRef>
					</siri:FramedVehicleJourneyRef>
					<siri:Monitored>true</siri:Monitored>
					<siri:MonitoredCall>
						<siri:StopPointRef>stopArea2</siri:StopPointRef>
						<siri:VehicleAtStop>false</siri:VehicleAtStop>
						<siri:ActualArrivalTime>1984-04-04T02:00:00.000Z</siri:ActualArrivalTime>
					</siri:MonitoredCall>
				</siri:MonitoredVehicleJourney>
			</siri:MonitoredStopVisit>
		</siri:StopMonitoringDelivery>
	</Answer>
	<AnswerExtension/>
</sw:GetSiriServiceResponse>`

	// Check the response body is what we expect.
	if responseBody != expectedResponseBody {
		t.Errorf("Unexpected response body:\n expected: %v\n got: %v", expectedResponseBody, responseBody)
	}
}

func Test_SIRIHandler_NotifyStopMonitoring(t *testing.T) {
	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/notify-stop-monitoring.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")
	partner := referential.Partners().FindAll()[0]

	partner.Subscriptions().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	subscription := partner.Subscriptions().FindOrCreateByKind(core.StopMonitoringCollect)
	subscription.Save()

	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "stopArea1")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	objectid2 := model.NewObjectID("objectidKind", "stopArea2")
	stopArea2.SetObjectID(objectid2)
	stopArea2.Save()

	siriHandler_Request(server, buffer, t)

	// Some Tests

	if count := len(referential.Model().StopVisits().FindAll()); count != 3 {
		t.Errorf("Notify should have created 3 StopVisits, got: %v", count)
	}
}

func Test_SIRIHandler_NotifyGeneralMessage(t *testing.T) {
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/notify-general-message.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")
	partner := referential.Partners().FindAll()[0]

	partner.Subscriptions().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	subscription := partner.Subscriptions().FindOrCreateByKind(core.GeneralMessageCollect)
	subscription.Save()

	siriHandler_Request(server, buffer, t)

	if count := len(referential.Model().Situations().FindAll()); count != 2 {
		t.Errorf("Notify should have created 2 Situation, got: %v", count)
	}
}

func Test_SIRIHandler_EstimatedTimetable(t *testing.T) {
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/estimated_timetable_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(model.NewObjectID("objectidKind", "stopArea1"))
	stopArea.Monitored = true
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(model.NewObjectID("objectidKind", "stopArea2"))
	stopArea2.Monitored = true
	stopArea2.Save()

	line := referential.Model().Lines().New()
	line.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.Model().Lines().New()
	line2.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney"))
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	vehicleJourney2 := referential.Model().VehicleJourneys().New()
	vehicleJourney2.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney2"))
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	pastStopVisit := referential.Model().StopVisits().New()
	pastStopVisit.SetObjectID(model.NewObjectID("objectidKind", "pastStopVisit"))
	pastStopVisit.VehicleJourneyId = vehicleJourney.Id()
	pastStopVisit.StopAreaId = stopArea.Id()
	pastStopVisit.PassageOrder = 0
	pastStopVisit.ArrivalStatus = "onTime"
	pastStopVisit.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(-1*time.Minute))
	pastStopVisit.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(-1*time.Minute))
	pastStopVisit.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.SetObjectID(model.NewObjectID("objectidKind", "stopVisit"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.PassageOrder = 1
	stopVisit.ArrivalStatus = "onTime"
	stopVisit.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.SetObjectID(model.NewObjectID("objectidKind", "stopVisit2"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.StopAreaId = stopArea2.Id()
	stopVisit2.PassageOrder = 2
	stopVisit2.ArrivalStatus = "onTime"
	stopVisit2.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(2*time.Minute))
	stopVisit2.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(2*time.Minute))
	stopVisit2.Save()

	stopVisit3 := referential.Model().StopVisits().New()
	stopVisit3.SetObjectID(model.NewObjectID("objectidKind", "stopVisit3"))
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.StopAreaId = stopArea.Id()
	stopVisit3.PassageOrder = 1
	stopVisit3.ArrivalStatus = "onTime"
	stopVisit3.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(1*time.Minute))
	stopVisit3.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(1*time.Minute))
	stopVisit3.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	envelope, err := remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}
	responseBody := envelope.Body().String()

	expectedResponseBody := `<sw:GetEstimatedTimetableResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>Ara</siri:ProducerRef>
		<siri:Address>http://ara</siri:Address>
		<siri:ResponseMessageIdentifier>Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>EstimatedTimetable:Test:0</siri:RequestMessageRef>
			<siri:Status>true</siri:Status>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>1984-04-04T00:00:00.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>NINOXE:Line:2:LOC</siri:LineRef>
					<siri:DirectionRef/>
					<siri:OperatorRef/>
					<siri:DatedVehicleJourneyRef>vehicleJourney</siri:DatedVehicleJourneyRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopArea1</siri:StopPointRef>
							<siri:Order>1</siri:Order>
							<siri:VehicleAtStop>false</siri:VehicleAtStop>
							<siri:AimedArrivalTime>1984-04-04T00:01:00.000Z</siri:AimedArrivalTime>
							<siri:ExpectedArrivalTime>1984-04-04T00:01:00.000Z</siri:ExpectedArrivalTime>
							<siri:ArrivalStatus>onTime</siri:ArrivalStatus>
						</siri:EstimatedCall>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopArea2</siri:StopPointRef>
							<siri:Order>2</siri:Order>
							<siri:VehicleAtStop>false</siri:VehicleAtStop>
							<siri:AimedArrivalTime>1984-04-04T00:02:00.000Z</siri:AimedArrivalTime>
							<siri:ExpectedArrivalTime>1984-04-04T00:02:00.000Z</siri:ExpectedArrivalTime>
							<siri:ArrivalStatus>onTime</siri:ArrivalStatus>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>1984-04-04T00:00:00.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
					<siri:DirectionRef/>
					<siri:OperatorRef/>
					<siri:DatedVehicleJourneyRef>vehicleJourney2</siri:DatedVehicleJourneyRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopArea1</siri:StopPointRef>
							<siri:Order>1</siri:Order>
							<siri:VehicleAtStop>false</siri:VehicleAtStop>
							<siri:AimedArrivalTime>1984-04-04T00:01:00.000Z</siri:AimedArrivalTime>
							<siri:ExpectedArrivalTime>1984-04-04T00:01:00.000Z</siri:ExpectedArrivalTime>
							<siri:ArrivalStatus>onTime</siri:ArrivalStatus>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
		</siri:EstimatedTimetableDelivery>
	</Answer>
	<AnswerExtension/>
</sw:GetEstimatedTimetableResponse>`

	// Check the response body is what we expect.
	if responseBody != expectedResponseBody {
		t.Errorf("Unexpected response body:\n expected: %v\n got: %v", expectedResponseBody, responseBody)
	}
}

func Test_SIRIHandler_LinesDiscovery(t *testing.T) {
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/lines-discovery-request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")

	line := referential.Model().Lines().New()
	line.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.Model().Lines().New()
	line2.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	line3 := referential.Model().Lines().New()
	line3.SetObjectID(model.NewObjectID("objectidKind2", "NINOXE:Line:4:LOC"))
	line3.Name = "lineName3"
	line3.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	envelope, err := remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}
	responseBody := envelope.Body().String()

	expectedResponseBody := `<sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer version="2.0">
		<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
		<siri:Status>true</siri:Status>
		<siri:AnnotatedLineRef>
			<siri:LineRef>NINOXE:Line:2:LOC</siri:LineRef>
			<siri:LineName>lineName</siri:LineName>
			<siri:Monitored>true</siri:Monitored>
		</siri:AnnotatedLineRef>
		<siri:AnnotatedLineRef>
			<siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
			<siri:LineName>lineName2</siri:LineName>
			<siri:Monitored>true</siri:Monitored>
		</siri:AnnotatedLineRef>
	</Answer>
	<AnswerExtension/>
</sw:LinesDiscoveryResponse>`

	// Check the response body is what we expect.
	if responseBody != expectedResponseBody {
		t.Errorf("Unexpected response body:\n expected: %v\n got: %v", expectedResponseBody, responseBody)
	}
}
