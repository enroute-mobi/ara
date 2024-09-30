package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"golang.org/x/exp/maps"
)

type SituationExchangeUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner         *Partner
	remoteCodeSpace string

	MonitoringRefs map[string]struct{}
	LineRefs       map[string]struct{}
}

type affectedModels struct {
	affectedAllLines  *model.AffectedAllLines
	affectedLines     map[model.LineId]*model.AffectedLine
	affectedStopAreas map[model.StopAreaId]*model.AffectedStopArea
}

func NewSituationExchangeUpdateEventBuilder(partner *Partner) SituationExchangeUpdateEventBuilder {
	return SituationExchangeUpdateEventBuilder{
		partner:         partner,
		remoteCodeSpace: partner.RemoteCodeSpace(),
		MonitoringRefs:  make(map[string]struct{}),
		LineRefs:        make(map[string]struct{}),
	}
}

func (builder *SituationExchangeUpdateEventBuilder) SetSituationExchangeDeliveryUpdateEvents(events *CollectUpdateEvents, xmlSituationExchangeDelivery *sxml.XMLSituationExchangeDelivery, producerRef string) {
	for _, xmlSituation := range xmlSituationExchangeDelivery.Situations() {
		builder.buildSituationExchangeUpdateEvent(events, xmlSituation, producerRef)
	}
}

func (builder *SituationExchangeUpdateEventBuilder) SetSituationExchangeUpdateEvents(event *CollectUpdateEvents, xmlResponse *sxml.XMLSituationExchangeResponse) {
	xmlSituationExchangeDeliveries := xmlResponse.SituationExchangeDeliveries()
	if len(xmlSituationExchangeDeliveries) == 0 {
		return
	}

	for _, xmlSituationExchangeDelivery := range xmlSituationExchangeDeliveries {
		builder.SetSituationExchangeDeliveryUpdateEvents(event, xmlSituationExchangeDelivery, xmlResponse.ProducerRef())
	}
}

