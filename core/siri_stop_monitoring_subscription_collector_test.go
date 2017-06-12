package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

// func Test_StopMonitoringCancelled(t *testing.T) {
// 	partners := createTestPartnerManager()
// 	partner := &Partner{
// 		context: make(Context),
// 		Settings: map[string]string{
// 			"remote_url":           "Une Magnifique Url",
// 			"remote_objectid_kind": "test kind",
// 		},
// 		manager: partners,
// 	}

// 	partners.Save(partner)

// 	siriStopMonitoringSubscriptionCollector := NewSIRIStopMonitoringSubscriptionCollector(partner)

// 	fs := fakeBroadcaster{}
// 	siriStopMonitoringSubscriptionCollector.SetStopAreaUpdateSubscriber(fs.FakeBroadcaster)
// 	siriStopMonitoringSubscriptionCollector.SetClock(model.NewFakeClock())
// 	cancelStopVisitMonitoring := make(map[string][]string)
// 	cancelStopVisitMonitoring["STIF:StopPoint:Q:411415:"] = []string{"SNCF-ACCES:Item::411415_125212:LOC"}
// 	// siriStopMonitoringSubscriptionCollector.CancelStopVisitMonitoring(cancelStopVisitMonitoring)

// 	time.Sleep(42 * time.Millisecond)
// 	if len(fs.Events) != 1 {
// 		t.Error("Events should have a lenght of 1 but got: ", len(fs.Events))
// 	}

// 	if len(fs.Events[0].StopVisitNotCollectedEvents) != 1 {
// 		t.Error(".Events.StopVisitNotCollectedEvents should have a lenght of 1 but got: ", len(fs.Events[0].StopVisitNotCollectedEvents))
// 	}

// 	if fs.Events[0].StopVisitNotCollectedEvents[0].StopVisitObjectId.Kind() != "StopMonitoring" {
// 		t.Error("Kind of the event should be 'StopMonitoring' but got: ", fs.Events[0].StopVisitNotCollectedEvents[0].StopVisitObjectId.Kind())
// 	}

// 	if fs.Events[0].StopVisitNotCollectedEvents[0].StopVisitObjectId.Value() != "SNCF-ACCES:Item::411415_125212:LOC" {
// 		t.Error("Id of the event should be 'SNCF-ACCES:Item::411415_125212:LOC' but got:", fs.Events[0].StopVisitNotCollectedEvents[0].StopVisitObjectId.Value())
// 	}
// }

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
