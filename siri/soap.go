package siri

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type SOAPClient struct {
	url string
}

func NewSOAPClient(url string) *SOAPClient {
	return &SOAPClient{url: url}
}

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

// Handle SIRI CRITICAL errors
type SiriError struct {
	message string
}

func (e *SiriError) Error() string {
	return e.message
}

func newSiriError(message string) error {
	return &SiriError{message: message}
}

// Temp
func WrapSoap(s string) string {
	soap := strings.Join([]string{
		"<?xml version='1.0' encoding='utf-8'?>\n<S:Envelope xmlns:S=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\">\n\t<S:Body>\n",
		s,
		"\n\t</S:Body>\n</S:Envelope>"}, "")
	return soap
}

func (client *SOAPClient) CheckStatus(request *SIRICheckStatusRequest) (*XMLCheckStatusResponse, error) {
	// Wrap the request XML
	soapRequest := WrapSoap(request.BuildXML())

	// Create http request
	httpRequest, err := http.NewRequest("POST", client.url, strings.NewReader(soapRequest))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	httpRequest.Header.Set("Content-Type", "text/xml")

	// Send http request
	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, newSiriError(strings.Join([]string{"SIRI CRITICAL: HTTP status ", strconv.Itoa(response.StatusCode)}, ""))
	}

	// Check if response is gzip
	var responseReader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		responseReader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer responseReader.Close()
	} else {
		responseReader = response.Body
	}

	// Create XMLCheckStatusResponse
	envelope, err := NewSOAPEnvelope(responseReader)

	if envelope.BodyType() != "CheckStatusResponse" {
		return nil, newSiriError(fmt.Sprintf("SIRI CRITICAL: Wrong Soap from server: %v", envelope.BodyType()))
	}

	checkStatus := NewXMLCheckStatusResponse(envelope.Body())

	return checkStatus, nil
}