func (builder *SituationExchangeUpdateEventBuilder) buildSituationExchangeUpdateEvent(event *CollectUpdateEvents, xmlSituation *sxml.XMLPtSituationElement, producerRef string) {
	if len(xmlSituation.Affects()) == 0 {
		return
	}

	situationEvent := &model.SituationUpdateEvent{
		Origin:         string(builder.partner.Slug()),
		CreatedAt:      builder.Clock().Now(),
		RecordedAt:     xmlSituation.RecordedAtTime(),
		SituationCode:  model.NewCode(builder.remoteCodeSpace, xmlSituation.SituationNumber()),
		Version:        xmlSituation.Version(),
		ProducerRef:    producerRef,
		ParticipantRef: xmlSituation.ParticipantRef(),
		VersionedAt:    xmlSituation.VersionedAtTime(),
	}
	situationEvent.SetId(model.SituationUpdateRequestId(builder.NewUUID()))

	situationEvent.Keywords = append(situationEvent.Keywords, xmlSituation.Keywords()...)
	situationEvent.ReportType = model.ReportType(xmlSituation.ReportType())

	var progress model.SituationProgress
	if err := progress.FromString(xmlSituation.Progress()); err == nil {
		situationEvent.Progress = progress
	} else {
		logger.Log.Debugf("%v", err)
	}

	var severity model.SituationSeverity
	if err := severity.FromString(xmlSituation.Severity()); err == nil {
		situationEvent.Severity = severity
	} else {
		logger.Log.Debugf("%v", err)
	}

	var reality model.SituationReality
	if err := reality.FromString(xmlSituation.Reality()); err == nil {
		situationEvent.Reality = reality
	} else {
		logger.Log.Debugf("%v", err)
	}

	// Summary
	s := model.NewTranslatedStringFromMap(xmlSituation.Summaries())
	situationEvent.Summary = s

	// Description
	d := model.NewTranslatedStringFromMap(xmlSituation.Descriptions())
	situationEvent.Description = d

	var alertCause model.SituationAlertCause
	if err := alertCause.FromString(xmlSituation.AlertCause()); err == nil {
		situationEvent.AlertCause = alertCause
	} else {
		logger.Log.Debugf("%v", err)
	}

	situationEvent.InternalTags = append(situationEvent.InternalTags, builder.partner.CollectSituationsInternalTags()...)

	for _, validityPeriod := range xmlSituation.ValidityPeriods() {
		period := &model.TimeRange{
			StartTime: validityPeriod.StartTime(),
			EndTime:   validityPeriod.EndTime(),
		}

		situationEvent.ValidityPeriods = append(situationEvent.ValidityPeriods, period)
	}

	for _, publicationWindow := range xmlSituation.PublicationWindows() {
		window := &model.TimeRange{
			StartTime: publicationWindow.StartTime(),
			EndTime:   publicationWindow.EndTime(),
		}

		situationEvent.PublicationWindows = append(
			situationEvent.PublicationWindows,
			window)
	}

	for _, affect := range xmlSituation.Affects() {
		affectedModels := builder.buildAffect(affect)
		if affectedModels.affectedAllLines != nil {
			situationEvent.Affects = append(situationEvent.Affects, affectedModels.affectedAllLines)
		}
		for _, affectedLine := range affectedModels.affectedLines {
			situationEvent.Affects = append(situationEvent.Affects, affectedLine)
		}
		for _, affectedStopAreas := range affectedModels.affectedStopAreas {
			situationEvent.Affects = append(situationEvent.Affects, affectedStopAreas)
		}
	}

	for _, consequence := range xmlSituation.Consequences() {
		builder.setConsequence(situationEvent, consequence)
	}

	for _, xmlPublishToWeb := range xmlSituation.PublishToWebActions() {
		builder.setPublishToWebAction(situationEvent, xmlPublishToWeb)
	}

	for _, xmlPublishToMobile := range xmlSituation.PublishToMobileActions() {
		builder.setPublishToMobileAction(situationEvent, xmlPublishToMobile)
	}

	for _, xmlPublishToDisplay := range xmlSituation.PublishToDisplayActions() {
		builder.setPublishToDisplayAction(situationEvent, xmlPublishToDisplay)
	}

	for _, xmlInfoLink := range xmlSituation.InfoLinks() {
		builder.setInfoLink(situationEvent, xmlInfoLink)

	}
	event.Situations = append(event.Situations, situationEvent)
}

func (builder *SituationExchangeUpdateEventBuilder) setInfoLink(situationEvent *model.SituationUpdateEvent, xmlInfoLink *sxml.XMLInfoLink) {
	var infoLink model.InfoLink

	infoLink.Uri = xmlInfoLink.Uri()
	infoLink.Label = xmlInfoLink.Label()
	infoLink.ImageRef = xmlInfoLink.ImageRef()

	var linkContent model.SituationLinkContent
	if err := linkContent.FromString(xmlInfoLink.LinkContent()); err == nil {
		infoLink.LinkContent = string(linkContent)
	} else {
		logger.Log.Debugf("%v", err)
	}

	situationEvent.InfoLinks = append(situationEvent.InfoLinks, &infoLink)
}

