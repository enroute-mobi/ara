package remote

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/jbowtie/gokogiri/xml"
)

const (
	RAW_SIRI_ENVELOPE  = "raw"
	SOAP_SIRI_ENVELOPE = "soap"
)

type SIRIEnvelope struct {
	body xml.Node

	bodyType string
}

func NewSIRIEnvelope(body io.Reader, envelopeType string) (*SIRIEnvelope, error) {
	// Attempt to read the body
	content, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("empty body")
	}
	// Parse the XML and store the body
	doc, err := xml.Parse(content, xml.DefaultEncodingBytes, nil, xml.StrictParseOption, xml.DefaultEncodingBytes)
	if err != nil {
		return nil, err
	}

	switch envelopeType {
	case RAW_SIRI_ENVELOPE:
		return newRawEnvelope(doc)
	case SOAP_SIRI_ENVELOPE:
		return newSOAPEnvelope(doc)
	default:
		return newSOAPEnvelope(doc)
	}
}

func newRawEnvelope(doc *xml.XmlDocument) (*SIRIEnvelope, error) {
	return &SIRIEnvelope{body: doc.Root().XmlNode.FirstChild().NextSibling()}, nil
}

func newSOAPEnvelope(doc *xml.XmlDocument) (*SIRIEnvelope, error) {
	nodes, err := doc.Root().Search("//*[local-name()='Body']/*")
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, errors.New("unable to find body when parsing SOAP request")
	}

	return &SIRIEnvelope{body: nodes[0]}, nil
}

func NewAutodetectSIRIEnvelope(body io.Reader) (*SIRIEnvelope, error) {
	// Attempt to read the body
	content, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("empty body")
	}

	// Parse the XML and store the body
	doc, err := xml.Parse(content, xml.DefaultEncodingBytes, nil, xml.StrictParseOption, xml.DefaultEncodingBytes)
	if err != nil {
		return nil, err
	}

	nodes, err := doc.Root().Search("//*[local-name()='Body']/*")
	if err != nil {
		return nil, fmt.Errorf("cannot search for `Body` tag: %s", err)
	}

	if len(nodes) == 0 {
		node := doc.Root().FirstChild()
		// ensure we only get the first node
		for node != nil && node.NodeType() != xml.XML_ELEMENT_NODE {
			node = node.NextSibling()
		}
		if node == nil {
			return nil, errors.New("invalid raw xml: cannot find body")
		}
		return &SIRIEnvelope{body: node}, nil
	}

	return &SIRIEnvelope{body: nodes[0]}, nil
}

func (envelope *SIRIEnvelope) BodyType() string {
	if envelope.bodyType == "" {
		envelope.bodyType = envelope.body.Name()
		switch envelope.bodyType {
		case "ServiceDelivery":
			xpath := fmt.Sprintf(".//*[contains(local-name(), '%s')]", "Delivery")
			node, err := envelope.body.Search(xpath)
			if err != nil {
				return ""
			}
			if len(node) == 0 {
				return ""
			}
			envelope.bodyType = node[0].Name()
		case "ServiceRequest":
			xpath := fmt.Sprintf(".//*[substring(name(),string-length(name())-6) = '%s']", "Request")
			node, err := envelope.body.Search(xpath)
			if err != nil {
				return ""
			}
			if len(node) == 0 {
				return ""
			}
			envelope.bodyType = node[0].Name()
		}
	}

	return strings.Replace(envelope.bodyType, "Request", "", -1)
}

func (envelope *SIRIEnvelope) Body() xml.Node {
	return envelope.body
}

func (envelope *SIRIEnvelope) BodyOrError(expectedResponsess []string) (xml.Node, error) {
	if slices.Contains(expectedResponsess, envelope.BodyType()) {
		return envelope.body, nil
	}
	if envelope.BodyType() == "Fault" {
		se := sxml.NewXMLSiriError(envelope.body)
		return nil, siri.NewSiriError(fmt.Sprintf("SIRI Error: %v", se.Error()))
	}
	return nil, siri.NewSiriError(fmt.Sprintf("SIRI CRITICAL: Wrong xml from server: %v", envelope.BodyType()))
}
