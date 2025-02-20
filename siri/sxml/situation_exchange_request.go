package sxml

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetSituationExchange struct {
	XMLSituationExchangeRequest

	requestorRef string
}

type XMLSituationExchangeRequest struct {
	LightRequestXMLStructure

	lineRefs        []string
	stopPointRefs   []string
	previewInterval time.Duration
	startTime       time.Time
}

func NewXMLGetSituationExchange(node xml.Node) *XMLGetSituationExchange {
	xmlGetSituationExchange := &XMLGetSituationExchange{}
	xmlGetSituationExchange.node = NewXMLNode(node)
	return xmlGetSituationExchange
}

func NewXMLGetSituationExchangeFromContent(content []byte) (*XMLGetSituationExchange, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetSituationExchange(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetSituationExchange) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent(siri_attributes.RequestorRef)
	}
	return request.requestorRef
}

func (request *XMLSituationExchangeRequest) LineRef() []string {
	if len(request.lineRefs) == 0 {
		nodes := request.findNodes(siri_attributes.LineRef)
		for _, lineRef := range nodes {
			request.lineRefs = append(request.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return request.lineRefs
}

func (request *XMLSituationExchangeRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent(siri_attributes.PreviewInterval)
	}
	return request.previewInterval
}

func (request *XMLSituationExchangeRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent(siri_attributes.StartTime)
	}
	return request.startTime
}
