package core

import (
	"fmt"

	"github.com/af83/edwig/model"
)

type MessageIdentifierGenerator interface {
	NewMessageIdentifier() string
}

var defaultMessageIdentifierGenerator MessageIdentifierGenerator = NewFormatMessageIdentifierGenerator("RATPDev:Message::%s:LOC")

func DefaultMessageIdentifierGenerator() MessageIdentifierGenerator {
	return defaultMessageIdentifierGenerator
}

type FormatMessageIdentifierGenerator struct {
	model.UUIDConsumer

	format string
}

func NewFormatMessageIdentifierGenerator(format string) *FormatMessageIdentifierGenerator {
	return &FormatMessageIdentifierGenerator{format: format}
}

func (generator *FormatMessageIdentifierGenerator) NewMessageIdentifier() string {
	return fmt.Sprintf(generator.format, generator.NewUUID())
}

type MessageIdentifierConsumer struct {
	messageIdentifierGenerator MessageIdentifierGenerator
}

func (consumer *MessageIdentifierConsumer) SetMessageIdentifierGenerator(generator MessageIdentifierGenerator) {
	consumer.messageIdentifierGenerator = generator
}

func (consumer *MessageIdentifierConsumer) MessageIdentifierGenerator() MessageIdentifierGenerator {
	if consumer.messageIdentifierGenerator == nil {
		consumer.messageIdentifierGenerator = DefaultMessageIdentifierGenerator()
	}
	return consumer.messageIdentifierGenerator
}

func (consumer *MessageIdentifierConsumer) NewMessageIdentifier() string {
	return consumer.MessageIdentifierGenerator().NewMessageIdentifier()
}

type ResponseMessageIdentifierConsumer struct {
	messageIdentifierGenerator MessageIdentifierGenerator
}

func (consumer *ResponseMessageIdentifierConsumer) SetResponseMessageIdentifierGenerator(generator MessageIdentifierGenerator) {
	consumer.messageIdentifierGenerator = generator
}

func (consumer *ResponseMessageIdentifierConsumer) ResponseMessageIdentifierGenerator() MessageIdentifierGenerator {
	if consumer.messageIdentifierGenerator == nil {
		consumer.messageIdentifierGenerator = NewFormatMessageIdentifierGenerator("RATPDev:ResponseMessage::%s:LOC")
	}
	return consumer.messageIdentifierGenerator
}

func (consumer *ResponseMessageIdentifierConsumer) NewResponseMessageIdentifier() string {
	return consumer.ResponseMessageIdentifierGenerator().NewMessageIdentifier()
}
