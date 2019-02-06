package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSubscriptionRequest struct {
	RequestXMLStructure

	consumerAddress string

	changeBeforeUpdates string
	incrementalUpdates  string

	smEntries  []*XMLStopMonitoringSubscriptionRequestEntry
	gmEntries  []*XMLGeneralMessageSubscriptionRequestEntry
	ettEntries []*XMLEstimatedTimetableSubscriptionRequestEntry
}

func NewXMLSubscriptionRequestFromContent(content []byte) (*XMLSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func NewXMLSubscriptionRequest(node xml.Node) *XMLSubscriptionRequest {
	xmlSubscriptionRequest := &XMLSubscriptionRequest{}
	xmlSubscriptionRequest.node = NewXMLNode(node)
	return xmlSubscriptionRequest
}

func (request *XMLSubscriptionRequest) XMLSubscriptionSMEntries() []*XMLStopMonitoringSubscriptionRequestEntry {
	if len(request.smEntries) != 0 {
		return request.smEntries
	}
	nodes := request.findNodes("StopMonitoringSubscriptionRequest")
	if nodes != nil {
		for _, stopMonitoring := range nodes {
			request.smEntries = append(request.smEntries, NewXMLStopMonitoringSubscriptionRequestEntry(stopMonitoring))
		}
	}
	return request.smEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionETTEntries() []*XMLEstimatedTimetableSubscriptionRequestEntry {
	if len(request.ettEntries) != 0 {
		return request.ettEntries
	}
	nodes := request.findNodes("EstimatedTimetableSubscriptionRequest")
	if nodes != nil {
		for _, ett := range nodes {
			request.ettEntries = append(request.ettEntries, NewXMLEstimatedTimetableSubscriptionRequestEntry(ett))
		}
	}
	return request.ettEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionGMEntries() []*XMLGeneralMessageSubscriptionRequestEntry {
	if len(request.gmEntries) != 0 {
		return request.gmEntries
	}
	nodes := request.findNodes("GeneralMessageSubscriptionRequest")
	if nodes != nil {
		for _, generalMessage := range nodes {
			request.gmEntries = append(request.gmEntries, NewXMLGeneralMessageSubscriptionRequestEntry(generalMessage))
		}
	}
	return request.gmEntries
}

func (request *XMLSubscriptionRequest) ConsumerAddress() string {
	if request.consumerAddress == "" {
		request.consumerAddress = request.findStringChildContent("ConsumerAddress")
	}
	return request.consumerAddress
}

func (request *XMLSubscriptionRequest) ChangeBeforeUpdates() string {
	if request.changeBeforeUpdates == "" {
		request.changeBeforeUpdates = request.findStringChildContent("ChangeBeforeUpdates")
	}
	return request.changeBeforeUpdates
}

func (request *XMLSubscriptionRequest) IncrementalUpdates() string {
	if request.incrementalUpdates == "" {
		request.incrementalUpdates = request.findStringChildContent("IncrementalUpdates")
	}
	return request.incrementalUpdates
}
