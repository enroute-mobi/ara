package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type StopMonitoringRequestBroadcaster interface {
	RequestStopArea(*sxml.XMLGetStopMonitoring, *audit.BigQueryMessage) *siri.SIRIStopMonitoringResponse
}

type SIRIStopMonitoringRequestBroadcaster struct {
	state.Startable
	connector
}

type SIRIStopMonitoringRequestBroadcasterFactory struct{}

func NewSIRIStopMonitoringRequestBroadcaster(partner *Partner) *SIRIStopMonitoringRequestBroadcaster {
	connector := &SIRIStopMonitoringRequestBroadcaster{}
	connector.partner = partner
	return connector
}

func (connector *SIRIStopMonitoringRequestBroadcaster) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
}
func (connector *SIRIStopMonitoringRequestBroadcaster) getStopMonitoringDelivery(request *sxml.XMLStopMonitoringRequest) siri.SIRIStopMonitoringDelivery {
	code := model.NewCode(connector.remoteCodeSpace, request.MonitoringRef())
	stopArea, ok := connector.partner.Model().StopAreas().FindByCode(code)
	if !ok {
		return siri.SIRIStopMonitoringDelivery{
			RequestMessageRef:  request.MessageIdentifier(),
			Status:             false,
			ResponseTimestamp:  connector.Clock().Now(),
			ErrorType:          "InvalidDataReferencesError",
			ErrorText:          fmt.Sprintf("StopArea not found: '%s'", code.Value()),
			MonitoringRef:      request.MonitoringRef(),
			VehicleJourneyRefs: make(map[string]struct{}),
		}
	}

	if !stopArea.CollectedAlways {
		stopArea.CollectedUntil = connector.Clock().Now().Add(15 * time.Minute)
		logger.Log.Printf("StopArea %s will be collected until %v", stopArea.Id(), stopArea.CollectedUntil)
		stopArea.Save()
	}

	delivery := siri.SIRIStopMonitoringDelivery{
		RequestMessageRef:  request.MessageIdentifier(),
		Status:             true,
		ResponseTimestamp:  connector.Clock().Now(),
		MonitoringRef:      request.MonitoringRef(),
		VehicleJourneyRefs: make(map[string]struct{}),
		LineRefs:           make(map[string]struct{}),
	}

	// Prepare StopVisit Selectors
	selectors := []model.StopVisitSelector{}
	if request.LineRef() != "" {
		lineIds := connector.partner.Model().Lines().FindFamilyFromCode(model.NewCode(connector.remoteCodeSpace, request.LineRef()))
		if len(lineIds) != 0 {
			selectors = append(selectors, model.StopVisitSelectorByLines(lineIds))
		}
	}
	if request.PreviewInterval() != 0 {
		duration := request.PreviewInterval()
		now := connector.Clock().Now()
		if !request.StartTime().IsZero() {
			now = request.StartTime()
		}
		selectors = append(selectors, model.StopVisitSelectorByTime(now, now.Add(duration)))
	}
	selector := model.CompositeStopVisitSelector(selectors)

	// Prepare Id Array
	var stopVisitArray []string

	// Initialize builder
	stopMonitoringBuilder := NewBroadcastStopMonitoringBuilder(connector.Partner(), SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
	stopMonitoringBuilder.StopVisitTypes = request.StopVisitTypes()
	stopMonitoringBuilder.MonitoringRef = request.MonitoringRef()

	// Find Descendants
	stopAreas := connector.partner.Model().StopAreas().FindFamily(stopArea.Id())

	// Fill StopVisits
	svs := connector.partner.Model().StopVisits().FindFollowingByStopAreaIds(stopAreas)
	for i := range svs {
		if svs[i].Origin == string(connector.Partner().Slug()) {
			continue
		}
		if request.MaximumStopVisits() > 0 && len(stopVisitArray) >= request.MaximumStopVisits() {
			break
		}
		if !selector(svs[i]) {
			continue
		}

		monitoredStopVisit := stopMonitoringBuilder.BuildMonitoredStopVisit(svs[i])
		if monitoredStopVisit == nil {
			continue
		}
		stopVisitArray = append(stopVisitArray, monitoredStopVisit.ItemIdentifier)
		delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, monitoredStopVisit)
		delivery.VehicleJourneyRefs[monitoredStopVisit.DatedVehicleJourneyRef] = struct{}{}
		delivery.LineRefs[monitoredStopVisit.LineRef] = struct{}{}

	}

	return delivery
}

func (connector *SIRIStopMonitoringRequestBroadcaster) RequestStopArea(request *sxml.XMLGetStopMonitoring, message *audit.BigQueryMessage) *siri.SIRIStopMonitoringResponse {
	response := &siri.SIRIStopMonitoringResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
	}

	response.SIRIStopMonitoringDelivery = connector.getStopMonitoringDelivery(&request.XMLStopMonitoringRequest)

	if !response.SIRIStopMonitoringDelivery.Status {
		message.Status = "Error"
		message.ErrorDetails = response.SIRIStopMonitoringDelivery.ErrorString()
	}
	message.Lines = GetModelReferenceSlice(response.LineRefs)
	message.StopAreas = []string{request.MonitoringRef()}
	message.VehicleJourneys = GetModelReferenceSlice(response.VehicleJourneyRefs)
	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	return response
}

func (factory *SIRIStopMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIStopMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestBroadcaster(partner)
}
