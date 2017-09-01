package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
	"github.com/jbowtie/gokogiri/xml"
)

func Test_SIRIStopmonitoringSubscriptionsCollector_HandleNotifyStopMonitoring(t *testing.T) {
	collectManager := NewTestCollectManager()
	referential := &Referential{
		collectManager: collectManager,
		model:          model.NewMemoryModel(),
	}
	referential.Model().StopAreas().(*model.MemoryStopAreas).SetUUIDGenerator(model.NewFakeUUIDGenerator())

	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("_internal", "coicogn2")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	objectid2 := model.NewObjectID("_internal", "coicogn3")
	stopArea2.SetObjectID(objectid2)
	stopArea2.Save()

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.Settings["generators.subscription_identifier"] = "Subscription::%{id}::LOC"

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	file, err := os.Open("testdata/notify-stop-monitoring.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := xml.Parse(content, xml.DefaultEncodingBytes, nil, xml.StrictParseOption, xml.DefaultEncodingBytes)
	if err != nil {
		t.Fatal(err)
	}

	deliveries := siri.NewXMLNotifyStopMonitoring(doc.Root())

	partner.Subscriptions().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	subscription, _ := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")
	subscription.Save()

	connector.HandleNotifyStopMonitoring(deliveries)

	if len(collectManager.(*TestCollectManager).Events) != 2 {
		t.Errorf("Wrong number of events in collectManager, expected 2 got %v", len(collectManager.(*TestCollectManager).Events))
	}

	for _, event := range collectManager.(*TestCollectManager).Events {
		if event.StopAreaId == model.StopAreaId("6ba7b814-9dad-11d1-0-00c04fd430c8") && len(event.StopVisitUpdateEvents) != 2 {
			t.Errorf("StopArea 6ba7b814-9dad-11d1-0-00c04fd430c8 should have 2 StopVisitEvents, got %v", len(event.StopVisitUpdateEvents))
		} else if event.StopAreaId == model.StopAreaId("6ba7b814-9dad-11d1-1-00c04fd430c8") && len(event.StopVisitUpdateEvents) != 1 {
			t.Errorf("StopArea 6ba7b814-9dad-11d1-1-00c04fd430c8 should have 1 StopVisitEvent, got %v", len(event.StopVisitUpdateEvents))
		} else if event.StopAreaId != model.StopAreaId("6ba7b814-9dad-11d1-0-00c04fd430c8") && event.StopAreaId != model.StopAreaId("6ba7b814-9dad-11d1-1-00c04fd430c8") {
			t.Errorf("Wrong StopAreaId, want 6ba7b814-9dad-11d1-0-00c04fd430c8 or 6ba7b814-9dad-11d1-1-00c04fd430c8, got %v", event.StopAreaId)
		}
	}
}

func Test_SIRIStopmonitoringSubscriptionsCollector_AddtoRessource(t *testing.T) {

	response, err := os.Open("testdata/stopmonitoringsubscription-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Close()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		io.Copy(w, response)
	}))
	defer ts.Close()

	// Create a SIRIStopMonitoringRequestCollector
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url":           ts.URL,
			"remote_objectid_kind": "test_kind",
		},
		manager: partners,
		// connectors: make(map[string]Connector),
	}
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	objectid := model.NewObjectID("test_kind", "value")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	connector.SetStopMonitoringSubscriber(NewFakeStopMonitoringSubscriber(connector))
	connector.RequestStopAreaUpdate(stopAreaUpdateRequest)
	subscription, _ := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	if len(subscription.ResourcesByObjectID()) != 1 {
		t.Errorf("Response should have 1 ressource but got %v\n", len(subscription.ResourcesByObjectID()))
	}
}

