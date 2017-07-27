package siri

import (
	"strings"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetEstimatedTimetableRequest struct {
	RequestXMLStructure

	lines []string
}

func NewXMLGetEstimatedTimetableRequest(node xml.Node) *XMLGetEstimatedTimetableRequest {
	xmlCheckStatusRequest := &XMLGetEstimatedTimetableRequest{}
	xmlCheckStatusRequest.node = NewXMLNode(node)
	return xmlCheckStatusRequest
}

func NewXMLGetEstimatedTimetableRequestFromContent(content []byte) (*XMLGetEstimatedTimetableRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetEstimatedTimetableRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetEstimatedTimetableRequest) Lines() []string {
	if len(request.lines) == 0 {
		nodes := request.findNodes("LineRef")
		if nodes != nil {
			for _, node := range nodes {
				request.lines = append(request.lines, strings.TrimSpace(node.NativeNode().Content()))
			}
		}
	}
	return request.lines
}
