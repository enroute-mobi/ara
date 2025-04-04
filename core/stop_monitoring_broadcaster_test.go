package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/assert"
)

func Test_StopMonitoringBroadcaster_Create_Events(t *testing.T) {
	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewTestMemoryModel()

	referential.model.SetBroadcastSMChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	code := model.NewCode("internal", string(stopArea.Id()))
	stopArea.SetCode(code)

	reference := model.Reference{
		Code: &code,
		Type: "StopArea",
	}

	subs := partner.Subscriptions().New("kind")
	subs.Save()
	subs.CreateAndAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	time.Sleep(5 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(5 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	if len(connector.(*TestStopMonitoringSubscriptionBroadcaster).events) != 2 {
		t.Error("2 events should have been generated got: ", len(connector.(*TestStopMonitoringSubscriptionBroadcaster).events))
	}
}

func Test_StopMonitoringBroadcaster_HandleStopMonitoringBroadcastWithReferent(t *testing.T) {
	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Save()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("internal", string(stopArea.Id()))
	stopArea.SetCode(code)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.ReferentId = stopArea.Id()
	stopArea2.Save()

	reference := model.Reference{
		Code: &code,
		Type: "StopArea",
	}

	subs := partner.Subscriptions().New("StopMonitoringBroadcast")
	subs.CreateAndAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea2.Id()
	stopVisit.Save()

	event := &model.StopMonitoringBroadcastEvent{
		ModelId:   string(stopVisit.Id()),
		ModelType: "StopVisit",
	}

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(event)
	if len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast) != 1 {
		t.Error("1 events should have been generated got: ", len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast))
	}
}

func Test_StopMonitoringBroadcaster_HandleStopMonitoringBroadcastWithLineRefFilter(t *testing.T) {
	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Save()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("internal", string(stopArea.Id()))
	stopArea.SetCode(code)
	stopArea.Save()

	reference := model.Reference{
		Code: &code,
		Type: "StopArea",
	}

	sub := partner.Subscriptions().New("StopMonitoringBroadcast")
	sub.CreateAndAddNewResource(reference)
	sub.SetExternalId("externalId")
	sub.SetSubscriptionOption("LineRef", "incorrect:lineRef")
	sub.Save()

	line := referential.Model().Lines().New()
	line_code := model.NewCode("internal", "line")
	line.SetCode(line_code)
	line.Save()

	vj := referential.Model().VehicleJourneys().New()
	vj.LineId = line.Id()
	vj.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vj.Id()
	stopVisit.Save()

	event := &model.StopMonitoringBroadcastEvent{
		ModelId:   string(stopVisit.Id()),
		ModelType: "StopVisit",
	}

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(event)
	if len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast) != 0 {
		t.Error("0 events should have been generated got: ", len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast))
	}

	sub.SetSubscriptionOption("LineRef", "internal:line")
	sub.Save()

	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(event)
	if len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast) != 1 {
		t.Error("1 events should have been generated got: ", len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast))
	}
}

func Test_StopMonitoringBroadcaster_Receive_Notify(t *testing.T) {
	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)

	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = io.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test referential
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		s.BROADCAST_SIRI_SM_MULTIPLE_SUBSCRIPTIONS: "true",
		"remote_code_space":                        "internal",
		"local_credential":                         "external",
		"remote_url":                               ts.URL,
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster = NewFakeStopMonitoringBroadcaster(connector.(*SIRIStopMonitoringSubscriptionBroadcaster))

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	code := model.NewCode("internal", string(stopArea.Id()))
	stopArea.SetCode(code)

	reference := model.Reference{
		Code: &code,
		Type: "StopArea",
	}

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.Save()

	code2 := model.NewCode("internal", string(stopArea2.Id()))
	stopArea2.SetCode(code2)

	reference2 := model.Reference{
		Code: &code2,
		Type: "StopArea",
	}

	subscription := partner.Subscriptions().New("StopMonitoringBroadcast")
	subscription.SubscriberRef = "subscriber"
	subscription.SetExternalId("externalId")
	subscription.CreateAndAddNewResource(reference)
	subscription.CreateAndAddNewResource(reference2)
	subscription.subscriptionOptions["ChangeBeforeUpdates"] = "PT4M"
	// subscription.subscriptionOptions["MaximumStopVisits"] = "1"
	subscription.Save()

	line := referential.Model().Lines().New()
	line.SetCode(code)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetCode(code)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetCode(model.NewCode("internal", "svid"))
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.SetCode(model.NewCode("internal", "svid2"))
	stopVisit2.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	stopVisit3 := referential.Model().StopVisits().New()
	stopVisit3.StopAreaId = stopArea2.Id()
	stopVisit3.VehicleJourneyId = vehicleJourney.Id()
	stopVisit3.SetCode(model.NewCode("internal", "svid3"))
	stopVisit3.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()
	time.Sleep(10 * time.Millisecond)
	stopVisit2.Save()
	time.Sleep(10 * time.Millisecond)
	stopVisit3.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Start()

	notify, _ := sxml.NewXMLNotifyStopMonitoringFromContent(response)
	delivery := notify.StopMonitoringDeliveries()

	if len(delivery) != 2 {
		t.Errorf("Should have received 2 deliveries but got == %v", len(delivery))
	}

	if delivery[0].SubscriberRef() != "subscriber" {
		t.Errorf("SubscriberRef should be subscriber but got == %v", delivery[0].SubscriberRef())
	}

	if delivery[0].SubscriptionRef() != "externalId" {
		t.Errorf("SubscriptionRef should be externalId but got == %v", delivery[0].SubscriptionRef())
	}

	sv := delivery[0].XMLMonitoredStopVisits()
	var expected int
	if delivery[0].MonitoringRef() == string(stopArea.Id()) {
		expected = 2
	} else {
		expected = 1
	}
	if len(sv) != expected {
		t.Errorf("Should have received %v StopVisits but got == %v", expected, len(sv))
	}

	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond)
	if len := len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast); len != 0 {
		t.Errorf("No stopVisit should need to be broadcasted %v", len)
	}

	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond)
	if len := len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast); len != 1 {
		t.Errorf("1 stopVisit should need to be broadcasted %v", len)
	}
}

