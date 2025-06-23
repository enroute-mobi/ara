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
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		"remote_code_space":                      "codeSpace",
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
	require := require.New(t)

	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// Create a request
	request, err := http.NewRequest("POST", "/default/siri", buffer)
	require.NoError(err)

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{referential_slug}/siri", server.HandleSIRI)

	// Call ServeHTTP method and pass in our Request and ResponseRecorder.
	mux.ServeHTTP(responseRecorder, request)

	// Check the status code is what we expect.
	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("text/xml; charset=utf-8", responseRecorder.Header().Get("Content-Type"))

	return responseRecorder
}

func Test_SIRIHandler_SOAP(t *testing.T) {
	assert := assert.New(t)

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
	assert.NoError(err)
	assert.True(response.Status())
}

func Test_SIRIHandler_SOAPResponse(t *testing.T) {
	assert := assert.New(t)

	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	assert.NoError(err)

	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer(remote.SOAP_SIRI_ENVELOPE)
	responseRecorder := siriHandler_Request(server, buffer, t)

	_, err = remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	assert.NoError(err, "We should receive a SOAP response")
}

func Test_SIRIHandler_RawResponse(t *testing.T) {
	assert := assert.New(t)

	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	assert.NoError(err)

	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer(remote.RAW_SIRI_ENVELOPE)
	responseRecorder := siriHandler_Request(server, buffer, t)

	_, err = remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	assert.Error(err, "NewSIRIEnvelope with SOAP option should return an error")

	responseRecorder = siriHandler_Request(server, buffer, t) // Making the request again as the reader should be empty
	_, err = remote.NewSIRIEnvelope(responseRecorder.Body, remote.RAW_SIRI_ENVELOPE)
	assert.NoError(err, "We shouldn't get an error while trying to create a raw envelope")
}

func Test_SIRIHandler_Raw(t *testing.T) {
	assert := assert.New(t)

	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.RAW_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	assert.NoError(err)

	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer("")
	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	assert.NoError(err)
	assert.True(response.Status())
}

