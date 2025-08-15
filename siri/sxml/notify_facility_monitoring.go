package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyFacilityMonitoring struct {
	ResponseXMLStructureWithStatus

	deliveries []*XMLNotifyFacilityMonitoringDelivery
}

type XMLNotifyFacilityMonitoringDelivery struct {
	SubscriptionDeliveryXMLStructure

	facilityCondtions []*XMLFacilityCondition
}

func NewXMLNotifyFacilityMonitoring(node xml.Node) *XMLNotifyFacilityMonitoring {
	xmlFacilityMonitoringResponse := &XMLNotifyFacilityMonitoring{}
	xmlFacilityMonitoringResponse.node = NewXMLNode(node)
	return xmlFacilityMonitoringResponse
}

func NewXMLNotifyFacilityMonitoringFromContent(content []byte) (*XMLNotifyFacilityMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifyFacilityMonitoring(doc.Root().XmlNode)
	return response, nil
}

func NewXMLNotifyFacilityMonitoringDelivery(node XMLNode) *XMLNotifyFacilityMonitoringDelivery {
	delivery := &XMLNotifyFacilityMonitoringDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyFacilityMonitoring) FacilityMonitoringDeliveries() []*XMLNotifyFacilityMonitoringDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLNotifyFacilityMonitoringDelivery{}
		nodes := notify.findNodes(siri_attributes.FacilityMonitoringDelivery)
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLNotifyFacilityMonitoringDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLNotifyFacilityMonitoringDelivery) FacilityConditions() []*XMLFacilityCondition {
	if delivery.facilityCondtions == nil {
		facilityConditions := []*XMLFacilityCondition{}
		nodes := delivery.findNodes(siri_attributes.FacilityCondition)
		for _, node := range nodes {
			facilityConditions = append(facilityConditions, NewXMLFacilityCondition(node))
		}
		delivery.facilityCondtions = facilityConditions
	}
	return delivery.facilityCondtions
}