func Test_StopMonitoringBroadcaster_Receive_Two_Notifications(t *testing.T) {
	assert := assert.New(t)

	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)

	// Create a test http server
	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = io.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test referential
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	f := audit.NewFakeBigQuery()
	audit.SetCurrentBigQuery("Un Referential Plutot Cool", f)

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"broadcast.siri.stop_monitoring.multiple_subscriptions": "true",
		"remote_code_space": "internal",
		"local_credential":  "external",
		"remote_url":        ts.URL,
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster = NewFakeStopMonitoringBroadcaster(connector.(*SIRIStopMonitoringSubscriptionBroadcaster))

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	code := model.NewCode("internal", string(stopArea.Id()))
	stopArea.SetCode(code)

	reference := model.Reference{
		Code: &code,
		Type: "StopArea",
	}

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.Save()

	code2 := model.NewCode("internal", string(stopArea2.Id()))
	stopArea2.SetCode(code2)

	reference2 := model.Reference{
		Code: &code2,
		Type: "StopArea",
	}

	subscription := partner.Subscriptions().New("StopMonitoringBroadcast")
	subscription.SubscriberRef = "subscriber"
	subscription.SetExternalId("externalId")
	subscription.CreateAndAddNewResource(reference)
	subscription.subscriptionOptions["ChangeBeforeUpdates"] = "PT4M"
	subscription.Save()

	subscription2 := partner.Subscriptions().New("StopMonitoringBroadcast")
	subscription2.SubscriberRef = "subscriber"
	subscription2.SetExternalId("externalId2")
	subscription2.CreateAndAddNewResource(reference2)
	subscription2.subscriptionOptions["ChangeBeforeUpdates"] = "PT4M"
	subscription2.Save()

	line := referential.Model().Lines().New()
	line.SetCode(code)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetCode(code)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetCode(code)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea2.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.SetCode(model.NewCode("internal", string(stopArea2.Id())))
	stopVisit2.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()
	time.Sleep(10 * time.Millisecond)
	stopVisit2.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Start()

	notify, _ := sxml.NewXMLNotifyStopMonitoringFromContent(response)
	delivery := notify.StopMonitoringDeliveries()

	assert.Lenf(delivery, 2,
		"Should have received 2 deliveries but got %v", len(delivery))

	assert.Len(f.Messages(), 1, "1 message should be sent to BigQuery")
	assert.ElementsMatch([]string{"externalId", "externalId2"},
		f.Messages()[0].SubscriptionIdentifiers,
		"2 subscriptionIdentifiers should be logged in BigQuery")
}

func Test_StopMonitoringBroadcaster_Receive_Two_Notifications_With_MaxPerDelivery(t *testing.T) {
	assert := assert.New(t)

	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)

	// Create a test referential
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	f := audit.NewFakeBigQuery()
	audit.SetCurrentBigQuery("Un Referential Plutot Cool", f)

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		s.BROADCAST_SIRI_SM_MULTIPLE_SUBSCRIPTIONS:         "false",
		s.BROADCAST_SIRI_SM_MAXIMUM_RESOURCES_PER_DELIVERY: "1",
		"remote_code_space":                                "internal",
		"local_credential":                                 "external",
		"remote_url":                                       "URL",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster = NewFakeStopMonitoringBroadcaster(connector.(*SIRIStopMonitoringSubscriptionBroadcaster))

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	code := model.NewCode("internal", string(stopArea.Id()))
	stopArea.SetCode(code)

	reference := model.Reference{
		Code: &code,
		Type: "StopArea",
	}

	subscription := partner.Subscriptions().New("StopMonitoringBroadcast")
	subscription.SubscriberRef = "subscriber"
	subscription.SetExternalId("externalId")
	subscription.CreateAndAddNewResource(reference)
	subscription.subscriptionOptions["ChangeBeforeUpdates"] = "PT4M"
	subscription.Save()

	line := referential.Model().Lines().New()
	line.SetCode(code)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetCode(code)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetCode(code)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.SetCode(model.NewCode("internal", string(stopArea.Id())))
	stopVisit2.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()
	time.Sleep(10 * time.Millisecond)
	stopVisit2.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Start()

	assert.Len(f.Messages(), 2, "2 messages should be sent to BigQuery")
}
