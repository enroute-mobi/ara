package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type CheckStatusServer interface {
	CheckStatus(*siri.XMLCheckStatusRequest, *audit.BigQueryMessage) (*siri.SIRICheckStatusResponse, error)
}

type SIRICheckStatusServer struct {
	clock.ClockConsumer

	connector
}

type SIRICheckStatusServerFactory struct{}

func NewSIRICheckStatusServer(partner *Partner) *SIRICheckStatusServer {
	siriCheckStatusServer := &SIRICheckStatusServer{}
	siriCheckStatusServer.partner = partner
	return siriCheckStatusServer
}

func (connector *SIRICheckStatusServer) CheckStatus(request *siri.XMLCheckStatusRequest, message *audit.BigQueryMessage) (*siri.SIRICheckStatusResponse, error) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLCheckStatusRequest(logStashEvent, request)

	response := &siri.SIRICheckStatusResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		RequestMessageRef:         request.MessageIdentifier(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
		Status:                    true,
		ResponseTimestamp:         connector.Clock().Now(),
		ServiceStartedTime:        connector.Partner().StartedAt(),
	}

	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	logSIRICheckStatusResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRICheckStatusServer) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "CheckStatusServer"
	return event
}

func (factory *SIRICheckStatusServerFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRICheckStatusServerFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusServer(partner)
}

func logXMLCheckStatusRequest(logStashEvent audit.LogStashEvent, request *siri.XMLCheckStatusRequest) {
	logStashEvent["siriType"] = "CheckStatusResponse"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRICheckStatusResponse(logStashEvent audit.LogStashEvent, response *siri.SIRICheckStatusResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
	}
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
