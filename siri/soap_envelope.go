package siri

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"runtime"

	"github.com/af83/edwig/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type SOAPEnvelope struct {
	body xml.Node

	bodyType string
}

type SOAPEnvelopeBuffer struct {
	buffer bytes.Buffer
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
		logger.Log.Panicf("Error while parsing XML: %v", err)
	}

	if len(nodes) == 0 {
		return nil, errors.New("Unable to find body when parsing SOAP request")
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

func NewSOAPEnvelopeBuffer() *SOAPEnvelopeBuffer {
	return &SOAPEnvelopeBuffer{}
}

func (writer *SOAPEnvelopeBuffer) WriteXML(xml string) {
	writer.buffer.WriteString("<?xml version='1.0' encoding='utf-8'?>\n<S:Envelope xmlns:S=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\">\n<S:Body>\n")
	writer.buffer.WriteString(xml)
	writer.buffer.WriteString("\n</S:Body>\n</S:Envelope>")
}

func (writer *SOAPEnvelopeBuffer) Read(p []byte) (n int, err error) {
	n, err = writer.buffer.Read(p)
	return
}

func (writer *SOAPEnvelopeBuffer) WriteTo(w io.Writer) (n int64, err error) {
	n, err = writer.buffer.WriteTo(w)
	return
}

func (writer *SOAPEnvelopeBuffer) String() string {
	return writer.buffer.String()
}

func (writer *SOAPEnvelopeBuffer) Length() int64 {
	return int64(writer.buffer.Len())
}
