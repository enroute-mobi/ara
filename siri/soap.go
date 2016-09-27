package siri

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type SOAPClient struct {
	url string
}

func NewSOAPClient(url string) *SOAPClient {
	return &SOAPClient{url: url}
}

const SOAPRequestTemplate = `<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
	{{.BuildXML}}
	</S:Body>
</S:Envelope>`

func (client *SOAPClient) CheckStatus(request *SIRICheckStatusRequest) (*XMLCheckStatusResponse, error) {
	// Wrap soap Request
	var buffer bytes.Buffer
	var soapRequest = template.Must(template.New("soapRequest").Parse(SOAPRequestTemplate))
	if err := soapRequest.Execute(&buffer, request); err != nil {
		return nil, err
	}

	// Create http request
	httpRequest, err := http.NewRequest("POST", "http://server/siri", bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	httpRequest.Header.Set("Content-Type", "text/xml")

	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(strings.Join([]string{"Request error, response status code :", strconv.Itoa(resp.StatusCode)}, " "))
	}

	// Create XMLCheckStatusResponse
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	xmlRequest := NewXMLCheckStatusResponseFromContent(responseContent)

	return xmlRequest, nil
}
