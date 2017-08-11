package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_GeneralMessageBroadcaster_Create_Events(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Start()
	defer referential.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.ConnectorTypes = []string{TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)

	situation := referential.Model().Situations().New()
	//situation.Save()

	objectid := model.NewObjectID("internal", string(situation.Id()))
	situation.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Id:       string(situation.Id()),
		Type:     "Situation",
	}

	subs := partner.Subscriptions().New()
	subs.Save()
	subs.CreateAddNewResource(reference)
	subs.SetKind(string(subs.Id()))
	subs.Save()
	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...

	situation.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	if len(connector.(*TestGeneralMessageSubscriptionBroadcaster).events) != 1 {
		t.Error("1 event should have been generated got: ", len(connector.(*TestGeneralMessageSubscriptionBroadcaster).events))
	}
}

func Test_GeneralMessageBroadcaster_Receive_Notify(t *testing.T) {
	// Create a test http server

	fmt.Println("esh komment allet vous")
	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = ioutil.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test http server
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.Settings["remote_credential"] = "external"
	partner.Settings["remote_url"] = ts.URL

	partner.ConnectorTypes = []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
	connector.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster = NewFakeGeneralMessageBroadcaster(connector.(*SIRIGeneralMessageSubscriptionBroadcaster))

	situation := referential.Model().Situations().New()

	objectid := model.NewObjectID("internal", string(situation.Id()))
	situation.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Id:       string(situation.Id()),
		Type:     "Situation",
	}

	subscription, _ := partner.Subscriptions().FindOrCreateByKind("This Kind should normaly be the exterior Subscription Id")
	subscription.CreateAddNewResource(reference)

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	situation.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster.Start()

	notify, _ := siri.NewXMLNotifyGeneralMessageFromContent(response)
	delivery := notify.GeneralMessagesDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	if delivery[0].SubscriberRef() != "external" {
		t.Errorf("SubscriberRef should be external but got == %v", delivery[0].SubscriptionRef())
	}

	sv := delivery[0].XMLGeneralMessages()

	if len(sv) != 1 {
		t.Errorf("Should have received 1 GeneralMessage but got == %v\n%v", len(sv), sv)
	}
}
