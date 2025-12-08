package core

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/jbowtie/gokogiri/xml"
	"github.com/stretchr/testify/assert"
)

func Test_SIRIStopmonitoringSubscriptionsCollector_HandleNotifyStopMonitoring(t *testing.T) {
	collectManager := NewTestCollectManager()
	referential := &Referential{
		collectManager: collectManager,
		model:          model.NewTestMemoryModel(),
	}
	referential.Model().StopAreas().(*model.MemoryStopAreas).SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("_internal", "coicogn2")
	stopArea.SetCode(code)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	code2 := model.NewCode("_internal", "coicogn3")
	stopArea2.SetCode(code2)
	stopArea2.Save()

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_code_space":                  "_internal",
		"generators.subscription_identifier": "Subscription::%{id}::LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	file, err := os.Open("testdata/notify-stop-monitoring.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := xml.Parse(content, xml.DefaultEncodingBytes, nil, xml.StrictParseOption, xml.DefaultEncodingBytes)
	if err != nil {
		t.Fatal(err)
	}

	deliveries := sxml.NewXMLNotifyStopMonitoring(doc.Root())

	partner.Subscriptions().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	subscription := connector.partner.Subscriptions().FindOrCreateByKind(StopMonitoringCollect)
	subscription.Save()

	connector.HandleNotifyStopMonitoring(deliveries)

	// 2 StopAreas 1 Line 3 VehicleJourneys 3 StopVisits
	if len(collectManager.(*TestCollectManager).UpdateEvents) != 9 {
		t.Errorf("Wrong number of events in collectManager, expected 9 got %v", len(collectManager.(*TestCollectManager).UpdateEvents))
	}
}

func Test_SIRIStopmonitoringSubscriptionsCollector_AddtoResource(t *testing.T) {
	assert := assert.New(t)
	response, err := os.Open("testdata/stopmonitoringsubscription-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Close()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		io.Copy(w, response)
	}))
	defer ts.Close()

	// Create a SIRIStopMonitoringRequestCollector
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := partners.New("slug")
	settings := map[string]string{
		"remote_url":        ts.URL,
		"remote_code_space": "test_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	code := model.NewCode("test_kind", "value")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	connector.SetStopMonitoringSubscriber(NewFakeStopMonitoringSubscriber(connector))
	connector.RequestStopAreaUpdate(stopAreaUpdateRequest)

	subscriptions := connector.partner.Subscriptions().FindAll()
	assert.Len(subscriptions, 1)
	assert.Equal(subscriptions[0].ResourcesLen(), 1)

	subscriptionResource := subscriptions[0].UniqueResource()
	assert.NotNil(subscriptionResource)
	assert.Equal(subscriptionResource.Reference.Type, "StopArea")
	assert.Equal(subscriptionResource.Reference.Code.String(), "test_kind:value")

	// Adding a new subscription
	subscription := connector.partner.Subscriptions().FindOrCreateByKind(StopMonitoringCollect)
	assert.Len(connector.partner.Subscriptions().FindAll(), 2)
	assert.Equal(subscription.ResourcesLen(), 0)
}

