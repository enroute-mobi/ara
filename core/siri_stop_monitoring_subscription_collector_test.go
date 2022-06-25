package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/jbowtie/gokogiri/xml"
)

func Test_SIRIStopmonitoringSubscriptionsCollector_HandleNotifyStopMonitoring(t *testing.T) {
	collectManager := NewTestCollectManager()
	referential := &Referential{
		collectManager: collectManager,
		model:          model.NewMemoryModel(),
	}
	referential.Model().StopAreas().(*model.MemoryStopAreas).SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

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
	partner.SetSetting("remote_objectid_kind", "_internal")
	partner.SetSetting("generators.subscription_identifier", "Subscription::%{id}::LOC")

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

	deliveries := sxml.NewXMLNotifyStopMonitoring(doc.Root())

	partner.Subscriptions().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoringCollect")
	subscription.Save()

	connector.HandleNotifyStopMonitoring(deliveries)

	// 2 StopAreas 1 Line 3 VehicleJourneys 3 StopVisits
	if len(collectManager.(*TestCollectManager).UpdateEvents) != 9 {
		t.Errorf("Wrong number of events in collectManager, expected 9 got %v", len(collectManager.(*TestCollectManager).UpdateEvents))
	}
}

func Test_SIRIStopmonitoringSubscriptionsCollector_AddtoResource(t *testing.T) {

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

	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"remote_url":           ts.URL,
		"remote_objectid_kind": "test_kind",
	})
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
	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoringCollect")

	if subscription.ResourcesLen() != 1 {
		t.Errorf("Response should have 1 ressource but got %v\n", subscription.ResourcesLen())
	}
}

func Test_SIRIStopMonitoringSubscriptionCollector(t *testing.T) {

	request := &sxml.XMLSubscriptionRequest{}
	// Create a test http server

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		body, _ := ioutil.ReadAll(r.Body)
		request, _ = sxml.NewXMLSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	// Create a SIRIStopMonitoringRequestCollector
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"local_url":                          "http://example.com/test/siri",
		"remote_url":                         ts.URL,
		"remote_objectid_kind":               "test_kind",
		"generators.subscription_identifier": "Subscription::%{id}::LOC",
	})

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
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

	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoringCollect")

	if expected := "http://example.com/test/siri"; request.ConsumerAddress() != expected {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", request.ConsumerAddress(), expected)
	}

	expectedUuid := fmt.Sprintf("%v", subscription.Id())
	if request.XMLSubscriptionSMEntries()[0].SubscriptionIdentifier() != expectedUuid {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", request.XMLSubscriptionSMEntries()[0].SubscriptionIdentifier(), expectedUuid)
	}
}

func Test_SIRIStopMonitoringDeleteSubscriptionRequest(t *testing.T) {
	request := &sxml.XMLDeleteSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		request, _ = sxml.NewXMLDeleteSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"local_url":                          "http://example.com/test/siri",
		"remote_url":                         ts.URL,
		"remote_objectid_kind":               "test_kind",
		"generators.subscription_identifier": "Subscription::%{id}::LOC",
	})
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	file, _ := os.Open("testdata/notify-stop-monitoring.xml")
	content, _ := ioutil.ReadAll(file)

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)
	connector.deletedSubscriptions = NewDeletedSubscriptions()

	notify, _ := sxml.NewXMLNotifyStopMonitoringFromContent(content)

	connector.HandleNotifyStopMonitoring(notify)

	if expected := "Subscription::6ba7b814-9dad-11d1-0-00c04fd430c8::LOC"; request.SubscriptionRef() != expected {
		t.Errorf("Wrong SubscriptionRef want : %v  got %v :", expected, request.SubscriptionRef())
	}
}
