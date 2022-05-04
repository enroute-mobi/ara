package core

import (
	"fmt"

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

	connector
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

func (factory *TestCheckStatusClientFactory) Validate(apiPartner *APIPartner) {} // Always valid

func (factory *TestCheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewTestCheckStatusClient()
}

func NewSIRICheckStatusClient(partner *Partner) *SIRICheckStatusClient {
	siriCheckStatusClient := &SIRICheckStatusClient{}
	siriCheckStatusClient.partner = partner
	return siriCheckStatusClient
}

func (connector *SIRICheckStatusClient) Status() (PartnerStatus, error) {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	partnerStatus := PartnerStatus{}
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      connector.Partner().RequestorRef(),
		RequestTimestamp:  startTime,
		MessageIdentifier: connector.Partner().NewMessageIdentifier(),
	}

	logSIRICheckStatusRequest(message, request)

	response, err := connector.Partner().SIRIClient().CheckStatus(request)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during CheckStatus: %v", err)
		message.Status = "Error"
		message.ErrorDetails = e
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UNKNOWN
		return partnerStatus, err
	}

	logXMLCheckStatusResponse(message, response)

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

func (factory *SIRICheckStatusClientFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRICheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusClient(partner)
}

func logSIRICheckStatusRequest(message *audit.BigQueryMessage, request *siri.SIRICheckStatusRequest) {
	xml, err := request.BuildXML()
	if err != nil {
		return
	}

	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLCheckStatusResponse(message *audit.BigQueryMessage, response *siri.XMLCheckStatusResponse) {
	if !response.Status() {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	message.ResponseIdentifier = response.ResponseMessageIdentifier()
}
