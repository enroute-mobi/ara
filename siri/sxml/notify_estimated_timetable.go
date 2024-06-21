package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyEstimatedTimetable struct {
	ResponseXMLStructure

	deliveries []*XMLNotifyEstimatedTimetableDelivery
}

type XMLNotifyEstimatedTimetableDelivery struct {
	SubscriptionDeliveryXMLStructure

	estimatedJourneyVersionFrames []*XMLEstimatedJourneyVersionFrame
}

func NewXMLNotifyEstimatedTimetableDelivery(node XMLNode) *XMLNotifyEstimatedTimetableDelivery {
	delivery := &XMLNotifyEstimatedTimetableDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyEstimatedTimetable) EstimatedTimetableDeliveries() []*XMLNotifyEstimatedTimetableDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLNotifyEstimatedTimetableDelivery{}
		nodes := notify.findNodes(siri_attributes.EstimatedTimetableDelivery)
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLNotifyEstimatedTimetableDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
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

func (delivery *XMLNotifyEstimatedTimetableDelivery) EstimatedJourneyVersionFrames() []*XMLEstimatedJourneyVersionFrame {
	if delivery.estimatedJourneyVersionFrames == nil {
		estimatedJourneyVersionFrames := []*XMLEstimatedJourneyVersionFrame{}
		nodes := delivery.findNodes(siri_attributes.EstimatedJourneyVersionFrame)
		for _, node := range nodes {
			estimatedJourneyVersionFrames = append(estimatedJourneyVersionFrames, NewXMLEstimatedJourneyVersionFrame(node))
		}
		delivery.estimatedJourneyVersionFrames = estimatedJourneyVersionFrames
	}
	return delivery.estimatedJourneyVersionFrames
}