func (builder *SituationExchangeUpdateEventBuilder) setPublishActionCommon(xmlCommon sxml.XMLCommonPublishingAction) (actionCommon model.ActionCommon) {
	actionCommon.Name = xmlCommon.Name()
	actionCommon.ActionType = xmlCommon.Type()
	actionCommon.Value = xmlCommon.Value()

	// Prompt
	actionCommon.Prompt = model.NewTranslatedStringFromMap(xmlCommon.Prompt())

	// scopeType
	var scopeType model.SituationScopeType
	if err := scopeType.FromString(xmlCommon.ScopeType()); err == nil {
		actionCommon.ScopeType = scopeType
	} else {
		logger.Log.Debugf("%v", err)
	}

	// affects
	for _, affect := range xmlCommon.Affects() {
		affectedModels := builder.buildAffect(affect)
		if affectedModels.affectedAllLines != nil {
			actionCommon.Affects = append(actionCommon.Affects, affectedModels.affectedAllLines)
		}
		for _, affectedLine := range affectedModels.affectedLines {
			actionCommon.Affects = append(actionCommon.Affects, affectedLine)
		}
		for _, affectedStopAreas := range affectedModels.affectedStopAreas {
			actionCommon.Affects = append(actionCommon.Affects, affectedStopAreas)
		}
	}

	// ActionStatus
	var actionStatus model.SituationActionStatus
	if err := actionStatus.FromString(xmlCommon.ActionStatus()); err == nil {
		actionCommon.ActionStatus = actionStatus
	} else {
		logger.Log.Debugf("%v", err)
	}
	// Description
	d := model.NewTranslatedStringFromMap(xmlCommon.Descriptions())
	actionCommon.Description = d

	// PublicationWindows
	for _, publicationWindow := range xmlCommon.PublicationWindows() {
		window := &model.TimeRange{
			StartTime: publicationWindow.StartTime(),
			EndTime:   publicationWindow.EndTime(),
		}

		actionCommon.PublicationWindows = append(
			actionCommon.PublicationWindows,
			window)
	}
	return
}

func (builder *SituationExchangeUpdateEventBuilder) setPublishToDisplayAction(situationEvent *model.SituationUpdateEvent, xmlPublishToDisplay *sxml.XMLPublishToDisplayAction) {
	da := &model.PublishToDisplayAction{}

	da.ActionCommon = builder.setPublishActionCommon(xmlPublishToDisplay.XMLCommonPublishingAction)

	if xmlPublishToDisplay.OnBoard() != nil {
		da.OnBoard = xmlPublishToDisplay.OnBoard()
	}

	if xmlPublishToDisplay.OnPlace() != nil {
		da.OnPlace = xmlPublishToDisplay.OnPlace()
	}

	situationEvent.PublishToDisplayActions = append(situationEvent.PublishToDisplayActions, da)
}

func (builder *SituationExchangeUpdateEventBuilder) setPublishToMobileAction(situationEvent *model.SituationUpdateEvent, xmlPublishToMobile *sxml.XMLPublishToMobileAction) {
	ma := &model.PublishToMobileAction{}

	ma.ActionCommon = builder.setPublishActionCommon(xmlPublishToMobile.XMLCommonPublishingAction)

	if xmlPublishToMobile.Incidents() != nil {
		ma.Incidents = xmlPublishToMobile.Incidents()
	}

	if xmlPublishToMobile.HomePage() != nil {
		ma.HomePage = xmlPublishToMobile.HomePage()
	}

	situationEvent.PublishToMobileActions = append(situationEvent.PublishToMobileActions, ma)
}

func (builder *SituationExchangeUpdateEventBuilder) setPublishToWebAction(situationEvent *model.SituationUpdateEvent, xmlPublishToWeb *sxml.XMLPublishToWebAction) {
	pa := &model.PublishToWebAction{}

	pa.ActionCommon = builder.setPublishActionCommon(xmlPublishToWeb.XMLCommonPublishingAction)

	if xmlPublishToWeb.Incidents() != nil {
		pa.Incidents = xmlPublishToWeb.Incidents()
	}

	if xmlPublishToWeb.HomePage() != nil {
		pa.HomePage = xmlPublishToWeb.HomePage()
	}

	if xmlPublishToWeb.Ticker() != nil {
		pa.Ticker = xmlPublishToWeb.Ticker()
	}

	if len(xmlPublishToWeb.SocialNetworks()) != 0 {
		pa.SocialNetworks = xmlPublishToWeb.SocialNetworks()
	}

	situationEvent.PublishToWebActions = append(situationEvent.PublishToWebActions, pa)
}

