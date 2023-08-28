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

	situationEvent.SituationAttributes.Format = xmlGeneralMessageEvent.FormatRef()
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
	builder.setReferences(situationEvent, &content)

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

func (builder *GeneralMessageUpdateEventBuilder) setAffects(event *model.SituationUpdateEvent, content *sxml.IDFGeneralMessageStructure) {
	remoteObjectidKind := builder.remoteObjectidKind
	for _, stoppointref := range content.StopPointRef() {
		stopPointRefObjectId := model.NewObjectID(remoteObjectidKind, stoppointref)
		stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(stopPointRefObjectId)
		if !ok {
			continue
		}
		affect := model.NewAffectedStopArea()
		affect.StopAreaId = stopArea.Id()

		event.Affects = append(event.Affects, affect)
	}
}

func (builder *GeneralMessageUpdateEventBuilder) setReferences(event *model.SituationUpdateEvent, content *sxml.IDFGeneralMessageStructure) {
	remoteObjectidKind := builder.remoteObjectidKind

	for _, lineref := range content.LineRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, lineref))
		ref.Type = "LineRef"
		event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
	}

	for _, journeypatternref := range content.JourneyPatternRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, journeypatternref))
		ref.Type = "JourneyPatternRef"
		event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
	}
	for _, destinationref := range content.DestinationRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, destinationref))
		ref.Type = "DestinationRef"
		event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
	}
	for _, routeref := range content.RouteRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, routeref))
		ref.Type = "RouteRef"
		event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
	}
	for _, groupoflinesref := range content.GroupOfLinesRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, groupoflinesref))
		ref.Type = "GroupOfLinesRef"
		event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
	}
	for _, lineSection := range content.LineSections() {
		builder.handleLineSection(remoteObjectidKind, lineSection, event)
	}
}

func (builder *GeneralMessageUpdateEventBuilder) handleLineSection(remoteObjectidKind string, lineSection *sxml.IDFLineSectionStructure, event *model.SituationUpdateEvent) {
	references := model.NewReferences()

	lineRef := model.NewReference(model.NewObjectID(remoteObjectidKind, lineSection.LineRef()))
	references.Set("LineRef", *lineRef)

	firstStopRef := model.NewReference(model.NewObjectID(remoteObjectidKind, lineSection.FirstStop()))
	references.Set("FirstStop", *firstStopRef)

	lastStopRef := model.NewReference(model.NewObjectID(remoteObjectidKind, lineSection.LastStop()))
	references.Set("LastStop", *lastStopRef)

	event.SituationAttributes.LineSections = append(event.SituationAttributes.LineSections, &references)
}
