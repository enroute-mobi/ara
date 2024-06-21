package sxml

import (
	"strings"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetVehicleMonitoring struct {
	XMLVehicleMonitoringRequest

	requestorRef string
}

type XMLVehicleMonitoringRequest struct {
	LightRequestXMLStructure

	lineRef           string
	vehicleRef        string
	messageIdentifier string

	lines []string
}

func NewXMLGetVehicleMonitoring(node xml.Node) *XMLGetVehicleMonitoring {
	xmlGetVehicleMonitoring := &XMLGetVehicleMonitoring{}
	xmlGetVehicleMonitoring.node = NewXMLNode(node)
	return xmlGetVehicleMonitoring
}

func NewXMLGetVehicleMonitoringFromContent(content []byte) (*XMLGetVehicleMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetVehicleMonitoring(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLVehicleMonitoringRequest) Lines() []string {
	if len(request.lines) == 0 {
		nodes := request.findNodes(siri_attributes.LineRef)
		for _, node := range nodes {
			request.lines = append(request.lines, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return request.lines
}

func (request *XMLGetVehicleMonitoring) LineRef() string {
	if request.lineRef == "" {
		request.lineRef = request.findStringChildContent(siri_attributes.LineRef)
	}
	return request.lineRef
}

func (request *XMLGetVehicleMonitoring) VehicleRef() string {
	if request.vehicleRef == "" {
		request.vehicleRef = request.findStringChildContent(siri_attributes.VehicleRef)
	}
	return request.vehicleRef
}

func (request *XMLGetVehicleMonitoring) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent(siri_attributes.MessageIdentifier)
	}
	return request.messageIdentifier
}

func (request *XMLGetVehicleMonitoring) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent(siri_attributes.RequestorRef)
	}
	return request.requestorRef
}
