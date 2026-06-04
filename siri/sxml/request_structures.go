package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type RequestXMLStructure struct {
	LightRequestXMLStructure

	requestorRef *string
}

type LightRequestXMLStructure struct {
	XMLStructure

	messageIdentifier *string
	requestTimestamp  *time.Time
}

func (request *RequestXMLStructure) RequestorRef() string {
	if request.requestorRef == nil {
		s := request.findStringChildContent(siri_attributes.RequestorRef)
		request.requestorRef = &s
	}
	return *request.requestorRef
}

func (request *LightRequestXMLStructure) MessageIdentifier() string {
	if request.messageIdentifier == nil {
		s := request.findStringChildContent(siri_attributes.MessageIdentifier)
		request.messageIdentifier = &s
	}
	return *request.messageIdentifier
}

func (request *LightRequestXMLStructure) RequestTimestamp() time.Time {
	if request.requestTimestamp == nil {
		t := request.findTimeChildContent(siri_attributes.RequestTimestamp)
		request.requestTimestamp = &t
	}
	return *request.requestTimestamp
}
