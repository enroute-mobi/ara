package siri

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLDeleteSubscriptionResponse struct {
	XMLStructure

	responderRef      string
	requestMessageRef string

	responseTimestamp time.Time

	responseStatus []*XMLTerminationResponseStatus
}

type XMLTerminationResponseStatus struct {
	XMLStructure

	subscriberRef     string
	subscriptionRef   string
	responseTimestamp time.Time

	status           Bool
	errorType        string
	errorNumber      int
	errorText        string
	errorDescription string
}

type SIRIDeleteSubscriptionResponse struct {
	ResponderRef      string
	RequestMessageRef string
	ResponseTimestamp time.Time

	ResponseStatus []*SIRITerminationResponseStatus
}

type SIRITerminationResponseStatus struct {
	SubscriberRef     string
	SubscriptionRef   string
	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber string
	ErrorText   string
}

const deleteSubscriptionResponseTemplate = `<sw:DeleteSubscriptionResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<DeleteSubscriptionAnswerInfo>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ResponderRef>{{ .ResponderRef }}</siri:ResponderRef>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
	</DeleteSubscriptionAnswerInfo>
	<Answer>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ResponderRef>{{ .ResponderRef }}</siri:ResponderRef>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>{{ range .ResponseStatus }}
		<siri:TerminationResponseStatus>
			<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
			<siri:SubscriberRef>{{ .SubscriberRef }}</siri:SubscriberRef>
			<siri:SubscriptionRef>{{ .SubscriptionRef }}</siri:SubscriptionRef>
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
			<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<siri:OtherError number="{{.ErrorNumber}}">{{ else }}
				<siri:{{.ErrorType}}>{{ end }}
					<siri:ErrorText>{{.ErrorText}}</siri:ErrorText>
				</siri:{{.ErrorType}}>
			</siri:ErrorCondition>{{ end }}
		</siri:TerminationResponseStatus>{{ end }}
	</Answer>
	<AnswerExtension/>
</sw:DeleteSubscriptionResponse>`

func NewXMLDeleteSubscriptionResponse(node xml.Node) *XMLDeleteSubscriptionResponse {
	xmlDeleteSubscriptionResponse := &XMLDeleteSubscriptionResponse{}
	xmlDeleteSubscriptionResponse.node = NewXMLNode(node)
	return xmlDeleteSubscriptionResponse
}

func NewXMLDeleteSubscriptionResponseFromContent(content []byte) (*XMLDeleteSubscriptionResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLDeleteSubscriptionResponse(doc.Root().XmlNode)
	return request, nil
}

func NewXMLTerminationResponseStatus(node XMLNode) *XMLTerminationResponseStatus {
	responseStatus := &XMLTerminationResponseStatus{}
	responseStatus.node = node
	return responseStatus
}

func (response *XMLDeleteSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent("ResponderRef")
	}
	return response.responderRef
}

func (response *XMLDeleteSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLDeleteSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLDeleteSubscriptionResponse) ResponseStatus() []*XMLTerminationResponseStatus {
	if len(response.responseStatus) == 0 {
		nodes := response.findNodes("TerminationResponseStatus")
		if nodes == nil {
			return response.responseStatus
		}
		for _, responseStatus := range nodes {
			response.responseStatus = append(response.responseStatus, NewXMLTerminationResponseStatus(responseStatus))
		}
	}
	return response.responseStatus
}

func (response *XMLTerminationResponseStatus) SubscriptionRef() string {
	if response.subscriptionRef == "" {
		response.subscriptionRef = response.findStringChildContent("SubscriptionRef")
	}
	return response.subscriptionRef
}

func (response *XMLTerminationResponseStatus) SubscriberRef() string {
	if response.subscriberRef == "" {
		response.subscriberRef = response.findStringChildContent("SubscriberRef")
	}
	return response.subscriberRef
}

func (response *XMLTerminationResponseStatus) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLTerminationResponseStatus) Status() bool {
	if !response.status.Defined {
		response.status.SetValue(response.findBoolChildContent("Status"))
	}
	return response.status.Value
}

func (response *XMLTerminationResponseStatus) ErrorType() string {
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

func (response *XMLTerminationResponseStatus) ErrorNumber() int {
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

func (response *XMLTerminationResponseStatus) ErrorText() string {
	if !response.Status() && response.errorText == "" {
		response.errorText = response.findStringChildContent("ErrorText")
	}
	return response.errorText
}

func (response *XMLTerminationResponseStatus) ErrorDescription() string {
	if !response.Status() && response.errorDescription == "" {
		response.errorDescription = response.findStringChildContent("Description")
	}
	return response.errorDescription
}

func (notify *SIRIDeleteSubscriptionResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var deleteSubscriptionResponse = template.Must(template.New("deleteSubscriptionResponseTemplate").Parse(deleteSubscriptionResponseTemplate))
	if err := deleteSubscriptionResponse.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
