package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
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
	model.ClockConsumer

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
	startTime := connector.Clock().Now()

	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	partnerStatus := PartnerStatus{}
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  startTime,
		MessageIdentifier: connector.Partner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
	}

	logSIRICheckStatusRequest(logStashEvent, request)

	response, err := connector.SIRIPartner().SOAPClient().CheckStatus(request)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = fmt.Sprintf("Error during CheckStatus: %v", err)
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UNKNOWN
		return partnerStatus, err
	}

	logXMLCheckStatusResponse(logStashEvent, response)

	partnerStatus.ServiceStartedAt = response.ServiceStartedTime()
	if response.Status() {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UP
		return partnerStatus, nil
	} else {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_DOWN
		return partnerStatus, nil
	}
}

func (connector *SIRICheckStatusClient) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "CheckStatusClient"
	return event
}

func (factory *SIRICheckStatusClientFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRICheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusClient(partner)
}

func logSIRICheckStatusRequest(logStashEvent audit.LogStashEvent, request *siri.SIRICheckStatusRequest) {
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
}

func logXMLCheckStatusResponse(logStashEvent audit.LogStashEvent, response *siri.XMLCheckStatusResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		logStashEvent["errorType"] = response.ErrorType()
		if response.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		}
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
	}
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime().String()
	logStashEvent["responseXML"] = response.RawXML()
}
