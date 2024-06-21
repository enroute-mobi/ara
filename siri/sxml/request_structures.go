package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type RequestXMLStructure struct {
	LightRequestXMLStructure

	requestorRef string
}

type LightRequestXMLStructure struct {
	XMLStructure

	messageIdentifier string
	requestTimestamp  time.Time
}

func (request *RequestXMLStructure) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent(siri_attributes.RequestorRef)
	}
	return request.requestorRef
}

func (request *LightRequestXMLStructure) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent(siri_attributes.MessageIdentifier)
	}
	return request.messageIdentifier
}

func (request *LightRequestXMLStructure) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent(siri_attributes.RequestTimestamp)
	}
	return request.requestTimestamp
}
