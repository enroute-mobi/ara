package siri

import (
	"bitbucket.org/enroute-mobi/ara/logger"
	"bytes"
	"fmt"
	"strings"
	"time"
)

type SIRIGeneralMessageResponse struct {
	SIRIGeneralMessageDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIGeneralMessageDelivery struct {
	RequestMessageRef string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	GeneralMessages []*SIRIGeneralMessage
}

type SIRIGeneralMessageCancellation struct {
	RecordedAtTime        time.Time
	ItemIdentifier        string
	InfoMessageIdentifier string
}

type SIRIGeneralMessage struct {
	RecordedAtTime        time.Time
	ValidUntilTime        time.Time
	ItemIdentifier        string
	InfoMessageIdentifier string
	FormatRef             string
	InfoMessageVersion    int
	InfoChannelRef        string

	AffectedRefs []*SIRIAffectedRef
	LineSections []*SIRILineSection
	Messages     []*SIRIMessage
}

type SIRIAffectedRef struct {
	Kind string
	Id   string
}

type SIRILineSection struct {
	FirstStop string
	LastStop  string
	LineRef   string
}

type SIRIMessage struct {
	Content             string
	Type                string
	NumberOfLines       int
	NumberOfCharPerLine int
}

func (response *SIRIGeneralMessageResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("general_message_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIGeneralMessageDelivery) ErrorString() string {
	return fmt.Sprintf("%v: %v", delivery.errorType(), delivery.ErrorText)
}

func (delivery *SIRIGeneralMessageDelivery) errorType() string {
	if delivery.ErrorType == "OtherError" {
		return fmt.Sprintf("%v %v", delivery.ErrorType, delivery.ErrorNumber)
	}
	return delivery.ErrorType
}

func (delivery *SIRIGeneralMessageDelivery) BuildGeneralMessageDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (message *SIRIGeneralMessage) BuildGeneralMessageXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message.template", message); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (message *SIRIGeneralMessageCancellation) BuildGeneralMessageCancellationXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message_cancellation.template", message); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
