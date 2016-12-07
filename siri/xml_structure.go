package siri

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStructure struct {
	node xml.Node
}

type ResponseXMLStructure struct {
	XMLStructure

	address                   string
	producerRef               string
	requestMessageRef         string
	responseMessageIdentifier string
	responseTimestamp         time.Time
}

func (xmlStruct *XMLStructure) findNode(localName string) xml.Node {
	xpath := fmt.Sprintf("//*[local-name()='%s']", localName)
	nodes, err := xmlStruct.node.Search(xpath)
	if err != nil {
		logger.Log.Panicf("Error while parsing XML: %v", err)
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// TODO: See how to handle errors
func (xmlStruct *XMLStructure) findStringChildContent(localName string) string {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return ""
	}
	return strings.TrimSpace(node.Content())
}

func (xmlStruct *XMLStructure) findTimeChildContent(localName string) time.Time {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(node.Content()))
	if err != nil {
		logger.Log.Panicf("Error while parsing XML: %v", err)
	}
	return t
}

func (xmlStruct *XMLStructure) findBoolChildContent(localName string) bool {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return false
	}
	s, err := strconv.ParseBool(strings.TrimSpace(node.Content()))
	if err != nil {
		logger.Log.Panicf("Error while parsing XML: %v", err)
	}
	return s
}

func (xmlStruct *XMLStructure) RawXML() string {
	return xmlStruct.node.String()
}

func (response *ResponseXMLStructure) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *ResponseXMLStructure) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent("ProducerRef")
	}
	return response.producerRef
}
func (response *ResponseXMLStructure) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *ResponseXMLStructure) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent("ResponseMessageIdentifier")
	}
	return response.responseMessageIdentifier
}

func (response *ResponseXMLStructure) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}