func Test_SIRIStopMonitoringSubscriptionCollector(t *testing.T) {
	assert := assert.New(t)

	// Create a SIRIStopMonitoringRequestCollector
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	code := model.NewCode("test_kind", "value")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	code2 := model.NewCode("test_kind", "value2")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetCode(code2)
	stopArea2.Save()

	code3 := model.NewCode("test_kind", "value3")
	stopArea3 := referential.Model().StopAreas().New()
	stopArea3.SetCode(code3)
	stopArea3.Save()

	var TestCases = []struct {
		testNumber                    int
		settings                      map[string]string
		expectedNumberOfRequests      int
		expectedMonitoringRefsNumbers []int
		testMessage                   string
	}{
		{
			testNumber: 1,
			settings: map[string]string{
				"local_url":                          "http://example.com/test/siri",
				"remote_code_space":                  "test_kind",
				"generators.subscription_identifier": "Subscription::%{id}::LOC",
			},
			expectedNumberOfRequests:      1,
			expectedMonitoringRefsNumbers: []int{3},
			testMessage: `Without setting, should send
 3 StopMonitoringSubscriptionRequest in 1 request`,
		},
		{
			testNumber: 2,
			settings: map[string]string{
				"local_url":                          "http://example.com/test/siri",
				"remote_code_space":                  "test_kind",
				"generators.subscription_identifier": "Subscription::%{id}::LOC",
				"collect.siri.stop_monitoring.maximum_subscriptions_per_request": "1",
			},
			expectedNumberOfRequests:      3,
			expectedMonitoringRefsNumbers: []int{1, 1, 1},
			testMessage: `with setting
collect.siri.stop_monitoring.maximum_subscriptions_per_request = 3, should send 3 StopMonitoringSubscriptionRequest in 3 different Requests`,
		},
		{
			testNumber: 3,
			settings: map[string]string{
				"local_url":                          "http://example.com/test/siri",
				"remote_code_space":                  "test_kind",
				"generators.subscription_identifier": "Subscription::%{id}::LOC",
				"collect.siri.stop_monitoring.maximum_subscriptions_per_request": "2",
			},
			expectedNumberOfRequests:      2,
			expectedMonitoringRefsNumbers: []int{1, 2},
			testMessage: `with setting
collect.siri.stop_monitoring.maximum_subscriptions_per_request = 2, should send
2 Requests, one with 2 StopMonitoringSubscriptionRequest and the other with
1 StopMonitoringSubscriptionRequest`,
		},
		{
			testNumber: 4,
			settings: map[string]string{
				"local_url":                          "http://example.com/test/siri",
				"remote_code_space":                  "test_kind",
				"generators.subscription_identifier": "Subscription::%{id}::LOC",
				"collect.siri.stop_monitoring.maximum_subscriptions_per_request": "3",
			},
			expectedNumberOfRequests:      1,
			expectedMonitoringRefsNumbers: []int{3},
			testMessage: `with setting
collect.siri.stop_monitoring.maximum_subscriptions_per_request = 3, should send
1 Requests with 3 StopMonitoringSubscriptionRequest`,
		},
	}

	for _, test := range TestCases {
		request := &sxml.XMLSubscriptionRequest{}
		requests := []*sxml.XMLSubscriptionRequest{}
		// Create a test http server

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength <= 0 {
				t.Errorf("Request ContentLength should be zero")
			}
			body, _ := io.ReadAll(r.Body)
			request, _ = sxml.NewXMLSubscriptionRequestFromContent(body)
			requests = append(requests, request)
		}))
		defer ts.Close()

		partner := partners.New("slug")

		test.settings["remote_url"] = ts.URL
		partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, test.settings)

		partner.subscriptionManager = NewMemorySubscriptions(partner)
		partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
		partners.Save(partner)

		connector := NewSIRIStopMonitoringSubscriptionCollector(partner)

		stopAreaUpdateEvent := NewStopAreaUpdateRequest(stopArea.Id())
		stopAreaUpdateEvent2 := NewStopAreaUpdateRequest(stopArea2.Id())
		stopAreaUpdateEvent3 := NewStopAreaUpdateRequest(stopArea3.Id())
		connector.SetStopMonitoringSubscriber(NewFakeStopMonitoringSubscriber(connector))
		connector.RequestStopAreaUpdate(stopAreaUpdateEvent)
		connector.RequestStopAreaUpdate(stopAreaUpdateEvent2)
		connector.RequestStopAreaUpdate(stopAreaUpdateEvent3)

		connector.stopMonitoringSubscriber.Start()
		assert.Len(requests, test.expectedNumberOfRequests)

		var monitoringRefsLength []int
		for i := range requests {
			monitoringRefsLength = append(monitoringRefsLength, len(requests[i].XMLSubscriptionSMEntries()))
		}
		assert.ElementsMatch(monitoringRefsLength, test.expectedMonitoringRefsNumbers, fmt.Sprintf("test %d: %s", test.testNumber, test.testMessage))
		connector.stopMonitoringSubscriber.Stop()
	}
}

func Test_SIRIStopMonitoringDeleteSubscriptionRequest(t *testing.T) {
	request := &sxml.XMLDeleteSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		request, _ = sxml.NewXMLDeleteSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := partners.New("slug")

	settings := map[string]string{
		"local_url":                          "http://example.com/test/siri",
		"remote_url":                         ts.URL,
		"remote_code_space":                  "test_kind",
		"generators.subscription_identifier": "Subscription::%{id}::LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	file, _ := os.Open("testdata/notify-stop-monitoring.xml")
	content, _ := io.ReadAll(file)

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)
	connector.deletedSubscriptions = NewDeletedSubscriptions()

	notify, _ := sxml.NewXMLNotifyStopMonitoringFromContent(content)

	connector.HandleNotifyStopMonitoring(notify)

	if expected := "Subscription::6ba7b814-9dad-11d1-0-00c04fd430c8::LOC"; request.SubscriptionRef() != expected {
		t.Errorf("Wrong SubscriptionRef want : %v  got %v :", expected, request.SubscriptionRef())
	}
}
