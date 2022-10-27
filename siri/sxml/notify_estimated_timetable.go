package sxml

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyEstimatedTimetable struct {
	ResponseXMLStructure

	deliveries []*XMLNotifyEstimatedTimetableDelivery
}

type XMLNotifyEstimatedTimetableDelivery struct {
	SubscriptionDeliveryXMLStructure

	lineRef string

	// monitoredStopVisits             []*XMLMonitoredStopVisit
	// monitoredStopVisitCancellations []*XMLMonitoredStopVisitCancellation
}

func NewXMLNotifyEstimatedTimetableDelivery(node XMLNode) *XMLNotifyEstimatedTimetableDelivery {
	delivery := &XMLNotifyEstimatedTimetableDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyEstimatedTimetable) EstimatedTimetableDeliveries() []*XMLNotifyEstimatedTimetableDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLNotifyEstimatedTimetableDelivery{}
		nodes := notify.findNodes("EstimatedTimetableDelivery")
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLNotifyEstimatedTimetableDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLNotifyEstimatedTimetableDelivery) LineRef() string {
	if delivery.lineRef == "" {
		delivery.lineRef = delivery.findStringChildContent("LineRef")
	}
	return delivery.lineRef
}

func NewXMLNotifyEstimatedTimetable(node xml.Node) *XMLNotifyEstimatedTimetable {
	xmlEstimatedTimetableResponse := &XMLNotifyEstimatedTimetable{}
	xmlEstimatedTimetableResponse.node = NewXMLNode(node)
	return xmlEstimatedTimetableResponse
}

func NewXMLNotifyEstimatedTimetableFromContent(content []byte) (*XMLNotifyEstimatedTimetable, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifyEstimatedTimetable(doc.Root().XmlNode)
	return response, nil
}
