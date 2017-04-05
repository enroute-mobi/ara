package siri

import (
	"testing"
	"time"
)

func Test_SIRIGeneralMessageRequest(t *testing.T) {
	expectedXML := `<ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <ns2:RequestTimestamp>2016-09-21T20:14:46.000Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>ref</ns2:RequestorRef>
        <ns2:MessageIdentifier>MessageId</ns2:MessageIdentifier>
      </ServiceRequestInfo>
      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>2016-09-21T20:14:46.000Z/ns2:RequestTimestamp>
        <ns2:MessageIdentifier>MessageId</ns2:MessageIdentifier>
      </Request>
      <RequestExtension/>
      </ns7:GetGeneralMessage>`

	requestTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	gmRequest := &SIRIGeneralMessageRequest{
		RequestTimestamp:  requestTimestamp,
		RequestorRef:      "ref",
		MessageIdentifier: "MessageId",
	}

	xml, err := gmRequest.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}

}
