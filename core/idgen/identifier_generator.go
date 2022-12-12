package idgen

import (
	"regexp"
	"strings"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

const (
	MESSAGE_IDENTIFIER             = "message_identifier"
	RESPONSE_MESSAGE_IDENTIFIER    = "response_message_identifier"
	DATA_FRAME_IDENTIFIER          = "data_frame_identifier"
	REFERENCE_IDENTIFIER           = "reference_identifier"
	REFERENCE_STOP_AREA_IDENTIFIER = "reference_stop_area_identifier"
	SUBSCRIPTION_IDENTIFIER        = "subscription_identifier"
)

var re = regexp.MustCompile("%{uuid}")

var defaultIdentifierGenerators = map[string]string{
	MESSAGE_IDENTIFIER:             "%{uuid}",
	RESPONSE_MESSAGE_IDENTIFIER:    "%{uuid}",
	DATA_FRAME_IDENTIFIER:          "%{id}",
	REFERENCE_IDENTIFIER:           "%{type}:%{id}",
	REFERENCE_STOP_AREA_IDENTIFIER: "%{id}",
	SUBSCRIPTION_IDENTIFIER:        "%{id}",
}

type IdentifierGenerator struct {
	uuidGenerator uuid.UUIDGenerator
	formatString  string
}

type IdentifierAttributes struct {
	Type string
	Id   string
}

func DefaultIdentifierGenerator(k string) string {
	return defaultIdentifierGenerators[k]
}

func NewIdentifierGenerator(formatString string, uuidGenerator uuid.UUIDGenerator) *IdentifierGenerator {
	return &IdentifierGenerator{
		uuidGenerator: uuidGenerator,
		formatString:  formatString,
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

func (generator *IdentifierGenerator) FormatString() string {
	return generator.formatString
}

func (generator *IdentifierGenerator) handleuuids(s string) string {
	return re.ReplaceAllStringFunc(s, func(string) string { return generator.uuidGenerator.NewUUID() })
}
