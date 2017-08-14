package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SubscriptionRequest_Dispatch_SM(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST, SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	objectid := model.NewObjectID("_internal", "coicogn2")

	stopArea.SetObjectID(objectid)
	stopArea.Save()

	file, _ := os.Open("testdata/stopmonitoringsubscription-request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLSubscriptionRequestFromContent(body)

	response := connector.(*SIRISubscriptionRequest).Dispatch(request)

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

	if sub.ExternalId() != "6ba7b814-9dad-11d1-2-00c04fd430c8" {
		t.Errorf("Wrong ExternalId value want: 6ba7b814-9dad-11d1-2-00c04fd430c8 got: %v", sub.ExternalId())
	}

	if sub.partner.Id() != partner.Id() {
		t.Errorf("Wrong Partner Id value want: %v got: %v", partner.Id(), sub.partner.Id())
	}
}

func Test_SubscriptionRequest_Dispatch_GM(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "_internal"
	partner.ConnectorTypes = []string{SIRI_SUBSCRIPTION_REQUEST, SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_SUBSCRIPTION_REQUEST)

	file, _ := os.Open("testdata/generalmessagesubscription-request-soap.xml")
	body, _ := ioutil.ReadAll(file)
	request, _ := siri.NewXMLSubscriptionRequestFromContent(body)

	response := connector.(*SIRISubscriptionRequest).Dispatch(request)

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

	if sub.ExternalId() != "6ba7b814-9dad-11d1-2-00c04fd430c8" {
		t.Errorf("Wrong ExternalId value want: 6ba7b814-9dad-11d1-2-00c04fd430c8 got: %v", sub.ExternalId())
	}

	if sub.partner.Id() != partner.Id() {
		t.Errorf("Wrong Partner Id value want: %v got: %v", partner.Id(), sub.partner.Id())
	}
}
