package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_EstimatedTimeTableBroadcaster_Receive_Notify(t *testing.T) {
	// logger.Log.Debug = true

	fakeClock := model.NewFakeClock()
	model.SetDefaultClock(fakeClock)
	uuidGenerator := model.NewFakeUUIDGenerator()
	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = ioutil.ReadAll(r.Body)
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
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.Settings["remote_credential"] = "external"
	partner.Settings["remote_url"] = ts.URL

	partner.ConnectorTypes = []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)

	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).SIRIPartner().SetUUIDGenerator(uuidGenerator)
	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).estimatedTimeTableBroadcaster = NewFakeEstimatedTimeTableBroadcaster(connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster))

	referential.Model().Lines().SetUUIDGenerator(uuidGenerator)
	referential.Model().StopAreas().SetUUIDGenerator(uuidGenerator)
	referential.Model().StopVisits().SetUUIDGenerator(uuidGenerator)
	referential.Model().VehicleJourneys().SetUUIDGenerator(uuidGenerator)
	partner.Subscriptions().SetUUIDGenerator(uuidGenerator)

	line := referential.Model().Lines().New()
	line.Save()

	objectid := model.NewObjectID("internal", string(line.Id()))
	line.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "Line",
	}

	subscription := partner.Subscriptions().New("EstimatedTimeTableBroadcast")
	subscription.SetExternalId("externalId")
	subscription.CreateAddNewResource(reference)
	subscription.SubscriptionOptions()["MessageIdentifier"] = "MessageIdentifier"
	subscription.Save()

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	operatorObjectid := model.NewObjectID("test", "1234")
	operatorRef := model.Reference{
		ObjectId: &operatorObjectid,
		Type:     "OperatorRef",
	}

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetObjectID(objectid)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.References.SetReference("OperatorRef", operatorRef)
	stopVisit.StopAreaId = stopArea.Id()

	operator := referential.Model().Operators().New()
	operator.SetObjectID(operatorObjectid)
	operator.SetObjectID(model.NewObjectID("internal", "123456789"))
	operator.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work

	if l := len(connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).toBroadcast); l != 1 {
		t.Errorf("should have 1 line to broadcast got : %v", l)
	}

	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).estimatedTimeTableBroadcaster.Start()

	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>external</siri:ProducerRef>
		<siri:ResponseMessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>MessageIdentifier</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Notification>
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>1984-04-04T00:00:00.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>MessageIdentifier</siri:RequestMessageRef>
			<siri:SubscriberRef>external</siri:SubscriberRef>
			<siri:SubscriptionRef>externalId</siri:SubscriptionRef>
			<siri:Status>true</siri:Status>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>1984-04-04T00:00:00.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>6ba7b814-9dad-11d1-0-00c04fd430c8</siri:LineRef>
					<siri:DirectionRef/>
					<siri:OperatorRef>123456789</siri:OperatorRef>
					<siri:DatedVehicleJourneyRef>6ba7b814-9dad-11d1-0-00c04fd430c8</siri:DatedVehicleJourneyRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>6ba7b814-9dad-11d1-0-00c04fd430c8</siri:StopPointRef>
							<siri:Order>0</siri:Order>
							<siri:VehicleAtStop>false</siri:VehicleAtStop>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
		</siri:EstimatedTimetableDelivery>
	</Notification>
	<NotificationExtension/>
</sw:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>`

	if string(response) != expected {
		t.Errorf("Got diffrent xml than expected, got: %v\nwant :%v", string(response), expected)
	}
}
