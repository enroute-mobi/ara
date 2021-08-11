package siri

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNode interface {
	NativeNode() xml.Node
}

func NewXMLNode(nativeNode xml.Node) XMLNode {
	node := &RootXMLNode{rootNode: nativeNode}

	finalizer := func(node *RootXMLNode) {
		node.Free()
	}
	runtime.SetFinalizer(node, finalizer)

	return node
}

func NewXMLNodeFromContent(content []byte) (XMLNode, error) {
	document, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	return NewXMLNode(document.Root().XmlNode), nil
}

type RootXMLNode struct {
	rootNode xml.Node
}

func (node *RootXMLNode) NativeNode() xml.Node {
	return node.rootNode
}

func (node *RootXMLNode) Free() {
	if node.rootNode != nil {
		node.rootNode.MyDocument().Free()
		node.rootNode = nil
	}
}

type SubXMLNode struct {
	parent     XMLNode
	nativeNode xml.Node
}

func (node *SubXMLNode) NativeNode() xml.Node {
	return node.nativeNode
}

func NewSubXMLNode(nativeNode xml.Node) *SubXMLNode {
	return &SubXMLNode{nativeNode: nativeNode}
}

type XMLStructure struct {
	node XMLNode
}

// Test Method
func (xmlStruct *XMLStructure) Node() XMLNode {
	return xmlStruct.node
}

type RequestXMLStructure struct {
	LightRequestXMLStructure

	requestorRef string
}

type LightRequestXMLStructure struct {
	XMLStructure

	messageIdentifier string
	requestTimestamp  time.Time
}

type ResponseXMLStructure struct {
	XMLStructure

	address                   string
	producerRef               string
	requestMessageRef         string
	responseMessageIdentifier string
	responseTimestamp         time.Time
}

type ResponseXMLStructureWithStatus struct {
	DeliveryXMLStructure // To avoid code duplication as much as possible

	address                   string
	producerRef               string
	responseMessageIdentifier string
}

type DeliveryXMLStructure struct {
	LightDeliveryXMLStructure

	requestMessageRef string
}

type LightDeliveryXMLStructure struct {
	XMLStatus

	responseTimestamp time.Time
}

type SubscriptionDeliveryXMLStructure struct {
	LightSubscriptionDeliveryXMLStructure

	requestMessageRef string
}

type LightSubscriptionDeliveryXMLStructure struct {
	LightDeliveryXMLStructure

	subscriberRef   string
	subscriptionRef string
}

type XMLStatus struct {
	XMLStructure

	status           Bool
	errorType        string
	errorNumber      int
	errorText        string
	errorDescription string
}

func (xmlStruct *XMLStructure) findNodeWithNamespace(localName string) xml.Node {
	xpath := fmt.Sprintf(".//*[local-name()='%s']", localName)

	nodes, err := xmlStruct.node.NativeNode().Search(xpath)
	if err != nil {
		return nil
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// func (xmlStruct *XMLStructure) findXMLNode(localName string) XMLNode {
// 	xpath := fmt.Sprintf(".//*[local-name()='%s']", localName)
// 	nodes, err := xmlStruct.node.NativeNode().Search(xpath)
// 	if err != nil {
// 		return nil
// 	}
// 	if len(nodes) == 0 {
// 		return nil
// 	}

// 	subNode := NewSubXMLNode(nodes[0])
// 	subNode.parent = xmlStruct.node

// 	return subNode
// }

func (xmlStruct *XMLStructure) findNode(localName string) xml.Node {
	xpath := fmt.Sprintf(".//%s", localName)

	nodes, err := xmlStruct.node.NativeNode().Search(xpath)
	if err != nil || len(nodes) == 0 {
		return xmlStruct.findNodeWithNamespace(localName)
	}
	return nodes[0]
}

func (xmlStruct *XMLStructure) findNodes(localName string) []XMLNode {
	return xmlStruct.nodes(fmt.Sprintf(".//*[local-name()='%s']", localName))
}

func (xmlStruct *XMLStructure) findDirectChildrenNodes(localName string) []XMLNode {
	return xmlStruct.nodes(fmt.Sprintf("./*[local-name()='%s']", localName))
}

func (xmlStruct *XMLStructure) nodes(xpath string) []XMLNode {
	nodes, err := xmlStruct.node.NativeNode().Search(xpath)
	if err != nil {
		return nil
	}
	if len(nodes) == 0 {
		return nil
	}

	xmlNodes := make([]XMLNode, 0)
	for _, node := range nodes {
		subNode := NewSubXMLNode(node)
		subNode.parent = xmlStruct.node
		xmlNodes = append(xmlNodes, subNode)
	}

	return xmlNodes
}

// TODO: See how to handle errors
func (xmlStruct *XMLStructure) findStringChildContent(localName string) string {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return ""
	}
	return strings.TrimSpace(node.Content())
}

func (xmlStruct *XMLStructure) findChildAttribute(localName, attr string) string {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return ""
	}
	return strings.TrimSpace(node.Attr(attr))
}

func (xmlStruct *XMLStructure) containSelfClosing(localName string) bool {
	node := xmlStruct.findNode(localName)
	return node != nil
}

