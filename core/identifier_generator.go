package core

import (
	"regexp"
	"strings"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

var defaultIdentifierGenerators = map[string]string{
	"message_identifier":             "%{uuid}",
	"response_message_identifier":    "%{uuid}",
	"data_frame_identifier":          "%{id}",
	"reference_identifier":           "%{type}:%{id}",
	"reference_stop_area_identifier": "%{id}",
	"subscription_identifier":        "%{id}",
}

type IdentifierGenerator struct {
	uuid.UUIDConsumer

	formatString string
}

type IdentifierAttributes struct {
	Type string
	Id   string
}

func DefaultIdentifierGenerator(k string) string {
	return defaultIdentifierGenerators[k]
}

func NewIdentifierGenerator(formatString string, uuidGenerator uuid.UUIDConsumer) *IdentifierGenerator {
	return &IdentifierGenerator{
		UUIDConsumer: uuidGenerator,
		formatString: formatString,
	}
}

func (generator *IdentifierGenerator) NewIdentifier(attributes IdentifierAttributes) string {
	// default and objectid are legacy values, keep them for now for a smoother transition
	replacer := strings.NewReplacer("%{id}", attributes.Id, "%{type}", attributes.Type, "%{default}", attributes.Id, "%{objectid}", attributes.Id)
	return generator.handleuuids(replacer.Replace(generator.formatString))
}

func (generator *IdentifierGenerator) NewMessageIdentifier() string {
	return generator.handleuuids(generator.formatString)
}

func (generator *IdentifierGenerator) handleuuids(s string) string {
	re := regexp.MustCompile("%{uuid}")
	return re.ReplaceAllStringFunc(s, func(string) string { return generator.NewUUID() })
}
