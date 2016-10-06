package siri

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/af83/edwig/api"
)

type SOAPClient struct {
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
	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(strings.Join([]string{"Request error, response status code: ", strconv.Itoa(response.StatusCode)}, ""))
	}

	// Create XMLCheckStatusResponse
	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	xmlResponse, err := NewXMLCheckStatusResponseFromContent(responseContent)
	if err != nil {
		return nil, err
	}

	return xmlResponse, nil
}

func CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	requestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//Handle error
	}
	xmlRequest, err := NewXMLCheckStatusRequestFromContent(requestContent)
	if err != nil {
		// Handle error
	}

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