func (builder *SituationExchangeUpdateEventBuilder) setConsequence(situationEvent *model.SituationUpdateEvent, xmlConsequence *sxml.XMLConsequence) {
	consequence := &model.Consequence{}
	for _, xmlPeriod := range xmlConsequence.Periods() {
		period := &model.TimeRange{
			StartTime: xmlPeriod.StartTime(),
			EndTime:   xmlPeriod.EndTime(),
		}
		consequence.Periods = append(consequence.Periods, period)
	}

	var condition model.SituationCondition
	if err := condition.FromString(xmlConsequence.Condition()); err == nil {
		consequence.Condition = condition
	} else {
		logger.Log.Debugf("Condition: %v", err)
	}

	var severity model.SituationSeverity
	if err := severity.FromString(xmlConsequence.Severity()); err == nil {
		consequence.Severity = severity
	} else {
		logger.Log.Debugf("Consequence: %v", err)
	}

	for _, affect := range xmlConsequence.Affects() {
		affectedModels := builder.buildAffect(affect)
		if affectedModels.affectedAllLines != nil {
			consequence.Affects = append(consequence.Affects, affectedModels.affectedAllLines)
		}
		for _, affectedLine := range affectedModels.affectedLines {
			consequence.Affects = append(consequence.Affects, affectedLine)
		}
		for _, affectedStopAreas := range affectedModels.affectedStopAreas {
			consequence.Affects = append(consequence.Affects, affectedStopAreas)
		}
	}

	if xmlConsequence.HasBlocking() {
		blocking := &model.Blocking{
			JourneyPlanner: xmlConsequence.JourneyPlanner(),
			RealTime:       xmlConsequence.RealTime(),
		}
		consequence.Blocking = blocking
	}

	situationEvent.Consequences = append(situationEvent.Consequences, consequence)
}

