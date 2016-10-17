package siri

import (
	"io"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

// Need to check if all are needed
type SOAPEnvelope struct {
	body xml.Node

	bodyType string
}

func NewSOAPEnvelope(body io.Reader) (*SOAPEnvelope, error) {
	// Attempt to read the body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// Parse the XML and store the body
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	nodes, err := doc.Root().Search("//*[local-name()='Body']/*")
	if err != nil {
		log.Fatal(err)
	}

	soapEnvelope := &SOAPEnvelope{body: nodes[0]}
	finalizer := func(document *xml.XmlDocument) {
		document.Free()
	}
	runtime.SetFinalizer(doc, finalizer)

	return soapEnvelope, nil
}

func (envelope *SOAPEnvelope) BodyType() string {
	if envelope.bodyType == "" {
		envelope.bodyType = envelope.body.Name()
	}
	return envelope.bodyType
}

func (envelope *SOAPEnvelope) Body() xml.Node {
	return envelope.body
}
