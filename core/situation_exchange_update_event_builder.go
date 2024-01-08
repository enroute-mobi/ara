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

func NewSituationExchangeUpdateEventBuilder(partner *Partner) SituationExchangeUpdateEventBuilder {
	return SituationExchangeUpdateEventBuilder{
		partner:         partner,
		remoteCodeSpace: partner.RemoteCodeSpace(),
		MonitoringRefs:  make(map[string]struct{}),
		LineRefs:        make(map[string]struct{}),
	}
}

func (builder *SituationExchangeUpdateEventBuilder) SetSituationExchangeDeliveryUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *sxml.XMLSituationExchangeResponse) {
	xmlSituationExchangeDeliveries := xmlResponse.SituationExchangeDeliveries()
	if len(xmlSituationExchangeDeliveries) == 0 {
		return
	}

	for _, xmlSituationExchangeDelivery := range xmlSituationExchangeDeliveries {
		for _, xmlSituation := range xmlSituationExchangeDelivery.Situations() {
			builder.buildSituationExchangeUpdateEvent(event, xmlSituation, xmlResponse.ProducerRef())
		}

	}
}

func (builder *SituationExchangeUpdateEventBuilder) buildSituationExchangeUpdateEvent(event *[]*model.SituationUpdateEvent, xmlSituation *sxml.XMLPtSituationElement, producerRef string) {
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
		VersionedAt:    xmlSituation.VersionedAtime(),
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

	situationEvent.Summary = &model.SituationTranslatedString{
		DefaultValue: xmlSituation.Summary(),
	}
	situationEvent.Description = &model.SituationTranslatedString{
		DefaultValue: xmlSituation.Description(),
	}

	var alertCause model.SituationAlertCause
	if err := alertCause.FromString(xmlSituation.AlertCause()); err == nil {
		situationEvent.AlertCause = alertCause
	} else {
		logger.Log.Debugf("%v", err)
	}

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
		builder.setAffect(situationEvent, affect)
	}

	for _, consequence := range xmlSituation.Consequences() {
		builder.setConsequence(situationEvent, consequence)
	}

	*event = append(*event, situationEvent)
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
	situationEvent.Consequences = append(situationEvent.Consequences, consequence)
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedStopArea(event *model.SituationUpdateEvent, stopPointRef string) {
	stopPointRefCode := model.NewCode(builder.remoteCodeSpace, stopPointRef)
	stopArea, ok := builder.partner.Model().StopAreas().FindByCode(stopPointRefCode)
	if !ok {
		return
	}
	affect := model.NewAffectedStopArea()
	affect.StopAreaId = stopArea.Id()

	event.Affects = append(event.Affects, affect)

	// Logging
	builder.MonitoringRefs[stopPointRefCode.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedLine(lineRef string, affectedLines map[model.LineId]*model.AffectedLine) {
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

func (builder *SituationExchangeUpdateEventBuilder) setAffectedRoute(lineId model.LineId, route string, affectedLines map[model.LineId]*model.AffectedLine) {
	affectedRoute := model.AffectedRoute{RouteRef: route}
	affectedLines[lineId].AffectedRoutes =
		append(affectedLines[lineId].AffectedRoutes, &affectedRoute)
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedDestination(lineId model.LineId, destination string, affectedLines map[model.LineId]*model.AffectedLine) {
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

func (builder *SituationExchangeUpdateEventBuilder) setAffectedSection(lineId model.LineId, section *sxml.XMLAffectedSection, affectedLines map[model.LineId]*model.AffectedLine) {
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

func (builder *SituationExchangeUpdateEventBuilder) setAffect(event *model.SituationUpdateEvent, xmlAffect *sxml.XMLAffect) {
	affectedLines := make(map[model.LineId]*model.AffectedLine)
	for _, lineRef := range xmlAffect.LineRefs() {
		builder.setAffectedLine(lineRef, affectedLines)
	}

	if len(affectedLines) == 1 {
		// get the LineId
		lineId := maps.Keys(affectedLines)[0]

		for _, route := range xmlAffect.AffectedRoutes() {
			builder.setAffectedRoute(lineId, route, affectedLines)
		}
		for _, section := range xmlAffect.AffectedSections() {
			builder.setAffectedSection(lineId, section, affectedLines)
		}
		for _, destination := range xmlAffect.AffectedDestinations() {
			builder.setAffectedDestination(lineId, destination, affectedLines)
		}
	}

	// Fill affectedLines
	for _, affectedLine := range affectedLines {
		event.Affects = append(event.Affects, affectedLine)
	}

	// Fill affectedStopAreas
	for _, stopPointRef := range xmlAffect.StopPoints() {
		builder.setAffectedStopArea(event, stopPointRef)
	}

}
