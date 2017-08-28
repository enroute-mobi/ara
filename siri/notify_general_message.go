package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyGeneralMessage struct {
	ResponseXMLStructure

	deliveries []*XMLGeneralMessageDelivery
}

type XMLGeneralMessageDelivery struct {
	ResponseXMLStructure

	subscriptionRef string
	subscriberRef   string

	xmlGeneralMessages              []*XMLGeneralMessage
	xmlGeneralMessagesCancellations []*XMLGeneralMessageCancellation
}

type SIRINotifyGeneralMessage struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	SubscriberRef             string
	SubscriptionIdentifier    string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	GeneralMessages []*SIRIGeneralMessage
}

const generalMessageNotifyTemplate = `<ns8:NotifyGeneralMessage xmlns:ns3="http://www.siri.org.uk/siri"
															 xmlns:ns5="http://www.ifopt.org.uk/ifopt"
															 xmlns:ns4="http://www.ifopt.org.uk/acsb"
															 xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
															 xmlns:ns7="http://scma/siri"
															 xmlns:ns8="http://wsdl.siri.org.uk">
	<ServiceDeliveryInfo>
		<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
		<ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>{{ if .Address }}
		<ns3:Address>{{ .Address }}</ns3:Address>{{ end }}
		<ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Notification xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
		<ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
			<ns5:RequestMessageRef>{{ .RequestMessageRef }}</ns5:RequestMessageRef>
			<ns5:SubscriberRef>{{ .SubscriberRef }}</ns5:SubscriberRef>
			<ns5:SubscriptionRef>{{ .SubscriptionIdentifier }}</ns5:SubscriptionRef>
			<ns3:Status>{{ .Status }}</ns3:Status>{{ if not .Status }}
			<ns3:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<ns3:OtherError number="{{ .ErrorNumber }}">{{ else }}
				<ns3:{{ .ErrorType }}>{{ end }}
					<ns3:ErrorText>{{ .ErrorText }}</ns3:ErrorText>
				</ns3:{{ .ErrorType }}>
			</ns3:ErrorCondition>{{ else }}{{ range .GeneralMessages }}
			{{ .BuildGeneralMessageXML }}{{ end }}{{ end }}
		 </ns3:GeneralMessageDelivery>
		</Notification>
 <NotifyExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
</ns8:NotifyGeneralMessage>`

func NewXMLNotifyGeneralMessage(node xml.Node) *XMLNotifyGeneralMessage {
	xmlGeneralMessageResponse := &XMLNotifyGeneralMessage{}
	xmlGeneralMessageResponse.node = NewXMLNode(node)
	return xmlGeneralMessageResponse
}

func NewXMLNotifyGeneralMessageFromContent(content []byte) (*XMLNotifyGeneralMessage, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifyGeneralMessage(doc.Root().XmlNode)
	return response, nil
}

func NewXMLGeneralMessageDelivery(node XMLNode) *XMLGeneralMessageDelivery {
	delivery := &XMLGeneralMessageDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyGeneralMessage) GeneralMessagesDeliveries() []*XMLGeneralMessageDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLGeneralMessageDelivery{}
		nodes := notify.findNodes("GeneralMessageDelivery")
		if nodes != nil {
			for _, node := range nodes {
				deliveries = append(deliveries, NewXMLGeneralMessageDelivery(node))
			}
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLGeneralMessageDelivery) SubscriptionRef() string {
	if delivery.subscriptionRef == "" {
		delivery.subscriptionRef = delivery.findStringChildContent("SubscriptionRef")
	}
	return delivery.subscriptionRef
}

func (delivery *XMLGeneralMessageDelivery) SubscriberRef() string {
	if delivery.subscriberRef == "" {
		delivery.subscriberRef = delivery.findStringChildContent("SubscriberRef")
	}
	return delivery.subscriberRef
}

func (delivery *XMLGeneralMessageDelivery) XMLGeneralMessages() []*XMLGeneralMessage {
	if delivery.xmlGeneralMessages == nil {
		nodes := delivery.findNodes("GeneralMessage")
		if nodes != nil {
			for _, node := range nodes {
				delivery.xmlGeneralMessages = append(delivery.xmlGeneralMessages, NewXMLGeneralMessage(node))
			}
		}
	}
	return delivery.xmlGeneralMessages
}

func (delivery *XMLGeneralMessageDelivery) XMLGeneralMessagesCancellations() []*XMLGeneralMessageCancellation {
	if delivery.xmlGeneralMessagesCancellations == nil {
		cancellations := []*XMLGeneralMessageCancellation{}
		nodes := delivery.findNodes("GeneralMessageCancellation")
		if nodes != nil {
			for _, node := range nodes {
				cancellations = append(cancellations, NewXMLCancelledGeneralMessage(node))
			}
		}
		delivery.xmlGeneralMessagesCancellations = cancellations
	}
	return delivery.xmlGeneralMessagesCancellations
}

func (notify *SIRINotifyGeneralMessage) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var notifyDelivery = template.Must(template.New("generalMessageNotifyTemplate").Parse(generalMessageNotifyTemplate))
	if err := notifyDelivery.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
