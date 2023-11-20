package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"golang.org/x/exp/maps"
)

type GeneralMessageUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner            *Partner
	remoteObjectidKind string
	affectedLines      map[model.LineId]*model.AffectedLine

	MonitoringRefs map[string]struct{}
	LineRefs       map[string]struct{}
}

type LineSection struct {
	LineRef   string
	FirstStop string
	LastStop  string
}

func NewGeneralMessageUpdateEventBuilder(partner *Partner) GeneralMessageUpdateEventBuilder {
	return GeneralMessageUpdateEventBuilder{
		partner:            partner,
		remoteObjectidKind: partner.RemoteObjectIDKind(),
		affectedLines:      make(map[model.LineId]*model.AffectedLine),

		MonitoringRefs: make(map[string]struct{}),
		LineRefs:       make(map[string]struct{}),
	}
}

func (builder *GeneralMessageUpdateEventBuilder) SetGeneralMessageDeliveryUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *sxml.XMLGeneralMessageDelivery, producerRef string) {
	xmlGeneralMessageEvents := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessageEvents) == 0 {
		return
	}

	for _, xmlGeneralMessageEvents := range xmlGeneralMessageEvents {
		builder.buildGeneralMessageUpdateEvent(event, xmlGeneralMessageEvents, producerRef)
	}
}

func (builder *GeneralMessageUpdateEventBuilder) SetGeneralMessageResponseUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *sxml.XMLGeneralMessageResponse) {
	xmlGeneralMessageEvents := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessageEvents) == 0 {
		return
	}

	for _, xmlGeneralMessageEvents := range xmlGeneralMessageEvents {
		builder.buildGeneralMessageUpdateEvent(event, xmlGeneralMessageEvents, xmlResponse.ProducerRef())
	}
}

func (builder *GeneralMessageUpdateEventBuilder) buildGeneralMessageUpdateEvent(event *[]*model.SituationUpdateEvent, xmlGeneralMessageEvent *sxml.XMLGeneralMessage, producerRef string) {
	if xmlGeneralMessageEvent.Content() == nil {
		return
	}

	situationEvent := &model.SituationUpdateEvent{
		Origin:            string(builder.partner.Slug()),
		CreatedAt:         builder.Clock().Now(),
		RecordedAt:        xmlGeneralMessageEvent.RecordedAtTime(),
		SituationObjectID: model.NewObjectID(builder.remoteObjectidKind, xmlGeneralMessageEvent.InfoMessageIdentifier()),
		Version:           xmlGeneralMessageEvent.InfoMessageVersion(),
		ProducerRef:       producerRef,
	}
	situationEvent.SetId(model.SituationUpdateRequestId(builder.NewUUID()))

	situationEvent.Format = xmlGeneralMessageEvent.FormatRef()
	situationEvent.Keywords = append(situationEvent.Keywords, xmlGeneralMessageEvent.InfoChannelRef())
	situationEvent.ReportType = builder.setReportType(xmlGeneralMessageEvent.InfoChannelRef())

	timeRange := &model.TimeRange{
		StartTime: xmlGeneralMessageEvent.RecordedAtTime(),
		EndTime:   xmlGeneralMessageEvent.ValidUntilTime(),
	}
	situationEvent.ValidityPeriods = []*model.TimeRange{timeRange}

	content := xmlGeneralMessageEvent.Content().(sxml.IDFGeneralMessageStructure)

	builder.buildSituationAndDescriptionFromMessages(content.Messages(), situationEvent)

	builder.setAffects(situationEvent, &content)

	*event = append(*event, situationEvent)
}

func (builder *GeneralMessageUpdateEventBuilder) buildSituationAndDescriptionFromMessages(messages []*sxml.XMLMessage, event *model.SituationUpdateEvent) {
	for _, xmlMessage := range messages {
		builder.buildSituationAndDescriptionFromMessage(xmlMessage.MessageType(), xmlMessage.MessageText(), event)
	}
}

func (builder *GeneralMessageUpdateEventBuilder) buildSituationAndDescriptionFromMessage(messageType, messageText string, event *model.SituationUpdateEvent) {
	switch messageType {
	case "shortMessage":
		event.Summary = &model.SituationTranslatedString{
			DefaultValue: messageText,
		}
	case "longMessage":
		event.Description = &model.SituationTranslatedString{
			DefaultValue: messageText,
		}
	default:
		if event.Summary == nil && len(messageText) < 160 {
			event.Summary = &model.SituationTranslatedString{
				DefaultValue: messageText,
			}
		} else {
			event.Description = &model.SituationTranslatedString{
				DefaultValue: messageText,
			}
		}
	}
}

func (builder *GeneralMessageUpdateEventBuilder) setReportType(infoChannelRef string) model.ReportType {
	switch infoChannelRef {
	case "Perturbation":
		return model.SituationReportTypeIncident
	default:
		return model.SituationReportTypeGeneral
	}

}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedStopArea(event *model.SituationUpdateEvent, stopPointRef string) {
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

func (builder *GeneralMessageUpdateEventBuilder) setAffectedLine(lineRef string) {
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

func (builder *GeneralMessageUpdateEventBuilder) setAffectedRoute(lineId model.LineId, route string) {
	affectedRoute := model.AffectedRoute{RouteRef: route}
	builder.affectedLines[model.LineId(lineId)].AffectedRoutes =
		append(builder.affectedLines[model.LineId(lineId)].AffectedRoutes, &affectedRoute)
}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedDestination(lineId model.LineId, destination string) {
	destinationObjectId := model.NewObjectID(builder.remoteObjectidKind, destination)
	stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(destinationObjectId)
	if !ok {
		return
	}

	affectedDestination := model.AffectedDestination{StopAreaId: stopArea.Id()}
	builder.affectedLines[model.LineId(lineId)].AffectedDestinations =
		append(builder.affectedLines[model.LineId(lineId)].AffectedDestinations, &affectedDestination)

	// Logging
	builder.MonitoringRefs[destinationObjectId.Value()] = struct{}{}
}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedSection(section LineSection) {
	LineRefObjectId := model.NewObjectID(builder.remoteObjectidKind, section.LineRef)
	line, ok := builder.partner.Model().Lines().FindByObjectId(LineRefObjectId)
	if !ok {
		return
	}

	firstStopRef := section.FirstStop
	firstStopObjectId := model.NewObjectID(builder.remoteObjectidKind, firstStopRef)
	firstStopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(firstStopObjectId)
	if !ok {
		return
	}
	lastStopRef := section.LastStop
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

func (builder *GeneralMessageUpdateEventBuilder) setAffects(event *model.SituationUpdateEvent, content *sxml.IDFGeneralMessageStructure) {

	for _, lineRef := range content.LineRef() {
		builder.setAffectedLine(lineRef)
	}

	if len(builder.affectedLines) == 1 {
		// get the LineId
		lineId := maps.Keys(builder.affectedLines)[0]

		for _, destination := range content.DestinationRef() {
			builder.setAffectedDestination(lineId, destination)
		}
		for _, route := range content.RouteRef() {
			builder.setAffectedRoute(lineId, route)
		}
	}

	for _, section := range content.LineSections() {
		lineSection := LineSection{
			LineRef:   section.LineRef(),
			FirstStop: section.FirstStop(),
			LastStop:  section.LastStop(),
		}

		builder.setAffectedSection(lineSection)
	}

	// Fill affectedLines
	for _, affectedLine := range builder.affectedLines {
		event.Affects = append(event.Affects, affectedLine)
	}

	for _, stopPointRef := range content.StopPointRef() {
		builder.setAffectedStopArea(event, stopPointRef)
	}

}