func (builder *SituationExchangeUpdateEventBuilder) buildAffectedStopArea(xmlAffectedStopPoint *sxml.XMLAffectedStopPoint, affectedStopAreas map[model.StopAreaId]*model.AffectedStopArea) {
	stopPointRef := xmlAffectedStopPoint.StopPointRef()
	stopPointRefCode := model.NewCode(builder.remoteCodeSpace, stopPointRef)
	stopArea, ok := builder.partner.Model().StopAreas().FindByCode(stopPointRefCode)
	if !ok {
		return
	}
	affect := model.NewAffectedStopArea()
	affect.StopAreaId = stopArea.Id()

	for _, lineRef := range xmlAffectedStopPoint.LineRefs() {
		LineRefCode := model.NewCode(builder.remoteCodeSpace, lineRef)
		line, ok := builder.partner.Model().Lines().FindByCode(LineRefCode)
		if !ok {
			logger.Log.Debugf("Unknown Line with code %s", LineRefCode.String())
			continue
		}
		builder.LineRefs[LineRefCode.Value()] = struct{}{}
		affect.LineIds = append(affect.LineIds, line.Id())
	}

	affectedStopAreas[affect.StopAreaId] = affect

	// Logging
	builder.MonitoringRefs[stopPointRefCode.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) buildAffectedLine(lineRef string, affectedLines map[model.LineId]*model.AffectedLine) {
	LineRefCode := model.NewCode(builder.remoteCodeSpace, lineRef)
	line, ok := builder.partner.Model().Lines().FindByCode(LineRefCode)
	if !ok {
		return
	}
	affect := model.NewAffectedLine()
	affect.LineId = line.Id()
	affectedLines[affect.LineId] = affect
	builder.LineRefs[LineRefCode.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) buildAffectedRoute(lineId model.LineId, xmlAffectedRoute *sxml.XMLAffectedRoute, affectedLines map[model.LineId]*model.AffectedLine) {
	affectedRoute := &model.AffectedRoute{}
	if xmlAffectedRoute.RouteRef() != "" {
		affectedRoute.RouteRef = xmlAffectedRoute.RouteRef()
	}
	for _, xmlAffectedStopPoint := range xmlAffectedRoute.AffectedStopPoints() {
		stopPointRef := xmlAffectedStopPoint.StopPointRef()
		stopPointRefCode := model.NewCode(builder.remoteCodeSpace, stopPointRef)
		stopArea, ok := builder.partner.Model().StopAreas().FindByCode(stopPointRefCode)
		if !ok {
			continue

		}
		affectedRoute.StopAreaIds = append(affectedRoute.StopAreaIds, stopArea.Id())
		builder.MonitoringRefs[stopPointRefCode.Value()] = struct{}{}
	}

	if affectedRoute != (&model.AffectedRoute{}) {
		affectedLines[lineId].AffectedRoutes =
			append(affectedLines[lineId].AffectedRoutes, affectedRoute)
	}
}

func (builder *SituationExchangeUpdateEventBuilder) buildAffectedDestination(lineId model.LineId, destination string, affectedLines map[model.LineId]*model.AffectedLine) {
	destinationCode := model.NewCode(builder.remoteCodeSpace, destination)
	stopArea, ok := builder.partner.Model().StopAreas().FindByCode(destinationCode)
	if !ok {
		return
	}

	affectedDestination := model.AffectedDestination{StopAreaId: stopArea.Id()}
	affectedLines[lineId].AffectedDestinations =
		append(affectedLines[lineId].AffectedDestinations, &affectedDestination)

	// Logging
	builder.MonitoringRefs[destinationCode.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) buildAffectedSection(lineId model.LineId, section *sxml.XMLAffectedSection, affectedLines map[model.LineId]*model.AffectedLine) {
	firstStopRef := section.FirstStop()
	firstStopCode := model.NewCode(builder.remoteCodeSpace, firstStopRef)
	firstStopArea, ok := builder.partner.Model().StopAreas().FindByCode(firstStopCode)
	if !ok {
		return
	}
	lastStopRef := section.LastStop()
	lastStopCode := model.NewCode(builder.remoteCodeSpace, lastStopRef)
	lastStopArea, ok := builder.partner.Model().StopAreas().FindByCode(lastStopCode)
	if !ok {
		return
	}

	affectedSection := &model.AffectedSection{
		FirstStop: firstStopArea.Id(),
		LastStop:  lastStopArea.Id(),
	}

	// Fill already existing AffectedLine
	affectedLine, ok := affectedLines[lineId]
	if ok {
		affectedLine.AffectedSections = append(affectedLine.AffectedSections, affectedSection)
		builder.MonitoringRefs[firstStopCode.Value()] = struct{}{}
		builder.MonitoringRefs[lastStopCode.Value()] = struct{}{}
		return
	}

	// Logging
	builder.MonitoringRefs[firstStopCode.Value()] = struct{}{}
	builder.MonitoringRefs[lastStopCode.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) buildAffect(xmlAffect *sxml.XMLAffect) (affects *affectedModels) {
	models := affectedModels{
		affectedLines:     make(map[model.LineId]*model.AffectedLine),
		affectedStopAreas: make(map[model.StopAreaId]*model.AffectedStopArea),
	}

	for _, xmlAffectedNetwork := range xmlAffect.AffectedNetworks() {
		if xmlAffectedNetwork.AllLines() {
			models.affectedAllLines = model.NewAffectedAllLines()
			continue
		}
		for _, lineRef := range xmlAffectedNetwork.LineRefs() {
			builder.buildAffectedLine(lineRef, models.affectedLines)
		}

		if len(models.affectedLines) == 1 {
			// get the LineId
			lineId := maps.Keys(models.affectedLines)[0]

			for _, xmlAffectedRoute := range xmlAffectedNetwork.AffectedRoutes() {
				builder.buildAffectedRoute(lineId, xmlAffectedRoute, models.affectedLines)
			}
			for _, section := range xmlAffectedNetwork.AffectedSections() {
				builder.buildAffectedSection(lineId, section, models.affectedLines)
			}
			for _, destination := range xmlAffectedNetwork.AffectedDestinations() {
				builder.buildAffectedDestination(lineId, destination, models.affectedLines)
			}
		}
	}

	for _, xmlAffectedStopPoint := range xmlAffect.AffectedStopPoints() {
		builder.buildAffectedStopArea(xmlAffectedStopPoint, models.affectedStopAreas)
	}

	affects = &models
	return affects
}
