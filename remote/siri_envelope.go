package remote

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"bitbucket.org/enroute-mobi/ara/siri"
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
	content, err := ioutil.ReadAll(body)
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
	content, err := ioutil.ReadAll(body)
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
		return nil, err
	}

	if len(nodes) == 0 {
		return &SIRIEnvelope{body: doc.Root().XmlNode.FirstChild().NextSibling()}, nil
	}

	return &SIRIEnvelope{body: nodes[0]}, nil
}

func (envelope *SIRIEnvelope) BodyType() string {
	if envelope.bodyType == "" {
		envelope.bodyType = envelope.body.Name()
	}
	return envelope.bodyType
}

func (envelope *SIRIEnvelope) Body() xml.Node {
	return envelope.body
}

func (envelope *SIRIEnvelope) BodyOrError(expectedResponse string) (xml.Node, error) {
	if envelope.BodyType() == expectedResponse {
		return envelope.body, nil
	}
	if envelope.BodyType() == "Fault" {
		se := siri.NewXMLSiriError(envelope.body)
		return nil, siri.NewSiriError(fmt.Sprintf("SIRI Error: %v", se.Error()))
	}
	return nil, siri.NewSiriError(fmt.Sprintf("SIRI CRITICAL: Wrong xml from server: %v", envelope.BodyType()))
}
