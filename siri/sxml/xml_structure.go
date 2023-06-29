package sxml

/*
#cgo pkg-config: libxml-2.0
#include <libxml/xmlerror.h>

static inline void libxml_nilErrorHandler(void *ctx, const char *msg, ...) {}

static inline void libxml_SilenceParseErrors() {
  xmlSetGenericErrorFunc(NULL, libxml_nilErrorHandler);
  xmlThrDefSetGenericErrorFunc(NULL, libxml_nilErrorHandler);
}
*/
import "C"
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

var durationRegex = regexp.MustCompile(`P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?`)

type XMLNode interface {
	NativeNode() xml.Node
}

func init() {
	C.libxml_SilenceParseErrors()
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
	s, _ := strconv.ParseBool(strings.TrimSpace(node.Content()))
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
