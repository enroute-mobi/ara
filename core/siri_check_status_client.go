package core

import (
	"fmt"
	"strconv"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type CheckStatusClient interface {
	Status() (OperationnalStatus, error)
}

type TestCheckStatusClient struct {
	status OperationnalStatus
	Done   chan bool
}

type TestCheckStatusClientFactory struct{}

type SIRICheckStatusClient struct {
	model.ClockConsumer

	SIRIConnector
}

type SIRICheckStatusClientFactory struct{}

func NewTestCheckStatusClient() *TestCheckStatusClient {
	return &TestCheckStatusClient{
		status: OPERATIONNAL_STATUS_UP,
		Done:   make(chan bool, 1),
	}
}

func (connector *TestCheckStatusClient) Status() (OperationnalStatus, error) {
	connector.Done <- true
	return connector.status, nil
}

func (connector *TestCheckStatusClient) SetStatus(status OperationnalStatus) {
	connector.status = status
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

func (connector *SIRICheckStatusClient) Status() (OperationnalStatus, error) {
	logStashEvent := make(audit.LogStashEvent)
	startTime := connector.Clock().Now()

	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  startTime,
		MessageIdentifier: connector.SIRIPartner().NewMessageIdentifier(),
	}

	logRequest(logStashEvent, request)

	response, err := connector.SIRIPartner().SOAPClient().CheckStatus(request)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["response"] = fmt.Sprintf("Error during CheckStatus: %v", err)
		return OPERATIONNAL_STATUS_UNKNOWN, err
	}

	logResponse(logStashEvent, response)

	if response.Status() {
		return OPERATIONNAL_STATUS_UP, nil
	} else {
		return OPERATIONNAL_STATUS_DOWN, nil
	}
}

func (factory *SIRICheckStatusClientFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRICheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusClient(partner)
}

func logRequest(logStashEvent audit.LogStashEvent, request *siri.SIRICheckStatusRequest) {
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	logStashEvent["requestXML"] = request.BuildXML()
}

func logResponse(logStashEvent audit.LogStashEvent, response *siri.XMLCheckStatusResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime().String()
	logStashEvent["responseXML"] = response.RawXML()
}
