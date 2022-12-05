package sxml

import (
	"fmt"
	"strings"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopPointsDiscoveryResponse struct {
	LightDeliveryXMLStructure

	annotatedStopPointRefs []*XMLAnnotatedStopPointRef
}

type XMLAnnotatedStopPointRef struct {
	XMLStructure

	stopPointRef string
	stopName     string

	lineRefs []string

	monitored Bool
}

func NewXMLStopPointsDiscoveryResponse(node xml.Node) *XMLStopPointsDiscoveryResponse {
	xmlStopPointsDiscoveryResponse := &XMLStopPointsDiscoveryResponse{}
	xmlStopPointsDiscoveryResponse.node = NewXMLNode(node)
	return xmlStopPointsDiscoveryResponse
}

func NewXMLStopPointsDiscoveryResponseFromContent(content []byte) (*XMLStopPointsDiscoveryResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopPointsDiscoveryResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLAnnotatedStopPointRef(node XMLNode) *XMLAnnotatedStopPointRef {
	annotatedStopPoint := &XMLAnnotatedStopPointRef{}
	annotatedStopPoint.node = node
	return annotatedStopPoint
}

func (response *XMLStopPointsDiscoveryResponse) ErrorString() string {
	return fmt.Sprintf("%v: %v", response.errorType(), response.ErrorText())
}

func (response *XMLStopPointsDiscoveryResponse) errorType() string {
	if response.ErrorType() == "OtherError" {
		return fmt.Sprintf("%v %v", response.ErrorType(), response.ErrorNumber())
	}
	return response.ErrorType()
}

func (response *XMLStopPointsDiscoveryResponse) AnnotatedStopPointRefs() []*XMLAnnotatedStopPointRef {
	if response.annotatedStopPointRefs == nil {
		annotatedStopPointRefs := []*XMLAnnotatedStopPointRef{}
		nodes := response.findNodes("AnnotatedStopPointRef")
		for _, node := range nodes {
			annotatedStopPointRefs = append(annotatedStopPointRefs, NewXMLAnnotatedStopPointRef(node))
		}
		response.annotatedStopPointRefs = annotatedStopPointRefs
	}
	return response.annotatedStopPointRefs
}

func (annotatedStopPoint *XMLAnnotatedStopPointRef) StopPointRef() string {
	if annotatedStopPoint.stopPointRef == "" {
		annotatedStopPoint.stopPointRef = annotatedStopPoint.findStringChildContent("StopPointRef")
	}
	return annotatedStopPoint.stopPointRef
}

func (annotatedStopPoint *XMLAnnotatedStopPointRef) StopName() string {
	if annotatedStopPoint.stopName == "" {
		annotatedStopPoint.stopName = annotatedStopPoint.findStringChildContent("StopName")
	}
	return annotatedStopPoint.stopName
}

func (annotatedStopPoint *XMLAnnotatedStopPointRef) LineRefs() []string {
	if len(annotatedStopPoint.lineRefs) == 0 {
		nodes := annotatedStopPoint.findNodes("LineRef")
		for _, node := range nodes {
			annotatedStopPoint.lineRefs = append(annotatedStopPoint.lineRefs, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return annotatedStopPoint.lineRefs
}

func (annotatedStopPoint *XMLAnnotatedStopPointRef) Monitored() bool {
	if !annotatedStopPoint.monitored.Defined {
		annotatedStopPoint.monitored.SetValue(annotatedStopPoint.findBoolChildContent("Monitored"))
	}
	return annotatedStopPoint.monitored.Value
}
