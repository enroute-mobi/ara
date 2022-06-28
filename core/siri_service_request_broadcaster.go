package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type ServiceRequestBroadcaster interface {
	HandleRequests(*sxml.XMLSiriServiceRequest, *audit.BigQueryMessage) *siri.SIRIServiceResponse
}

type SIRIServiceRequestBroadcaster struct {
	connector
}

type SIRIServiceRequestBroadcasterFactory struct{}

func NewSIRIServiceRequestBroadcaster(partner *Partner) *SIRIServiceRequestBroadcaster {
	siriServiceRequestBroadcaster := &SIRIServiceRequestBroadcaster{}
	siriServiceRequestBroadcaster.partner = partner
	return siriServiceRequestBroadcaster
}

func (connector *SIRIServiceRequestBroadcaster) HandleRequests(request *sxml.XMLSiriServiceRequest, message *audit.BigQueryMessage) *siri.SIRIServiceResponse {
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

	return response
}

func (connector *SIRIServiceRequestBroadcaster) handleStopMonitoringRequests(requests []*sxml.XMLStopMonitoringRequest, response *siri.SIRIServiceResponse) (stopIds []string) {
	for _, stopMonitoringRequest := range requests {
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
			delivery = stopMonitoringConnector.(*SIRIStopMonitoringRequestBroadcaster).getStopMonitoringDelivery(stopMonitoringRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		response.StopMonitoringDeliveries = append(response.StopMonitoringDeliveries, &delivery)

		stopIds = append(stopIds, stopMonitoringRequest.MonitoringRef())
	}
	return
}

func (connector *SIRIServiceRequestBroadcaster) handleGeneralMessageRequests(requests []*sxml.XMLGeneralMessageRequest, response *siri.SIRIServiceResponse) {
	for _, generalMessageRequest := range requests {
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
			delivery = generalMessageConnector.(*SIRIGeneralMessageRequestBroadcaster).getGeneralMessageDelivery(generalMessageRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		response.GeneralMessageDeliveries = append(response.GeneralMessageDeliveries, &delivery)
	}
}

func (connector *SIRIServiceRequestBroadcaster) handleEstimatedTimetableRequests(requests []*sxml.XMLEstimatedTimetableRequest, response *siri.SIRIServiceResponse) (lineIds []string) {
	for _, estimatedTimetableRequest := range requests {
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
			delivery = estimatedTimetabeConnector.(*SIRIEstimatedTimetableBroadcaster).getEstimatedTimetableDelivery(estimatedTimetableRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		response.EstimatedTimetableDeliveries = append(response.EstimatedTimetableDeliveries, &delivery)

		lineIds = append(lineIds, estimatedTimetableRequest.Lines()...)
	}
	return
}

func (factory *SIRIServiceRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIServiceRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIServiceRequestBroadcaster(partner)
}
