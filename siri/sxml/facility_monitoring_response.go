package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLFacilityMonitoringResponse struct {
	ResponseXMLStructure

	deliveries []*XMLFacilityMonitoringDelivery
}

type XMLFacilityMonitoringDelivery struct {
	DeliveryXMLStructure

	facilityConditions []*XMLFacilityCondition
}

type XMLFacilityCondition struct {
	XMLStructure

	facilityRef    string
	facilityStatus string
}

func NewXMLFacilityMonitoringResponse(node xml.Node) *XMLFacilityMonitoringResponse {
	xmlFacilityMonitoringResponse := &XMLFacilityMonitoringResponse{}
	xmlFacilityMonitoringResponse.node = NewXMLNode(node)
	return xmlFacilityMonitoringResponse
}

func NewXMLFacilityMonitoringResponseFromContent(content []byte) (*XMLFacilityMonitoringResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLFacilityMonitoringResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLFacilityMonitoringDelivery(node XMLNode) *XMLFacilityMonitoringDelivery {
	delivery := &XMLFacilityMonitoringDelivery{}
	delivery.node = node
	return delivery
}

func NewXMLFacilityCondition(node XMLNode) *XMLFacilityCondition {
	facilityCondition := &XMLFacilityCondition{}
	facilityCondition.node = node
	return facilityCondition
}

func (response *XMLFacilityMonitoringResponse) FacilityMonitoringDeliveries() []*XMLFacilityMonitoringDelivery {
	if response.deliveries == nil {
		deliveries := []*XMLFacilityMonitoringDelivery{}
		nodes := response.findNodes(siri_attributes.FacilityMonitoringDelivery)
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLFacilityMonitoringDelivery(node))
		}
		response.deliveries = deliveries
	}
	return response.deliveries
}

func (delivery *XMLFacilityMonitoringDelivery) FacilityConditions() []*XMLFacilityCondition {
	if delivery.facilityConditions == nil {
		facilityConditions := []*XMLFacilityCondition{}
		nodes := delivery.findNodes(siri_attributes.FacilityCondition)
		for _, node := range nodes {
			facilityConditions = append(facilityConditions, NewXMLFacilityCondition(node))
		}
		delivery.facilityConditions = facilityConditions
	}
	return delivery.facilityConditions
}

func (delivery *XMLFacilityCondition) FacilityRef() string {
	if delivery.facilityRef == "" {
		delivery.facilityRef = delivery.findStringChildContent(siri_attributes.FacilityRef)
	}
	return delivery.facilityRef
}

func (delivery *XMLFacilityCondition) FacilityStatus() string {
	if delivery.facilityStatus == "" {
		delivery.facilityStatus = delivery.findStringChildContent(siri_attributes.FacilityStatus)
	}
	return delivery.facilityStatus
}
