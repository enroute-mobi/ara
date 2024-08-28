package core

import (
	"bytes"
	"testing"

	"bitbucket.org/enroute-mobi/ara/audit"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func Test_GetSubscriptionRequest(t *testing.T) {
	assert := assert.New(t)

	partners := createTestPartnerManager()
	partner := partners.New("slug")

	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRISubscriptionRequestDispatcher(partner)
	partner.connectors["test-startable-connector-connector"] = connector

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	c, _ := partner.Connector("test-startable-connector-connector")
	subscriber := NewCollectSubcriber(c, "kind")
	assert.Len(subscriber.GetSubscriptionRequest(), 0, `No subscriptionRequest
without Subscription`)

	// Create a Subscription
	subscription := partner.Subscriptions().FindOrCreateByKind("kind")
	subscription.Save()

	subscriptionRequest := subscriber.GetSubscriptionRequest()
	assert.Empty(subscriptionRequest, "No subscriptionRequest with Subscription without Resource")

	// Create and add Resource to Subscription
	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
		Type: "type",
	}
	subscription.CreateAndAddNewResource(reference)
	subscriptionRequest = subscriber.GetSubscriptionRequest()

	assert.Len(subscriptionRequest, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequest)[0])

	modelsToRequest := subscriptionRequest[subscription.Id()].modelsToRequest
	assert.Len(modelsToRequest, 1)
	assert.Equal(obj, modelsToRequest[0].code)
	assert.Equal("type", modelsToRequest[0].kind)

	// Add another Resource to the Subscription
	obj2 := model.NewCode("internal", "AnotherValue")
	reference2 := model.Reference{
		Code: &obj2,
		Type: "type2",
	}
	subscription.CreateAndAddNewResource(reference2)
	subscriptionRequest = subscriber.GetSubscriptionRequest()

	assert.Len(subscriptionRequest, 1, "1 subscriptionRequest with Subscription with 2 Resources")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequest)[0])

	modelsToRequest = subscriptionRequest[subscription.Id()].modelsToRequest
	assert.Len(modelsToRequest, 2, "2 Resources for the Subscription")
	var codes []model.Code
	var kinds []string
	for i := range modelsToRequest {
		codes = append(codes, modelsToRequest[i].code)
		kinds = append(kinds, modelsToRequest[i].kind)
	}
	assert.ElementsMatch(codes, []model.Code{obj, obj2})
	assert.ElementsMatch(kinds, []string{"type", "type2"})

	// Force Subscription 1st Resource RetryCount above 10
	resource := subscription.Resource(obj)
	resource.RetryCount += 11

	subscriptionRequest = subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequest, 1, `1 subscriptionRequest with Subscription without 1 Resource
having 1 RetryCount > 10`)

	// Force Subscription 2nd Resource RetryCount above 10
	resource = subscription.Resource(obj2)
	resource.RetryCount += 11

	subscriptionRequest = subscriber.GetSubscriptionRequest()
	assert.Empty(subscriptionRequest)
}

func Test_HandleResponse_BadResponse(t *testing.T) {
	assert := assert.New(t)

	subscriber, subscription, responseTemplate := testSetup()

	// Create and add Resource to Subscription
	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}

	subscription.CreateAndAddNewResource(reference)

	subscriptionRequests := subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequests, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequests)[0])

	requestMessageRef := maps.Values(subscriptionRequests)[0].requestMessageRef

	var TestCases = []struct {
		requestMessageRef     string
		subscriptionIdentifer string
		status                string
		message               string
	}{
		{
			"WRONG",
			"WRONG",
			"true",
			`when the requestMessageRef is unknown
		if there is only one request the resource should not be subscribed and the subscriptionRequests must not be emtpy`,
		},
		{
			requestMessageRef,
			"WRONG",
			"true",
			`when the requestMessageRef is OK and subscriptionRef is unknown
		the resource should not be subscribed and the subscriptionRequests must not be emtpy`,
		},
	}

	for _, test := range TestCases {
		response := bytes.Replace(responseTemplate, []byte("{request_message_ref}"), []byte(test.requestMessageRef), 1)
		response = bytes.Replace(response, []byte("{subscription_ref}"), []byte(test.subscriptionIdentifer), 1)
		response = bytes.Replace(response, []byte("{status}"), []byte(test.status), 1)
		resp, err := sxml.NewXMLSubscriptionResponseFromContent(response)
		assert.NoError(err, "cannot parse xml")

		message := &audit.BigQueryMessage{}
		subRequests := subscriptionRequests
		subscriber.HandleResponse(subRequests, message, resp)
		resource := subscription.Resource(obj)

		assert.Zero(resource.RetryCount)
		assert.Zero(resource.SubscribedAt())
		assert.NotEmpty(subRequests, test.message)
	}
}

