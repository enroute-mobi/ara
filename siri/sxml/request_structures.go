package sxml

import (
	"time"
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
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *LightRequestXMLStructure) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *LightRequestXMLStructure) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}
