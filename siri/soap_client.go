package siri

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SOAPClient struct {
	url string
}

func NewSOAPClient(url string) *SOAPClient {
	return &SOAPClient{url: url}
}

func (client *SOAPClient) CheckStatus(request *SIRICheckStatusRequest) (*XMLCheckStatusResponse, error) {
	// Wrap the request XML
	soapEnvelope := NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(request.BuildXML())

	// Create http request
	httpRequest, err := http.NewRequest("POST", client.url, soapEnvelope)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
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

func (client *SOAPClient) StopMonitoring(request *SIRIStopMonitoringRequest) (*XMLStopMonitoringResponse, error) {
	return nil, nil
}
