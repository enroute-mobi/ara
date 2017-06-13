package core

import (
	"io/ioutil"
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
	}

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
		if event.StopAreaId == model.StopAreaId("coicogn2") && len(event.StopVisitUpdateEvents) != 2 {
			t.Errorf("StopArea coicogn2 should have 2 StopVisitEvents, got %v", len(event.StopVisitUpdateEvents))
		} else if event.StopAreaId == model.StopAreaId("coicogn3") && len(event.StopVisitUpdateEvents) != 1 {
			t.Errorf("StopArea coicogn3 should have 1 StopVisitEvent, got %v", len(event.StopVisitUpdateEvents))
		} else if event.StopAreaId != model.StopAreaId("coicogn2") && event.StopAreaId != model.StopAreaId("coicogn3") {
			t.Errorf("Wrong StopAreaId, want coicogn2 or coicogn2, got %v", event.StopAreaId)
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
