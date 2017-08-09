package core

import (
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_GeneralMessageBroadcaster_Create_Events(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")

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
