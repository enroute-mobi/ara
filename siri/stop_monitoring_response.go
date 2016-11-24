package siri

import (
	"runtime"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringResponse struct {
	ResponseXMLStructure
}

func NewXMLStopMonitoringResponse(node xml.Node) *XMLStopMonitoringResponse {
	xmlStopMonitoringResponse := &XMLStopMonitoringResponse{}
	xmlStopMonitoringResponse.node = node
	return xmlStopMonitoringResponse
}

func NewXMLStopMonitoringResponseFromContent(content []byte) (*XMLStopMonitoringResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringResponse(doc.Root().XmlNode)
	finalizer := func(response *XMLStopMonitoringResponse) {
		doc.Free()
	}
	runtime.SetFinalizer(response, finalizer)
	return response, nil
}
