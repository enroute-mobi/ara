package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
)

func Test_SIRIGeneralMessageRequest(t *testing.T) {
	expectedXML := `<sw:GetGeneralMessage xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri" xmlns:sws="http://wsdl.siri.org.uk/siri">
	<ServiceRequestInfo>
		<siri:RequestTimestamp>2016-09-21T20:14:46.000Z</siri:RequestTimestamp>
		<siri:RequestorRef>ref</siri:RequestorRef>
		<siri:MessageIdentifier>MessageId</siri:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		<siri:RequestTimestamp>2016-09-21T20:14:46.000Z</siri:RequestTimestamp>
		<siri:MessageIdentifier>MessageId</siri:MessageIdentifier>
		<siri:Extensions>
			<sws:IDFGeneralMessageRequestFilter>
				<siri:LineRef>lineRef</siri:LineRef>
			</sws:IDFGeneralMessageRequestFilter>
		</siri:Extensions>
	</Request>
	<RequestExtension/>
</sw:GetGeneralMessage>`

	requestTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	gmRequest := &siri.SIRIGetGeneralMessageRequest{
		RequestorRef: "ref",
	}
	gmRequest.RequestTimestamp = requestTimestamp
	gmRequest.MessageIdentifier = "MessageId"
	gmRequest.LineRef = []string{"lineRef"}

	xml, err := gmRequest.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}

}
