package siri

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/af83/edwig/version"
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

func (client *SOAPClient) responseFromFormat(body io.Reader, contentType string) io.Reader {
	r, _ := regexp.Compile("^text/xml;charset=([ -~]+)")
	s := r.FindStringSubmatch(contentType)
	if len(s) == 0 {
		return body
	}
	if s[1] == "ISO-8859-1" {
		return charmap.ISO8859_1.NewDecoder().Reader(body)
	}
	return body
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
	httpRequest.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpRequest.Header.Set("User-Agent", version.ApplicationName())
	httpRequest.ContentLength = soapEnvelope.Length()

	// Send http request
	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, NewSiriError(strings.Join([]string{"SIRI CRITICAL: HTTP status ", strconv.Itoa(response.StatusCode)}, ""))
	}

	if !strings.Contains(response.Header.Get("Content-Type"), "text/xml") {
		return nil, NewSiriError(fmt.Sprintf("SIRI CRITICAL: HTTP Content-Type %v", response.Header.Get("Content-Type")))
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
		responseReader = client.responseFromFormat(response.Body, response.Header.Get("Content-Type"))
	}

	// Create SOAPEnvelope and check body type
	envelope, err := NewSOAPEnvelope(responseReader)
	if err != nil {
		return nil, err
	}
	if envelope.BodyType() != resource {
		return nil, NewSiriError(fmt.Sprintf("SIRI CRITICAL: Wrong Soap from server: %v", envelope.BodyType()))
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

func (client *SOAPClient) StopMonitoring(request *SIRIGetStopMonitoringRequest) (*XMLStopMonitoringResponse, error) {
	// WIP
	node, err := client.prepareAndSendRequest(request, "GetStopMonitoringResponse", false)
	if err != nil {
		return nil, err
	}

	stopMonitoring := NewXMLStopMonitoringResponse(node)
	return stopMonitoring, nil
}

// WIP
func (client *SOAPClient) StopMonitoringSubscription(request *SIRIStopMonitoringSubscriptionRequest) (*XMLStopMonitoringSubscriptionResponse, error) {

	node, err := client.prepareAndSendRequest(request, "SubscribeResponse", false)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringSubscriptionResponse(node)
	return response, nil
}

func (client *SOAPClient) GeneralMessageSubscription(request *SIRIGeneralMessageSubscriptionRequest) (*XMLGeneralMessageSubscriptionResponse, error) {
	node, err := client.prepareAndSendRequest(request, "SubscribeResponse", false)
	if err != nil {
		return nil, err
	}
	response := NewXMLGeneralMessageSubscriptionResponse(node)
	return response, nil
}

func (client *SOAPClient) SituationMonitoring(request *SIRIGetGeneralMessageRequest) (*XMLGeneralMessageResponse, error) {
	// WIP
	node, err := client.prepareAndSendRequest(request, "GetGeneralMessageResponse", false)
	if err != nil {
		return nil, err
	}

	generalMessage := NewXMLGeneralMessageResponse(node)
	return generalMessage, nil
}

func (client *SOAPClient) DeleteSubscription(request *SIRIDeleteSubscriptionRequest) (*XMLDeleteSubscriptionResponse, error) {
	// WIP
	node, err := client.prepareAndSendRequest(request, "DeleteSubscription", false)
	if err != nil {
		return nil, err
	}

	terminatedSub := NewXMLDeleteSubscriptionResponse(node)
	return terminatedSub, nil
}

func (client *SOAPClient) NotifyStopMonitoring(request *SIRINotifyStopMonitoring) error {
	_, err := client.prepareAndSendRequest(request, "NotifyStopMonitoring", false)
	if err != nil {
		return err
	}
	return nil
}

func (client *SOAPClient) NotifyGeneralMessage(request *SIRINotifyGeneralMessage) error {
	_, err := client.prepareAndSendRequest(request, "NotifyGeneralMessage", false)
	if err != nil {
		return err
	}
	return nil
}

func (client *SOAPClient) NotifyEstimatedTimeTable(request *SIRINotifyEstimatedTimeTable) error {
	_, err := client.prepareAndSendRequest(request, "NotifyEstimatedTimeTable", false)
	if err != nil {
		return err
	}
	return nil
}
