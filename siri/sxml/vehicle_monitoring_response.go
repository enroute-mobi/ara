package sxml

import (
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLVehicleMonitoringResponse struct {
	ResponseXMLStructure

	deliveries []*XMLVehicleMonitoringDelivery
}

type XMLVehicleMonitoringDelivery struct {
	DeliveryXMLStructure

	vehicleActivities []*XMLVehicleActivity
}

type XMLVehicleActivity struct {
	XMLMonitoredVehicleJourney

	itemIdentifier       string
	linkDistance         string
	percentage           string
	vehicleMonitoringRef string
	recordedAtTime       time.Time
	validUntilTime       time.Time
}

func NewXMLVehicleMonitoringResponse(node xml.Node) *XMLVehicleMonitoringResponse {
	xmlVehicleMonitoringResponse := &XMLVehicleMonitoringResponse{}
	xmlVehicleMonitoringResponse.node = NewXMLNode(node)
	return xmlVehicleMonitoringResponse
}

func NewXMLVehicleMonitoringResponseFromContent(content []byte) (*XMLVehicleMonitoringResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLVehicleMonitoringResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLVehicleMonitoringDelivery(node XMLNode) *XMLVehicleMonitoringDelivery {
	delivery := &XMLVehicleMonitoringDelivery{}
	delivery.node = node
	return delivery
}

func NewXMLVehicleActivity(node XMLNode) *XMLVehicleActivity {
	activity := &XMLVehicleActivity{}
	activity.node = node
	return activity
}

func (response *XMLVehicleMonitoringResponse) VehicleMonitoringDeliveries() []*XMLVehicleMonitoringDelivery {
	if response.deliveries == nil {
		deliveries := []*XMLVehicleMonitoringDelivery{}
		nodes := response.findNodes("VehicleMonitoringDelivery")
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLVehicleMonitoringDelivery(node))
		}
		response.deliveries = deliveries
	}
	return response.deliveries
}

func (delivery *XMLVehicleMonitoringDelivery) VehicleActivities() []*XMLVehicleActivity {
	if delivery.vehicleActivities == nil {
		activities := []*XMLVehicleActivity{}
		nodes := delivery.findNodes("VehicleActivity")
		for _, node := range nodes {
			activities = append(activities, NewXMLVehicleActivity(node))
		}
		delivery.vehicleActivities = activities
	}
	return delivery.vehicleActivities
}

func (va *XMLVehicleActivity) ItemIdentifier() string {
	if va.itemIdentifier == "" {
		va.itemIdentifier = va.findStringChildContent("ItemIdentifier")
	}
	return va.itemIdentifier
}

func (va *XMLVehicleActivity) LinkDistance() string {
	if va.linkDistance == "" {
		va.linkDistance = va.findStringChildContent("LinkDistance")
	}
	return va.linkDistance
}

func (va *XMLVehicleActivity) Percentage() string {
	if va.percentage == "" {
		va.percentage = va.findStringChildContent("Percentage")
	}
	return va.percentage
}

func (va *XMLVehicleActivity) VehicleMonitoringRef() string {
	if va.vehicleMonitoringRef == "" {
		va.vehicleMonitoringRef = va.findStringChildContent("VehicleMonitoringRef")
	}
	return va.vehicleMonitoringRef
}

func (va *XMLVehicleActivity) RecordedAtTime() time.Time {
	if va.recordedAtTime.IsZero() {
		va.recordedAtTime = va.findTimeChildContent("RecordedAtTime")
	}
	return va.recordedAtTime
}

func (va *XMLVehicleActivity) ValidUntilTime() time.Time {
	if va.validUntilTime.IsZero() {
		va.validUntilTime = va.findTimeChildContent("RecordedAtTime")
	}
	return va.validUntilTime
}
