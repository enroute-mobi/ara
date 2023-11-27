package core

import (
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type SituationExchangeRequestBroadcaster interface {
	Situations(*sxml.XMLGetSituationExchange, *audit.BigQueryMessage) *siri.SIRISituationExchangeResponse
}

type SIRISituationExchangeRequestBroadcaster struct {
	state.Startable

	connector
}

type SIRISituationExchangeRequestBroadcasterFactory struct{}

func NewSIRISituationExchangeRequestBroadcaster(partner *Partner) *SIRISituationExchangeRequestBroadcaster {
	connector := &SIRISituationExchangeRequestBroadcaster{}

	connector.partner = partner
	return connector
}

func (connector *SIRISituationExchangeRequestBroadcaster) Start() {
	connector.remoteObjectidKind = connector.partner.RemoteObjectIDKind(SIRI_SITUATION_EXCHANGE_REQUEST_BROADCASTER)
}

func (connector *SIRISituationExchangeRequestBroadcaster) Situations(request *sxml.XMLGetSituationExchange, message *audit.BigQueryMessage) *siri.SIRISituationExchangeResponse {
	response := &siri.SIRISituationExchangeResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
	}

	delivery := &siri.SIRISituationExchangeDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
		LineRefs:          make(map[string]struct{}),
		MonitoringRefs:    make(map[string]struct{}),
	}

	connector.getSituationExchangeDelivery(delivery, &request.XMLSituationExchangeRequest)

	message.Lines = GetModelReferenceSlice(delivery.LineRefs)
	message.StopAreas = GetModelReferenceSlice(delivery.MonitoringRefs)

	response.SIRISituationExchangeDelivery = *delivery
	return response
}

func (connector *SIRISituationExchangeRequestBroadcaster) getSituationExchangeDelivery(delivery *siri.SIRISituationExchangeDelivery, _ *sxml.XMLSituationExchangeRequest) {
	situations := connector.partner.Model().Situations().FindAll()
	for i := range situations {
		connector.buildSituation(delivery, situations[i])
	}
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildSituation(delivery *siri.SIRISituationExchangeDelivery, situation model.Situation) {
	if situation.Origin == string(connector.partner.Slug()) || situation.GMValidUntil().Before(connector.Clock().Now()) {
		return
	}

	var situationNumber string
	objectid, present := situation.ObjectID(connector.remoteObjectidKind)
	if present {
		situationNumber = objectid.Value()
	} else {
		objectid, present = situation.ObjectID("_default")
		if !present {
			logger.Log.Debugf("Unknown ObjectId for Situation %s", situation.Id())
			return
		}
		situationNumber = connector.Partner().ReferenceIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: objectid.Value()})
	}

	ptSituationElement := &siri.SIRIPtSituationElement{
		SituationNumber: situationNumber,
		CreationTime:    situation.RecordedAt,
		Version:         situation.Version,
		ValidityPeriods: situation.ValidityPeriods,
		Keywords:        strings.Join(situation.Keywords, " "),
		ReportType:      situation.ReportType,
	}

	if situation.Description != nil {
		ptSituationElement.Description = situation.Description.DefaultValue
	}

	if situation.Summary != nil {
		ptSituationElement.Summary = situation.Summary.DefaultValue
	}

	for _, affect := range situation.Affects {
		switch affect.GetType() {
		case model.SituationTypeStopArea:
			connector.buildAffectedStopArea(affect, ptSituationElement, delivery)
		case model.SituationTypeLine:
			connector.buildAffectedLine(affect, ptSituationElement, delivery)
		}
	}
	delivery.Situations = append(delivery.Situations, ptSituationElement)
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffectedStopArea(affect model.Affect, ptSituationElement *siri.SIRIPtSituationElement, delivery *siri.SIRISituationExchangeDelivery) {
	affect, _ = affect.(*model.AffectedStopArea)
	affectedStopAreaRef, ok := connector.resolveStopAreaRef(model.StopAreaId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown StopArea %s", affect.GetId())
		return
	}

	affectedStopPoint := siri.AffectedStopPoint{
		StopPointRef: affectedStopAreaRef,
	}

	ptSituationElement.AffectedStopPoints = append(ptSituationElement.AffectedStopPoints, &affectedStopPoint)

	// Logging
	delivery.MonitoringRefs[affectedStopAreaRef] = struct{}{}
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffectedLine(affect model.Affect, ptSituationElement *siri.SIRIPtSituationElement, delivery *siri.SIRISituationExchangeDelivery) {
	affect, _ = affect.(*model.AffectedLine)
	line, ok := connector.partner.Model().Lines().Find(model.LineId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown Line %s", affect.GetId())
		return
	}
	lineObjectId, ok := line.ObjectID(connector.remoteObjectidKind)
	if !ok {
		logger.Log.Debugf("Unknown Line ObjectId %s", connector.remoteObjectidKind)
		return
	}

	affectedLine := siri.AffectedLine{
		LineRef: lineObjectId.Value(),
	}
	delivery.LineRefs[lineObjectId.Value()] = struct{}{}

	for _, affectedDestination := range affect.(*model.AffectedLine).AffectedDestinations {
		affectedDestinationRef, ok := connector.resolveStopAreaRef(model.StopAreaId(affectedDestination.StopAreaId))
		if !ok {
			logger.Log.Debugf("Cannot find destination %s", affectedDestination.StopAreaId)
			continue
		}
		destination := &siri.SIRIAffectedDestination{
			StopPlaceRef: affectedDestinationRef,
		}
		delivery.MonitoringRefs[affectedDestinationRef] = struct{}{}
		affectedLine.Destinations = append(affectedLine.Destinations, *destination)
	}

	for _, affectedSection := range affect.(*model.AffectedLine).AffectedSections {
		firstStopRef, ok := connector.resolveStopAreaRef(model.StopAreaId(affectedSection.FirstStop))
		if !ok {
			logger.Log.Debugf("Cannot find firstStop  %s", affectedSection.FirstStop)
			continue
		}
		lastStopRef, ok := connector.resolveStopAreaRef(model.StopAreaId(affectedSection.LastStop))
		if !ok {
			logger.Log.Debugf("Cannot find lastStop  %s", affectedSection.LastStop)
			continue
		}
		section := &siri.SIRIAffectedSection{
			FirstStopPointRef: firstStopRef,
			LastStopPointRef:  lastStopRef,
		}
		delivery.MonitoringRefs[firstStopRef] = struct{}{}
		delivery.MonitoringRefs[lastStopRef] = struct{}{}
		affectedLine.Sections = append(affectedLine.Sections, *section)
	}

	for _, affectedRoute := range affect.(*model.AffectedLine).AffectedRoutes {
		route := &siri.SIRIAffectedRoute{
			RouteRef: affectedRoute.RouteRef,
		}
		affectedLine.Routes = append(affectedLine.Routes, *route)
	}

	ptSituationElement.AffectedLines = append(ptSituationElement.AffectedLines, &affectedLine)
}

func (connector *SIRISituationExchangeRequestBroadcaster) resolveStopAreaRef(stopAreaId model.StopAreaId) (string, bool) {
	stopArea, ok := connector.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return "", false
	}
	stopAreaObjectId, ok := stopArea.ReferentOrSelfObjectId(connector.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return stopAreaObjectId.Value(), true
}

func (factory *SIRISituationExchangeRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISituationExchangeRequestBroadcaster(partner)
}

func (factory *SIRISituationExchangeRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}
