package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

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
