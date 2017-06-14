package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

	delivery := siri.NewXMLStopMonitoringResponse(doc.Root())
	connector.HandleNotifyStopMonitoring(delivery)

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
	partners := createTestPartnerManager()

	partner := &Partner{
		context: make(Context),
		manager: partners,
	}

	partners.Save(partner)
	partner.subscriptionManager = NewMemorySubscriptions(partner)

	file, err := os.Open("testdata/stopmonitoringdeliveries-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	xmlRequest, _ := siri.NewXMLStopMonitoringResponseFromContent(content)
	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	stopvisits := xmlRequest.XMLMonitoredStopVisits()
	stopvisit := stopvisits[0]
	stopAreaUpdateRequest := NewStopAreaUpdateRequest(model.StopAreaId(stopvisit.StopPointRef()))
	connector.RequestStopAreaUpdate(stopAreaUpdateRequest)
	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")
	if len(subscription.ResourcesByObjectID()) != 1 {
		t.Errorf("Response should have 1 ressource but got %v\n", len(subscription.ResourcesByObjectID()))
	}

	connector.RequestStopAreaUpdate(stopAreaUpdateRequest)
	subscription = connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")
	if len(subscription.ResourcesByObjectID()) != 1 {
		t.Errorf("Response should have 1 ressource but got %v\n", len(subscription.ResourcesByObjectID()))
	}

}

func Test_SIRIStopMonitoringSubscriptionCollector(t *testing.T) {

	request := &siri.XMLStopMonitoringSubscriptionRequest{}
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		body, _ := ioutil.ReadAll(r.Body)
		request, _ = siri.NewXMLStopMonitoringSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	connectors := make(map[string]Connector)

	// Create a SIRIStopMonitoringRequestCollector
	partners := createTestPartnerManager()
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url":           ts.URL,
			"remote_objectid_kind": "test kind",
		},
		ConnectorTypes: []string{"siri-stop-monitoring-deliveries-response-collector"},
		manager:        partners,
		connectors:     connectors,
	}

	connectors[SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR] = NewSIRIStopMonitoringSubscriptionCollector(partner)
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	stopAreaUpdateEvent := NewStopAreaUpdateRequest(model.StopAreaId("NINOXE:StopPoint:SP:24:LOC"))
	connectors[SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR].(StopMonitoringSubscriptionCollector).RequestStopAreaUpdate(stopAreaUpdateEvent)

	if request.MonitoringRef() != "NINOXE:StopPoint:SP:24:LOC" {
		t.Errorf("Wrong MonitoringRef:\n got: %v\nwant: %v", request.MonitoringRef(), "NINOXE:StopPoint:SP:24:LOC")
	}

	if request.ConsumerAddress() != "https://edwig-staging.af83.io/test/siri" {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", request.ConsumerAddress(), "https://edwig-staging.af83.io/test/siri")
	}

	if request.SubscriptionIdentifier() != "Edwig:Subscription::NINOXE:StopPoint:SP:24:LOC:LOC" {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", request.SubscriptionIdentifier(), "Edwig:Subscription::NINOXE:StopPoint:SP:24:LOC:LOC")
	}
}
