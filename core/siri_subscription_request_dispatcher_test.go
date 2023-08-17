package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SubscriptionRequest_Dispatch_ETT(t *testing.T) {
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")

	settings := map[string]string{
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)

	line := referential.Model().Lines().New()
	objectid := model.NewObjectID("_internal", "6ba7b814-9dad-11d1-1-00c04fd430c8")
	line.SetObjectID(objectid)
	line.Save()

	file, _ := os.Open("testdata/estimatedtimetable-request-soap.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLSubscriptionRequestFromContent(body)

	response, err := connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request, &audit.BigQueryMessage{})

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

func Test_SubscriptionRequest_Dispatch_PTT(t *testing.T) {
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)

	line := referential.Model().Lines().New()
	objectid := model.NewObjectID("_internal", "6ba7b814-9dad-11d1-1-00c04fd430c8")
	line.SetObjectID(objectid)
	line.Save()

	file, _ := os.Open("testdata/productiontimetable-request.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLSubscriptionRequestFromContent(body)

	response, err := connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request, &audit.BigQueryMessage{})
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
	settings := map[string]string{
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
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
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLSubscriptionRequestFromContent(body)

	response, err := connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request, &audit.BigQueryMessage{})
	if err != nil {
		t.Fatalf("Error while handling subscription request: %v", err)
	}

	if len(response.ResponseStatus) != 1 {
		t.Errorf("Wrong ResponseStatus size want 1 got : %v", len(response.ResponseStatus))
	}

	if !response.ResponseStatus[0].Status {
		t.Errorf("Wrong first ResponseStatus status want true got : %v", response.ResponseStatus[0].Status)
	}

	subs := partner.Subscriptions().FindByResourceId(objectid.String(), "StopMonitoringBroadcast")

	if len(subs) == 0 {
		t.Fatalf("Should have been able to find the stopArea ressource : %v", objectid.String())
	}

	sub := subs[0]
	externalId := "Ara:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"
	if sub.ExternalId() != externalId {
		t.Errorf("Wrong ExternalId value want: %v got: %v", externalId, sub.ExternalId())
	}
}

func Test_SubscriptionRequest_Dispatch_GM(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)

	file, _ := os.Open("testdata/generalmessagesubscription-request-soap.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLSubscriptionRequestFromContent(body)

	response, err := connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request, &audit.BigQueryMessage{})
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
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST_DISPATCHER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	sub := partner.Subscriptions().New("Test")
	sub.SetExternalId("Ara:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC")

	partner.Subscriptions().Save(sub)

	file, _ := os.Open("testdata/terminated_subscription_request-soap.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLDeleteSubscriptionRequestFromContent(body)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	response := connector.(*SIRISubscriptionRequestDispatcher).CancelSubscription(request, &audit.BigQueryMessage{})

	if len(response.ResponseStatus) != 1 {
		t.Fatalf("Response should have 1 responseStatus, got: %v", len(response.ResponseStatus))
	}
	if !response.ResponseStatus[0].Status {
		t.Errorf("Status should be true but got false")
	}

	if _, ok := partner.Subscriptions().FindByExternalId("Ara:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC"); ok {
		t.Errorf("Subscription shouldn't exist")
	}
}

