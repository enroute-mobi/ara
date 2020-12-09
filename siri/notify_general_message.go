package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyGeneralMessage struct {
	ResponseXMLStructure

	deliveries []*XMLGeneralMessageDelivery
}

type XMLGeneralMessageDelivery struct {
	SubscriptionDeliveryXMLStructure

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
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLGeneralMessageDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLGeneralMessageDelivery) XMLGeneralMessages() []*XMLGeneralMessage {
	if delivery.xmlGeneralMessages == nil {
		nodes := delivery.findNodes("GeneralMessage")
		for _, node := range nodes {
			delivery.xmlGeneralMessages = append(delivery.xmlGeneralMessages, NewXMLGeneralMessage(node))
		}
	}
	return delivery.xmlGeneralMessages
}

func (delivery *XMLGeneralMessageDelivery) XMLGeneralMessagesCancellations() []*XMLGeneralMessageCancellation {
	if delivery.xmlGeneralMessagesCancellations == nil {
		cancellations := []*XMLGeneralMessageCancellation{}
		nodes := delivery.findNodes("GeneralMessageCancellation")
		for _, node := range nodes {
			cancellations = append(cancellations, NewXMLCancelledGeneralMessage(node))
		}
		delivery.xmlGeneralMessagesCancellations = cancellations
	}
	return delivery.xmlGeneralMessagesCancellations
}

func (notify *SIRINotifyGeneralMessage) ErrorString() string {
	return fmt.Sprintf("%v: %v", notify.errorType(), notify.ErrorText)
}

func (notify *SIRINotifyGeneralMessage) errorType() string {
	if notify.ErrorType == "OtherError" {
		return fmt.Sprintf("%v %v", notify.ErrorType, notify.ErrorNumber)
	}
	return notify.ErrorType
}

func (notify *SIRINotifyGeneralMessage) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message_notify.template", notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
