package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type CheckStatusClient interface {
	Status() (PartnerStatus, error)
}

type TestCheckStatusClient struct {
	partnerStatus PartnerStatus
	Done          chan bool
}

type TestCheckStatusClientFactory struct{}

type SIRICheckStatusClient struct {
	clock.ClockConsumer

	siriConnector
}

type SIRICheckStatusClientFactory struct{}

func NewTestCheckStatusClient() *TestCheckStatusClient {
	return &TestCheckStatusClient{
		partnerStatus: PartnerStatus{
			OperationnalStatus: OPERATIONNAL_STATUS_UP,
		},
		Done: make(chan bool, 1),
	}
}

func (connector *TestCheckStatusClient) Status() (PartnerStatus, error) {
	connector.Done <- true

	return connector.partnerStatus, nil
}

func (connector *TestCheckStatusClient) SetStatus(status OperationnalStatus) {
	connector.partnerStatus.OperationnalStatus = status
}

func (factory *TestCheckStatusClientFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestCheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewTestCheckStatusClient()
}

func NewSIRICheckStatusClient(partner *Partner) *SIRICheckStatusClient {
	siriCheckStatusClient := &SIRICheckStatusClient{}
	siriCheckStatusClient.partner = partner
	return siriCheckStatusClient
}

func (connector *SIRICheckStatusClient) Status() (PartnerStatus, error) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	partnerStatus := PartnerStatus{}
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  startTime,
		MessageIdentifier: connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier(),
	}

	logSIRICheckStatusRequest(logStashEvent, message, request)

	response, err := connector.SIRIPartner().SOAPClient().CheckStatus(request)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during CheckStatus: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = e
		message.Status = "Error"
		message.ErrorDetails = e
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UNKNOWN
		return partnerStatus, err
	}

	logXMLCheckStatusResponse(logStashEvent, message, response)

	partnerStatus.ServiceStartedAt = response.ServiceStartedTime()
	if response.Status() {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UP
		return partnerStatus, nil
	} else {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_DOWN
		return partnerStatus, nil
	}
}

func (connector *SIRICheckStatusClient) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "CheckStatusRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (connector *SIRICheckStatusClient) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "CheckStatusClient"
	return event
}

func (factory *SIRICheckStatusClientFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
	return ok
}

func (factory *SIRICheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusClient(partner)
}

func logSIRICheckStatusRequest(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, request *siri.SIRICheckStatusRequest) {
	logStashEvent["siriType"] = "CheckStatusRequest"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml

	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLCheckStatusResponse(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, response *siri.XMLCheckStatusResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		message.Status = "Error"
		logStashEvent["errorType"] = response.ErrorType()
		if response.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		}
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
		message.ErrorDetails = response.ErrorString()
	}
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime().String()
	logStashEvent["responseXML"] = response.RawXML()
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	message.ResponseIdentifier = response.ResponseMessageIdentifier()
}
