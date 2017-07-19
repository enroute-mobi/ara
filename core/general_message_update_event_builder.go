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

func newGeneralMessageUpdateEventBuilder(partner *Partner) GeneralMessageUpdateEventBuilder {
	return GeneralMessageUpdateEventBuilder{partner: partner}
}

func (builder *GeneralMessageUpdateEventBuilder) buildGeneralMessageUpdateEvent(event *[]*model.SituationUpdateEvent, xmlGeneralMessageEvent *siri.XMLGeneralMessage, producerRef string) {
	situationEvent := model.SituationUpdateEvent{
		CreatedAt:         builder.Clock().Now(),
		RecordedAt:        builder.Clock().Now(),
		SituationObjectID: model.NewObjectID(builder.partner.Setting("remote_objectid_kind"), xmlGeneralMessageEvent.InfoMessageIdentifier()),
		Version:           int64(xmlGeneralMessageEvent.InfoMessageVersion()),
		ProducerRef:       producerRef,
	}
	situationEvent.SetId(model.SituationUpdateRequestId(builder.NewUUID()))
	if xmlGeneralMessageEvent.Content() != nil {
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
	}
	situationEvent.SituationAttributes.Format = xmlGeneralMessageEvent.FormatRef()
	situationEvent.SituationAttributes.Channel = xmlGeneralMessageEvent.InfoChannelRef()
	situationEvent.SituationAttributes.ValidUntil = xmlGeneralMessageEvent.ValidUntilTime()

	*event = append(*event, &situationEvent)
}

func (builder *GeneralMessageUpdateEventBuilder) SetGeneralMessageUpdateEvents(event *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageDelivery) {
	xmlGeneralMessageEvents := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessageEvents) == 0 {
		return
	}

	for _, xmlGeneralMessageEvents := range xmlGeneralMessageEvents {
		builder.buildGeneralMessageUpdateEvent(event, xmlGeneralMessageEvents, xmlResponse.ProducerRef())
	}
}
