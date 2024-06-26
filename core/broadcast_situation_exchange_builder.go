package core

import (
	"slices"
	"strings"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type BroadcastSituationExchangeBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner         *Partner
	remoteCodeSpace string
	lineRef         map[string]struct{}
	stopPointRef    map[string]struct{}
}

func NewBroadcastSituationExchangeBuilder(partner *Partner, connector string) *BroadcastSituationExchangeBuilder {
	return &BroadcastSituationExchangeBuilder{
		partner:         partner,
		remoteCodeSpace: partner.RemoteCodeSpace(connector),
		lineRef:         make(map[string]struct{}),
		stopPointRef:    make(map[string]struct{}),
	}
}

func (builder *BroadcastSituationExchangeBuilder) BuildSituationExchange(situation model.Situation) (situationExchangeDelivery *siri.SIRISituationExchangeDelivery) {
	if !builder.canBroadcast(situation) {
		return nil
	}

	var situationNumber string
	code, present := situation.Code(builder.remoteCodeSpace)
	if present {
		situationNumber = code.Value()
	} else {
		code, present = situation.Code("_default")
		if !present {
			logger.Log.Debugf("Unknown Code for Situation %s", situation.Id())
			return nil
		}
		situationNumber = builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: code.Value()})
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

	delivery := &siri.SIRISituationExchangeDelivery{
		Status:         true,
		LineRefs:       make(map[string]struct{}),
		MonitoringRefs: make(map[string]struct{}),
	}

	for _, affect := range situation.Affects {
		switch affect.GetType() {
		case model.SituationTypeStopArea:
			affectedStopArea, ok := builder.buildAffectedStopArea(affect, ptSituationElement, delivery)
			if ok {
				ptSituationElement.AffectedStopPoints = append(
					ptSituationElement.AffectedStopPoints,
					affectedStopArea,
				)
			}
		case model.SituationTypeLine:
			affectedLine, ok := builder.buildAffectedLine(affect, ptSituationElement, delivery)
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
				affectedStopArea, ok := builder.buildAffectedStopArea(affect, ptSituationElement, delivery)
				if ok {
					c.AffectedStopPoints = append(c.AffectedStopPoints, affectedStopArea)
				}
			case model.SituationTypeLine:
				affectedLine, ok := builder.buildAffectedLine(affect, ptSituationElement, delivery)
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

	situationExchangeDelivery = delivery
	return
}

func (builder *BroadcastSituationExchangeBuilder) buildAffectedStopArea(affect model.Affect, ptSituationElement *siri.SIRIPtSituationElement, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedStopPoint, bool) {
	affect, _ = affect.(*model.AffectedStopArea)

	affectedStopAreaRef, ok := builder.resolveStopAreaRef(model.StopAreaId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown StopArea %s", affect.GetId())
		return nil, false
	}

	// Logging
	delivery.MonitoringRefs[affectedStopAreaRef] = struct{}{}

	affectedStopPoint := &siri.AffectedStopPoint{StopPointRef: affectedStopAreaRef}
	for _, lineId := range affect.(*model.AffectedStopArea).LineIds {
		line, ok := builder.partner.Model().Lines().Find(lineId)
		if !ok {
			logger.Log.Debugf("Unknown Line %s", affect.GetId())
			continue
		}
		lineCode, ok := line.Code(builder.remoteCodeSpace)
		if !ok {
			logger.Log.Debugf("Unknown Line Code %s", builder.remoteCodeSpace)
			continue
		}
		affectedStopPoint.LineRefs = append(affectedStopPoint.LineRefs, lineCode.Value())
	}

	return affectedStopPoint, true
}

func (builder *BroadcastSituationExchangeBuilder) buildAffectedLine(affect model.Affect, ptSituationElement *siri.SIRIPtSituationElement, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedLine, bool) {
	affect, _ = affect.(*model.AffectedLine)
	line, ok := builder.partner.Model().Lines().Find(model.LineId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown Line %s", affect.GetId())
		return nil, false
	}
	lineCode, ok := line.Code(builder.remoteCodeSpace)
	if !ok {
		logger.Log.Debugf("Unknown Line Code %s", builder.remoteCodeSpace)
		return nil, false
	}

	affectedLine := siri.AffectedLine{
		LineRef: lineCode.Value(),
	}

	delivery.LineRefs[lineCode.Value()] = struct{}{}

	for _, affectedDestination := range affect.(*model.AffectedLine).AffectedDestinations {
		affectedDestinationRef, ok := builder.resolveStopAreaRef(model.StopAreaId(affectedDestination.StopAreaId))
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
		firstStopRef, ok := builder.resolveStopAreaRef(model.StopAreaId(affectedSection.FirstStop))
		if !ok {
			logger.Log.Debugf("Cannot find firstStop  %s", affectedSection.FirstStop)
			continue
		}
		lastStopRef, ok := builder.resolveStopAreaRef(model.StopAreaId(affectedSection.LastStop))
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
			stopAreaRef, ok := builder.resolveStopAreaRef(stopArea)
			if ok {
				route.StopPointRefs = append(route.StopPointRefs, stopAreaRef)
			}
		}
		affectedLine.Routes = append(affectedLine.Routes, *route)
	}

	return &affectedLine, true
}

func (builder *BroadcastSituationExchangeBuilder) resolveStopAreaRef(stopAreaId model.StopAreaId) (string, bool) {
	stopArea, ok := builder.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return "", false
	}
	stopAreaCode, ok := stopArea.ReferentOrSelfCode(builder.remoteCodeSpace)
	if !ok {
		return "", false
	}
	return stopAreaCode.Value(), true
}

func (builder *BroadcastSituationExchangeBuilder) canBroadcast(situation model.Situation) bool {
	if situation.Origin == string(builder.partner.Slug()) {
		return false
	}

	if !situation.GMValidUntil().IsZero() &&
		situation.GMValidUntil().Before(builder.Clock().Now()) {
		return false
	}

	tagsToBroadcast := builder.partner.BroadcastSituationsInternalTags()
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
