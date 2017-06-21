package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyStopMonitoring struct {
	XMLStructure

	deliveries []*XMLStopMonitoringDelivery
}

func (notify *XMLNotifyStopMonitoring) StopMonitoringDeliveries() []*XMLStopMonitoringDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLStopMonitoringDelivery{}
		nodes := notify.findNodes("StopMonitoringDelivery")
		if nodes != nil {
			for _, node := range nodes {
				deliveries = append(deliveries, NewXMLStopMonitoringDelivery(node))
			}
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
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
