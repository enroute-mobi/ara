package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type GeneralMessageUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner            *Partner
	remoteObjectidKind string
}

func NewGeneralMessageUpdateEventBuilder(partner *Partner) GeneralMessageUpdateEventBuilder {
	return GeneralMessageUpdateEventBuilder{
		partner:            partner,
		remoteObjectidKind: partner.RemoteObjectIDKind(),
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
}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedLine(event *model.SituationUpdateEvent, lineRef string) {
	LineRefObjectId := model.NewObjectID(builder.remoteObjectidKind, lineRef)
	line, ok := builder.partner.Model().Lines().FindByObjectId(LineRefObjectId)
	if !ok {
		return
	}
	affect := model.NewAffectedLine()
	affect.LineId = line.Id()

	event.Affects = append(event.Affects, affect)
}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedRoute(event *model.SituationUpdateEvent, route string, affectedLine *model.AffectedLine) {
	affectedRoute := model.AffectedRoute{RouteRef: route}
	affectedLine.AffectedRoutes = append(affectedLine.AffectedRoutes, &affectedRoute)
}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedDestination(event *model.SituationUpdateEvent, destination string, affectedLine *model.AffectedLine) {
	destinationObjectId := model.NewObjectID(builder.remoteObjectidKind, destination)
	stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(destinationObjectId)
	if !ok {
		return
	}

	affectedDestination := model.AffectedDestination{StopAreaId: stopArea.Id()}
	affectedLine.AffectedDestinations = append(affectedLine.AffectedDestinations, &affectedDestination)
}

func (builder *GeneralMessageUpdateEventBuilder) setAffectedSection(event *model.SituationUpdateEvent, section *sxml.IDFLineSectionStructure, affectedLine *model.AffectedLine) {
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
	affectedSection := model.AffectedSection{
		FirstStop: firstStopArea.Id(),
		LastStop:  lastStopArea.Id(),
	}
	affectedLine.AffectedSections = append(affectedLine.AffectedSections, &affectedSection)
}

func (builder *GeneralMessageUpdateEventBuilder) setAffects(event *model.SituationUpdateEvent, content *sxml.IDFGeneralMessageStructure) {

	for _, lineRef := range content.LineRef() {
		builder.setAffectedLine(event, lineRef)
	}

	if len(event.Affects) == 1 && event.Affects[0].GetType() == "Line" {
		for _, destination := range content.DestinationRef() {
			builder.setAffectedDestination(event, destination, event.Affects[0].(*model.AffectedLine))
		}
		for _, section := range content.LineSections() {
			builder.setAffectedSection(event, section, event.Affects[0].(*model.AffectedLine))
		}
		for _, route := range content.RouteRef() {
			builder.setAffectedRoute(event, route, event.Affects[0].(*model.AffectedLine))
		}
	}

	for _, stopPointRef := range content.StopPointRef() {
		builder.setAffectedStopArea(event, stopPointRef)
	}

}
