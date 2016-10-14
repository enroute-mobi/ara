package siri

import (
	"compress/gzip"
	"io"
	"io/ioutil"
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
	responseContent, err := ioutil.ReadAll(responseReader)
	if err != nil {
		return nil, err
	}
	xmlResponse, err := NewXMLCheckStatusResponseFromContent(responseContent)
	if err != nil {
		return nil, err
	}

	return xmlResponse, nil
}
