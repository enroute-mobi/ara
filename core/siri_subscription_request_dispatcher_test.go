package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SubscriptionRequest_Dispatch_SM(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	objectid := model.NewObjectID("_internal", "coicogn2")

	stopArea.SetObjectID(objectid)
	stopArea.Save()

	file, _ := os.Open("testdata/stopmonitoringsubscription-request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLSubscriptionRequestFromContent(body)

	response, err := connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request)
	if err != nil {
		t.Fatalf("Error while handling subscription request: %v", err)
	}

	if len(response.ResponseStatus) != 1 {
		t.Errorf("Wrong ResponseStatus size want 1 got : %v", len(response.ResponseStatus))
	}

	if !response.ResponseStatus[0].Status {
		t.Errorf("Wrong first ResponseStatus status want true got : %v", response.ResponseStatus[0].Status)
	}

	sub, ok := partner.Subscriptions().FindByRessourceId(objectid.String())
	if !ok {
		t.Errorf("Should have been able to find the stopArea ressource : %v", objectid.String())
	}
	externalId := "Edwig:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"
	if sub.ExternalId() != externalId {
		t.Errorf("Wrong ExternalId value want: %v got: %v", externalId, sub.ExternalId())
	}
}

func Test_SubscriptionRequest_Dispatch_GM(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.ConnectorTypes = []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)

	file, _ := os.Open("testdata/generalmessagesubscription-request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLSubscriptionRequestFromContent(body)

	response, err := connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request)
	if err != nil {
		t.Fatalf("Error while handling subscription request: %v", err)
	}

	if len(response.ResponseStatus) != 1 {
		t.Errorf("Wrong ResponseStatus size want 1 got : %v", len(response.ResponseStatus))
	}

	if !response.ResponseStatus[0].Status {
		t.Errorf("Wrong first ResponseStatus status want true got : %v", response.ResponseStatus[0].Status)
	}

	sub, ok := partner.Subscriptions().FindByKind("Situation")
	if !ok {
		t.Errorf("Could not find a subscription with kind of Situation")
	}

	externalId := "NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"
	if sub.ExternalId() != externalId {
		t.Errorf("Wrong ExternalId value want: %v got: %v", externalId, sub.ExternalId())
	}
}

func Test_CancelSubscription(t *testing.T) {
	model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST_DISPATCHER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	sub := partner.Subscriptions().New("Test")
	sub.SetExternalId("Edwig:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC")

	partner.Subscriptions().Save(sub)

	file, _ := os.Open("testdata/terminated_subscription_request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLTerminatedSubscriptionRequestFromContent(body)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	response := connector.(*SIRISubscriptionRequestDispatcher).CancelSubscription(request)

	if response.Status != true {
		t.Errorf("Status should be true but got false")
	}

	if _, ok := partner.Subscriptions().FindByExternalId("Edwig:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC"); ok {
		t.Errorf("Subscription shouldn't exist")
	}
}

func Test_ReceiveStateSM(t *testing.T) {
	fakeClock := model.NewFakeClock()
	model.SetDefaultClock(fakeClock)

	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = ioutil.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	referential.model.SetBroadcastSMChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.Settings["remote_credential"] = "external"
	partner.Settings["remote_url"] = ts.URL
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST_DISPATCHER, SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	connector2, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

	connector2.(*SIRIStopMonitoringSubscriptionBroadcaster).SetClock(fakeClock)
	connector2.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster = NewFakeStopMonitoringBroadcaster(connector2.(*SIRIStopMonitoringSubscriptionBroadcaster))

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	objectid := model.NewObjectID("_internal", "coicogn2")

	stopArea.SetObjectID(objectid)
	stopArea.Save()

	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.Save()

	objectid = model.NewObjectID("_internal", "value")

	sv1 := referential.Model().StopVisits().New()
	sv1.SetObjectID(objectid)
	sv1.StopAreaId = stopArea.Id()
	sv1.VehicleJourneyId = vehicleJourney.Id()
	sv1.Schedules.SetArrivalTime("actual", fakeClock.Now().Add(5*time.Minute))
	sv1.Save()

	objectid = model.NewObjectID("_internal", "value2")
	sv2 := referential.Model().StopVisits().New()
	sv2.SetObjectID(objectid)
	sv2.StopAreaId = stopArea.Id()
	sv2.VehicleJourneyId = vehicleJourney.Id()
	sv2.Schedules.SetArrivalTime("actual", fakeClock.Now().Add(5*time.Minute))
	sv2.Save()

	file, _ := os.Open("testdata/stopmonitoringsubscription-request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLSubscriptionRequestFromContent(body)

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	sv1.Save()
	sv2.Save()

	connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request)
	connector2.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Start()
	time.Sleep(10 * time.Millisecond)

	notify, _ := siri.NewXMLNotifyStopMonitoringFromContent(response)
	delivery := notify.StopMonitoringDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	if len(delivery[0].XMLMonitoredStopVisits()) != 2 {
		t.Errorf("Should have received 2 Monitored stop visit but got %v", len(delivery[0].XMLMonitoredStopVisits()))
	}
}
