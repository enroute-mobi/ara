package siri

import "time"

type XMLGeneralMessageSubscriptionResponse struct {
	XMLStructure

	address           string
	requestMessageRef string
	responderRef      string

	responseTimestamp  time.Time
	serviceStartedTime time.Time

	responseStatus XMLResponseStatus
}

func (response *XMLGeneralMessageSubscriptionResponse) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *XMLGeneralMessageSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent("ResponderRef")
	}
	return response.responderRef
}

func (response *XMLGeneralMessageSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLGeneralMessageSubscriptionResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent("ServiceStartedTime")
	}
	return response.serviceStartedTime
}

func (response *XMLGeneralMessageSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLGeneralMessageSubscriptionResponse) ResponseStatus() XMLResponseStatus {
	if response.responseStatus == (XMLResponseStatus{}) {
		node := response.findXMLNode("ResponseStatus")
		if node == nil {
			return response.responseStatus
		}
		xmlResponseStatus := XMLResponseStatus{}
		xmlResponseStatus.node = node
		response.responseStatus = xmlResponseStatus
	}
	return response.responseStatus
}
