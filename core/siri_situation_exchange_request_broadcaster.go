package core

import (
	"slices"
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
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(SIRI_SITUATION_EXCHANGE_REQUEST_BROADCASTER)
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
	if !connector.canBroadcast(situation) {
		return
	}

	var situationNumber string
	code, present := situation.Code(connector.remoteCodeSpace)
	if present {
		situationNumber = code.Value()
	} else {
		code, present = situation.Code("_default")
		if !present {
			logger.Log.Debugf("Unknown Code for Situation %s", situation.Id())
			return
		}
		situationNumber = connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: code.Value()})
	}

	ptSituationElement := &siri.SIRIPtSituationElement{
		SituationNumber:    situationNumber,
		CreationTime:       situation.RecordedAt,
		Version:            situation.Version,
		ValidityPeriods:    situation.ValidityPeriods,
		PublicationWindows: situation.PublicationWindows,
		Keywords:           strings.Join(situation.Keywords, " "),
		ReportType:         situation.ReportType,
		ParticipantRef:     situation.ParticipantRef,
		VersionedAtTime:    situation.VersionedAt,
		Progress:           situation.Progress,
		Reality:            situation.Reality,
		AlertCause:         situation.AlertCause,
		Severity:           situation.Severity,
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
			affectedStopArea, ok := connector.buildAffectedStopArea(affect, ptSituationElement, delivery)
			if ok {
				ptSituationElement.AffectedStopPoints = append(
					ptSituationElement.AffectedStopPoints,
					affectedStopArea,
				)
			}
		case model.SituationTypeLine:
			affectedLine, ok := connector.buildAffectedLine(affect, ptSituationElement, delivery)
			if ok {
				ptSituationElement.AffectedLines = append(
					ptSituationElement.AffectedLines,
					affectedLine,
				)
			}
		}
		if ptSituationElement.AffectedLines != nil || ptSituationElement.AffectedStopPoints != nil {
			ptSituationElement.HasAffects = true
		}
	}

	for _, consequence := range situation.Consequences {
		c := &siri.Consequence{
			Periods:  consequence.Periods,
			Severity: consequence.Severity,
		}
		for _, affect := range consequence.Affects {
			switch affect.GetType() {
			case model.SituationTypeStopArea:
				affectedStopArea, ok := connector.buildAffectedStopArea(affect, ptSituationElement, delivery)
				if ok {
					c.AffectedStopPoints = append(c.AffectedStopPoints, affectedStopArea)
				}
			case model.SituationTypeLine:
				affectedLine, ok := connector.buildAffectedLine(affect, ptSituationElement, delivery)
				if ok {
					c.AffectedLines = append(c.AffectedLines, affectedLine)
				}
			}
		}

		if c.AffectedLines != nil || c.AffectedStopPoints != nil {
			c.HasAffects = true
		}
		if consequence.Blocking != nil {
			c.Blocking = consequence.Blocking
		}

		ptSituationElement.Consequences = append(ptSituationElement.Consequences, c)
	}

	delivery.Situations = append(delivery.Situations, ptSituationElement)
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffectedStopArea(affect model.Affect, ptSituationElement *siri.SIRIPtSituationElement, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedStopPoint, bool) {
	affect, _ = affect.(*model.AffectedStopArea)

	affectedStopAreaRef, ok := connector.resolveStopAreaRef(model.StopAreaId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown StopArea %s", affect.GetId())
		return nil, false
	}

	// Logging
	delivery.MonitoringRefs[affectedStopAreaRef] = struct{}{}

	affectedStopPoint := &siri.AffectedStopPoint{StopPointRef: affectedStopAreaRef}
	for _, lineId := range affect.(*model.AffectedStopArea).LineIds {
		line, ok := connector.partner.Model().Lines().Find(lineId)
		if !ok {
			logger.Log.Debugf("Unknown Line %s", affect.GetId())
			continue
		}
		lineCode, ok := line.Code(connector.remoteCodeSpace)
		if !ok {
			logger.Log.Debugf("Unknown Line Code %s", connector.remoteCodeSpace)
			continue
		}
		affectedStopPoint.LineRefs = append(affectedStopPoint.LineRefs, lineCode.Value())
	}

	return affectedStopPoint, true
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffectedLine(affect model.Affect, ptSituationElement *siri.SIRIPtSituationElement, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedLine, bool) {
	affect, _ = affect.(*model.AffectedLine)
	line, ok := connector.partner.Model().Lines().Find(model.LineId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown Line %s", affect.GetId())
		return nil, false
	}
	lineCode, ok := line.Code(connector.remoteCodeSpace)
	if !ok {
		logger.Log.Debugf("Unknown Line Code %s", connector.remoteCodeSpace)
		return nil, false
	}

	affectedLine := siri.AffectedLine{
		LineRef: lineCode.Value(),
	}
	delivery.LineRefs[lineCode.Value()] = struct{}{}

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

		for _, stopArea := range affectedRoute.StopAreaIds {
			stopAreaRef, ok := connector.resolveStopAreaRef(stopArea)
			if ok {
				route.StopPointRefs = append(route.StopPointRefs, stopAreaRef)
			}
		}
		affectedLine.Routes = append(affectedLine.Routes, *route)
	}

	return &affectedLine, true
}

func (connector *SIRISituationExchangeRequestBroadcaster) resolveStopAreaRef(stopAreaId model.StopAreaId) (string, bool) {
	stopArea, ok := connector.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return "", false
	}
	stopAreaCode, ok := stopArea.ReferentOrSelfCode(connector.remoteCodeSpace)
	if !ok {
		return "", false
	}
	return stopAreaCode.Value(), true
}

func (connector *SIRISituationExchangeRequestBroadcaster) canBroadcast(situation model.Situation) bool {
	if situation.Origin == string(connector.partner.Slug()) {
		return false
	}

	if !situation.GMValidUntil().IsZero() &&
		situation.GMValidUntil().Before(connector.Clock().Now()) {
		return false
	}

	tagsToBroadcast := connector.partner.BroadcastSituationsInternalTags()
	if len(tagsToBroadcast) != 0 {
		for _, tag := range situation.InternalTags {
			if slices.Contains(tagsToBroadcast, tag) {
				return true
			}
		}
		return false
	}

	return true
}

func (factory *SIRISituationExchangeRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISituationExchangeRequestBroadcaster(partner)
}

func (factory *SIRISituationExchangeRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}
