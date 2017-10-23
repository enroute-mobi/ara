package siri

import (
	"testing"
	"time"
)

func Test_DeleteSubscriptionResponse_BuildXML(t *testing.T) {
	expectedXML := `<sw:DeleteSubscriptionResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<DeleteSubscriptionAnswerInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ResponderRef>responder</siri:ResponderRef>
		<siri:RequestMessageRef>requestref</siri:RequestMessageRef>
	</DeleteSubscriptionAnswerInfo>
	<Answer>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ResponderRef>responder</siri:ResponderRef>
		<siri:RequestMessageRef>requestref</siri:RequestMessageRef>
		<siri:TerminationResponseStatus>
			<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
			<siri:SubscriberRef>subscriber</siri:SubscriberRef>
			<siri:SubscriptionRef>subscription</siri:SubscriptionRef>
			<siri:Status>true</siri:Status>
		</siri:TerminationResponseStatus>
	</Answer>
	<AnswerExtension/>
</sw:DeleteSubscriptionResponse>`

	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	rs := &SIRITerminationResponseStatus{
		SubscriberRef:     "subscriber",
		SubscriptionRef:   "subscription",
		ResponseTimestamp: responseTimestamp,
		Status:            true,
	}

	dsResponse := &SIRIDeleteSubscriptionResponse{
		ResponderRef:      "responder",
		RequestMessageRef: "requestref",
		ResponseTimestamp: responseTimestamp,
		ResponseStatus:    []*SIRITerminationResponseStatus{rs},
	}

	xml, err := dsResponse.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
