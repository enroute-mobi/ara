package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
)

func Test_DeleteSubscriptionRequest_BuildXML(t *testing.T) {
	expectedXML := `<sw:DeleteSubscription xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<DeleteSubscriptionInfo>
		<siri:RequestTimestamp>2016-09-21T20:14:46.000Z</siri:RequestTimestamp>
		<siri:RequestorRef>requestor</siri:RequestorRef>
		<siri:MessageIdentifier>mid</siri:MessageIdentifier>
	</DeleteSubscriptionInfo>
	<Request>
		<siri:All/>
	</Request>
	<RequestExtension/>
</sw:DeleteSubscription>`

	requestTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	dsRequest := &siri.SIRIDeleteSubscriptionRequest{
		RequestorRef:      "requestor",
		RequestTimestamp:  requestTimestamp,
		MessageIdentifier: "mid",
		CancelAll:         true,
	}

	xml, err := dsRequest.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
