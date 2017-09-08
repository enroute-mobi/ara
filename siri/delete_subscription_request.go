package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type SIRIDeleteSubscriptionRequest struct {
	RequestorRef     string
	RequestTimestamp time.Time

	MessageIdentifier string
	SubscriptionRef   string
	CancelAll         bool
}

type XMLDeleteSubscriptionRequest struct {
	RequestXMLStructure

	cancelAll       Bool
	subscriptionRef string
}

const deleteSubscriptionRequestTemplate = `<sw:DeleteSubscription xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<DeleteSubscriptionInfo>
		<siri:RequestTimestamp>{{ .RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{ .MessageIdentifier }}</siri:MessageIdentifier>
	</DeleteSubscriptionInfo>
	<Request version="2.0:FR-IDF-2.4">{{ if .CancelAll }}
		<siri:All/>{{ else }}
		<siri:SubscriptionRef>{{.SubscriptionRef}}</siri:SubscriptionRef>{{ end }}
	</Request>
</sw:DeleteSubscription>`

func NewXMLDeleteSubscriptionRequest(node xml.Node) *XMLDeleteSubscriptionRequest {
	xmlDeleteSubscriptionRequest := &XMLDeleteSubscriptionRequest{}
	xmlDeleteSubscriptionRequest.node = NewXMLNode(node)
	return xmlDeleteSubscriptionRequest
}

func NewXMLDeleteSubscriptionRequestFromContent(content []byte) (*XMLDeleteSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLDeleteSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLDeleteSubscriptionRequest) SubscriptionRef() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionRef
}

func (request *XMLDeleteSubscriptionRequest) CancelAll() bool {
	if !request.cancelAll.Defined {
		request.cancelAll.SetValue(request.containSelfClosing("All"))
	}
	return request.cancelAll.Value
}

func (request *SIRIDeleteSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var deleteSubscriptionRequest = template.Must(template.New("deleteSubscriptionRequestTemplate").Parse(deleteSubscriptionRequestTemplate))
	if err := deleteSubscriptionRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
