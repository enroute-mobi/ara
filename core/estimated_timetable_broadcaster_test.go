package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EstimatedTimetableBroadcaster_Receive_Notify(t *testing.T) {
	assert := assert.New(t)

	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)
	uuidGenerator := uuid.NewFakeUUIDGenerator()
	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = io.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test http server
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
		"remote_credential": "external",
		"local_credential":  "local",
		"remote_url":        ts.URL,
	}
	partner.SetUUIDGenerator(uuidGenerator)
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)

	connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster).Partner().SetUUIDGenerator(uuidGenerator)
	connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster).estimatedTimetableBroadcaster = NewFakeSIRIEstimatedTimetableBroadcaster(connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster))

	referential.Model().Lines().SetUUIDGenerator(uuidGenerator)
	referential.Model().StopAreas().SetUUIDGenerator(uuidGenerator)
	referential.Model().StopVisits().SetUUIDGenerator(uuidGenerator)
	referential.Model().VehicleJourneys().SetUUIDGenerator(uuidGenerator)
	partner.Subscriptions().SetUUIDGenerator(uuidGenerator)

	line := referential.Model().Lines().New()
	line.Save()

	code := model.NewCode("internal", string(line.Id()))
	line.SetCode(code)

	reference := model.Reference{
		Code: &code,
		Type: "Line",
	}

	subscription := partner.Subscriptions().New(EstimatedTimetableBroadcast)
	subscription.SubscriberRef = "subscriber"
	subscription.SetExternalId("externalId")
	subscription.CreateAndAddNewResource(reference)
	subscription.SetSubscriptionOption("MessageIdentifier", "MessageIdentifier")
	subscription.Save()

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	operatorCode := model.NewCode("test", "1234")
	operatorRef := model.Reference{
		Code: &operatorCode,
		Type: "OperatorRef",
	}

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetCode(code)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetCode(code)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.References.SetReference("OperatorRef", operatorRef)
	stopVisit.StopAreaId = stopArea.Id()

	operator := referential.Model().Operators().New()
	operator.SetCode(operatorCode)
	operator.SetCode(model.NewCode("internal", "123456789"))
	operator.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work

	if l := len(connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster).toBroadcast); l != 1 {
		t.Errorf("should have 1 line to broadcast got : %v", l)
	}

	connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster).estimatedTimetableBroadcaster.Start()

	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>external</siri:ProducerRef>
		<siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:ResponseMessageIdentifier>
	</ServiceDeliveryInfo>
	<Notification>
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
			<siri:SubscriberRef>subscriber</siri:SubscriberRef>
			<siri:SubscriptionRef>externalId</siri:SubscriptionRef>
			<siri:Status>true</siri:Status>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>1984-04-04T00:00:00.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>6ba7b814-9dad-11d1-0-00c04fd430c8</siri:LineRef>
					<siri:DirectionRef>unknown</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>6ba7b814-9dad-11d1-0-00c04fd430c8</siri:DatedVehicleJourneyRef>
					<siri:OperatorRef>123456789</siri:OperatorRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>6ba7b814-9dad-11d1-0-00c04fd430c8</siri:StopPointRef>
							<siri:Order>0</siri:Order>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
		</siri:EstimatedTimetableDelivery>
	</Notification>
	<SiriExtension />
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>`

	assert.Equal(expected, string(response))
}
