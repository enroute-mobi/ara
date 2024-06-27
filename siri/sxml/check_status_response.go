package sxml

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusResponse struct {
	ResponseXMLStructureWithStatus

	serviceStartedTime time.Time
}

func NewXMLCheckStatusResponse(node xml.Node) *XMLCheckStatusResponse {
	xmlCheckStatusResponse := &XMLCheckStatusResponse{}
	xmlCheckStatusResponse.node = NewXMLNode(node)
	return xmlCheckStatusResponse
}

func NewXMLCheckStatusResponseFromContent(content []byte) (*XMLCheckStatusResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLCheckStatusResponse(doc.Root().XmlNode)
	return response, nil
}

func (response *XMLCheckStatusResponse) ErrorString() string {
	return fmt.Sprintf("%v: %v", response.errorType(), response.ErrorText())
}

func (response *XMLCheckStatusResponse) errorType() string {
	if response.ErrorType() == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", response.ErrorType(), response.ErrorNumber())
	}
	return response.ErrorType()
}

func (response *XMLCheckStatusResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent(siri_attributes.ServiceStartedTime)
	}
	return response.serviceStartedTime
}
