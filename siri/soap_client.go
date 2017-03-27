package siri

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jbowtie/gokogiri/xml"
	"golang.org/x/text/encoding/charmap"
)

type Request interface {
	BuildXML() (string, error)
}

type SOAPClient struct {
	url string
}

func NewSOAPClient(url string) *SOAPClient {
	return &SOAPClient{url: url}
}
func (client *SOAPClient) URL() string {
	return client.url
}

func (client *SOAPClient) prepareAndSendRequest(request Request, resource string, acceptGzip bool) (xml.Node, error) {
	// Wrap the request XML
	soapEnvelope := NewSOAPEnvelopeBuffer()
	xml, err := request.BuildXML()
	if err != nil {
		return nil, err
	}

	soapEnvelope.WriteXML(xml)
	// Create http request
	httpRequest, err := http.NewRequest("POST", client.url, soapEnvelope)
	if err != nil {
		return nil, err
	}
	if acceptGzip {
		httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	}
	httpRequest.Header.Set("Content-Type", "text/xml")
	httpRequest.ContentLength = soapEnvelope.Length()

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

	if !strings.Contains(response.Header.Get("Content-Type"), "text/xml") {
		return nil, newSiriError(fmt.Sprintf("SIRI CRITICAL: HTTP Content-Type %v", response.Header.Get("Content-Type")))
	}

	// Check if response is gzip
	var responseReader io.Reader
	if acceptGzip && response.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()
		responseReader = gzipReader
	} else {
		responseReader = response.Body
	}

	if response.Header.Get("Content-Type") == "text/xml;charset=ISO-8859-1" {
		responseReader = charmap.ISO8859_1.NewDecoder().Reader(responseReader)
	}

	// Create SOAPEnvelope and check body type
	envelope, err := NewSOAPEnvelope(responseReader)
	if err != nil {
		return nil, err
	}
	if envelope.BodyType() != resource {
		return nil, newSiriError(fmt.Sprintf("SIRI CRITICAL: Wrong Soap from server: %v", envelope.BodyType()))
	}

	return envelope.Body(), nil
}

func (client *SOAPClient) CheckStatus(request *SIRICheckStatusRequest) (*XMLCheckStatusResponse, error) {
	node, err := client.prepareAndSendRequest(request, "CheckStatusResponse", true)
	if err != nil {
		return nil, err
	}

	checkStatus := NewXMLCheckStatusResponse(node)
	return checkStatus, nil
}

func (client *SOAPClient) StopMonitoring(request *SIRIStopMonitoringRequest) (*XMLStopMonitoringResponse, error) {
	// WIP
	node, err := client.prepareAndSendRequest(request, "GetStopMonitoringResponse", false)
	if err != nil {
		return nil, err
	}

	stopMonitoring := NewXMLStopMonitoringResponse(node)
	return stopMonitoring, nil
}
