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
		Id:       string(line.Id()),
		Type:     "Line",
	}

	subscription := partner.Subscriptions().New("sub")
	subscription.SetExternalId("externalId")
	subscription.CreateAddNewResource(reference)
	subscription.Save()

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetObjectID(objectid)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.StopAreaId = stopArea.Id()

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
<ns1:NotifyEstimatedTimetable
	xmlns:ns1="http://wsdl.siri.org.uk"
	xmlns:ns2="http://www.ifopt.org.uk/acsb"
	xmlns:ns3="http://www.ifopt.org.uk/ifopt"
	xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
	xmlns:ns5="http://www.siri.org.uk/siri"
	xmlns:ns6="http://wsdl.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<ns5:ResponseTimestamp>1984-04-04T00:00:00.000Z</ns5:ResponseTimestamp>
		<ns5:ProducerRef>external</ns5:ProducerRef>
		<ns5:ResponseMessageIdentifier>6ba7b814-9dad-11d1-5-00c04fd430c8</ns5:ResponseMessageIdentifier>
	</ServiceDeliveryInfo>
	<Notification>
		<ns3:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>1984-04-04T00:00:00.000Z</ns3:ResponseTimestamp>
			<ns5:SubscriberRef>external</ns5:SubscriberRef>
			<ns5:SubscriptionRef>6ba7b814-9dad-11d1-1-00c04fd430c8</ns5:SubscriptionRef>
			<ns3:Status>true</ns3:Status>
			<ns3:EstimatedJourneyVersionFrame>
				<ns3:RecordedAtTime>2017-01-01T12:00:00.000Z</ns3:RecordedAtTime>
				<ns3:EstimatedVehicleJourney>
					<ns3:LineRef>6ba7b814-9dad-11d1-0-00c04fd430c8</ns3:LineRef>
					<ns3:DatedVehicleJourneyRef>6ba7b814-9dad-11d1-0-00c04fd430c8</ns3:DatedVehicleJourneyRef>
					<ns3:EstimatedCalls>
						<ns3:EstimatedCall>
							<ns3:StopPointRef>6ba7b814-9dad-11d1-0-00c04fd430c8</ns3:StopPointRef>
							<ns3:Order>0</ns3:Order>
							<ns3:VehicleAtStop>false</ns3:VehicleAtStop>
						</ns3:EstimatedCall>
					</ns3:EstimatedCalls>
				</ns3:EstimatedVehicleJourney>
			</ns3:EstimatedJourneyVersionFrame>
		</ns3:EstimatedTimetableDelivery>
	</Notification>
	<NotifyExtension
		xmlns:ns2="http://www.ifopt.org.uk/acsb"
		xmlns:ns3="http://www.ifopt.org.uk/ifopt"
		xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
		xmlns:ns5="http://www.siri.org.uk/siri"
		xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
	</ns1:NotifyEstimatedTimetable>
</S:Body>
</S:Envelope>`

	if string(response) != expected {
		t.Errorf("Got diffrent xml than expected, got: %v want :%v", string(response), expected)
	}
}
