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

func Test_SubscriptionRequest_Dispatch_ETT(t *testing.T) {
	model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.ConnectorTypes = []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)

	line := referential.Model().Lines().New()
	objectid := model.NewObjectID("_internal", "6ba7b814-9dad-11d1-1-00c04fd430c8")
	line.SetObjectID(objectid)
	line.Save()

	file, _ := os.Open("testdata/estimatedtimetable-request-soap.xml")
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
}

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

	subs := partner.Subscriptions().FindByRessourceId(objectid.String(), "StopMonitoringBroadcast")

	if len(subs) == 0 {
		t.Errorf("Should have been able to find the stopArea ressource : %v", objectid.String())
	}

	sub := subs[0]
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

	sub, ok := partner.Subscriptions().FindByKind("GeneralMessageBroadcast")
	if !ok {
		t.Fatalf("Could not find a subscription with kind of Situation")
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
	request, _ := siri.NewXMLDeleteSubscriptionRequestFromContent(body)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	response := connector.(*SIRISubscriptionRequestDispatcher).CancelSubscription(request)

	if len(response.ResponseStatus) != 1 {
		t.Fatalf("Response should have 1 responseStatus, got: %v", len(response.ResponseStatus))
	}
	if !response.ResponseStatus[0].Status {
		t.Errorf("Status should be true but got false")
	}

	if _, ok := partner.Subscriptions().FindByExternalId("Edwig:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC"); ok {
		t.Errorf("Subscription shouldn't exist")
	}
}

func Test_CancelSubscriptionAll(t *testing.T) {
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

	sub = partner.Subscriptions().New("Test")
	sub.SetExternalId("Edwig:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c9:LOC")

	partner.Subscriptions().Save(sub)

	file, _ := os.Open("testdata/terminated_subscription_request_all-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLDeleteSubscriptionRequestFromContent(body)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	response := connector.(*SIRISubscriptionRequestDispatcher).CancelSubscription(request)

	if len(response.ResponseStatus) != 2 {
		t.Fatalf("Response should have 1 responseStatus, got: %v", len(response.ResponseStatus))
	}
	if !response.ResponseStatus[0].Status || !response.ResponseStatus[1].Status {
		t.Errorf("Status should be true but got false: %v %v", response.ResponseStatus[0].Status, response.ResponseStatus[1].Status)
	}

	if len(partner.Subscriptions().FindAll()) != 0 {
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

func Test_ReceiveStateGM(t *testing.T) {
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

	referential.model.SetBroadcastGMChan(referential.broacasterManager.GetGeneralMessageBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.Settings["remote_credential"] = "external"
	partner.Settings["remote_url"] = ts.URL
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST_DISPATCHER, SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	connector2, _ := partner.Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)

	connector2.(*SIRIGeneralMessageSubscriptionBroadcaster).SetClock(fakeClock)
	connector2.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster = NewFakeGeneralMessageBroadcaster(connector2.(*SIRIGeneralMessageSubscriptionBroadcaster))
	situation := referential.Model().Situations().New()

	line := referential.Model().Lines().New()
	objectid0 := model.NewObjectID("_internal", "line1")
	line.SetObjectID(objectid0)
	line.Save()

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()
	objectid1 := model.NewObjectID("_internal", "coicogn1")

	situation.References = append(situation.References, &model.Reference{ObjectId: &objectid1, Type: "StopPointRef"})

	stopArea.SetObjectID(objectid1)
	stopArea.Save()

	stopArea = referential.Model().StopAreas().New()
	stopArea.Save()
	objectid2 := model.NewObjectID("_internal", "coicogn2")

	situation.References = append(situation.References, &model.Reference{ObjectId: &objectid2, Type: "StopPointRef"})

	stopArea.SetObjectID(objectid2)
	stopArea.Save()

	objectid3 := model.NewObjectID("_internal", string(situation.Id()))
	situation.Channel = "Mondial"
	situation.ValidUntil = fakeClock.Now().Add(10 * time.Minute)
	message := &model.Message{
		Content:             "Le content",
		Type:                "Le Type",
		NumberOfLines:       1,
		NumberOfCharPerLine: 10,
	}
	situation.Messages = append(situation.Messages, message)
	situation.SetObjectID(objectid3)

	lineSectionReferences := model.NewReferences()
	lineSectionReferences["FirstStop"] = model.Reference{ObjectId: &objectid2, Type: "StopPointRef"}
	lineSectionReferences["LastStop"] = model.Reference{ObjectId: &objectid1, Type: "StopPointRef"}
	lineSectionReferences["LinesRef"] = model.Reference{ObjectId: &objectid0, Type: "LineRef"}

	situation.LineSections = append(situation.LineSections, &lineSectionReferences)
	situation.Save()

	file, _ := os.Open("testdata/generalmessagesubscription-request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLSubscriptionRequestFromContent(body)
	time.Sleep(10 * time.Millisecond)

	connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request)
	connector2.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster.Start()
	time.Sleep(10 * time.Millisecond)

	notify, _ := siri.NewXMLNotifyGeneralMessageFromContent(response)
	delivery := notify.GeneralMessagesDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	gms := delivery[0].XMLGeneralMessages()
	if len(gms) != 1 {
		t.Errorf("Should have received 1 GeneralMessage but got == %v", len(gms))
	}

	content := gms[0].Content().(siri.IDFGeneralMessageStructure)

	messages := content.Messages()

	if len(messages) != 1 {
		t.Errorf("Should have received 1 messages but got == %v", len(messages))
	}

}