func Test_HandleResponse_GoodResponse(t *testing.T) {
	assert := assert.New(t)

	subscriber, subscription, responseTemplate := testSetup()

	// Create and add Resource to Subscription
	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}

	subscription.CreateAndAddNewResource(reference)

	subscriptionRequests := subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequests, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequests)[0])

	requestMessageRef := maps.Values(subscriptionRequests)[0].requestMessageRef

	response := bytes.Replace(responseTemplate, []byte("{request_message_ref}"), []byte(requestMessageRef), 1)
	response = bytes.Replace(response, []byte("{subscription_ref}"), []byte("6ba7b814-9dad-11d1-0-00c04fd430c8"), 1)
	response = bytes.Replace(response, []byte("{status}"), []byte("true"), 1)
	resp, err := sxml.NewXMLSubscriptionResponseFromContent(response)
	assert.NoError(err, "cannot parse xml")

	message := &audit.BigQueryMessage{}
	subscriber.HandleResponse(subscriptionRequests, message, resp)
	resource := subscription.Resource(obj)

	testMessage := `when the requestMessageRef is OK and subscriptionRef is OK
and Status is true, the resource should be subscribed and the subscriptionRequests must be emtpy`

	assert.Zero(resource.RetryCount)
	assert.NotZero(resource.SubscribedAt(), testMessage)
	assert.Empty(subscriptionRequests, testMessage)
}

func Test_HandleResponse_GoodResponse_With_GeneralMessageCollect_all(t *testing.T) {
	assert := assert.New(t)

	subscriber, subscription, responseTemplate := testSetup()

	// Create and add Resource to Subscription
	obj := model.NewCode("GeneralMessageCollect", "all")
	reference := model.Reference{
		Code: &obj,
	}

	subscription.CreateAndAddNewResource(reference)

	subscriptionRequests := subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequests, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequests)[0])

	requestMessageRef := maps.Values(subscriptionRequests)[0].requestMessageRef

	response := bytes.Replace(responseTemplate, []byte("{request_message_ref}"), []byte(requestMessageRef), 1)
	response = bytes.Replace(response, []byte("{subscription_ref}"), []byte("6ba7b814-9dad-11d1-0-00c04fd430c8"), 1)
	response = bytes.Replace(response, []byte("{status}"), []byte("true"), 1)
	resp, err := sxml.NewXMLSubscriptionResponseFromContent(response)
	assert.NoError(err, "cannot parse xml")

	message := &audit.BigQueryMessage{}
	subscriber.HandleResponse(subscriptionRequests, message, resp)
	resource := subscription.Resource(obj)

	testMessage := `when the requestMessageRef is OK and subscriptionRef is OK
and Status is true, the resource should be subscribed and the subscriptionRequests must be emtpy`

	assert.Zero(resource.RetryCount)
	assert.NotZero(resource.SubscribedAt(), testMessage)
	assert.Empty(subscriptionRequests, testMessage)
}

func Test_HandleResponse_GoodResponse_With_SituationExchangeCollect_all(t *testing.T) {
	assert := assert.New(t)

	subscriber, subscription, responseTemplate := testSetup()

	// Create and add Resource to Subscription
	obj := model.NewCode("SituationExchangeCollect", "all")
	reference := model.Reference{
		Code: &obj,
	}

	subscription.CreateAndAddNewResource(reference)

	subscriptionRequests := subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequests, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequests)[0])

	requestMessageRef := maps.Values(subscriptionRequests)[0].requestMessageRef

	response := bytes.Replace(responseTemplate, []byte("{request_message_ref}"), []byte(requestMessageRef), 1)
	response = bytes.Replace(response, []byte("{subscription_ref}"), []byte("6ba7b814-9dad-11d1-0-00c04fd430c8"), 1)
	response = bytes.Replace(response, []byte("{status}"), []byte("true"), 1)
	resp, err := sxml.NewXMLSubscriptionResponseFromContent(response)
	assert.NoError(err, "cannot parse xml")

	message := &audit.BigQueryMessage{}
	subscriber.HandleResponse(subscriptionRequests, message, resp)
	resource := subscription.Resource(obj)

	testMessage := `when the requestMessageRef is OK and subscriptionRef is OK
and Status is true, the resource should be subscribed and the subscriptionRequests must be emtpy`

	assert.Zero(resource.RetryCount)
	assert.NotZero(resource.SubscribedAt(), testMessage)
	assert.Empty(subscriptionRequests, testMessage)
}

