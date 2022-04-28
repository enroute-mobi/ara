package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type ServiceRequestBroadcaster interface {
	HandleRequests(*siri.XMLSiriServiceRequest, *audit.BigQueryMessage) *siri.SIRIServiceResponse
}

type SIRIServiceRequestBroadcaster struct {
	clock.ClockConsumer

	connector
}

type SIRIServiceRequestBroadcasterFactory struct{}

func NewSIRIServiceRequestBroadcaster(partner *Partner) *SIRIServiceRequestBroadcaster {
	siriServiceRequestBroadcaster := &SIRIServiceRequestBroadcaster{}
	siriServiceRequestBroadcaster.partner = partner
	return siriServiceRequestBroadcaster
}

func (connector *SIRIServiceRequestBroadcaster) HandleRequests(request *siri.XMLSiriServiceRequest, message *audit.BigQueryMessage) *siri.SIRIServiceResponse {
	logStashEvent := connector.newLogStashEvent("")
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLSiriServiceRequest(logStashEvent, request)

	response := &siri.SIRIServiceResponse{
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
		Status:                    true,
		RequestMessageRef:         request.MessageIdentifier(),
		ResponseTimestamp:         connector.Clock().Now(),
	}

	var stopIds, lineIds []string
	if smRequests := request.StopMonitoringRequests(); len(smRequests) != 0 {
		stopIds = connector.handleStopMonitoringRequests(smRequests, response)
	}
	if gmRequests := request.GeneralMessageRequests(); len(gmRequests) != 0 {
		connector.handleGeneralMessageRequests(gmRequests, response)
	}
	if ettRequests := request.EstimatedTimetableRequests(); len(ettRequests) != 0 {
		lineIds = connector.handleEstimatedTimetableRequests(ettRequests, response)
	}

	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier
	if !response.Status {
		message.Status = "Error"
	}
	message.StopAreas = stopIds
	message.Lines = lineIds

	logSIRIServiceResponse(logStashEvent, response)

	return response
}

func (connector *SIRIServiceRequestBroadcaster) handleStopMonitoringRequests(requests []*siri.XMLStopMonitoringRequest, response *siri.SIRIServiceResponse) (stopIds []string) {
	for _, stopMonitoringRequest := range requests {
		SMLogStashEvent := connector.newLogStashEvent("StopMonitoringRequestBroadcaster")
		logXMLStopMonitoringRequest(SMLogStashEvent, stopMonitoringRequest)
		SMLogStashEvent["siriType"] = "StopMonitoringDelivery for GetSiriServiceResponse"

		var delivery siri.SIRIStopMonitoringDelivery

		stopMonitoringConnector, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
		if !ok {
			delivery = siri.SIRIStopMonitoringDelivery{
				RequestMessageRef: stopMonitoringRequest.MessageIdentifier(),
				Status:            false,
				ResponseTimestamp: connector.Clock().Now(),
				ErrorType:         "CapabilityNotSupportedError",
				ErrorText:         "Can't find a StopMonitoringRequestBroadcaster connector",
			}
		} else {
			delivery = stopMonitoringConnector.(*SIRIStopMonitoringRequestBroadcaster).getStopMonitoringDelivery(SMLogStashEvent, stopMonitoringRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		logSIRIStopMonitoringDelivery(SMLogStashEvent, delivery)
		audit.CurrentLogStash().WriteEvent(SMLogStashEvent)

		response.StopMonitoringDeliveries = append(response.StopMonitoringDeliveries, &delivery)

		stopIds = append(stopIds, stopMonitoringRequest.MonitoringRef())
	}
	return
}

func (connector *SIRIServiceRequestBroadcaster) handleGeneralMessageRequests(requests []*siri.XMLGeneralMessageRequest, response *siri.SIRIServiceResponse) {
	for _, generalMessageRequest := range requests {
		GMLogStashEvent := connector.newLogStashEvent("GeneralMessageRequestBroadcaster")
		logXMLGeneralMessageRequest(GMLogStashEvent, generalMessageRequest)
		GMLogStashEvent["siriType"] = "GeneralMessageDelivery for GetSiriServiceResponse"

		var delivery siri.SIRIGeneralMessageDelivery

		generalMessageConnector, ok := connector.Partner().Connector(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER)
		if !ok {
			delivery = siri.SIRIGeneralMessageDelivery{
				RequestMessageRef: generalMessageRequest.MessageIdentifier(),
				Status:            false,
				ResponseTimestamp: connector.Clock().Now(),
				ErrorType:         "CapabilityNotSupportedError",
				ErrorText:         "Can't find a GeneralMessageRequestBroadcaster connector",
			}
		} else {
			delivery = generalMessageConnector.(*SIRIGeneralMessageRequestBroadcaster).getGeneralMessageDelivery(GMLogStashEvent, generalMessageRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		logSIRIGeneralMessageDelivery(GMLogStashEvent, delivery)
		audit.CurrentLogStash().WriteEvent(GMLogStashEvent)

		response.GeneralMessageDeliveries = append(response.GeneralMessageDeliveries, &delivery)
	}
}

func (connector *SIRIServiceRequestBroadcaster) handleEstimatedTimetableRequests(requests []*siri.XMLEstimatedTimetableRequest, response *siri.SIRIServiceResponse) (lineIds []string) {
	for _, estimatedTimetableRequest := range requests {
		ETTLogStashEvent := connector.newLogStashEvent("EstimatedTimetableRequestBroadcaster")
		logXMLEstimatedTimetableRequest(ETTLogStashEvent, estimatedTimetableRequest)
		ETTLogStashEvent["siriType"] = "EstimatedTimetableDelivery for GetSiriServiceResponse"

		var delivery siri.SIRIEstimatedTimetableDelivery

		estimatedTimetabeConnector, ok := connector.Partner().Connector(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER)
		if !ok {
			delivery = siri.SIRIEstimatedTimetableDelivery{
				RequestMessageRef: estimatedTimetableRequest.MessageIdentifier(),
				Status:            false,
				ResponseTimestamp: connector.Clock().Now(),
				ErrorType:         "CapabilityNotSupportedError",
				ErrorText:         "Can't find a EstimatedTimetableBroadcaster connector",
			}
		} else {
			delivery = estimatedTimetabeConnector.(*SIRIEstimatedTimetableBroadcaster).getEstimatedTimetableDelivery(estimatedTimetableRequest, ETTLogStashEvent)
		}

		if !delivery.Status {
			response.Status = false
		}

		audit.CurrentLogStash().WriteEvent(ETTLogStashEvent)

		response.EstimatedTimetableDeliveries = append(response.EstimatedTimetableDeliveries, &delivery)

		lineIds = append(lineIds, estimatedTimetableRequest.Lines()...)
	}
	return
}

func (connector *SIRIServiceRequestBroadcaster) newLogStashEvent(connectorName string) audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	if connectorName != "" {
		event["connector"] = fmt.Sprintf("%v for SIRIServiceRequestBroadcaster", connectorName)
		return event
	}
	event["connector"] = "SIRIServiceRequestBroadcaster"
	return event
}

func (factory *SIRIServiceRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIServiceRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIServiceRequestBroadcaster(partner)
}

func logXMLSiriServiceRequest(logStashEvent audit.LogStashEvent, request *siri.XMLSiriServiceRequest) {
	logStashEvent["siriType"] = "GetSIRIServiceResponse"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIServiceResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIServiceResponse) {
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
