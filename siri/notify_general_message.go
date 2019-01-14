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

const generalMessageNotifyTemplate = `<sw:NotifyGeneralMessage xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{ .ProducerRef }}</siri:ProducerRef>{{ if .Address }}
		<siri:Address>{{ .Address }}</siri:Address>{{ end }}
		<siri:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Notification>
		<siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
			<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
			<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
			<siri:SubscriberRef>{{ .SubscriberRef }}</siri:SubscriberRef>
			<siri:SubscriptionRef>{{ .SubscriptionIdentifier }}</siri:SubscriptionRef>
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
			<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<siri:OtherError number="{{ .ErrorNumber }}">{{ else }}
				<siri:{{ .ErrorType }}>{{ end }}
					<siri:ErrorText>{{ .ErrorText }}</siri:ErrorText>
				</siri:{{ .ErrorType }}>
			</siri:ErrorCondition>{{ else }}{{ range .GeneralMessages }}
			{{ .BuildGeneralMessageXML }}{{ end }}{{ end }}
		 </siri:GeneralMessageDelivery>
	</Notification>
	<NotifyExtension />
</sw:NotifyGeneralMessage>`

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
