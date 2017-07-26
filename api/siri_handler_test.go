package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/core"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SIRIHandler_LogSIRIError(t *testing.T) {
	logstash := audit.NewFakeLogStash()
	audit.SetCurrentLogstash(logstash)
	defer audit.SetCurrentLogstash(audit.NewNullLogStash())

	siriError("errCode", "errDescription", httptest.NewRecorder())

	if len(logstash.Events()) != 1 {
		t.Errorf("Logstash should have one event, got %d", len(logstash.Events()))
	}
}

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
		"remote_url":                             "",
		"remote_credential":                      "",
		"remote_objectid_kind":                   "objectidKind",
		"local_credential":                       "Edwig",
		"local_url":                              "http://edwig",
		"generators.message_identifier":          "Edwig:Message::%{uuid}:LOC",
		"generators.response_message_identifier": "Edwig:ResponseMessage::%{uuid}:LOC",
		"generators.data_frame_identifier":       "RATPDev:DataFrame::%{id}:LOC",
	}
	partner.ConnectorTypes = []string{
		"siri-check-status-server",
		"siri-stop-monitoring-request-broadcaster",
		"siri-service-request-broadcaster",
		"siri-stop-monitoring-subscription-collector",
		"siri-general-message-subscription-collector",
	}
	partner.RefreshConnectors()
	siriPartner := core.NewSIRIPartner(partner)

	siriPartner.SetUUIDGenerator(model.NewFakeUUIDGenerator())
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

	if expected := "http://edwig"; response.Address() != expected {
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

func Test_SIRIHandler_CheckStatus_Gzip(t *testing.T) {
	server, _ := siriHandler_PrepareServer()

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
	response, err := siri.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if !response.Status() {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status())
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

	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.LineId = line.Id()
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

	responseRecorder := siriHandler_Request(server, soapEnvelope, t)

	// Check the response body is what we expect.
	response, err := siri.NewXMLStopMonitoringResponseFromContent(responseRecorder.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if len(response.XMLMonitoredStopVisits()) != 2 {
		t.Errorf("Past StopVisit should be ignored")
	}

	if !response.XMLMonitoredStopVisits()[1].ActualArrivalTime().After(response.XMLMonitoredStopVisits()[0].ActualArrivalTime()) {
		t.Errorf("Stop visits are not chronollogicaly ordered ")
	}

	if expected := "http://edwig"; response.Address() != expected {
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

func Test_SIRIHandler_SiriService(t *testing.T) {
	// Generate the request Body
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()

	file, err := os.Open("testdata/siri-service-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	soapEnvelope.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer()
	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "stopArea1")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.LineId = line.Id()
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
	stopArea2.Save()

	line2 := referential.Model().Lines().New()
	line2.SetObjectID(objectid2)
	line2.Save()

	vehicleJourney2 := referential.Model().VehicleJourneys().New()
	vehicleJourney2.SetObjectID(objectid2)
	vehicleJourney2.LineId = line2.Id()
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

	responseRecorder := siriHandler_Request(server, soapEnvelope, t)

	// responseRecorder.Body.String()
	envelope, err := siri.NewSOAPEnvelope(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	responseBody := envelope.Body().String()

	// TEMP: Find a better way to test?
	expectedResponseBody := `<ns1:GetSiriServiceResponse xmlns:ns1="http://wsdl.siri.org.uk">
	<Answer xmlns:ns3="http://www.siri.org.uk/siri" xmlns:ns4="http://www.ifopt.org.uk/acsb" xmlns:ns5="http://www.ifopt.org.uk/ifopt" xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns7="http://scma/siri" xmlns:ns8="http://wsdl.siri.org.uk" xmlns:ns9="http://wsdl.siri.org.uk/siri">
		<ns3:ResponseTimestamp>1984-04-04T00:00:00.000Z</ns3:ResponseTimestamp>
		<ns3:ProducerRef>Edwig</ns3:ProducerRef>
		<ns3:ResponseMessageIdentifier>Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns3:RequestMessageRef>
		<ns3:Status>true</ns3:Status>
		<ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>1984-04-04T00:00:00.000Z</ns3:ResponseTimestamp>
			<ns3:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns3:RequestMessageRef>
			<ns3:Status>true</ns3:Status>
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>first</ns3:ItemIdentifier>
				<ns3:MonitoringRef>stopArea1</ns3:MonitoringRef>
				<ns3:MonitoredVehicleJourney>
					<ns3:LineRef>stopArea1</ns3:LineRef>
					<ns3:FramedVehicleJourneyRef>
						<ns3:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</ns3:DataFrameRef>
						<ns3:DatedVehicleJourneyRef>stopArea1</ns3:DatedVehicleJourneyRef>
					</ns3:FramedVehicleJourneyRef>
					<ns3:MonitoredCall>
						<ns3:StopPointRef>stopArea1</ns3:StopPointRef>
						<ns3:VehicleAtStop>false</ns3:VehicleAtStop>
						<ns3:ActualArrivalTime>1984-04-04T01:00:00.000Z</ns3:ActualArrivalTime>
					</ns3:MonitoredCall>
				</ns3:MonitoredVehicleJourney>
			</ns3:MonitoredStopVisit>
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>second</ns3:ItemIdentifier>
				<ns3:MonitoringRef>stopArea1</ns3:MonitoringRef>
				<ns3:MonitoredVehicleJourney>
					<ns3:LineRef>stopArea1</ns3:LineRef>
					<ns3:FramedVehicleJourneyRef>
						<ns3:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</ns3:DataFrameRef>
						<ns3:DatedVehicleJourneyRef>stopArea1</ns3:DatedVehicleJourneyRef>
					</ns3:FramedVehicleJourneyRef>
					<ns3:MonitoredCall>
						<ns3:StopPointRef>stopArea1</ns3:StopPointRef>
						<ns3:VehicleAtStop>false</ns3:VehicleAtStop>
						<ns3:ActualArrivalTime>1984-04-04T02:00:00.000Z</ns3:ActualArrivalTime>
					</ns3:MonitoredCall>
				</ns3:MonitoredVehicleJourney>
			</ns3:MonitoredStopVisit>
		</ns3:StopMonitoringDelivery>
		<ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>1984-04-04T00:00:00.000Z</ns3:ResponseTimestamp>
			<ns3:RequestMessageRef>GetSIRIStopMonitoring:Test:0</ns3:RequestMessageRef>
			<ns3:Status>true</ns3:Status>
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>fourth</ns3:ItemIdentifier>
				<ns3:MonitoringRef>stopArea2</ns3:MonitoringRef>
				<ns3:MonitoredVehicleJourney>
					<ns3:LineRef>stopArea2</ns3:LineRef>
					<ns3:FramedVehicleJourneyRef>
						<ns3:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</ns3:DataFrameRef>
						<ns3:DatedVehicleJourneyRef>stopArea2</ns3:DatedVehicleJourneyRef>
					</ns3:FramedVehicleJourneyRef>
					<ns3:MonitoredCall>
						<ns3:StopPointRef>stopArea2</ns3:StopPointRef>
						<ns3:VehicleAtStop>false</ns3:VehicleAtStop>
						<ns3:ActualArrivalTime>1984-04-04T01:00:00.000Z</ns3:ActualArrivalTime>
					</ns3:MonitoredCall>
				</ns3:MonitoredVehicleJourney>
			</ns3:MonitoredStopVisit>
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>third</ns3:ItemIdentifier>
				<ns3:MonitoringRef>stopArea2</ns3:MonitoringRef>
				<ns3:MonitoredVehicleJourney>
					<ns3:LineRef>stopArea2</ns3:LineRef>
					<ns3:FramedVehicleJourneyRef>
						<ns3:DataFrameRef>RATPDev:DataFrame::1984-04-04:LOC</ns3:DataFrameRef>
						<ns3:DatedVehicleJourneyRef>stopArea2</ns3:DatedVehicleJourneyRef>
					</ns3:FramedVehicleJourneyRef>
					<ns3:MonitoredCall>
						<ns3:StopPointRef>stopArea2</ns3:StopPointRef>
						<ns3:VehicleAtStop>false</ns3:VehicleAtStop>
						<ns3:ActualArrivalTime>1984-04-04T02:00:00.000Z</ns3:ActualArrivalTime>
					</ns3:MonitoredCall>
				</ns3:MonitoredVehicleJourney>
			</ns3:MonitoredStopVisit>
		</ns3:StopMonitoringDelivery>
	</Answer>
</ns1:GetSiriServiceResponse>`

	// Check the response body is what we expect.
	if responseBody != expectedResponseBody {
		t.Errorf("Unexpected response body:\n expected: %v\n got: %v", expectedResponseBody, responseBody)
	}
}

func Test_SIRIHandler_NotifyStopMonitoring(t *testing.T) {
	// Generate the request Body
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()

	file, err := os.Open("testdata/notify-stop-monitoring.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	soapEnvelope.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer()
	partner := referential.Partners().FindAll()[0]

	partner.Subscriptions().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	subscription, _ := partner.Subscriptions().FindOrCreateByKind("StopMonitoring")
	subscription.Save()

	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "stopArea1")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	objectid2 := model.NewObjectID("objectidKind", "stopArea2")
	stopArea2.SetObjectID(objectid2)
	stopArea2.Save()

	siriHandler_Request(server, soapEnvelope, t)

	// Some Tests

	if count := len(referential.Model().StopVisits().FindAll()); count != 3 {
		t.Errorf("Notify should have created 3 StopVisits, got: %v", count)
	}
}

func Test_SIRIHandler_NotifyGeneralMessage(t *testing.T) {
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()

	file, err := os.Open("../siri/testdata/notify-general-message.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	soapEnvelope.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer()
	partner := referential.Partners().FindAll()[0]

	partner.Subscriptions().FindOrCreateByKind("GeneralMessage")

	siriHandler_Request(server, soapEnvelope, t)

	if count := len(referential.Model().Situations().FindAll()); count != 2 {
		t.Errorf("Notify should have created 2 Situation, got: %v", count)
	}
}
