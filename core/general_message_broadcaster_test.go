package core

import (
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

	objectid := model.NewObjectID("internal", string(situation.Id()))
	situation.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "Situation",
	}

	subs := partner.Subscriptions().New("kind")
	subs.Save()
	subs.CreateAddNewResource(reference)
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

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = ioutil.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test http server
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Start()
	defer referential.Stop()

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
	situation.ValidUntil = referential.Clock().Now().Add(5 * time.Minute)

	objectid := model.NewObjectID("internal", string(situation.Id()))
	situation.SetObjectID(objectid)
	routeReference := model.NewReference(model.NewObjectID("internal", "value"))
	routeReference.Type = "RouteRef"
	situation.References = append(situation.References, routeReference)

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "Situation",
	}

	subscription := partner.Subscriptions().FindOrCreateByKind("situation")
	subscription.CreateAddNewResource(reference)

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	situation.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster.Start()
	defer connector.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster.Stop()

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
