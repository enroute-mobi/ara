package siri

import (
	"testing"
	"time"
)

func Test_SIRISubscriptionTeminated_BUILDXML(t *testing.T) {
	st := SIRISubscriptionTerminated{
		ResponseMessageIdentifier: "un response message identifier",
		RequestMessageRef:         "un request message ref",
		ProducerRef:               "un producer ref",
		SubscriberRef:             "un subscriber ref",
		SubscriptionRef:           "une subscription ref",
		ResponseTimestamp:         time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}

	xml, err := st.BuildXML()
	if err != nil {
		t.Errorf(err.Error())
	}

	xmlSTR, err := NewXMLStopMonitoringSubscriptionTerminatedResponseFromContent([]byte(xml))
	if err != nil {
		t.Errorf(err.Error())
	}

	if xmlSTR.ProducerRef() != st.ProducerRef {
		t.Errorf("Wrong ProducerRef want : %v, got: %v", st.ProducerRef, xmlSTR.ProducerRef())
	}

	if xmlSTR.ResponseMessageIdentifier() != st.ResponseMessageIdentifier {
		t.Errorf("Wrong ResponseMessageIdentifier want : %v, got: %v", st.ProducerRef, xmlSTR.ProducerRef())
	}

	xmlSTs := xmlSTR.XMLSubscriptionTerminateds()

	if len(xmlSTs) != 1 {
		t.Errorf("Should have received one subscriptionTerminated but got : %v", len(xmlSTs))
	}

	xmlST := xmlSTs[0]
	if xmlST.SubscriberRef() != st.SubscriberRef {
		t.Errorf("Wrong SubscriberRef want : %v, got: %v", st.SubscriberRef, xmlST.ProducerRef())
	}

	if xmlST.SubscriptionRef() != st.SubscriptionRef {
		t.Errorf("Wrong SubscriptionReff want : %v, got: %v", st.SubscriptionRef, xmlST.ProducerRef())
	}
}
