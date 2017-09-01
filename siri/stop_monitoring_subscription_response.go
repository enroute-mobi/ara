package siri

import (
	"bytes"
	"text/template"
	"strconv"
	"strings"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

//WIP
type XMLStopMonitoringSubscriptionResponse struct {
	XMLStructure

	address           string
	requestMessageRef string
	responderRef      string

	responseTimestamp  time.Time
	serviceStartedTime time.Time

	responseStatus []*XMLResponseStatus
}

type XMLResponseStatus struct {
	XMLStructure

	requestMessageRef string
	subscriberRef     string
	subscriptionRef   string

	status           bool
	errorType        string
	errorNumber      int
	errorText        string
	errorDescription string

	responseTimestamp time.Time
	validUntil        time.Time
}

type SIRIStopMonitoringSubscriptionResponse struct {
	Address           string
	ResponderRef      string
	RequestMessageRef string

	ResponseTimestamp  time.Time
	ServiceStartedTime time.Time

	ResponseStatus []SIRIResponseStatus
}

type SIRIResponseStatus struct {
	RequestMessageRef string
	SubscriberRef     string
	SubscriptionRef   string

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	ResponseTimestamp time.Time
	ValidUntil        time.Time
}

const stopMonitoringSubscriptionResponseTemplate = `<ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
    <SubscriptionAnswerInfo
			xmlns:ns2="http://www.ifopt.org.uk/acsb"
			xmlns:ns3="http://www.ifopt.org.uk/ifopt"
			xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
			xmlns:ns5="http://www.siri.org.uk/siri"
			xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
        <ns5:Address>{{.Address}}</ns5:Address>
        <ns5:ResponderRef>{{.ResponderRef}}</ns5:ResponderRef>
        <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">{{.RequestMessageRef}}</ns5:RequestMessageRef>
    </SubscriptionAnswerInfo>
    <Answer
			xmlns:ns2="http://www.ifopt.org.uk/acsb"
			xmlns:ns3="http://www.ifopt.org.uk/ifopt"
			xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
			xmlns:ns5="http://www.siri.org.uk/siri"
			xmlns:ns6="http://wsdl.siri.org.uk/siri">{{ range .ResponseStatus }}
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>{{.RequestMessageRef}}</ns5:RequestMessageRef>
            <ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
            <ns5:SubscriptionRef>{{.SubscriptionRef}}</ns5:SubscriptionRef>
            <ns5:Status>{{.Status}}</ns5:Status>{{ if not .Status }}
						<ns5:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
							<ns5:OtherError number="{{.ErrorNumber}}">{{ else }}
							<ns5:{{.ErrorType}}>{{ end }}
								<ns5:ErrorText>{{.ErrorText}}</ns5:ErrorText>
							</ns5:{{.ErrorType}}>
						</ns5:ErrorCondition>{{ end }}{{ if not .ValidUntil.IsZero }}
            <ns5:ValidUntil>{{.ValidUntil}}</ns5:ValidUntil>{{ end }}
        </ns5:ResponseStatus>{{ end }}
        <ns5:ServiceStartedTime>{{.ServiceStartedTime}}</ns5:ServiceStartedTime>
    </Answer>
    <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
</ns1:SubscribeResponse>`

func NewXMLStopMonitoringSubscriptionResponse(node xml.Node) *XMLStopMonitoringSubscriptionResponse {
	xmlStopMonitoringSubscriptionResponse := &XMLStopMonitoringSubscriptionResponse{}
	xmlStopMonitoringSubscriptionResponse.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionResponse
}

func NewXMLStopMonitoringSubscriptionResponseFromContent(content []byte) (*XMLStopMonitoringSubscriptionResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringSubscriptionResponse(doc.Root().XmlNode)
	return response, nil
}

func (response *XMLStopMonitoringSubscriptionResponse) ResponseStatus() []*XMLResponseStatus {
	if len(response.responseStatus) == 0 {
		nodes := response.findNodes("ResponseStatus")
		if nodes == nil {
			return response.responseStatus
		}
		for _, responseStatusNode := range nodes {
			xmlResponseStatus := &XMLResponseStatus{}
			xmlResponseStatus.node = responseStatusNode
			response.responseStatus = append(response.responseStatus, xmlResponseStatus)
		}
	}
	return response.responseStatus
}

func (response *XMLStopMonitoringSubscriptionResponse) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *XMLStopMonitoringSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent("ResponderRef")
	}
	return response.responderRef
}

func (response *XMLStopMonitoringSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLStopMonitoringSubscriptionResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent("ServiceStartedTime")
	}
	return response.serviceStartedTime
}

func (response *XMLStopMonitoringSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLResponseStatus) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLResponseStatus) SubscriberRef() string {
	if response.subscriberRef == "" {
		response.subscriberRef = response.findStringChildContent("SubscriberRef")
	}
	return response.subscriberRef
}

func (response *XMLResponseStatus) SubscriptionRef() string {
	if response.subscriptionRef == "" {
		response.subscriptionRef = response.findStringChildContent("SubscriptionRef")
	}
	return response.subscriptionRef
}

func (response *XMLResponseStatus) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLResponseStatus) ValidUntil() time.Time {
	if response.validUntil.IsZero() {
		response.validUntil = response.findTimeChildContent("ValidUntil")
	}
	return response.validUntil
}

func (response *XMLResponseStatus) Status() bool {
	if !response.status {
		response.status = response.findBoolChildContent("Status")
	}
	return response.status
}

func (response *XMLResponseStatus) ErrorType() string {
	if !response.Status() && response.errorType == "" {
		node := response.findNode("ErrorText")
		if node != nil {
			response.errorType = node.Parent().Name()
			// Find errorText and errorNumber to avoir too much parsing
			response.errorText = strings.TrimSpace(node.Content())
			if response.errorType == "OtherError" {
				n, err := strconv.Atoi(node.Parent().Attr("number"))
				if err != nil {
					return ""
				}
				response.errorNumber = n
			}
		}
	}
	return response.errorType
}

func (response *XMLResponseStatus) ErrorNumber() int {
	if !response.Status() && response.ErrorType() == "OtherError" && response.errorNumber == 0 {
		node := response.findNode("ErrorText")
		n, err := strconv.Atoi(node.Parent().Attr("number"))
		if err != nil {
			return -1
		}
		response.errorNumber = n
	}
	return response.errorNumber
}

func (response *XMLResponseStatus) ErrorText() string {
	if !response.Status() && response.errorText == "" {
		response.errorText = response.findStringChildContent("ErrorText")
	}
	return response.errorText
}

func (response *XMLResponseStatus) ErrorDescription() string {
	if !response.Status() && response.errorDescription == "" {
		response.errorDescription = response.findStringChildContent("Description")
	}
	return response.errorDescription
}

func (response *SIRIStopMonitoringSubscriptionResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("SubscribeResponse").Parse(stopMonitoringSubscriptionResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
