package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type ServiceRequestBroadcaster interface {
	HandleRequests(*sxml.XMLSiriServiceRequest, *audit.BigQueryMessage) *siri.SIRIServiceResponse
}

type SIRIServiceRequestBroadcaster struct {
	state.Startable

	connector
}

type SIRIServiceRequestBroadcasterFactory struct{}

func NewSIRIServiceRequestBroadcaster(partner *Partner) *SIRIServiceRequestBroadcaster {
	siriServiceRequestBroadcaster := &SIRIServiceRequestBroadcaster{}
	siriServiceRequestBroadcaster.partner = partner
	return siriServiceRequestBroadcaster
}

func (connector *SIRIServiceRequestBroadcaster) Start() {
	stopMonitoringConnector, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
	if ok {
		stopMonitoringConnector.(*SIRIStopMonitoringRequestBroadcaster).Start()
	}

	estimatedTimetableConnector, ok := connector.Partner().Connector(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER)
	if ok {
		estimatedTimetableConnector.(*SIRIEstimatedTimetableRequestBroadcaster).Start()
	}
}

func (connector *SIRIServiceRequestBroadcaster) HandleRequests(request *sxml.XMLSiriServiceRequest, message *audit.BigQueryMessage) *siri.SIRIServiceResponse {
	response := &siri.SIRIServiceResponse{
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
		Status:                    true,
		RequestMessageRef:         request.MessageIdentifier(),
		ResponseTimestamp:         connector.Clock().Now(),
		LineRefs:                  make(map[string]struct{}),
		VehicleJourneyRefs:        make(map[string]struct{}),
		MonitoringRefs:            make(map[string]struct{}),
	}

	if smRequests := request.StopMonitoringRequests(); len(smRequests) != 0 {
		connector.handleStopMonitoringRequests(smRequests, response)
	}
	if gmRequests := request.GeneralMessageRequests(); len(gmRequests) != 0 {
		connector.handleGeneralMessageRequests(gmRequests, response)
	}
	if ettRequests := request.EstimatedTimetableRequests(); len(ettRequests) != 0 {
		connector.handleEstimatedTimetableRequests(ettRequests, response)
	}

	// log models
	message.StopAreas = GetModelReferenceSlice(response.MonitoringRefs)
	message.Lines = GetModelReferenceSlice(response.LineRefs)
	message.VehicleJourneys = GetModelReferenceSlice(response.VehicleJourneyRefs)

	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier
	if !response.Status {
		message.Status = "Error"
	}

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
		response.MonitoringRefs[stopMonitoringRequest.MonitoringRef()] = struct{}{}
		for line := range delivery.LineRefs {
			response.LineRefs[line] = struct{}{}
		}

		for vehicleJourney := range delivery.VehicleJourneyRefs {
			response.VehicleJourneyRefs[vehicleJourney] = struct{}{}
		}
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

func (connector *SIRIServiceRequestBroadcaster) handleEstimatedTimetableRequests(requests []*sxml.XMLEstimatedTimetableRequest, response *siri.SIRIServiceResponse) {
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
			delivery = estimatedTimetabeConnector.(*SIRIEstimatedTimetableRequestBroadcaster).getEstimatedTimetableDelivery(estimatedTimetableRequest)
		}

		if !delivery.Status {
			response.Status = false
		}

		response.EstimatedTimetableDeliveries = append(response.EstimatedTimetableDeliveries, &delivery)
		for _, line := range estimatedTimetableRequest.Lines() {
			response.LineRefs[line] = struct{}{}
		}
		for stopArea := range delivery.MonitoringRefs {
			response.MonitoringRefs[stopArea] = struct{}{}
		}
		for vehicleJourney := range delivery.VehicleJourneyRefs {
			response.VehicleJourneyRefs[vehicleJourney] = struct{}{}
		}
	}
}

func (factory *SIRIServiceRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIServiceRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIServiceRequestBroadcaster(partner)
}
