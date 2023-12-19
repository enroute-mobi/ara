package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SituationExchangeUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner            *Partner
	remoteObjectidKind string
	affectedLines      map[model.LineId]*model.AffectedLine

	MonitoringRefs map[string]struct{}
	LineRefs       map[string]struct{}
}

func NewSituationExchangeUpdateEventBuilder(partner *Partner) SituationExchangeUpdateEventBuilder {
	return SituationExchangeUpdateEventBuilder{
		partner:            partner,
		remoteObjectidKind: partner.RemoteObjectIDKind(),
		affectedLines:      make(map[model.LineId]*model.AffectedLine),

		MonitoringRefs: make(map[string]struct{}),
		LineRefs:       make(map[string]struct{}),
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
		Origin:            string(builder.partner.Slug()),
		CreatedAt:         builder.Clock().Now(),
		RecordedAt:        xmlSituation.RecordedAtTime(),
		SituationObjectID: model.NewObjectID(builder.remoteObjectidKind, xmlSituation.SituationNumber()),
		Version:           xmlSituation.Version(),
		ProducerRef:       producerRef,
	}
	situationEvent.SetId(model.SituationUpdateRequestId(builder.NewUUID()))

	situationEvent.Keywords = append(situationEvent.Keywords, xmlSituation.Keywords()...)
	situationEvent.ReportType = model.ReportType(xmlSituation.ReportType())
	situationEvent.Summary = &model.SituationTranslatedString{
		DefaultValue: xmlSituation.Summary(),
	}
	situationEvent.Description = &model.SituationTranslatedString{
		DefaultValue: xmlSituation.Description(),
	}

	for _, validityPeriod := range xmlSituation.ValidityPeriods() {
		period := &model.TimeRange{
			StartTime: validityPeriod.StartTime(),
			EndTime:   validityPeriod.EndTime(),
		}

		situationEvent.ValidityPeriods = append(situationEvent.ValidityPeriods, period)
	}

	for _, affect := range xmlSituation.Affects() {
		builder.setAffect(situationEvent, affect)
	}

	*event = append(*event, situationEvent)
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedStopArea(event *model.SituationUpdateEvent, stopPointRef string) {
	stopPointRefObjectId := model.NewObjectID(builder.remoteObjectidKind, stopPointRef)
	stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(stopPointRefObjectId)
	if !ok {
		return
	}
	affect := model.NewAffectedStopArea()
	affect.StopAreaId = stopArea.Id()

	event.Affects = append(event.Affects, affect)

	// Logging
	builder.MonitoringRefs[stopPointRefObjectId.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedLine(lineRef string) {
	LineRefObjectId := model.NewObjectID(builder.remoteObjectidKind, lineRef)
	line, ok := builder.partner.Model().Lines().FindByObjectId(LineRefObjectId)
	if !ok {
		return
	}
	affect := model.NewAffectedLine()
	affect.LineId = line.Id()
	builder.affectedLines[affect.LineId] = affect
	builder.LineRefs[LineRefObjectId.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedRoute(lineRef string, route string) {
	LineRefObjectId := model.NewObjectID(builder.remoteObjectidKind, lineRef)
	line, ok := builder.partner.Model().Lines().FindByObjectId(LineRefObjectId)
	if !ok {
		return
	}
	affectedRoute := model.AffectedRoute{RouteRef: route}
	builder.affectedLines[line.Id()].AffectedRoutes =
		append(builder.affectedLines[line.Id()].AffectedRoutes, &affectedRoute)
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedDestination(lineRef string, destination string) {
	LineRefObjectId := model.NewObjectID(builder.remoteObjectidKind, lineRef)
	line, ok := builder.partner.Model().Lines().FindByObjectId(LineRefObjectId)
	if !ok {
		return
	}

	destinationObjectId := model.NewObjectID(builder.remoteObjectidKind, destination)
	stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(destinationObjectId)
	if !ok {
		return
	}

	affectedDestination := model.AffectedDestination{StopAreaId: stopArea.Id()}
	builder.affectedLines[model.LineId(line.Id())].AffectedDestinations =
		append(builder.affectedLines[model.LineId(line.Id())].AffectedDestinations, &affectedDestination)

	// Logging
	builder.MonitoringRefs[destinationObjectId.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) setAffectedSection(lineRef string, section *sxml.XMLAffectedSection) {
	LineRefObjectId := model.NewObjectID(builder.remoteObjectidKind, lineRef)
	line, ok := builder.partner.Model().Lines().FindByObjectId(LineRefObjectId)
	if !ok {
		return
	}

	firstStopRef := section.FirstStop()
	firstStopObjectId := model.NewObjectID(builder.remoteObjectidKind, firstStopRef)
	firstStopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(firstStopObjectId)
	if !ok {
		return
	}
	lastStopRef := section.LastStop()
	lastStopObjectId := model.NewObjectID(builder.remoteObjectidKind, lastStopRef)
	lastStopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(lastStopObjectId)
	if !ok {
		return
	}

	affectedSection := &model.AffectedSection{
		FirstStop: firstStopArea.Id(),
		LastStop:  lastStopArea.Id(),
	}

	// Fill already existing AffectedLine if exists
	affectedLine, ok := builder.affectedLines[line.Id()]
	if ok {
		affectedLine.AffectedSections = append(affectedLine.AffectedSections, affectedSection)
		builder.MonitoringRefs[firstStopObjectId.Value()] = struct{}{}
		builder.MonitoringRefs[lastStopObjectId.Value()] = struct{}{}
		return
	}

	// otherwise create new AffectedLine
	affectedLine = model.NewAffectedLine()
	affectedLine.LineId = line.Id()
	affectedLine.AffectedSections = append(affectedLine.AffectedSections, affectedSection)
	builder.affectedLines[line.Id()] = affectedLine

	// Logging
	builder.LineRefs[LineRefObjectId.Value()] = struct{}{}
	builder.MonitoringRefs[firstStopObjectId.Value()] = struct{}{}
	builder.MonitoringRefs[lastStopObjectId.Value()] = struct{}{}
}

func (builder *SituationExchangeUpdateEventBuilder) setAffect(event *model.SituationUpdateEvent, xmlAffect *sxml.XMLAffect) {
	lineRef := xmlAffect.LineRef()

	builder.setAffectedLine(lineRef)
	if len(builder.affectedLines) != 0 {
		for _, route := range xmlAffect.AffectedRoutes() {
			builder.setAffectedRoute(lineRef, route)
		}
		for _, section := range xmlAffect.AffectedSections() {
			builder.setAffectedSection(lineRef, section)
		}
		for _, destination := range xmlAffect.AffectedDestinations() {
			builder.setAffectedDestination(lineRef, destination)
		}
	}

	for _, affectedLine := range builder.affectedLines {
		event.Affects = append(event.Affects, affectedLine)
	}

	for _, stopPointRef := range xmlAffect.StopPoints() {
		builder.setAffectedStopArea(event, stopPointRef)
	}

}