func Test_SIRIHandler_CheckStatus(t *testing.T) {
	assert := assert.New(t)

	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRICheckStatusRequest("Ara",
		clock.DefaultClock().Now(),
		"Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC").BuildXML()
	assert.NoError(err)

	buffer.WriteXML(request)

	server, _ := siriHandler_PrepareServer("")
	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	assert.NoError(err)
	assert.Equal("http://ara", response.Address())
	assert.Equal("Ara", response.ProducerRef())
	assert.Equal("Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.RequestMessageRef())
	assert.Equal("Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier())
	assert.True(response.Status())

	expectedDate := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	assert.True(response.ResponseTimestamp().Equal(expectedDate))
	assert.True(response.ServiceStartedTime().Equal(expectedDate))
}

func Test_SIRIHandler_CheckStatus_Gzip(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	server, _ := siriHandler_PrepareServer("")

	// Create a request
	file, err := os.Open("testdata/checkstatus-soap-request.xml.gz")
	require.NoError(err)
	defer file.Close()

	request, err := http.NewRequest("POST", "/default/siri", file)
	require.NoError(err)

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{referential_slug}/siri", server.HandleSIRI)

	// Call ServeHTTP method and pass in our Request and ResponseRecorder.
	mux.ServeHTTP(responseRecorder, request)
	assert.Equal(http.StatusOK, responseRecorder.Code)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLCheckStatusResponseFromContent(responseRecorder.Body.Bytes())
	assert.NoError(err)
	assert.True(response.Status())
}

func Test_SIRIHandler_StopMonitoring(t *testing.T) {
	assert := assert.New(t)

	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)
	request, err := siri.NewSIRIGetStopMonitoringRequest("Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC",
		"codeValue",
		"Ara",
		clock.DefaultClock().Now()).BuildXML()
	assert.NoError(err)

	buffer.WriteXML(request)

	server, referential := siriHandler_PrepareServer("")
	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("codeSpace", "codeValue")
	stopArea.SetCode(code)
	stopArea.Monitored = true
	stopArea.Save()

	line := referential.Model().Lines().New()
	line.SetCode(code)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetCode(code)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Monitored = true
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(2*time.Hour))
	stopVisit.SetCode(model.NewCode("codeSpace", "second"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Save()

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(1*time.Hour))
	stopVisit2.SetCode(model.NewCode("codeSpace", "first"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.Save()

	pastStopVisit := referential.Model().StopVisits().New()
	pastStopVisit.StopAreaId = stopArea.Id()
	pastStopVisit.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(-1*time.Hour))
	pastStopVisit.SetCode(model.NewCode("codeSpace", "past"))
	pastStopVisit.VehicleJourneyId = vehicleJourney.Id()
	pastStopVisit.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	// Check the response body is what we expect.
	response, err := sxml.NewXMLStopMonitoringResponseFromContent(responseRecorder.Body.Bytes())
	assert.NoError(err)

	delivery := response.StopMonitoringDeliveries()[0]
	assert.Len(delivery.XMLMonitoredStopVisits(), 2, "Past StopVisit should be ignored")
	assert.True(delivery.XMLMonitoredStopVisits()[1].ActualArrivalTime().After(delivery.XMLMonitoredStopVisits()[0].ActualArrivalTime()),
		"Stop visits are not chronollogicaly ordered ")

	assert.Equal("http://ara", response.Address())
	assert.Equal("Ara", response.ProducerRef())
	assert.Equal("Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.RequestMessageRef())
	assert.Equal("Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier())

	expectedDate := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	assert.True(response.ResponseTimestamp().Equal(expectedDate))
}

func Test_SIRIHandler_SiriService(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Generate the request Body
	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/siri-service-request-soap.xml")
	require.NoError(err)
	defer file.Close()

	content, err := io.ReadAll(file)
	require.NoError(err)

	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")
	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("codeSpace", "stopArea1")
	stopArea.SetCode(code)
	stopArea.Monitored = true
	stopArea.Save()

	line := referential.Model().Lines().New()
	line.SetCode(code)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetCode(code)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Monitored = true
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(2*time.Hour))
	stopVisit.SetCode(model.NewCode("codeSpace", "second"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Save()

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(1*time.Hour))
	stopVisit2.SetCode(model.NewCode("codeSpace", "first"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.Save()

	pastStopVisit := referential.Model().StopVisits().New()
	pastStopVisit.StopAreaId = stopArea.Id()
	pastStopVisit.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(-1*time.Hour))
	pastStopVisit.SetCode(model.NewCode("codeSpace", "past"))
	pastStopVisit.VehicleJourneyId = vehicleJourney.Id()
	pastStopVisit.Save()

	stopArea2 := referential.Model().StopAreas().New()
	code2 := model.NewCode("codeSpace", "stopArea2")
	stopArea2.SetCode(code2)
	stopArea2.Monitored = true
	stopArea2.Save()

	line2 := referential.Model().Lines().New()
	line2.SetCode(code2)
	line2.Save()

	vehicleJourney2 := referential.Model().VehicleJourneys().New()
	vehicleJourney2.SetCode(code2)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Monitored = true
	vehicleJourney2.Save()

	stopVisit3 := referential.Model().StopVisits().New()
	stopVisit3.StopAreaId = stopArea2.Id()
	stopVisit3.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(2*time.Hour))
	stopVisit3.SetCode(model.NewCode("codeSpace", "third"))
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.Save()

	stopVisit4 := referential.Model().StopVisits().New()
	stopVisit4.StopAreaId = stopArea2.Id()
	stopVisit4.Schedules.SetArrivalTime(schedules.Actual, referential.Clock().Now().Add(1*time.Hour))
	stopVisit4.SetCode(model.NewCode("codeSpace", "fourth"))
	stopVisit4.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit4.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	// responseRecorder.Body.String()
	envelope, err := remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	assert.NoError(err)

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
			<siri:Status>true</siri:Status>
			<siri:MonitoringRef>stopArea1</siri:MonitoringRef>
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
			<siri:Status>true</siri:Status>
			<siri:MonitoringRef>stopArea2</siri:MonitoringRef>
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
</sw:GetSiriServiceResponse>`

	assert.Equal(expectedResponseBody, responseBody)
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
	code := model.NewCode("codeSpace", "stopArea1")
	stopArea.SetCode(code)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	code2 := model.NewCode("codeSpace", "stopArea2")
	stopArea2.SetCode(code2)
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
	assert := assert.New(t)
	require := require.New(t)

	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/estimated_timetable_request.xml")
	require.NoError(err)
	defer file.Close()

	content, err := io.ReadAll(file)
	require.NoError(err)

	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(model.NewCode("codeSpace", "stopArea1"))
	stopArea.Monitored = true
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetCode(model.NewCode("codeSpace", "stopArea2"))
	stopArea2.Monitored = true
	stopArea2.Save()

	line := referential.Model().Lines().New()
	line.SetCode(model.NewCode("codeSpace", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.Model().Lines().New()
	line2.SetCode(model.NewCode("codeSpace", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.SetCode(model.NewCode("codeSpace", "vehicleJourney"))
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	vehicleJourney2 := referential.Model().VehicleJourneys().New()
	vehicleJourney2.SetCode(model.NewCode("codeSpace", "vehicleJourney2"))
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	pastStopVisit := referential.Model().StopVisits().New()
	pastStopVisit.SetCode(model.NewCode("codeSpace", "pastStopVisit"))
	pastStopVisit.VehicleJourneyId = vehicleJourney.Id()
	pastStopVisit.StopAreaId = stopArea.Id()
	pastStopVisit.PassageOrder = 0
	pastStopVisit.ArrivalStatus = "onTime"
	pastStopVisit.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(-1*time.Minute))
	pastStopVisit.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(-1*time.Minute))
	pastStopVisit.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.SetCode(model.NewCode("codeSpace", "stopVisit"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.PassageOrder = 1
	stopVisit.ArrivalStatus = "onTime"
	stopVisit.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.SetCode(model.NewCode("codeSpace", "stopVisit2"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.StopAreaId = stopArea2.Id()
	stopVisit2.PassageOrder = 2
	stopVisit2.ArrivalStatus = "onTime"
	stopVisit2.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(2*time.Minute))
	stopVisit2.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(2*time.Minute))
	stopVisit2.Save()

	stopVisit3 := referential.Model().StopVisits().New()
	stopVisit3.SetCode(model.NewCode("codeSpace", "stopVisit3"))
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.StopAreaId = stopArea.Id()
	stopVisit3.PassageOrder = 1
	stopVisit3.ArrivalStatus = "onTime"
	stopVisit3.Schedules.SetArrivalTime("aimed", referential.Clock().Now().Add(1*time.Minute))
	stopVisit3.Schedules.SetArrivalTime("expected", referential.Clock().Now().Add(1*time.Minute))
	stopVisit3.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	envelope, err := remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	assert.NoError(err)

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
					<siri:DirectionRef>unknown</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>vehicleJourney</siri:DatedVehicleJourneyRef>
					<siri:PublishedLineName>lineName</siri:PublishedLineName>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopArea1</siri:StopPointRef>
							<siri:Order>1</siri:Order>
							<siri:AimedArrivalTime>1984-04-04T00:01:00.000Z</siri:AimedArrivalTime>
							<siri:ExpectedArrivalTime>1984-04-04T00:01:00.000Z</siri:ExpectedArrivalTime>
							<siri:ArrivalStatus>onTime</siri:ArrivalStatus>
						</siri:EstimatedCall>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopArea2</siri:StopPointRef>
							<siri:Order>2</siri:Order>
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
					<siri:DirectionRef>unknown</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>vehicleJourney2</siri:DatedVehicleJourneyRef>
					<siri:PublishedLineName>lineName2</siri:PublishedLineName>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopArea1</siri:StopPointRef>
							<siri:Order>1</siri:Order>
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

	assert.Equal(expectedResponseBody, responseBody)
}

func Test_SIRIHandler_LinesDiscovery(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	buffer := remote.NewSIRIBuffer(remote.SOAP_SIRI_ENVELOPE)

	file, err := os.Open("testdata/lines-discovery-request.xml")
	require.NoError(err)
	defer file.Close()

	content, err := io.ReadAll(file)
	require.NoError(err)

	buffer.WriteXML(string(content))

	server, referential := siriHandler_PrepareServer("")

	line := referential.Model().Lines().New()
	line.SetCode(model.NewCode("codeSpace", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.Model().Lines().New()
	line2.SetCode(model.NewCode("codeSpace", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	line3 := referential.Model().Lines().New()
	line3.SetCode(model.NewCode("codeSpace2", "NINOXE:Line:4:LOC"))
	line3.Name = "lineName3"
	line3.Save()

	responseRecorder := siriHandler_Request(server, buffer, t)

	envelope, err := remote.NewSIRIEnvelope(responseRecorder.Body, remote.SOAP_SIRI_ENVELOPE)
	assert.NoError(err)

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

	assert.Equal(expectedResponseBody, responseBody)
}
