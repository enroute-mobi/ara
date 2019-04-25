package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyStopMonitoring struct {
	ResponseXMLStructure

	deliveries []*XMLNotifyStopMonitoringDelivery
}

type XMLNotifyStopMonitoringDelivery struct {
	SubscriptionDeliveryXMLStructure

	monitoringRef string

	monitoredStopVisits             []*XMLMonitoredStopVisit
	monitoredStopVisitCancellations []*XMLMonitoredStopVisitCancellation
}

func NewXMLNotifyStopMonitoringDelivery(node XMLNode) *XMLNotifyStopMonitoringDelivery {
	delivery := &XMLNotifyStopMonitoringDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyStopMonitoring) StopMonitoringDeliveries() []*XMLNotifyStopMonitoringDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLNotifyStopMonitoringDelivery{}
		nodes := notify.findNodes("StopMonitoringDelivery")
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLNotifyStopMonitoringDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLNotifyStopMonitoringDelivery) MonitoringRef() string {
	if delivery.monitoringRef == "" {
		delivery.monitoringRef = delivery.findStringChildContent("MonitoringRef")
	}
	return delivery.monitoringRef
}

func (delivery *XMLNotifyStopMonitoringDelivery) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if delivery.monitoredStopVisits == nil {
		stopVisits := []*XMLMonitoredStopVisit{}
		nodes := delivery.findNodes("MonitoredStopVisit")
		for _, node := range nodes {
			stopVisits = append(stopVisits, NewXMLMonitoredStopVisit(node))
		}
		delivery.monitoredStopVisits = stopVisits
	}
	return delivery.monitoredStopVisits
}

func (delivery *XMLNotifyStopMonitoringDelivery) XMLMonitoredStopVisitCancellations() []*XMLMonitoredStopVisitCancellation {
	if delivery.monitoredStopVisitCancellations == nil {
		cancellations := []*XMLMonitoredStopVisitCancellation{}
		nodes := delivery.findNodes("MonitoredStopVisitCancellation")
		for _, node := range nodes {
			cancellations = append(cancellations, NewXMLCancelledStopVisit(node))
		}
		delivery.monitoredStopVisitCancellations = cancellations
	}
	return delivery.monitoredStopVisitCancellations
}

func NewXMLNotifyStopMonitoring(node xml.Node) *XMLNotifyStopMonitoring {
	xmlStopMonitoringResponse := &XMLNotifyStopMonitoring{}
	xmlStopMonitoringResponse.node = NewXMLNode(node)
	return xmlStopMonitoringResponse
}

func NewXMLNotifyStopMonitoringFromContent(content []byte) (*XMLNotifyStopMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifyStopMonitoring(doc.Root().XmlNode)
	return response, nil
}
