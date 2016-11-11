package model

import (
	"fmt"
)

type MessageIdentifierGenerator interface {
	NewMessageIdentifier() string
}

var defaultMessageIdentifierGenerator MessageIdentifierGenerator = NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")

func DefaultMessageIdentifierGenerator() MessageIdentifierGenerator {
	return defaultMessageIdentifierGenerator
}

type FormatMessageIdentifierGenerator struct {
	UUIDConsumer

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
