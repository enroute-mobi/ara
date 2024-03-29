package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/assert"
)

func Test_GeneralMessageBroadcaster_Create_Events(t *testing.T) {
	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Start()
	defer referential.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)

	situation := referential.Model().Situations().New()

	code := model.NewCode("internal", string(situation.Id()))
	situation.SetCode(code)

	reference := model.Reference{
		Code: &code,
		Type: "Situation",
	}

	subs := partner.Subscriptions().New("")
	subs.Save()
	subs.CreateAndAddNewResource(reference)
	subs.Save()
	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...

	situation.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	if len(connector.(*TestGeneralMessageSubscriptionBroadcaster).events) != 1 {
		t.Error("1 event should have been generated got: ", len(connector.(*TestGeneralMessageSubscriptionBroadcaster).events))
	}
}

func Test_GeneralMessageBroadcaster_Receive_Notify(t *testing.T) {
	assert := assert.New(t)
	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = io.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test http server
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Start()
	defer referential.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
		"local_credential":  "external",
		"remote_url":        ts.URL,
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
	connector.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster = NewFakeGeneralMessageBroadcaster(connector.(*SIRIGeneralMessageSubscriptionBroadcaster))

	situation := referential.Model().Situations().New()
	period := &model.TimeRange{EndTime: referential.Clock().Now().Add(5 * time.Minute)}
	situation.ValidityPeriods = []*model.TimeRange{period}
	situation.Keywords = []string{"Perturbation"}

	code := model.NewCode("internal", string(situation.Id()))
	situation.SetCode(code)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()
	code2 := model.NewCode("internal", "value")
	stopArea.SetCode(code2)
	stopArea.Save()

	affectedStopArea := model.NewAffectedStopArea()
	affectedStopArea.StopAreaId = stopArea.Id()
	situation.Affects = append(situation.Affects, affectedStopArea)

	code3 := model.NewCode("SituationResource", "Situation")

	reference := model.Reference{
		Code: &code3,
		Type: "Situation",
	}

	subscription := partner.Subscriptions().FindOrCreateByKind("GeneralMessageBroadcast")
	subscription.SubscriberRef = "subscriber"
	subscription.CreateAndAddNewResource(reference)

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	situation.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIGeneralMessageSubscriptionBroadcaster).generalMessageBroadcaster.Start()

	notify, err := sxml.NewXMLNotifyGeneralMessageFromContent(response)
	assert.Nil(err)
	delivery := notify.GeneralMessagesDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	if delivery[0].SubscriberRef() != "subscriber" {
		t.Errorf("SubscriberRef should be subscriber but got == %v", delivery[0].SubscriptionRef())
	}

	sv := delivery[0].XMLGeneralMessages()

	if len(sv) != 1 {
		t.Errorf("Should have received 1 GeneralMessage but got == %v\n%v", len(sv), sv)
	}
}
