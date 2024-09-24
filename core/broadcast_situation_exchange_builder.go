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

func (builder *BroadcastSituationExchangeBuilder) BuildSituationExchange(situation model.Situation, delivery *siri.SIRISituationExchangeDelivery) {
	if !builder.canBroadcast(situation) {
		return
	}

	var situationNumber string
	code, present := situation.Code(builder.remoteCodeSpace)
	if present {
		situationNumber = code.Value()
	} else {
		code, present = situation.Code(model.Default)
		if !present {
			logger.Log.Debugf("Unknown Code for Situation %s", situation.Id())
			return
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
		d := &siri.SIRITranslatedString{
			Tag:              "Description",
			TranslatedString: *situation.Description,
		}

		ptSituationElement.Description = d
	}

	if situation.Summary != nil {
		s := &siri.SIRITranslatedString{
			Tag:              "Summary",
			TranslatedString: *situation.Summary,
		}

		ptSituationElement.Summary = s
	}

	builder.buildAffects(situation.Affects, &ptSituationElement.SIRIAffects, delivery)

	if ptSituationElement.AffectedLines != nil || ptSituationElement.AffectedStopPoints != nil {
		ptSituationElement.HasAffects = true
	}

	for _, consequence := range situation.Consequences {
		c := &siri.Consequence{
			Periods:   consequence.Periods,
			Severity:  consequence.Severity,
			Condition: consequence.Condition,
		}

		builder.buildAffects(consequence.Affects, &c.SIRIAffects, delivery)

		if c.AffectedLines != nil || c.AffectedStopPoints != nil || c.AffectedAllLines {
			c.HasAffects = true
		}
		if consequence.Blocking != nil {
			c.Blocking = consequence.Blocking
		}

		ptSituationElement.Consequences = append(ptSituationElement.Consequences, c)
	}

	for _, publishToWebAction := range situation.PublishToWebActions {
		wa := &siri.PublishToWebAction{}

		wa.Incidents = publishToWebAction.Incidents
		wa.HomePage = publishToWebAction.HomePage

		wa.SocialNetworks = append(wa.SocialNetworks, publishToWebAction.SocialNetworks...)

		builder.buildActionCommon(publishToWebAction.ActionCommon, &wa.SIRIPublishActionCommon, delivery)
		ptSituationElement.PublishToWebActions = append(ptSituationElement.PublishToWebActions, wa)
	}

	for _, publishToMobileAction := range situation.PublishToMobileActions {
		ma := &siri.PublishToMobileAction{}

		ma.Incidents = publishToMobileAction.Incidents
		ma.HomePage = publishToMobileAction.HomePage

		builder.buildActionCommon(publishToMobileAction.ActionCommon, &ma.SIRIPublishActionCommon, delivery)
		ptSituationElement.PublishToMobileActions = append(ptSituationElement.PublishToMobileActions, ma)
	}

	for _, publishToDisplayAction := range situation.PublishToDisplayActions {
		da := &siri.PublishToDisplayAction{}

		da.OnBoard = publishToDisplayAction.OnBoard
		da.OnPlace = publishToDisplayAction.OnPlace

		builder.buildActionCommon(publishToDisplayAction.ActionCommon, &da.SIRIPublishActionCommon, delivery)
		ptSituationElement.PublishToDisplayActions = append(ptSituationElement.PublishToDisplayActions, da)
	}

	if len(ptSituationElement.PublishToWebActions) != 0 ||
		len(ptSituationElement.PublishToMobileActions) != 0 ||
		len(ptSituationElement.PublishToDisplayActions) != 0 {
		ptSituationElement.HasPublishingActions = true
	}

	delivery.Situations = append(delivery.Situations, ptSituationElement)
}

func (connector *BroadcastSituationExchangeBuilder) buildActionCommon(actionCommon model.ActionCommon, siriActionCommon *siri.SIRIPublishActionCommon, delivery *siri.SIRISituationExchangeDelivery) {
	siriActionCommon.Name = actionCommon.Name
	siriActionCommon.ActionType = actionCommon.ActionType
	siriActionCommon.Value = actionCommon.Value
	siriActionCommon.ScopeType = actionCommon.ScopeType
	siriActionCommon.ActionStatus = actionCommon.ActionStatus
	siriActionCommon.PublicationWindows = actionCommon.PublicationWindows

	if actionCommon.Prompt != nil {
		p := &siri.SIRITranslatedString{
			Tag:              "Prompt",
			TranslatedString: *actionCommon.Prompt,
		}
		siriActionCommon.Prompt = p
	}

	if actionCommon.Description != nil {
		d := &siri.SIRITranslatedString{
			Tag:              "Description",
			TranslatedString: *actionCommon.Description,
		}
		siriActionCommon.Description = d
	}

	connector.buildAffects(actionCommon.Affects, &siriActionCommon.SIRIAffects, delivery)

	if siriActionCommon.AffectedLines != nil || siriActionCommon.AffectedStopPoints != nil || siriActionCommon.AffectedAllLines{
		siriActionCommon.HasAffects = true
	}

	if siriActionCommon.ScopeType != "" && siriActionCommon.HasAffects {
		siriActionCommon.HasPublishAtScope = true
	}
}

func (builder *BroadcastSituationExchangeBuilder) buildAffects(affects model.Affects, siriAffects *siri.SIRIAffects, delivery *siri.SIRISituationExchangeDelivery) {
	for _, affect := range affects {
		switch affect.GetType() {
		case model.SituationTypeStopArea:
			affectedStopArea, ok := builder.buildAffectedStopArea(affect, delivery)
			if ok {
				siriAffects.AffectedStopPoints = append(siriAffects.AffectedStopPoints, affectedStopArea)
			}
		case model.SituationTypeLine:
			affectedLine, ok := builder.buildAffectedLine(affect, delivery)
			if ok {
				siriAffects.AffectedLines = append(siriAffects.AffectedLines, affectedLine)
			}
		case model.SituationTypeAllLines:
			siriAffects.AffectedAllLines = true
		}
	}
}

func (builder *BroadcastSituationExchangeBuilder) buildAffectedStopArea(affect model.Affect, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedStopPoint, bool) {
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

func (builder *BroadcastSituationExchangeBuilder) buildAffectedLine(affect model.Affect, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedLine, bool) {
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
