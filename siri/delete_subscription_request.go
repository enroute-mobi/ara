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

const deleteSubscriptionRequestTemplate = `<ns1:DeleteSubscription xmlns:ns1="http://wsdl.siri.org.uk" xmlns:ns5="http://www.siri.org.uk/siri">
	<DeleteSubscriptionInfo
	 xmlns:ns2="http://www.ifopt.org.uk/acsb"
	 xmlns:ns3="http://www.ifopt.org.uk/ifopt"
	 xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
	 xmlns:ns6="http://wsdl.siri.org.uk/siri">
		<ns5:RequestTimestamp>{{ .RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:RequestTimestamp>
		<ns5:RequestorRef>{{.RequestorRef}}</ns5:RequestorRef>
		<ns5:MessageIdentifier>{{ .MessageIdentifier }}</ns5:MessageIdentifier>
	</DeleteSubscriptionInfo>
	<Request version="2.0:FR-IDF-2.4">{{ if .CancelAll }}
		<ns5:All/>{{ else }}
		<ns5:SubscriptionRef>{{.SubscriptionRef}}</ns5:SubscriptionRef>{{ end }}
	</Request>
	<RequestExtension/>
</ns1:DeleteSubscription>`

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
