package core

import (
	"regexp"
	"strings"

	"github.com/af83/edwig/model"
)

var defaultIdentifierGenerators = map[string]string{
	"message_identifier":             "%{uuid}",
	"response_message_identifier":    "%{uuid}",
	"data_frame_identifier":          "%{id}",
	"reference_identifier":           "%{type}:%{default}",
	"reference_stop_area_identifier": "%{default}",
	"subscription_identifier":        "%{id}",
}

type IdentifierGenerator struct {
	model.UUIDConsumer

	formatString string
}

type IdentifierAttributes struct {
	Default string
	Id      string
	Type    string
}

func NewIdentifierGenerator(formatString string) *IdentifierGenerator {
	return &IdentifierGenerator{formatString: formatString}
}

func NewIdentifierGeneratorWithUUID(formatString string, uuidGenerator model.UUIDConsumer) *IdentifierGenerator {
	return &IdentifierGenerator{formatString: formatString, UUIDConsumer: uuidGenerator}
}

func (generator *IdentifierGenerator) NewIdentifier(attributes IdentifierAttributes) string {
	replacer := strings.NewReplacer("%{id}", attributes.Id, "%{type}", attributes.Type, "%{default}", attributes.Default)
	return generator.handleuuids(replacer.Replace(generator.formatString))
}

func (generator *IdentifierGenerator) NewMessageIdentifier() string {
	return generator.handleuuids(generator.formatString)
}

func (generator *IdentifierGenerator) handleuuids(s string) string {
	re := regexp.MustCompile("%{uuid}")
	return re.ReplaceAllStringFunc(s, func(string) string { return generator.NewUUID() })
}