func Test_CancelSubscriptionAll(t *testing.T) {
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST_DISPATCHER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	sub := partner.Subscriptions().New("Test")
	sub.kind = "Broadcast"
	sub.SetExternalId("Ara:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC")

	partner.Subscriptions().Save(sub)

	sub = partner.Subscriptions().New("Test")
	sub.kind = "Broadcast"
	sub.SetExternalId("Ara:Subscription::6ba7b814-9dad-11d1-1-00c04fd430c9:LOC")

	partner.Subscriptions().Save(sub)

	file, _ := os.Open("testdata/terminated_subscription_request_all-soap.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLDeleteSubscriptionRequestFromContent(body)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER)
	response := connector.(*SIRISubscriptionRequestDispatcher).CancelSubscription(request, &audit.BigQueryMessage{})

	if len(response.ResponseStatus) != 2 {
		t.Fatalf("Response should have 2 responseStatus, got: %v", len(response.ResponseStatus))
	}
	if !response.ResponseStatus[0].Status || !response.ResponseStatus[1].Status {
		t.Errorf("Status should be true but got false: %v %v", response.ResponseStatus[0].Status, response.ResponseStatus[1].Status)
	}

	if len(partner.Subscriptions().FindAll()) != 0 {
		t.Errorf("Subscription shouldn't exist")
	}
}

func Test_ReceiveStateSM(t *testing.T) {
	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)

	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = io.ReadAll(r.Body)
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

	settings := map[string]string{
		"remote_objectid_kind": "_internal",
		"remote_credential":    "external",
		"remote_url":           ts.URL,
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

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
	sv2.ArrivalStatus = model.STOP_VISIT_ARRIVAL_CANCELLED
	sv2.Save()

	file, _ := os.Open("testdata/stopmonitoringsubscription-request-soap.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLSubscriptionRequestFromContent(body)

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...

	connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request, &audit.BigQueryMessage{})
	connector2.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Start()
	time.Sleep(10 * time.Millisecond)

	notify, _ := sxml.NewXMLNotifyStopMonitoringFromContent(response)
	delivery := notify.StopMonitoringDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	if len(delivery[0].XMLMonitoredStopVisits()) != 1 {
		t.Errorf("Should have received 1 Monitored stop visit but got %v", len(delivery[0].XMLMonitoredStopVisits()))
	}

	if len(delivery[0].XMLMonitoredStopVisitCancellations()) != 1 {
		t.Errorf("Should have received 1 Cancelled stop visit but got %v", len(delivery[0].XMLMonitoredStopVisitCancellations()))
	}
}

func Test_ReceiveStateGM(t *testing.T) {
	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)

	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = io.ReadAll(r.Body)
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

	settings := map[string]string{
		"remote_objectid_kind": "_internal",
		"remote_credential":    "external",
		"remote_url":           ts.URL,
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

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
	situation.Keywords = []string{"Perturbation"}
	period := &model.TimeRange{EndTime: fakeClock.Now().Add(10 * time.Minute)}
	situation.ValidityPeriods = []*model.TimeRange{period}

	message := &model.Message{
		Content:             "Le content",
		Type:                "Le Type",
		NumberOfLines:       1,
		NumberOfCharPerLine: 10,
	}
	situation.Messages = append(situation.Messages, message)
	situation.SetObjectID(objectid3)

	lineSectionReferences := model.NewReferences()
	lineSectionReferences.SetReference("FirstStop", model.Reference{ObjectId: &objectid2, Type: "StopPointRef"})
	lineSectionReferences.SetReference("LastStop", model.Reference{ObjectId: &objectid1, Type: "StopPointRef"})
	lineSectionReferences.SetReference("LinesRef", model.Reference{ObjectId: &objectid0, Type: "LineRef"})

	situation.LineSections = append(situation.LineSections, &lineSectionReferences)
	situation.Save()

	file, _ := os.Open("testdata/generalmessagesubscription-request-soap.xml")
	body, _ := io.ReadAll(file)
	request, _ := sxml.NewXMLSubscriptionRequestFromContent(body)
	time.Sleep(10 * time.Millisecond)

	connector.(*SIRISubscriptionRequestDispatcher).Dispatch(request, &audit.BigQueryMessage{})
	connector2.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster.Start()
	time.Sleep(10 * time.Millisecond)

	notify, _ := sxml.NewXMLNotifyGeneralMessageFromContent(response)
	delivery := notify.GeneralMessagesDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	gms := delivery[0].XMLGeneralMessages()
	if len(gms) != 1 {
		t.Errorf("Should have received 1 GeneralMessage but got == %v", len(gms))
	}

	content := gms[0].Content().(sxml.IDFGeneralMessageStructure)

	messages := content.Messages()

	if len(messages) != 1 {
		t.Errorf("Should have received 1 messages but got == %v", len(messages))
	}

}

func Test_HandleSubscriptionTerminatedNotification(t *testing.T) {
	file, err := os.Open("testdata/subscription_terminated_notification-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLSubscriptionTerminatedNotificationFromContent(content)

	partners := createTestPartnerManager()
	partner := partners.New("slug")

	settings := map[string]string{
		"remote_url":           "une url",
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRISubscriptionRequestDispatcher(partner)
	partner.connectors[SIRI_SUBSCRIPTION_REQUEST_DISPATCHER] = connector
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	subscription := partner.Subscriptions().FindOrCreateByKind(StopMonitoringCollect)
	subscription.Save()

	if _, ok := partner.Subscriptions().Find("6ba7b814-9dad-11d1-0-00c04fd430c8"); !ok {
		t.Fatalf("Subscriptions should be found")
	}

	connector.HandleSubscriptionTerminatedNotification(response)

	if _, ok := partner.Subscriptions().Find("6ba7b814-9dad-11d1-0-00c04fd430c8"); ok {
		t.Errorf("Subscriptions should not be found")
	}
}

func Test_HandleNotifySubscriptionTerminated(t *testing.T) {
	file, err := os.Open("testdata/notify-subscription-terminated-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLNotifySubscriptionTerminatedFromContent(content)

	partners := createTestPartnerManager()

	partner := partners.New("slug")

	settings := map[string]string{
		"remote_url":           "une url",
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRISubscriptionRequestDispatcher(partner)
	partner.connectors[SIRI_SUBSCRIPTION_REQUEST_DISPATCHER] = connector

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	subscription := partner.Subscriptions().FindOrCreateByKind(StopMonitoringCollect)
	subscription.Save()

	if _, ok := partner.Subscriptions().Find("6ba7b814-9dad-11d1-0-00c04fd430c8"); !ok {
		t.Fatalf("Subscriptions should be found")
	}

	connector.HandleNotifySubscriptionTerminated(response)

	if _, ok := partner.Subscriptions().Find("6ba7b814-9dad-11d1-0-00c04fd430c8"); ok {
		t.Errorf("Subscriptions should not be found")
	}
}
