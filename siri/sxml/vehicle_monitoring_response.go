package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
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

	vehicleActivityRawAttributes map[string]string

	itemIdentifier       *string
	linkDistance         *string
	percentage           *string
	vehicleMonitoringRef *string
	vehicleActivityNote  *string
	recordedAtTime       *time.Time
	validUntilTime       *time.Time
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
		nodes := response.findNodes(siri_attributes.VehicleMonitoringDelivery)
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
		nodes := delivery.findNodes(siri_attributes.VehicleActivity)
		for _, node := range nodes {
			activities = append(activities, NewXMLVehicleActivity(node))
		}
		delivery.vehicleActivities = activities
	}
	return delivery.vehicleActivities
}

func (va *XMLVehicleActivity) ItemIdentifier() string {
	if va.itemIdentifier == nil {
		s := va.findStringChildContent(siri_attributes.ItemIdentifier)
		va.itemIdentifier = &s
	}
	return *va.itemIdentifier
}

func (va *XMLVehicleActivity) LinkDistance() string {
	if va.linkDistance == nil {
		s := va.findStringChildContent(siri_attributes.LinkDistance)
		va.linkDistance = &s
	}
	return *va.linkDistance
}

func (va *XMLVehicleActivity) Percentage() string {
	if va.percentage == nil {
		s := va.findStringChildContent(siri_attributes.Percentage)
		va.percentage = &s
	}
	return *va.percentage
}

func (va *XMLVehicleActivity) VehicleMonitoringRef() string {
	if va.vehicleMonitoringRef == nil {
		s := va.findStringChildContent(siri_attributes.VehicleMonitoringRef)
		va.vehicleMonitoringRef = &s
	}
	return *va.vehicleMonitoringRef
}

func (va *XMLVehicleActivity) VehicleActivityNote() string {
	if va.vehicleActivityNote == nil {
		s := va.findStringChildContent(siri_attributes.VehicleActivityNote)
		va.vehicleActivityNote = &s
	}
	return *va.vehicleActivityNote
}

func (va *XMLVehicleActivity) RecordedAtTime() time.Time {
	if va.recordedAtTime == nil {
		t := va.findTimeChildContent(siri_attributes.RecordedAtTime)
		va.recordedAtTime = &t
	}
	return *va.recordedAtTime
}

func (va *XMLVehicleActivity) ValidUntilTime() time.Time {
	if va.validUntilTime == nil {
		t := va.findTimeChildContent(siri_attributes.RecordedAtTime)
		va.validUntilTime = &t
	}
	return *va.validUntilTime
}

func (va *XMLVehicleActivity) RawAttributes() map[string]string {
	if va.vehicleActivityRawAttributes != nil {
		return va.vehicleActivityRawAttributes
	}
	attrs := make(map[string]string)
	if v := va.VehicleActivityNote(); v != "" {
		attrs[siri_attributes.VehicleActivityNote] = v
	}
	va.vehicleActivityRawAttributes = attrs
	return attrs
}
