package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type GeneralMessageUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner *Partner
}

func NewGeneralMessageUpdateEventBuilder(partner *Partner) GeneralMessageUpdateEventBuilder {
	return GeneralMessageUpdateEventBuilder{partner: partner}
}

func (builder *GeneralMessageUpdateEventBuilder) SetGeneralMessageDeliveryUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageDelivery, producerRef string) {
	xmlGeneralMessageEvents := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessageEvents) == 0 {
		return
	}

	for _, xmlGeneralMessageEvents := range xmlGeneralMessageEvents {
		builder.buildGeneralMessageUpdateEvent(event, xmlGeneralMessageEvents, producerRef)
	}
}

func (builder *GeneralMessageUpdateEventBuilder) SetGeneralMessageResponseUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageResponse) {
	xmlGeneralMessageEvents := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessageEvents) == 0 {
		return
	}

	for _, xmlGeneralMessageEvents := range xmlGeneralMessageEvents {
		builder.buildGeneralMessageUpdateEvent(event, xmlGeneralMessageEvents, xmlResponse.ProducerRef())
	}
}

func (builder *GeneralMessageUpdateEventBuilder) buildGeneralMessageUpdateEvent(event *[]*model.SituationUpdateEvent, xmlGeneralMessageEvent *siri.XMLGeneralMessage, producerRef string) {
	if xmlGeneralMessageEvent.Content() == nil {
		return
	}

	situationEvent := &model.SituationUpdateEvent{
		Origin:            string(builder.partner.Slug()),
		CreatedAt:         builder.Clock().Now(),
		RecordedAt:        xmlGeneralMessageEvent.RecordedAtTime(),
		SituationObjectID: model.NewObjectID(builder.partner.RemoteObjectIDKind(), xmlGeneralMessageEvent.InfoMessageIdentifier()),
		Version:           xmlGeneralMessageEvent.InfoMessageVersion(),
		ProducerRef:       producerRef,
	}
	situationEvent.SetId(model.SituationUpdateRequestId(builder.NewUUID()))

	situationEvent.SituationAttributes.Format = xmlGeneralMessageEvent.FormatRef()
	situationEvent.SituationAttributes.Channel = xmlGeneralMessageEvent.InfoChannelRef()
	situationEvent.SituationAttributes.ValidUntil = xmlGeneralMessageEvent.ValidUntilTime()

	content := xmlGeneralMessageEvent.Content().(siri.IDFGeneralMessageStructure)
	for _, xmlMessage := range content.Messages() {
		message := &model.Message{
			Content:             xmlMessage.MessageText(),
			Type:                xmlMessage.MessageType(),
			NumberOfLines:       xmlMessage.NumberOfLines(),
			NumberOfCharPerLine: xmlMessage.NumberOfCharPerLine(),
		}
		situationEvent.SituationAttributes.Messages = append(situationEvent.SituationAttributes.Messages, message)
	}

	builder.setReferences(situationEvent, &content)

	*event = append(*event, situationEvent)
}

func (builder *GeneralMessageUpdateEventBuilder) setReferences(event *model.SituationUpdateEvent, content *siri.IDFGeneralMessageStructure) {
	tx := builder.partner.Referential().NewTransaction()
	defer tx.Close()

	remoteObjectidKind := builder.partner.RemoteObjectIDKind()

	for _, lineref := range content.LineRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, lineref))
		ref.Type = "LineRef"
		event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
	}
	for _, stoppointref := range content.StopPointRef() {
		ref := model.NewReference(model.NewObjectID(remoteObjectidKind, stoppointref))
		ref.Type = "StopPointRef"
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
		builder.handleLineSection(tx, remoteObjectidKind, lineSection, event)
	}
}

func (builder *GeneralMessageUpdateEventBuilder) handleLineSection(tx *model.Transaction, remoteObjectidKind string, lineSection *siri.IDFLineSectionStructure, event *model.SituationUpdateEvent) {
	references := model.NewReferences()

	lineRef := model.NewReference(model.NewObjectID(remoteObjectidKind, lineSection.LineRef()))
	references.Set("LineRef", *lineRef)

	firstStopRef := model.NewReference(model.NewObjectID(remoteObjectidKind, lineSection.FirstStop()))
	references.Set("FirstStop", *firstStopRef)

	lastStopRef := model.NewReference(model.NewObjectID(remoteObjectidKind, lineSection.LastStop()))
	references.Set("LastStop", *lastStopRef)

	event.SituationAttributes.LineSections = append(event.SituationAttributes.LineSections, &references)
}