func (xmlStruct *XMLStructure) findTimeChildContent(localName string) time.Time {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", strings.TrimSpace(node.Content()))
	// t, err := time.Parse(time.RFC3339, strings.TrimSpace(node.Content()))
	if err != nil {
		return time.Time{}
	}
	return t
}

func (xmlStruct *XMLStructure) findDurationChildContent(localName string) time.Duration {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return 0
	}
	durationRegex := regexp.MustCompile(`P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?`)
	matches := durationRegex.FindStringSubmatch(strings.TrimSpace(node.Content()))

	if len(matches) == 0 {
		return 0
	}
	years := parseDuration(matches[1]) * 24 * 365 * time.Hour
	months := parseDuration(matches[2]) * 30 * 24 * time.Hour
	days := parseDuration(matches[3]) * 24 * time.Hour
	hours := parseDuration(matches[4]) * time.Hour
	minutes := parseDuration(matches[5]) * time.Minute
	seconds := parseDuration(matches[6]) * time.Second

	return time.Duration(years + months + days + hours + minutes + seconds)
}

func parseDuration(value string) time.Duration {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return time.Duration(parsed)
}

func (xmlStruct *XMLStructure) findBoolChildContent(localName string) bool {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return false
	}
	s, err := strconv.ParseBool(strings.TrimSpace(node.Content()))
	if err != nil {
		return false
	}
	return s
}

func (xmlStruct *XMLStructure) findIntChildContent(localName string) int {
	node := xmlStruct.findNode(localName)
	if node == nil {
		return 0
	}
	s, err := strconv.Atoi(strings.TrimSpace(node.Content()))
	if err != nil {
		return 0
	}
	return s
}

func (xmlStruct *XMLStructure) RawXML() string {
	return xmlStruct.node.NativeNode().String()
}

func (request *RequestXMLStructure) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *LightRequestXMLStructure) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *LightRequestXMLStructure) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
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

func (response *ResponseXMLStructure) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent("ResponseMessageIdentifier")
	}
	return response.responseMessageIdentifier
}

func (response *ResponseXMLStructure) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *ResponseXMLStructure) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *ResponseXMLStructureWithStatus) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *ResponseXMLStructureWithStatus) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent("ProducerRef")
	}
	return response.producerRef
}

func (response *ResponseXMLStructureWithStatus) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent("ResponseMessageIdentifier")
	}
	return response.responseMessageIdentifier
}

func (delivery *DeliveryXMLStructure) RequestMessageRef() string {
	if delivery.requestMessageRef == "" {
		delivery.requestMessageRef = delivery.findStringChildContent("RequestMessageRef")
	}
	return delivery.requestMessageRef
}

func (delivery *LightDeliveryXMLStructure) ResponseTimestamp() time.Time {
	if delivery.responseTimestamp.IsZero() {
		delivery.responseTimestamp = delivery.findTimeChildContent("ResponseTimestamp")
	}
	return delivery.responseTimestamp
}

func (delivery *SubscriptionDeliveryXMLStructure) RequestMessageRef() string {
	if delivery.requestMessageRef == "" {
		delivery.requestMessageRef = delivery.findStringChildContent("RequestMessageRef")
	}
	return delivery.requestMessageRef
}

func (delivery *LightSubscriptionDeliveryXMLStructure) SubscriberRef() string {
	if delivery.subscriberRef == "" {
		delivery.subscriberRef = delivery.findStringChildContent("SubscriberRef")
	}
	return delivery.subscriberRef
}

func (delivery *LightSubscriptionDeliveryXMLStructure) SubscriptionRef() string {
	if delivery.subscriptionRef == "" {
		delivery.subscriptionRef = delivery.findStringChildContent("SubscriptionRef")
	}
	return delivery.subscriptionRef
}

func (response *XMLStatus) Status() bool {
	if !response.status.Defined {
		response.status.SetValue(response.findBoolChildContent("Status"))
	}
	return response.status.Value
}

func (response *XMLStatus) ErrorType() string {
	if !response.Status() && response.errorType == "" {
		node := response.findNode("ErrorText")
		if node != nil {
			response.errorType = node.Parent().Name()
			// Find errorText and errorNumber to avoir too much parsing
			response.errorText = strings.TrimSpace(node.Content())
			if response.errorType == "OtherError" {
				n, err := strconv.Atoi(node.Parent().Attr("number"))
				if err != nil {
					return ""
				}
				response.errorNumber = n
			}
		}
	}
	return response.errorType
}

func (response *XMLStatus) ErrorNumber() int {
	if !response.Status() && response.ErrorType() == "OtherError" && response.errorNumber == 0 {
		node := response.findNode("ErrorText")
		n, err := strconv.Atoi(node.Parent().Attr("number"))
		if err != nil {
			return -1
		}
		response.errorNumber = n
	}
	return response.errorNumber
}

func (response *XMLStatus) ErrorText() string {
	if !response.Status() && response.errorText == "" {
		response.errorText = response.findStringChildContent("ErrorText")
	}
	return response.errorText
}

func (response *XMLStatus) ErrorDescription() string {
	if !response.Status() && response.errorDescription == "" {
		response.errorDescription = response.findStringChildContent("Description")
	}
	return response.errorDescription
}
