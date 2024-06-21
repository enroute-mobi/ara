package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
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
	sxEntries  []*XMLSituationExchangeSubscriptionRequestEntry
	ettEntries []*XMLEstimatedTimetableSubscriptionRequestEntry
	pttEntries []*XMLProductionTimetableSubscriptionRequestEntry
	vmEntries  []*XMLVehicleMonitoringSubscriptionRequestEntry
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
	nodes := request.findNodes(siri_attributes.StopMonitoringSubscriptionRequest)
	for _, stopMonitoring := range nodes {
		request.smEntries = append(request.smEntries, NewXMLStopMonitoringSubscriptionRequestEntry(stopMonitoring))
	}
	return request.smEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionETTEntries() []*XMLEstimatedTimetableSubscriptionRequestEntry {
	if len(request.ettEntries) != 0 {
		return request.ettEntries
	}
	nodes := request.findNodes(siri_attributes.EstimatedTimetableSubscriptionRequest)
	for _, ett := range nodes {
		request.ettEntries = append(request.ettEntries, NewXMLEstimatedTimetableSubscriptionRequestEntry(ett))
	}
	return request.ettEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionVMEntries() []*XMLVehicleMonitoringSubscriptionRequestEntry {
	if len(request.vmEntries) != 0 {
		return request.vmEntries
	}
	nodes := request.findNodes(siri_attributes.VehicleMonitoringSubscriptionRequest)
	for _, vm := range nodes {
		request.vmEntries = append(request.vmEntries, NewXMLVehicleMonitoringSubscriptionRequestEntry(vm))
	}
	return request.vmEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionPTTEntries() []*XMLProductionTimetableSubscriptionRequestEntry {
	if len(request.pttEntries) != 0 {
		return request.pttEntries
	}
	nodes := request.findNodes(siri_attributes.ProductionTimetableSubscriptionRequest)
	for _, ett := range nodes {
		request.pttEntries = append(request.pttEntries, NewXMLProductionTimetableSubscriptionRequestEntry(ett))
	}
	return request.pttEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionGMEntries() []*XMLGeneralMessageSubscriptionRequestEntry {
	if len(request.gmEntries) != 0 {
		return request.gmEntries
	}
	nodes := request.findNodes(siri_attributes.GeneralMessageSubscriptionRequest)
	for _, generalMessage := range nodes {
		request.gmEntries = append(request.gmEntries, NewXMLGeneralMessageSubscriptionRequestEntry(generalMessage))
	}
	return request.gmEntries
}

func (request *XMLSubscriptionRequest) XMLSubscriptionSXEntries() []*XMLSituationExchangeSubscriptionRequestEntry {
	if len(request.sxEntries) != 0 {
		return request.sxEntries
	}
	nodes := request.findNodes(siri_attributes.SituationExchangeSubscriptionRequest)
	for _, situationExchangeMessage := range nodes {
		request.sxEntries = append(request.sxEntries, NewXMLSituationExchangeSubscriptionRequestEntry(situationExchangeMessage))
	}
	return request.sxEntries
}

func (request *XMLSubscriptionRequest) ConsumerAddress() string {
	if request.consumerAddress == "" {
		request.consumerAddress = request.findStringChildContent(siri_attributes.ConsumerAddress)
	}
	return request.consumerAddress
}

func (request *XMLSubscriptionRequest) ChangeBeforeUpdates() string {
	if request.changeBeforeUpdates == "" {
		request.changeBeforeUpdates = request.findStringChildContent(siri_attributes.ChangeBeforeUpdates)
	}
	return request.changeBeforeUpdates
}

func (request *XMLSubscriptionRequest) IncrementalUpdates() string {
	if request.incrementalUpdates == "" {
		request.incrementalUpdates = request.findStringChildContent(siri_attributes.IncrementalUpdates)
	}
	return request.incrementalUpdates
}
