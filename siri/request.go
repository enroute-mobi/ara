package siri

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/af83/edwig/api"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusRequest struct {
	node *xml.XmlNode

	messageIdentifier string
	requestorRef      string
	requestTimestamp  time.Time
}

type SIRICheckStatusRequest struct {
	api.UUIDConsumer

	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time
}

const SIRIRequestTemplate = `<ns7:CheckStatus xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
	<Request>
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}</ns2:RequestTimestamp>
		<ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
	</Request>
	<RequestExtension/>
</ns7:CheckStatus>`

func NewXMLCheckStatusRequest(node *xml.XmlNode) *XMLCheckStatusRequest {
	return &XMLCheckStatusRequest{node: node}
}

func NewXMLCheckStatusRequestFromContent(content []byte) (*XMLCheckStatusRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLCheckStatusRequest(doc.Root().XmlNode)
	finalizer := func(request *XMLCheckStatusRequest) {
		doc.Free()
	}
	runtime.SetFinalizer(request, finalizer)
	return request, nil
}

func NewSIRICheckStatusRequest(RequestorRef string, RequestTimestamp time.Time, MessageIdentifier string) *SIRICheckStatusRequest {
	return &SIRICheckStatusRequest{RequestorRef: RequestorRef, RequestTimestamp: RequestTimestamp, MessageIdentifier: MessageIdentifier}
}

// TODO : Handle errors
func (request *XMLCheckStatusRequest) RequestorRef() string {
	if request.requestorRef == "" {
		nodes, _ := request.node.Search("//*[local-name()='RequestorRef']")
		request.requestorRef = strings.TrimSpace(nodes[0].Content())
	}
	return request.requestorRef
}

// TODO : Handle errors
func (request *XMLCheckStatusRequest) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		nodes, _ := request.node.Search("//*[local-name()='RequestTimestamp']")
		t, _ := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		request.requestTimestamp = t
	}
	return request.requestTimestamp
}

// TODO : Handle errors
func (request *XMLCheckStatusRequest) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		nodes, _ := request.node.Search("//*[local-name()='MessageIdentifier']")
		request.messageIdentifier = strings.TrimSpace(nodes[0].Content())
	}
	return request.messageIdentifier
}

// TODO : Handle errors
func (request *SIRICheckStatusRequest) BuildXML() string {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(SIRIRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}

func (request *SIRICheckStatusRequest) GenerateMessageIdentifier() {
	request.MessageIdentifier = fmt.Sprintf("Edwig:Message::%s:LOC", request.NewUUID())
}
