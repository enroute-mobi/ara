package siri

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/af83/edwig/api"
)

type SOAPClient struct {
	api.ClockConsumer

	url string
}

func NewSOAPClient(url string) *SOAPClient {
	return &SOAPClient{url: url}
}

// Temp
func WrapSoap(s string) string {
	soap := strings.Join([]string{
		"<S:Envelope xmlns:S=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\">\n\t<S:Body>\n",
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
	startTime := client.Clock().Now()
	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseTime := client.Clock().Since(startTime)

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(strings.Join([]string{"Request error, response status code: ", strconv.Itoa(response.StatusCode)}, ""))
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
	responseContent, err := ioutil.ReadAll(responseReader)
	if err != nil {
		return nil, err
	}
	xmlResponse, err := NewXMLCheckStatusResponseFromContent(responseContent)
	if err != nil {
		return nil, err
	}

	// Log
	var logMessage []byte
	if xmlResponse.Status() {
		logMessage = []byte("SIRI OK - status true - ")
	} else {
		logMessage = []byte("SIRI CRITICAL: status false - ")
		if xmlResponse.ErrorType() == "OtherError" {
			logMessage = append(logMessage, fmt.Sprintf("%s %d %s - ", xmlResponse.ErrorType(), xmlResponse.ErrorNumber(), xmlResponse.ErrorText())...)
		} else {
			logMessage = append(logMessage, fmt.Sprintf("%s %s - ", xmlResponse.ErrorType(), xmlResponse.ErrorText())...)
		}
	}
	logMessage = append(logMessage, fmt.Sprintf("%.3f seconds response time", responseTime.Seconds())...)
	log.Println(string(logMessage[:]))

	return xmlResponse, nil
}

func CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Try to read and parse request body
	requestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request: can't read content", 500)
		return
	}
	xmlRequest, err := NewXMLCheckStatusRequestFromContent(requestContent)
	if err != nil {
		http.Error(w, "Invalid request: can't parse content", 500)
		return
	}

	// Set Content-Type header and create a SIRICheckStatusResponse
	w.Header().Set("Content-Type", "text/xml")

	response := new(SIRICheckStatusResponse)
	response.Address = strings.Join([]string{r.URL.Host, r.URL.Path}, "")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = xmlRequest.MessageIdentifier()
	response.GenerateMessageIdentifier()
	response.Status = true // Temp
	response.SetResponseTimestamp()
	response.ServiceStartedTime = api.DefaultClock().Now() //Temp

	fmt.Fprintf(w, response.BuildXML())
}
