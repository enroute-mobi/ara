package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyVehicleMonitoring struct {
	ResponseXMLStructureWithStatus

	deliveries []*XMLNotifyVehicleMonitoringDelivery
}

type XMLNotifyVehicleMonitoringDelivery struct {
	SubscriptionDeliveryXMLStructure

	vehicleActivities []*XMLVehicleActivity
}

func NewXMLNotifyVehicleMonitoring(node xml.Node) *XMLNotifyVehicleMonitoring {
	xmlVehicleMonitoringResponse := &XMLNotifyVehicleMonitoring{}
	xmlVehicleMonitoringResponse.node = NewXMLNode(node)
	return xmlVehicleMonitoringResponse
}

func NewXMLNotifyVehicleMonitoringFromContent(content []byte) (*XMLNotifyVehicleMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifyVehicleMonitoring(doc.Root().XmlNode)
	return response, nil
}

func NewXMLNotifyVehicleMonitoringDelivery(node XMLNode) *XMLNotifyVehicleMonitoringDelivery {
	delivery := &XMLNotifyVehicleMonitoringDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyVehicleMonitoring) VehicleMonitoringDeliveries() []*XMLNotifyVehicleMonitoringDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLNotifyVehicleMonitoringDelivery{}
		nodes := notify.findNodes(siri_attributes.VehicleMonitoringDelivery)
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLNotifyVehicleMonitoringDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLNotifyVehicleMonitoringDelivery) VehicleActivities() []*XMLVehicleActivity {
	if delivery.vehicleActivities == nil {
		vas := []*XMLVehicleActivity{}
		nodes := delivery.findNodes(siri_attributes.VehicleActivity)
		for _, node := range nodes {
			vas = append(vas, NewXMLVehicleActivity(node))
		}
		delivery.vehicleActivities = vas
	}
	return delivery.vehicleActivities
}
