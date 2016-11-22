package siri

import (
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringResponse struct {
	XMLStructure
}

func NewXMLStopMonitoringResponse(node xml.Node) *XMLStopMonitoringResponse {
	return &XMLStopMonitoringResponse{XMLStructure: XMLStructure{node: node}}
}
