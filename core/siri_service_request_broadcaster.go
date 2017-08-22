package core

import (
	"fmt"
	"strconv"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type ServiceRequestBroadcaster interface {
	HandleRequests(request *siri.XMLSiriServiceRequest) *siri.SIRIServiceResponse
}

type SIRIServiceRequestBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRIServiceRequestBroadcasterFactory struct{}

func NewSIRIServiceRequestBroadcaster(partner *Partner) *SIRIServiceRequestBroadcaster {
	siriServiceRequestBroadcaster := &SIRIServiceRequestBroadcaster{}
	siriServiceRequestBroadcaster.partner = partner
	return siriServiceRequestBroadcaster
}

func (connector *SIRIServiceRequestBroadcaster) HandleRequests(request *siri.XMLSiriServiceRequest) *siri.SIRIServiceResponse {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent("")
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLSiriServiceRequest(logStashEvent, request)

	response := &siri.SIRIServiceResponse{
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
		Status:            true,
		RequestMessageRef: request.MessageIdentifier(),
		ResponseTimestamp: connector.Clock().Now(),
	}

	if smRequests := request.StopMonitoringRequests(); len(smRequests) != 0 {
		connector.handleStopMonitoringRequests(tx, smRequests, response)
	}
	if gmRequests := request.GeneralMessageRequests(); len(gmRequests) != 0 {
		connector.handleGeneralMessageRequests(tx, gmRequests, response)
	}
	if ettRequests := request.EstimatedTimetableRequests(); len(ettRequests) != 0 {
		connector.handleEstimatedTimetableRequests(tx, ettRequests, response)
	}

	logSIRIServiceResponse(logStashEvent, response)

	return response
}

func (connector *SIRIServiceRequestBroadcaster) handleStopMonitoringRequests(tx *model.Transaction, requests []*siri.XMLStopMonitoringRequest, response *siri.SIRIServiceResponse) {
	for _, stopMonitoringRequest := range requests {
		SMLogStashEvent := connector.newLogStashEvent("StopMonitoringRequestBroadcaster")
		logXMLStopMonitoringRequest(SMLogStashEvent, stopMonitoringRequest)

		var delivery siri.SIRIStopMonitoringDelivery

		stopMonitoringConnector, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
		if !ok {
			delivery = siri.SIRIStopMonitoringDelivery{
				RequestMessageRef: stopMonitoringRequest.MessageIdentifier(),
				Status:            false,
				ResponseTimestamp: connector.Clock().Now(),
				ErrorType:         "NotFound",
				ErrorText:         "Can't find a SIRIStopMonitoringRequestBroadcaster connector",
			}
		} else {
			delivery = stopMonitoringConnector.(*SIRIStopMonitoringRequestBroadcaster).getStopMonitoringDelivery(tx, SMLogStashEvent, stopMonitoringRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		logSIRIStopMonitoringDelivery(SMLogStashEvent, delivery)
		audit.CurrentLogStash().WriteEvent(SMLogStashEvent)

		response.StopMonitoringDeliveries = append(response.StopMonitoringDeliveries, &delivery)
	}
}

func (connector *SIRIServiceRequestBroadcaster) handleGeneralMessageRequests(tx *model.Transaction, requests []*siri.XMLGeneralMessageRequest, response *siri.SIRIServiceResponse) {
	for _, generalMessageRequest := range requests {
		GMLogStashEvent := connector.newLogStashEvent("GeneralMessageRequestBroadcaster")
		logXMLGeneralMessageRequest(GMLogStashEvent, generalMessageRequest)

		var delivery siri.SIRIGeneralMessageDelivery

		generalMessageConnector, ok := connector.Partner().Connector(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER)
		if !ok {
			delivery = siri.SIRIGeneralMessageDelivery{
				RequestMessageRef: generalMessageRequest.MessageIdentifier(),
				Status:            false,
				ResponseTimestamp: connector.Clock().Now(),
				ErrorType:         "NotFound",
				ErrorText:         "Can't find a SIRIGeneralMessageRequestBroadcaster connector",
			}
		} else {
			delivery = generalMessageConnector.(*SIRIGeneralMessageRequestBroadcaster).getGeneralMessageDelivery(tx, GMLogStashEvent, generalMessageRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		logSIRIGeneralMessageDelivery(GMLogStashEvent, delivery)
		audit.CurrentLogStash().WriteEvent(GMLogStashEvent)

		response.GeneralMessageDeliveries = append(response.GeneralMessageDeliveries, &delivery)
	}
}

func (connector *SIRIServiceRequestBroadcaster) handleEstimatedTimetableRequests(tx *model.Transaction, requests []*siri.XMLEstimatedTimetableRequest, response *siri.SIRIServiceResponse) {
	for _, estimatedTimetableRequest := range requests {
		ETTLogStashEvent := connector.newLogStashEvent("EstimatedTimetableRequestBroadcaster")
		logXMLEstimatedTimetableRequest(ETTLogStashEvent, estimatedTimetableRequest)

		var delivery siri.SIRIEstimatedTimetableDelivery

		estimatedTimetabeConnector, ok := connector.Partner().Connector(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER)
		if !ok {
			delivery = siri.SIRIEstimatedTimetableDelivery{
				RequestMessageRef: estimatedTimetableRequest.MessageIdentifier(),
				Status:            false,
				ResponseTimestamp: connector.Clock().Now(),
				ErrorType:         "NotFound",
				ErrorText:         "Can't find a SIRIEstimatedTimetableBroadcaster connector",
			}
		} else {
			delivery = estimatedTimetabeConnector.(*SIRIEstimatedTimetableBroadcaster).getEstimatedTimetableDelivery(tx, estimatedTimetableRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		logSIRIEstimatedTimetableDelivery(ETTLogStashEvent, delivery)
		audit.CurrentLogStash().WriteEvent(ETTLogStashEvent)

		response.EstimatedTimetableDeliveries = append(response.EstimatedTimetableDeliveries, &delivery)
	}
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

func (factory *SIRIServiceRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIServiceRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIServiceRequestBroadcaster(partner)
}

func logXMLSiriServiceRequest(logStashEvent audit.LogStashEvent, request *siri.XMLSiriServiceRequest) {
	logStashEvent["Connector"] = "SIRIServiceRequestBroadcaster"
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
