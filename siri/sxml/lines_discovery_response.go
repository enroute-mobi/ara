package sxml

import (
	"fmt"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLLinesDiscoveryResponse struct {
	LightDeliveryXMLStructure

	annotatedLineRefs []*XMLAnnotatedLineRef
}

type XMLAnnotatedLineRef struct {
	XMLStructure

	lineRef  string
	lineName string

	monitored Bool
}

func NewXMLLinesDiscoveryResponse(node xml.Node) *XMLLinesDiscoveryResponse {
	xmlLinesDiscoveryResponse := &XMLLinesDiscoveryResponse{}
	xmlLinesDiscoveryResponse.node = NewXMLNode(node)
	return xmlLinesDiscoveryResponse
}

func NewXMLLinesDiscoveryResponseFromContent(content []byte) (*XMLLinesDiscoveryResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLLinesDiscoveryResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLAnnotatedLineRef(node XMLNode) *XMLAnnotatedLineRef {
	annotatedLine := &XMLAnnotatedLineRef{}
	annotatedLine.node = node
	return annotatedLine
}

func (response *XMLLinesDiscoveryResponse) ErrorString() string {
	return fmt.Sprintf("%v: %v", response.errorType(), response.ErrorText())
}

func (response *XMLLinesDiscoveryResponse) errorType() string {
	if response.ErrorType() == "OtherError" {
		return fmt.Sprintf("%v %v", response.ErrorType(), response.ErrorNumber())
	}
	return response.ErrorType()
}

func (response *XMLLinesDiscoveryResponse) AnnotatedLineRefs() []*XMLAnnotatedLineRef {
	if response.annotatedLineRefs == nil {
		annotatedLineRefs := []*XMLAnnotatedLineRef{}
		nodes := response.findNodes("AnnotatedLineRef")
		for _, node := range nodes {
			annotatedLineRefs = append(annotatedLineRefs, NewXMLAnnotatedLineRef(node))
		}
		response.annotatedLineRefs = annotatedLineRefs
	}
	return response.annotatedLineRefs
}

func (annotatedLine *XMLAnnotatedLineRef) LineRef() string {
	if annotatedLine.lineRef == "" {
		annotatedLine.lineRef = annotatedLine.findStringChildContent("LineRef")
	}
	return annotatedLine.lineRef
}

func (annotatedLine *XMLAnnotatedLineRef) LineName() string {
	if annotatedLine.lineName == "" {
		annotatedLine.lineName = annotatedLine.findStringChildContent("LineName")
	}
	return annotatedLine.lineName
}

func (annotatedLine *XMLAnnotatedLineRef) Monitored() bool {
	if !annotatedLine.monitored.Defined {
		annotatedLine.monitored.SetValue(annotatedLine.findBoolChildContent("Monitored"))
	}
	return annotatedLine.monitored.Value
}