func Test_HandleResponse_GoodResponse_With_MultipleResources(t *testing.T) {
	assert := assert.New(t)

	subscriber, subscription, responseTemplate := testSetup()

	// Create and add Resources to Subscription
	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}
	subscription.CreateAndAddNewResource(reference)

	obj2 := model.NewCode("internal", "Value2")
	reference = model.Reference{
		Code: &obj2,
	}
	subscription.CreateAndAddNewResource(reference)

	subscriptionRequests := subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequests, 1, "1 subscriptionRequest with Subscription with Resources")
	assert.Len(subscriptionRequests[SubscriptionId("6ba7b814-9dad-11d1-0-00c04fd430c8")].modelsToRequest, 2)
	assert.Equal(SubscriptionId("6ba7b814-9dad-11d1-0-00c04fd430c8"), maps.Keys(subscriptionRequests)[0])

	requestMessageRef := maps.Values(subscriptionRequests)[0].requestMessageRef

	response := bytes.Replace(responseTemplate, []byte("{request_message_ref}"), []byte(requestMessageRef), 1)
	response = bytes.Replace(response, []byte("{subscription_ref}"), []byte("6ba7b814-9dad-11d1-0-00c04fd430c8"), 1)
	response = bytes.Replace(response, []byte("{status}"), []byte("true"), 1)
	resp, err := sxml.NewXMLSubscriptionResponseFromContent(response)
	assert.NoError(err, "cannot parse xml")

	message := &audit.BigQueryMessage{}
	subscriber.HandleResponse(subscriptionRequests, message, resp)

	testMessage := `when the requestMessageRef is OK and subscriptionRef is OK
and Status is true, all Resources of the subscription should be subscribed
and the subscriptionRequests must be emtpy`

	// Testing 1st Subscription Resource
	resource := subscription.Resource(obj)
	assert.Zero(resource.RetryCount, testMessage)
	assert.NotZero(resource.SubscribedAt(), testMessage)

	// Testing 2nd Subscription Resource
	resource = subscription.Resource(obj2)
	assert.Zero(resource.RetryCount, testMessage)
	assert.NotZero(resource.SubscribedAt(), testMessage)

	assert.Empty(subscriptionRequests, testMessage)
}

func Test_HandleResponse_StatusFalse(t *testing.T) {
	assert := assert.New(t)

	subscriber, subscription, responseTemplate := testSetup()

	// Create and add Resource to Subscription
	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}

	subscription.CreateAndAddNewResource(reference)

	subscriptionRequests := subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequests, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequests)[0])

	requestMessageRef := maps.Values(subscriptionRequests)[0].requestMessageRef

	response := bytes.Replace(responseTemplate, []byte("{request_message_ref}"), []byte(requestMessageRef), 1)
	response = bytes.Replace(response, []byte("{subscription_ref}"), []byte("6ba7b814-9dad-11d1-0-00c04fd430c8"), 1)
	response = bytes.Replace(response, []byte("{status}"), []byte("false"), 1)
	resp, err := sxml.NewXMLSubscriptionResponseFromContent(response)
	assert.NoError(err, "cannot parse xml")

	message := &audit.BigQueryMessage{}
	subscriber.HandleResponse(subscriptionRequests, message, resp)
	resource := subscription.Resource(obj)

	testMessage := `when the requestMessageRef is OK and subscriptionRef is OK
and Status is false, the resource should not subscribed and the subscriptionRequests must be emtpy`

		assert.Equal(1, resource.RetryCount)
	assert.Zero(resource.SubscribedAt(), testMessage)
	assert.Empty(subscriptionRequests, testMessage)
}

func testSetup() (subscriber *CollectSubscriber, subscription *Subscription, response []byte) {
	response = []byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
    <soap:Body>
        <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <SubscriptionAnswerInfo xmlns:ns2="http://www.ifopt.org.uk/acsb"
									xmlns:ns3="http://www.ifopt.org.uk/ifopt"
									xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
									xmlns:ns5="http://www.siri.org.uk/siri"
									xmlns:ns6="http://wsdl.siri.org.uk/siri">
                <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
                <ns5:Address>http://sqybus-siri:8080/ProfilSiriKidf2_4Producer-Sqybus/SiriServices</ns5:Address>
                <ns5:ResponderRef>SQYBUS</ns5:ResponderRef>
                <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">28679112-9dad-11d1-2-00c04fd430c8</ns5:RequestMessageRef>
            </SubscriptionAnswerInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb"
	 				xmlns:ns3="http://www.ifopt.org.uk/ifopt"
					xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
					xmlns:ns5="http://www.siri.org.uk/siri"
					xmlns:ns6="http://wsdl.siri.org.uk/siri">
                <ns5:ResponseStatus>
                    <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
                    <ns5:RequestMessageRef>{request_message_ref}</ns5:RequestMessageRef>
                    <ns5:SubscriberRef>RATPDEV:Concerto</ns5:SubscriberRef>
                    <ns5:SubscriptionRef>{subscription_ref}</ns5:SubscriptionRef>
                    <ns5:Status>{status}</ns5:Status>
                    <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
                </ns5:ResponseStatus>
                <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
        </ns1:SubscribeResponse>
    </soap:Body>
</soap:Envelope>`)

	partners := createTestPartnerManager()
	partner := partners.New("slug")

	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRISubscriptionRequestDispatcher(partner)
	connector.remoteCodeSpace = "internal"
	partner.connectors["test-startable-connector-connector"] = connector

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	c, _ := partner.Connector("test-startable-connector-connector")

	subscriber = NewCollectSubcriber(c, "kind")

	// Create a Subscription
	subscription = partner.Subscriptions().FindOrCreateByKind("kind")
	subscription.Save()

	return subscriber, subscription, response
}