func Test_SIRIStopMonitoringSubscriptionTerminationCollector(t *testing.T) {
	file, err := os.Open("../siri/testdata/subscription_terminated_notification-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := siri.NewXMLStopMonitoringSubscriptionTerminatedResponseFromContent(content)
	connectors := make(map[string]Connector)

	partners := createTestPartnerManager()
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url":           "une url",
			"remote_objectid_kind": "_internal",
		},
		ConnectorTypes: []string{"siri-stop-monitoring-subscription-collector"},
		manager:        partners,
		connectors:     connectors,
	}

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)
	connectors[SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR] = connector

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(model.NewFakeUUIDGenerator())

	referential := partner.Referential()
	stopArea := referential.Model().StopAreas().New()
	stopArea.CollectedAlways = false
	objectid := model.NewObjectID("_internal", "coicogn2")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.Collected(time.Now())
	objectid = model.NewObjectID("_internal", "stopvisit1")
	stopVisit.SetObjectID(objectid)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Save()

	objId := model.NewObjectID("_internal", "coicogn2")
	ref := model.Reference{
		ObjectId: &objId,
		Id:       string(stopArea.Id()),
		Type:     "StopArea",
	}

	subscription, _ := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")
	subscription.CreateAddNewResource(ref)
	subscription.Save()

	connector.HandleTerminatedNotification(response)

	if _, ok := connector.partner.Subscriptions().Find("6ba7b814-9dad-11d1-0-00c04fd430c8"); ok {
		t.Errorf("Subscriptions should not be found \n")
	}
}

func Test_SIRIStopMonitoringSubscriptionCollector(t *testing.T) {

	request := &siri.XMLSubscriptionRequest{}
	// Create a test http server

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		body, _ := ioutil.ReadAll(r.Body)
		request, _ = siri.NewXMLSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	// Create a SIRIStopMonitoringRequestCollector
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"local_url":                          "http://example.com/test/siri",
			"remote_url":                         ts.URL,
			"remote_objectid_kind":               "test_kind",
			"generators.subscription_identifier": "Subscription::%{id}::LOC",
		},
		manager: partners,
	}

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partner.subscriptionManager.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	partners.Save(partner)

	objectid := model.NewObjectID("test_kind", "value")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	objectid2 := model.NewObjectID("test_kind", "value2")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(objectid2)
	stopArea2.Save()

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	stopAreaUpdateEvent := NewStopAreaUpdateRequest(stopArea.Id())
	connector.SetStopMonitoringSubscriber(NewFakeStopMonitoringSubscriber(connector))
	connector.RequestStopAreaUpdate(stopAreaUpdateEvent)
	connector.stopMonitoringSubscriber.Start()

	subscription, _ := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	if expected := "http://example.com/test/siri"; request.ConsumerAddress() != expected {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", request.ConsumerAddress(), expected)
	}

	expectedUuid := fmt.Sprintf("%v", subscription.Id())
	if request.XMLSubscriptionSMEntries()[0].SubscriptionIdentifier() != expectedUuid {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", request.XMLSubscriptionSMEntries()[0].SubscriptionIdentifier(), expectedUuid)
	}
}

func Test_SIRIStopMonitoringTerminatedSubscriptionRequest(t *testing.T) {
	request := &siri.XMLTerminatedSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		request, _ = siri.NewXMLTerminatedSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"local_url":                          "http://example.com/test/siri",
			"remote_url":                         ts.URL,
			"remote_objectid_kind":               "test_kind",
			"generators.subscription_identifier": "Subscription::%{id}::LOC",
		},
		manager: partners,
	}
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	file, _ := os.Open("testdata/notify-stop-monitoring.xml")
	content, _ := ioutil.ReadAll(file)

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	notify, _ := siri.NewXMLNotifyStopMonitoringFromContent(content)

	connector.HandleNotifyStopMonitoring(notify)

	if expected := "Subscription::6ba7b814-9dad-11d1-0-00c04fd430c8::LOC"; request.SubscriptionRef() != expected {
		t.Errorf("Wrong SubscriptionRef want : %v  got %v :", expected, request.SubscriptionRef())
	}
}
