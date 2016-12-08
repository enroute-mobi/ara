package siri

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jbowtie/gokogiri/xml"
)

type Request interface {
	BuildXML() (xml string)
}

type SOAPClient struct {
	url string
}

func NewSOAPClient(url string) *SOAPClient {
	return &SOAPClient{url: url}
}

func (client *SOAPClient) prepareAndSendRequest(request Request, resource string, acceptGzip bool) (xml.Node, error) {
	// Wrap the request XML
	soapEnvelope := NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(request.BuildXML())

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

	if response.Header.Get("Content-Type") != "text/xml" {
		return nil, newSiriError(fmt.Sprintf("SIRI CRITICAL: HTTP Content-Type %v", response.Header.Get("Content-Type")))
	}

	// Check if response is gzip
	var responseReader io.ReadCloser
	if acceptGzip && response.Header.Get("Content-Encoding") == "gzip" {
		responseReader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer responseReader.Close()
	} else {
		responseReader = response.Body
	}

	// Create SOAPEnvelope and check body type
	envelope, err := NewSOAPEnvelope(responseReader)
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
