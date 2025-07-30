package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type FacilityMonitoringRequestBroadcaster interface {
	RequestFacility(*sxml.XMLGetFacilityMonitoring, *audit.BigQueryMessage) *siri.SIRIFacilityMonitoringResponse
}

type SIRIFacilityMonitoringRequestBroadcaster struct {
	state.Startable
	connector
}

type SIRIFacilityMonitoringRequestBroadcasterFactory struct{}

func NewSIRIFacilityMonitoringRequestBroadcaster(partner *Partner) *SIRIFacilityMonitoringRequestBroadcaster {
	connector := &SIRIFacilityMonitoringRequestBroadcaster{}
	connector.partner = partner
	return connector
}

func (connector *SIRIFacilityMonitoringRequestBroadcaster) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(SIRI_FACILITY_MONITORING_REQUEST_BROADCASTER)
}
func (connector *SIRIFacilityMonitoringRequestBroadcaster) getFacilityMonitoringDelivery(request *sxml.XMLFacilityMonitoringRequest) siri.SIRIFacilityMonitoringDelivery {
	code := model.NewCode(connector.remoteCodeSpace, request.FacilityRef())
	facility, ok := connector.partner.Model().Facilities().FindByCode(code)
	if !ok {
		return siri.SIRIFacilityMonitoringDelivery{
			RequestMessageRef: request.MessageIdentifier(),
			Status:            false,
			ResponseTimestamp: connector.Clock().Now(),
			ErrorType:         "InvalidDataReferencesError",
			ErrorText:         fmt.Sprintf("Facility not found: '%s'", code.Value()),
		}
	}

	// if !facility.CollectedAlways {
	// 	facility.CollectedUntil = connector.Clock().Now().Add(15 * time.Minute)
	// 	logger.Log.Printf("Facility %s will be collected until %v", faciility.Id(), faciility.CollectedUntil)
	// 	faciility.Save()
	// }

	delivery := siri.SIRIFacilityMonitoringDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	delivery.FacilityRef = request.FacilityRef()
	delivery.FacilityStatus = string(facility.Status)

	// facilityBuilder := NewBroadcastFacilityMonitoringBuilder(connector.Partner(), SIRI_FACILITY_MONITORING_REQUEST_BROADCASTER)
	return delivery
}

func (connector *SIRIFacilityMonitoringRequestBroadcaster) RequestFacility(request *sxml.XMLGetFacilityMonitoring, message *audit.BigQueryMessage) *siri.SIRIFacilityMonitoringResponse {
	response := &siri.SIRIFacilityMonitoringResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
	}

	response.SIRIFacilityMonitoringDelivery = connector.getFacilityMonitoringDelivery(&request.XMLFacilityMonitoringRequest)

	if !response.SIRIFacilityMonitoringDelivery.Status {
		message.Status = "Error"
		message.ErrorDetails = response.SIRIFacilityMonitoringDelivery.ErrorString()
	}
	// message.Facilities = []string{request.MonitoringRef()}
	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	return response
}

func (factory *SIRIFacilityMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIFacilityMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIFacilityMonitoringRequestBroadcaster(partner)
}
