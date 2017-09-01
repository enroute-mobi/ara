package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type GeneralMessageUpdateEventBuilder struct {
	model.ClockConsumer
	model.UUIDConsumer

	partner *Partner
}

func NewGeneralMessageUpdateEventBuilder(partner *Partner) GeneralMessageUpdateEventBuilder {
	return GeneralMessageUpdateEventBuilder{partner: partner}
}

func (builder *GeneralMessageUpdateEventBuilder) SetGeneralMessageDeliveryUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageDelivery) {
	xmlGeneralMessageEvents := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessageEvents) == 0 {
		return
	}

	for _, xmlGeneralMessageEvents := range xmlGeneralMessageEvents {
		builder.buildGeneralMessageUpdateEvent(event, xmlGeneralMessageEvents, xmlResponse.ProducerRef())
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
		CreatedAt:         builder.Clock().Now(),
		RecordedAt:        xmlGeneralMessageEvent.RecordedAtTime(),
		SituationObjectID: model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlGeneralMessageEvent.InfoMessageIdentifier()),
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

	remoteObjectidKind := builder.partner.Setting("remote_objectid_kind")

	if lineRefs := content.LineRef(); len(lineRefs) != 0 {
		for i, _ := range lineRefs {
			ref := model.NewReference(model.NewObjectID(remoteObjectidKind, lineRefs[i]))
			ref.Type = "LineRef"
			event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
		}
	}
	if stopPointRefs := content.StopPointRef(); len(stopPointRefs) != 0 {
		for i, _ := range stopPointRefs {
			ref := model.NewReference(model.NewObjectID(remoteObjectidKind, stopPointRefs[i]))
			ref.Type = "StopPointRef"
			event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
		}
	}
	if journeyPatternRefs := content.JourneyPatternRef(); len(journeyPatternRefs) != 0 {
		for i, _ := range journeyPatternRefs {
			ref := model.NewReference(model.NewObjectID(remoteObjectidKind, journeyPatternRefs[i]))
			ref.Type = "JourneyPatternRef"
			event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
		}
	}
	if destinationRefs := content.DestinationRef(); len(destinationRefs) != 0 {
		for i, _ := range destinationRefs {
			ref := model.NewReference(model.NewObjectID(remoteObjectidKind, destinationRefs[i]))
			ref.Type = "DestinationRef"
			event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
		}
	}
	if routeRefs := content.RouteRef(); len(routeRefs) != 0 {
		for i, _ := range routeRefs {
			ref := model.NewReference(model.NewObjectID(remoteObjectidKind, routeRefs[i]))
			ref.Type = "RouteRef"
			event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
		}
	}
	if groupOfLinesRefs := content.GroupOfLinesRef(); len(groupOfLinesRefs) != 0 {
		for i, _ := range groupOfLinesRefs {
			ref := model.NewReference(model.NewObjectID(remoteObjectidKind, groupOfLinesRefs[i]))
			ref.Type = "GroupOfLinesRef"
			event.SituationAttributes.References = append(event.SituationAttributes.References, ref)
		}
	}
	if lineSections := content.LineSections(); len(lineSections) != 0 {
		for _, lineSection := range lineSections {
			builder.handleLineSection(tx, remoteObjectidKind, lineSection, event)
		}
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
