package core

import (
	"slices"
	"strings"
	"time"

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

func (connector *SIRISituationExchangeRequestBroadcaster) getSituationExchangeDelivery(delivery *siri.SIRISituationExchangeDelivery, request *sxml.XMLSituationExchangeRequest) {
	situations := connector.partner.Model().Situations().FindAll()
	requestPeriod := connector.getBroadcastPeriod(request)
	for i := range situations {
		connector.buildSituation(delivery, situations[i], requestPeriod)
	}
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildSituation(delivery *siri.SIRISituationExchangeDelivery, situation model.Situation, requestPeriod *model.TimeRange) {
	if !connector.canBroadcast(situation, requestPeriod) {
		return
	}

	var situationNumber string
	code, present := situation.Code(connector.remoteCodeSpace)
	if present {
		situationNumber = code.Value()
	} else {
		code, present = situation.Code(model.Default)
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

	connector.buildAffects(situation.Affects, &ptSituationElement.SIRIAffects, delivery)

	if ptSituationElement.AffectedLines != nil || ptSituationElement.AffectedStopPoints != nil || ptSituationElement.AffectedAllLines {
		ptSituationElement.HasAffects = true
	}

	for _, consequence := range situation.Consequences {
		c := &siri.Consequence{
			Periods:   consequence.Periods,
			Severity:  consequence.Severity,
			Condition: consequence.Condition,
		}

		connector.buildAffects(consequence.Affects, &c.SIRIAffects, delivery)

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

		connector.buildActionCommon(publishToWebAction.ActionCommon, &wa.SIRIPublishActionCommon, delivery)
		ptSituationElement.PublishToWebActions = append(ptSituationElement.PublishToWebActions, wa)
	}

	for _, publishToMobileAction := range situation.PublishToMobileActions {
		ma := &siri.PublishToMobileAction{}

		ma.Incidents = publishToMobileAction.Incidents
		ma.HomePage = publishToMobileAction.HomePage

		connector.buildActionCommon(publishToMobileAction.ActionCommon, &ma.SIRIPublishActionCommon, delivery)
		ptSituationElement.PublishToMobileActions = append(ptSituationElement.PublishToMobileActions, ma)
	}

	for _, publishToDisplayAction := range situation.PublishToDisplayActions {
		da := &siri.PublishToDisplayAction{}

		da.OnBoard = publishToDisplayAction.OnBoard
		da.OnPlace = publishToDisplayAction.OnPlace

		connector.buildActionCommon(publishToDisplayAction.ActionCommon, &da.SIRIPublishActionCommon, delivery)
		ptSituationElement.PublishToDisplayActions = append(ptSituationElement.PublishToDisplayActions, da)
	}

	if len(ptSituationElement.PublishToWebActions) != 0 ||
		len(ptSituationElement.PublishToMobileActions) != 0 ||
		len(ptSituationElement.PublishToDisplayActions) != 0 {
		ptSituationElement.HasPublishingActions = true
	}

	for _, infoLink := range situation.InfoLinks {
		connector.buildInfoLink(ptSituationElement, infoLink)
	}

	delivery.Situations = append(delivery.Situations, ptSituationElement)
}

func (connector *SIRISituationExchangeRequestBroadcaster) getBroadcastPeriod(request *sxml.XMLSituationExchangeRequest) *model.TimeRange {
	period := model.TimeRange{}

	period.StartTime = connector.Clock().Now().Add(-1 * time.Hour)
	if !request.StartTime().IsZero() {
		period.StartTime = request.StartTime()
	}

	if request.PreviewInterval() != 0 {
		period.EndTime = request.StartTime().Add(request.PreviewInterval())
	}

	return &period
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildInfoLink(ptSituationElement *siri.SIRIPtSituationElement, infoLink *model.InfoLink) {
	link := &siri.InfoLink{
		Uri:         infoLink.Uri,
		Label:       infoLink.Label,
		ImageRef:    infoLink.ImageRef,
		LinkContent: infoLink.LinkContent,
	}

	ptSituationElement.InfoLinks = append(ptSituationElement.InfoLinks, link)

	if len(ptSituationElement.InfoLinks) != 0 {
		ptSituationElement.HasInfoLinks = true
	}
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildActionCommon(actionCommon model.ActionCommon, siriActionCommon *siri.SIRIPublishActionCommon, delivery *siri.SIRISituationExchangeDelivery) {
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

	if siriActionCommon.AffectedLines != nil || siriActionCommon.AffectedStopPoints != nil || siriActionCommon.AffectedAllLines {
		siriActionCommon.HasAffects = true
	}

	if siriActionCommon.ScopeType != "" && siriActionCommon.HasAffects {
		siriActionCommon.HasPublishAtScope = true
	}
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffects(affects model.Affects, siriAffects *siri.SIRIAffects, delivery *siri.SIRISituationExchangeDelivery) {
	for _, affect := range affects {
		switch affect.GetType() {
		case model.SituationTypeStopArea:
			affectedStopArea, ok := connector.buildAffectedStopArea(affect, delivery)
			if ok {
				siriAffects.AffectedStopPoints = append(siriAffects.AffectedStopPoints, affectedStopArea)
			}
		case model.SituationTypeLine:
			affectedLine, ok := connector.buildAffectedLine(affect, delivery)
			if ok {
				siriAffects.AffectedLines = append(siriAffects.AffectedLines, affectedLine)
			}
		case model.SituationTypeAllLines:
			siriAffects.AffectedAllLines = true
		}
	}
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffectedStopArea(affect model.Affect, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedStopPoint, bool) {
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
		lineCode, ok := line.ReferentOrSelfCode(connector.remoteCodeSpace)
		if !ok {
			logger.Log.Debugf("Unknown Line Code %s", connector.remoteCodeSpace)
			continue
		}
		affectedStopPoint.LineRefs = append(affectedStopPoint.LineRefs, lineCode.Value())
	}

	return affectedStopPoint, true
}

func (connector *SIRISituationExchangeRequestBroadcaster) buildAffectedLine(affect model.Affect, delivery *siri.SIRISituationExchangeDelivery) (*siri.AffectedLine, bool) {
	affect, _ = affect.(*model.AffectedLine)
	line, ok := connector.partner.Model().Lines().Find(model.LineId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown Line %s", affect.GetId())
		return nil, false
	}
	lineCode, ok := line.ReferentOrSelfCode(connector.remoteCodeSpace)
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

func (connector *SIRISituationExchangeRequestBroadcaster) canBroadcast(situation model.Situation, requestPeriod *model.TimeRange) bool {
	if situation.Origin == string(connector.partner.Slug()) {
		return false
	}

	if !situation.BroadcastPeriod().Overlaps(requestPeriod) {
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
