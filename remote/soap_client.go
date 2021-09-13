package remote

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/version"
	"github.com/jbowtie/gokogiri/xml"
	"golang.org/x/text/encoding/charmap"
)

type requestType int

const (
	DEFAULT requestType = iota
	SUBSCRIPTION
	NOTIFICATION
	CHECK_STATUS
)

type Request interface {
	BuildXML() (string, error)
}

type SOAPClient struct {
	httpClient *HTTPClient
}

type soapClientArguments struct {
	request          Request
	requestType      requestType
	expectedResponse string
	acceptGzip       bool
}

func NewSOAPClient(c *HTTPClient) *SOAPClient {
	return &SOAPClient{
		httpClient: c,
	}
}

func (client *SOAPClient) remoteClient() *http.Client {
	return client.httpClient.HTTPClient()
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

func (client *SOAPClient) prepareAndSendRequest(args soapClientArguments) (xml.Node, error) {
	// Wrap the request XML
	soapEnvelope := NewSOAPEnvelopeBuffer()
	xml, err := args.request.BuildXML()
	if err != nil {
		return nil, err
	}

	soapEnvelope.WriteXML(xml)

	// For tests
	// logger.Log.Debugf("%v", soapEnvelope.String())

	// Create http request
	ctx, cncl := context.WithTimeout(context.Background(), getTimeOut(args.requestType))
	defer cncl()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.getURL(args.requestType), soapEnvelope)
	if err != nil {
		return nil, err
	}
	if args.acceptGzip {
		httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	}
	httpRequest.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpRequest.Header.Set("User-Agent", version.ApplicationName())
	httpRequest.ContentLength = soapEnvelope.Length()

	// Send http request
	response, err := client.remoteClient().Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	// Do nothing if request is a notification
	if args.requestType == NOTIFICATION {
		return nil, nil
	}

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, siri.NewSiriError(strings.Join([]string{"SIRI CRITICAL: HTTP status ", strconv.Itoa(response.StatusCode)}, ""))
	}

	if !strings.Contains(response.Header.Get("Content-Type"), "text/xml") {
		return nil, siri.NewSiriError(fmt.Sprintf("SIRI CRITICAL: HTTP Content-Type %v", response.Header.Get("Content-Type")))
	}

	// Check if response is gzip
	var responseReader io.Reader
	if args.acceptGzip && response.Header.Get("Content-Encoding") == "gzip" {
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
	return envelope.BodyOrError(args.expectedResponse)
}

func (client *SOAPClient) getURL(requestType requestType) string {
	switch requestType {
	case SUBSCRIPTION:
		if client.httpClient.SubscriptionsUrl != "" {
			return client.httpClient.SubscriptionsUrl
		}
	case NOTIFICATION:
		if client.httpClient.NotificationsUrl != "" {
			return client.httpClient.NotificationsUrl
		}
	}
	return client.httpClient.Url
}

func getTimeOut(rt requestType) time.Duration {
	switch rt {
	case SUBSCRIPTION:
		return 30 * time.Second
	case CHECK_STATUS:
		return 9 * time.Second
	default:
		return 5 * time.Second
	}
}

func (client *SOAPClient) CheckStatus(request *siri.SIRICheckStatusRequest) (*siri.XMLCheckStatusResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		expectedResponse: "CheckStatusResponse",
		requestType:      CHECK_STATUS,
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}

	checkStatus := siri.NewXMLCheckStatusResponse(node)
	return checkStatus, nil
}

func (client *SOAPClient) StopDiscovery(request *siri.SIRIStopPointsDiscoveryRequest) (*siri.XMLStopPointsDiscoveryResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		expectedResponse: "StopPointsDiscoveryResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}

	stopDiscovery := siri.NewXMLStopPointsDiscoveryResponse(node)
	return stopDiscovery, nil
}

func (client *SOAPClient) StopMonitoring(request *siri.SIRIGetStopMonitoringRequest) (*siri.XMLStopMonitoringResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		expectedResponse: "GetStopMonitoringResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}

	stopMonitoring := siri.NewXMLStopMonitoringResponse(node)
	return stopMonitoring, nil
}

func (client *SOAPClient) SituationMonitoring(request *siri.SIRIGetGeneralMessageRequest) (*siri.XMLGeneralMessageResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		expectedResponse: "GetGeneralMessageResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}

	generalMessage := siri.NewXMLGeneralMessageResponse(node)
	return generalMessage, nil
}

func (client *SOAPClient) VehicleMonitoring(request *siri.SIRIGetVehicleMonitoringRequest) (*siri.XMLVehicleMonitoringResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		expectedResponse: "GetVehicleMonitoringResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}

	vehicleMonitoring := siri.NewXMLVehicleMonitoringResponse(node)
	return vehicleMonitoring, nil
}

func (client *SOAPClient) StopMonitoringSubscription(request *siri.SIRIStopMonitoringSubscriptionRequest) (*siri.XMLSubscriptionResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		requestType:      SUBSCRIPTION,
		expectedResponse: "SubscribeResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}
	response := siri.NewXMLSubscriptionResponse(node)
	return response, nil
}

func (client *SOAPClient) GeneralMessageSubscription(request *siri.SIRIGeneralMessageSubscriptionRequest) (*siri.XMLSubscriptionResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		requestType:      SUBSCRIPTION,
		expectedResponse: "SubscribeResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}
	response := siri.NewXMLSubscriptionResponse(node)
	return response, nil
}

func (client *SOAPClient) DeleteSubscription(request *siri.SIRIDeleteSubscriptionRequest) (*siri.XMLDeleteSubscriptionResponse, error) {
	node, err := client.prepareAndSendRequest(soapClientArguments{
		request:          request,
		requestType:      SUBSCRIPTION,
		expectedResponse: "DeleteSubscriptionResponse",
		acceptGzip:       true,
	})
	if err != nil {
		return nil, err
	}

	terminatedSub := siri.NewXMLDeleteSubscriptionResponse(node)
	return terminatedSub, nil
}

func (client *SOAPClient) NotifyStopMonitoring(request *siri.SIRINotifyStopMonitoring) error {
	_, err := client.prepareAndSendRequest(soapClientArguments{
		request:     request,
		requestType: NOTIFICATION,
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *SOAPClient) NotifyGeneralMessage(request *siri.SIRINotifyGeneralMessage) error {
	_, err := client.prepareAndSendRequest(soapClientArguments{
		request:     request,
		requestType: NOTIFICATION,
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *SOAPClient) NotifyEstimatedTimeTable(request *siri.SIRINotifyEstimatedTimeTable) error {
	_, err := client.prepareAndSendRequest(soapClientArguments{
		request:     request,
		requestType: NOTIFICATION,
	})
	if err != nil {
		return err
	}
	return nil
}
