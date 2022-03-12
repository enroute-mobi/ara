package remote

import (
	"bytes"
	"io"
)

type Buffer interface {
	io.Reader
	io.WriterTo
	WriteXML(string)
	String() string
	Length() int64
}

func NewSIRIBuffer(envelopeType string) Buffer {
	switch envelopeType {
	case RAW_SIRI_ENVELOPE:
		return newRawBuffer()
	case SOAP_SIRI_ENVELOPE:
		return newSOAPBuffer()
	default:
		return newSOAPBuffer()
	}
}

type RawBuffer struct {
	b bytes.Buffer
}

func newRawBuffer() *RawBuffer {
	return &RawBuffer{}
}

func (rb *RawBuffer) WriteXML(xml string) {
	rb.b.WriteString(xml)
}

func (rb *RawBuffer) Read(p []byte) (n int, err error) {
	n, err = rb.b.Read(p)
	return
}

func (rb *RawBuffer) WriteTo(w io.Writer) (n int64, err error) {
	n, err = rb.b.WriteTo(w)
	return
}

func (rb *RawBuffer) String() string {
	return rb.b.String()
}

func (rb *RawBuffer) Length() int64 {
	return int64(rb.b.Len())
}

type SOAPBuffer struct {
	b bytes.Buffer
}

func newSOAPBuffer() *SOAPBuffer {
	return &SOAPBuffer{}
}

func (sb *SOAPBuffer) WriteXML(xml string) {
	sb.b.WriteString("<?xml version='1.0' encoding='utf-8'?>\n<S:Envelope xmlns:S=\"http://schemas.xmlsoap.org/soap/envelope/\">\n<S:Body>\n")
	sb.b.WriteString(xml)
	sb.b.WriteString("\n</S:Body>\n</S:Envelope>")
}

func (sb *SOAPBuffer) Read(p []byte) (n int, err error) {
	n, err = sb.b.Read(p)
	return
}

func (sb *SOAPBuffer) WriteTo(w io.Writer) (n int64, err error) {
	n, err = sb.b.WriteTo(w)
	return
}

func (sb *SOAPBuffer) String() string {
	return sb.b.String()
}

func (sb *SOAPBuffer) Length() int64 {
	return int64(sb.b.Len())
}
